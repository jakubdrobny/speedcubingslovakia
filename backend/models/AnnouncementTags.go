package models

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateAnnouncementTags(announcement *AnnouncementState, db *pgxpool.Pool, tx pgx.Tx, envMap map[string]string) error {
	var err error
	var tag_ids []int

	if tag_ids, err = announcement.RemoveAllTags(db, tx); err != nil {
		return err
	}
	if err := announcement.AddTags(tx, tag_ids); err != nil {
		return err
	}

	return nil
}
