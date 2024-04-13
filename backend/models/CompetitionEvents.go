package models

import "github.com/jackc/pgx/v5/pgxpool"

type CompetitionEvents struct {
	Id int
	Competition_id string
	Event_id int
}

func UpdateCompetitionEvents(competition *CompetitionData, db *pgxpool.Pool) error {
	if err := competition.RemoveAllEvents(db); err != nil { return err }
	if err := competition.AddEvents(db); err != nil { return err }
	return nil;
}