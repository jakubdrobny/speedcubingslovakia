package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CompetitionEvent struct {
	Id int `json:"id"`
	Displayname string `json:"displayname"`
	Format string `json:"format"`
	Iconcode string `json:"iconcode"`
	Puzzlecode string `json:"puzzlecode"`
}

func GetEventById(db *pgxpool.Pool, eventId int) (CompetitionEvent, error) {
	rows, err := db.Query(context.Background(), "SELECT e.event_id, e.displayname, e.format, e.iconcode, e.puzzlecode FROM events e WHERE e.event_id = $1;", eventId);
	if err != nil { return CompetitionEvent{}, err }

	var event CompetitionEvent
	found := false
	for rows.Next() {
		err = rows.Scan(&event.Id, &event.Displayname, &event.Format, &event.Iconcode, &event.Puzzlecode)
		if err != nil { return CompetitionEvent{}, err }
		found = true
	}

	if !found { return CompetitionEvent{}, fmt.Errorf("event not found by id") }

	return event, nil
}

func GetAvailableEvents(db *pgxpool.Pool) ([]CompetitionEvent, error) {
	rows, err := db.Query(context.Background(), "SELECT e.event_id, e.displayname, e.format, e.iconcode, e.puzzlecode FROM events e ORDER BY e.event_id;");
	if err != nil { return []CompetitionEvent{}, err }

	var events []CompetitionEvent
	for rows.Next() {
		var event CompetitionEvent
		err = rows.Scan(&event.Id, &event.Displayname, &event.Format, &event.Iconcode, &event.Puzzlecode)
		if err != nil { return []CompetitionEvent{}, err }
		events = append(events, event)
	}

	return events, nil
}
