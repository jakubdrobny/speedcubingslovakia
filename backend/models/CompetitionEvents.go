package models

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompetitionEvents struct {
	Id int
	Competition_id string
	Event_id int
}

func UpdateCompetitionEvents(competition *CompetitionData, db *pgxpool.Pool, tx pgx.Tx, envMap map[string]string) error {
	var err error
	var event_ids []int

	if event_ids, err = competition.RemoveAllEvents(db, tx); err != nil { return err }
	if err := competition.AddEvents(db, tx, event_ids, envMap); err != nil { return err }

	return nil;
}