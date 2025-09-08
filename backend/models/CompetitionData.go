package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alexsergivan/transliterator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type CompetitionData struct {
	Id        string             `json:"id"`
	Name      string             `json:"name"`
	Startdate time.Time          `json:"startdate"`
	Enddate   time.Time          `json:"enddate"`
	Events    []CompetitionEvent `json:"events"`
	Scrambles []ScrambleSet      `json:"scrambles"`
	Results   ResultEntry        `json:"results"`
}

func (c *CompetitionData) RemoveAllEvents(db *pgxpool.Pool, tx pgx.Tx) ([]int, error) {
	rows, err := tx.Query(context.Background(), `SELECT event_id FROM competition_events WHERE competition_id = $1;`, c.Id)
	if err != nil {
		return []int{}, err
	}

	event_ids := make([]int, 0)
	for rows.Next() {
		var event_id int
		err = rows.Scan(&event_id)
		if err != nil {
			return []int{}, err
		}

		event_ids = append(event_ids, event_id)
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM competition_events WHERE competition_id = $1;`, c.Id)
	return event_ids, err
}

func (c *CompetitionData) AddEvents(db *pgxpool.Pool, tx pgx.Tx, event_ids []int, envMap map[string]string) error {
	for _, event := range c.Events {
		if event.Id < 0 {
			continue
		}
		_, err := tx.Exec(context.Background(), `INSERT INTO competition_events (competition_id, event_id) VALUES ($1, $2);`, c.Id, event.Id)
		if err != nil {
			return err
		}

		has, err := event.HasScrambles(db, tx, c.Id)
		if err != nil {
			return err
		}

		if !has {
			noOfSolves, err := utils.GetNoOfSolves(event.Format)
			if err != nil {
				return err
			}

			ismbld := event.Iconcode == "333mbf"

			scrambles, err := GenerateScramblesForEvent(event.Scramblingcode, noOfSolves, ismbld, envMap)
			if err != nil {
				return err
			}

			images, err := GenerateImagesForScrambles(scrambles, event.Scramblingcode, ismbld, envMap)
			if err != nil {
				return err
			}

			for scrambleIdx, scramble := range scrambles {
				_, err := tx.Exec(context.Background(), `INSERT INTO scrambles (scramble, event_id, competition_id, "order", img) VALUES ($1,$2,$3,$4,$5);`, scramble, event.Id, c.Id, scrambleIdx+1, images[scrambleIdx])
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (competition *CompetitionData) RecomputeCompetitionId() {
	trans := transliterator.NewTransliterator(nil)
	competition.Id = trans.Transliterate(strings.Join(strings.Split(competition.Name, " "), ""), "")
}

func (c *CompetitionData) GetScrambles(db *pgxpool.Pool) error {
	scrambleSets := make([]ScrambleSet, 0)

	for _, event := range c.Events {
		rows, err := db.Query(context.Background(), `SELECT s.scramble_id, s.scramble, e.event_id, e.displayname, e.format, e.iconcode, e.scramblingcode, s.img FROM scrambles s LEFT JOIN events e ON s.event_id = e.event_id WHERE s.competition_id = $1 AND s.event_id = $2 ORDER BY e.event_id, s."order";`, c.Id, event.Id)
		if err != nil {
			return err
		}

		var scrambleSet ScrambleSet
		for rows.Next() {
			var scramble Scramble
			var scrambleId int
			err := rows.Scan(&scrambleId, &scramble.Scramble, &scrambleSet.Event.Id, &scrambleSet.Event.Displayname, &scrambleSet.Event.Format, &scrambleSet.Event.Iconcode, &scrambleSet.Event.Scramblingcode, &scramble.Img)
			if err != nil {
				return err
			}

			if time.Now().Before(c.Startdate) {
				scramble.Scramble = "Competition has not started yet ;)"
				scramble.Img = ""
			}

			scrambleSet.AddScramble(scramble)
		}

		scrambleSets = append(scrambleSets, scrambleSet)
	}

	c.Scrambles = scrambleSets

	return nil
}

func (c *CompetitionData) GetEvents(db *pgxpool.Pool) error {
	events := make([]CompetitionEvent, 0)

	rows, err := db.Query(context.Background(), `SELECT e.event_id, e.displayname, e.format, e.iconcode, e.scramblingcode FROM competition_events ce JOIN events e ON ce.event_id = e.event_id WHERE ce.competition_id = $1 ORDER BY e.event_id`, c.Id)
	if err != nil {
		return err
	}

	for rows.Next() {
		var event CompetitionEvent
		err := rows.Scan(&event.Id, &event.Displayname, &event.Format, &event.Iconcode, &event.Scramblingcode)
		if err != nil {
			return err
		}
		events = append(events, event)
	}
	events = append(events, CompetitionEvent{-1, "", "Overall", "", "", ""})

	c.Events = events

	return nil
}

func GetCompetitionByIdObject(db *pgxpool.Pool, id string) (CompetitionData, error) {
	rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c WHERE c.competition_id = $1;`, id)
	if err != nil {
		return CompetitionData{}, err
	}

	var competition CompetitionData
	found := false
	for rows.Next() {
		err = rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
		if err != nil {
			return CompetitionData{}, err
		}
		found = true
	}

	if !found {
		return CompetitionData{}, err
	}

	return competition, nil
}

func GetAllCompetitions(db *pgxpool.Pool) ([]CompetitionData, error) {
	rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c;`)
	if err != nil {
		return []CompetitionData{}, err
	}

	competitions := make([]CompetitionData, 0)

	for rows.Next() {
		var competition CompetitionData
		err = rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
		if err != nil {
			return []CompetitionData{}, err
		}
		competition.GetEvents(db)
		competitions = append(competitions, competition)
	}

	return competitions, nil
}

func GetScrambles(scramblingcode string, noOfSolves int, envMap map[string]string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v0/scramble/%s?numScrambles=%d", envMap["SCRAMBLING_SERVICE_URL"], scramblingcode, noOfSolves)
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	scrambles := make([]string, 0)
	json.Unmarshal(respBody, &scrambles)

	return scrambles, nil
}

func GenerateScramblesForEvent(scramblingcode string, noOfSolves int, ismbld bool, envMap map[string]string) ([]string, error) {
	if !ismbld {
		return GetScrambles(scramblingcode, noOfSolves, envMap)
	}

	scrambles := make([]string, 0)
	for range noOfSolves {
		currentScrambles, err := GetScrambles(scramblingcode, constants.MBLD_MAX_CUBES_PER_ATTEMPT, envMap)
		if err != nil {
			return []string{}, err
		}

		scrambles = append(scrambles, strings.Join(currentScrambles, "\n"))
	}

	return scrambles, nil
}

func GenerateImagesForScrambles(scrambles []string, scramblingcode string, ismbld bool, envMap map[string]string) ([]string, error) {
	images := make([]string, 0)

	for _, scramble := range scrambles {
		if ismbld {
			scramble = ""
		}
		url := fmt.Sprintf("%s/api/v0/view/%s/svg?scramble=%s", envMap["SCRAMBLING_SERVICE_URL"], scramblingcode, url.QueryEscape(scramble))
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []string{}, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return []string{}, err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, err
		}

		img_id := utils.RandSeq(64) + ".svg" // with extension
		err = utils.SaveScrambleImg(img_id, string(respBody))
		if err != nil {
			return []string{}, err
		}

		images = append(images, img_id)
	}

	return images, nil
}

func (c *CompetitionData) GenerateScrambles(envMap map[string]string) error {
	for _, event := range c.Events {
		noOfSolves, err := utils.GetNoOfSolves(event.Format)
		if err != nil {
			return err
		}

		ismbld := event.Iconcode == "333mbf"

		fmt.Printf("GenerateScrambles:dbg1\nevent:%+v\n", event)

		scrambles, err := GenerateScramblesForEvent(event.Scramblingcode, noOfSolves, ismbld, envMap)
		if err != nil {
			return err
		}

		fmt.Println("GenerateScrambles:dbg2")
		images, err := GenerateImagesForScrambles(scrambles, event.Scramblingcode, ismbld, envMap)
		if err != nil {
			return err
		}

		fmt.Println("GenerateScrambles:dbg3")
		var scrambleSet ScrambleSet
		scrambleSet.Event = event
		for idx, scrambleText := range scrambles {
			var scramble Scramble
			scramble.Scramble = scrambleText
			scramble.Img = images[idx]
			scrambleSet.Scrambles = append(scrambleSet.Scrambles, scramble)
		}
		c.Scrambles = append(c.Scrambles, scrambleSet)
	}

	return nil
}
