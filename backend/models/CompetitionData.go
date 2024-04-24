package models

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexsergivan/transliterator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompetitionData struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Startdate time.Time `json:"startdate"`
	Enddate time.Time `json:"enddate"`
	Events []CompetitionEvent `json:"events"`
	Scrambles []ScrambleSet `json:"scrambles"`
	Results ResultEntry `json:"results"`
}

func (c *CompetitionData) RemoveAllEvents(db *pgxpool.Pool, tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), `DELETE FROM competition_events WHERE competition_id = $1;`, c.Id)
	return err
}

func (c *CompetitionData) AddEvents(db *pgxpool.Pool, tx pgx.Tx) error {
	for _, event := range c.Events {
		_, err := tx.Exec(context.Background(), `INSERT INTO competition_events (competition_id, event_id) VALUES ($1, $2);`, c.Id, event.Id)
		if err != nil { return err }
	}

	return nil;
}

func (competition *CompetitionData) RecomputeCompetitionId() {
	trans := transliterator.NewTransliterator(nil)
	competition.Id = trans.Transliterate(strings.Join(strings.Split(competition.Name, " "), ""), "")
}

func (c *CompetitionData) GetScrambles(db *pgxpool.Pool) (error) {
	scrambleSets := make([]ScrambleSet, 0)

	for _, event := range c.Events {
		rows, err := db.Query(context.Background(), `SELECT s.scramble, e.event_id, e.displayname, e.format, e.iconcode, e.puzzlecode FROM scrambles s LEFT JOIN events e ON s.event_id = e.event_id WHERE s.competition_id = $1 AND s.event_id = $2 ORDER BY e.event_id, s."order";`, c.Id, event.Id)
		if err != nil { return err }

		var scrambleSet ScrambleSet
		for rows.Next() {
			var scramble string
			err := rows.Scan(&scramble, &scrambleSet.Event.Id, &scrambleSet.Event.Displayname, &scrambleSet.Event.Format, &scrambleSet.Event.Iconcode, &scrambleSet.Event.Puzzlecode)
			if err != nil { return err }
			scrambleSet.AddScramble(scramble)
		}

		scrambleSets = append(scrambleSets, scrambleSet)
	}

	c.Scrambles = scrambleSets

	return nil
}

func (c *CompetitionData) GetEvents(db *pgxpool.Pool) (error) {
	events := make([]CompetitionEvent, 0)
	
	rows, err := db.Query(context.Background(), `SELECT e.event_id, e.displayname, e.format, e.iconcode, e.puzzlecode FROM competition_events ce JOIN events e ON ce.event_id = e.event_id WHERE ce.competition_id = $1 ORDER BY e.event_id`, c.Id)
	if err != nil { return err }
	
	for rows.Next() {
		var event CompetitionEvent
		err := rows.Scan(&event.Id, &event.Displayname, &event.Format, &event.Iconcode, &event.Puzzlecode)
		if err != nil { return err }
		events = append(events, event)
	}
	events = append(events, CompetitionEvent{-1, "Overall", "", "", ""})

	c.Events = events

	return nil
}

func GetCompetitionByIdObject(db *pgxpool.Pool, id string) (CompetitionData, error) {
	rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c WHERE c.competition_id = $1;`, id)
	if err != nil { return CompetitionData{}, err }

	var competition CompetitionData
	found := false
	for rows.Next() {
		err = rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
		if err != nil { return CompetitionData{}, err }
		found = true
	}

	if !found { return CompetitionData{}, err }

	return competition, nil
}


func GetAllCompetitions(db *pgxpool.Pool) ([]CompetitionData, error) {
	rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c;`)
	if err != nil { return []CompetitionData{}, err }

	competitions := make([]CompetitionData, 0)

	for rows.Next() {
		var competition CompetitionData
		err = rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
		if err != nil { return []CompetitionData{}, err }
		competition.GetEvents(db)
		competitions = append(competitions, competition)
	}

	return competitions, nil
}

func (c *CompetitionData) GenerateScrambles() (error) {
	req, err := http.NewRequest("POST", "http://localhost:2014/frontend/puzzle/pyram/scramble", bytes.NewBuffer([]byte("{\"B\":\"#0000ff\",\"R\":\"#ff0000\",\"D\":\"#ffff00\",\"U\":\"#ffffff\",\"F\":\"#00ff00\",\"L\":\"#ff8000\"}")))
	if err != nil { return err }

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, */*")

    resp, err := http.DefaultClient.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))

	return nil
}