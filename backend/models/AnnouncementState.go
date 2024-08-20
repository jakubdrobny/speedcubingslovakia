package models

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnnouncementState struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	AuthorId int `json:"authorId"`
	AuthorWcaId string `json:"authorWcaId"`
	AuthorUsername string `json:"authorUsername"`
	Tags []Tag `json:"tags"`
}

func (a *AnnouncementState) GetTags(db *pgxpool.Pool) (error) {
	rows, err := db.Query(context.Background(), `SELECT t.tag_id, t.label, t.color FROM tags t JOIN announcement_tags at ON t.tag_id = at.tag_id WHERE at.announcement_id = $1;`, a.Id)
	if err != nil { return err }

	a.Tags = make([]Tag, 0)

	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.Id, &tag.Label, &tag.Color)
		if err != nil { return err }

		a.Tags = append(a.Tags, tag)
	}

	return nil
}

func (a *AnnouncementState) RemoveAllTags(db *pgxpool.Pool, tx pgx.Tx) ([]int, error) {
	rows, err := tx.Query(context.Background(), `SELECT tag_id FROM announcement_tags WHERE announcement_id = $1;`, a.Id)
	if err != nil { return []int{}, err }

	tag_ids := make([]int, 0)
	for rows.Next() {
		var tag_id int
		err = rows.Scan(&tag_id)
		if err != nil { return []int{}, err }

		tag_ids = append(tag_ids, tag_id)
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM announcement_tags WHERE announcement_id = $1;`, a.Id)
	return tag_ids, err
}

func (a *AnnouncementState) AddTags(tx pgx.Tx, tag_ids []int) error {
	for _, tag := range a.Tags {
		_, err := tx.Exec(context.Background(), `INSERT INTO announcement_tags (announcement_id, tag_id) VALUES ($1, $2);`, a.Id, tag.Id)
		if err != nil { return err }
	}

	return nil;
}

func (a *AnnouncementState) Create(db *pgxpool.Pool, envMap map[string]string) (string, string) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return "ERR db.Begin in AnnouncementState.Create: " + err.Error(), "Failed to start transaction."
	}

	_, err = tx.Exec(context.Background(), `INSERT INTO announcements (title, content, author_id) VALUES ($1,$2,$3);`, a.Title, a.Content, a.AuthorId)
	if err != nil {
		tx.Rollback(context.Background())
		return "ERR tx.Exec INSERT INTO annoucements in AnnouncementState.Create: " + err.Error(), "Failed inserting announcement into database."
	}

	for _, tag := range a.Tags {
		_, err := tx.Exec(context.Background(), `INSERT INTO announcement_tags (announcement_id, tag_id) VALUES ($1,$2);`, a.Id, tag.Id)
		if err != nil {
			tx.Rollback(context.Background())
			return "ERR tx.Exec INSERT INTO announcement_tags in AnnouncementState.Create: " + err.Error(), "Failed to insert announcement tag connections into database."
		}
	}

	err = tx.Commit(context.Background())
	if err != nil { return "ERR tx.commit in AnnouncementState.Create: " + err.Error(), "Failed to finish transaction." }

	return "", ""
}