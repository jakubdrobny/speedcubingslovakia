package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompetitionEvent struct {
	Id              int    `json:"id"`
	Fulldisplayname string `json:"fulldisplayname"`
	Displayname     string `json:"displayname"`
	Format          string `json:"format"`
	Iconcode        string `json:"iconcode"`
	Scramblingcode  string `json:"scramblingcode"`
}

func GetEventById(db *pgxpool.Pool, eventId int) (CompetitionEvent, error) {
	rows, err := db.Query(context.Background(), "SELECT e.event_id, e.displayname, e.format, e.iconcode, e.scramblingcode FROM events e WHERE e.event_id = $1;", eventId)
	if err != nil {
		return CompetitionEvent{}, err
	}

	var event CompetitionEvent
	found := false
	for rows.Next() {
		err = rows.Scan(&event.Id, &event.Displayname, &event.Format, &event.Iconcode, &event.Scramblingcode)
		if err != nil {
			return CompetitionEvent{}, err
		}
		found = true
	}

	if !found {
		return CompetitionEvent{}, fmt.Errorf("event not found by id")
	}

	return event, nil
}

func GetAvailableEvents(db *pgxpool.Pool) ([]CompetitionEvent, error) {
	rows, err := db.Query(context.Background(), "SELECT e.event_id, e.fulldisplayname, e.displayname, e.format, e.iconcode, e.scramblingcode FROM events e ORDER BY e.event_id;")
	if err != nil {
		return []CompetitionEvent{}, err
	}

	var events []CompetitionEvent
	for rows.Next() {
		var event CompetitionEvent
		err = rows.Scan(&event.Id, &event.Fulldisplayname, &event.Displayname, &event.Format, &event.Iconcode, &event.Scramblingcode)
		if err != nil {
			return []CompetitionEvent{}, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (e *CompetitionEvent) HasScrambles(db *pgxpool.Pool, tx pgx.Tx, cid string) (bool, error) {
	rows, err := tx.Query(context.Background(), `SELECT COUNT(*) FROM scrambles WHERE competition_id = $1 AND event_id = $2;`, cid, e.Id)
	if err != nil {
		return true, err
	}

	var cnt int
	for rows.Next() {
		err = rows.Scan(&cnt)
		if err != nil {
			return true, err
		}
	}

	return cnt > 0, nil
}
