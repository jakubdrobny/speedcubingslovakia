package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmojiCounter struct {
	Id int `json:"-"`
	AnnouncementId int `json:"announcementId"`
	Emoji string `json:"emoji"`
	By string `json:"by"`
	Set bool `json:"-"`
}

// checks for (emoji, by) combination
func (ec *EmojiCounter) Exists(db *pgxpool.Pool) (bool, error) {
	rows, err := db.Query(context.Background(), `SELECT ar.announcement_reaction_id, ar."set" FROM announcement_reaction ar WHERE ar.emoji = $1 AND ar."by" = $2 AND ar.announcement_id = $3;`, ec.Emoji, ec.By, ec.AnnouncementId)
	if err != nil { return false, err }

	found := false

	if rows.Next() {
		err = rows.Scan(&ec.Id, &ec.Set)
		if err != nil { return false, err }

		found = true
	}

	return found, nil
}

func (ec *EmojiCounter) Update(db *pgxpool.Pool, userId int) error {
	exists, err := ec.Exists(db)
	if err != nil { return err }

	if !exists {
		_, err = db.Exec(context.Background(), `INSERT INTO announcement_reaction (announcement_id, user_id, emoji, "by", "set") VALUES ($1,$2,$3,$4,TRUE);`, ec.AnnouncementId, userId, ec.Emoji, ec.By)
		ec.Set = true
	} else {
		ec.Set = !ec.Set
		_, err = db.Exec(context.Background(), `UPDATE announcement_reaction SET "set" = $1 WHERE emoji = $2 AND "by" = $3;`, ec.Set, ec.Emoji, ec.By)
	}

	if err != nil { return err }

	return nil
}