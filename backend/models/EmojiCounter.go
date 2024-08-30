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

func (ec *EmojiCounter) Update(conn *pgxpool.Conn, userId int) error {
	err := conn.QueryRow(context.Background(), `INSERT INTO announcement_reaction (announcement_id, user_id, emoji, "by", "set") VALUES ($1,$2,$3,$4,TRUE) ON CONFLICT (announcement_id, user_id, emoji) DO UPDATE SET "set" = not announcement_reaction."set" RETURNING "set";`, ec.AnnouncementId, userId, ec.Emoji, ec.By).Scan(&ec.Set)
	if err != nil { return err }

	return nil
}