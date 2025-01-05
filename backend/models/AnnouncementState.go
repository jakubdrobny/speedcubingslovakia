package models

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnnouncementState struct {
	Id             int            `json:"id"`
	Title          string         `json:"title"`
	Content        string         `json:"content"`
	AuthorId       int            `json:"authorId"`
	AuthorWcaId    string         `json:"authorWcaId"`
	AuthorUsername string         `json:"authorUsername"`
	Tags           []Tag          `json:"tags"`
	Read           bool           `json:"read"`
	EmojiCounters  []EmojiCounter `json:"emojiCounters"`
}

func GetCompetitionAnnouncementTag(db *pgxpool.Pool) (Tag, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT tag_id, label, color FROM tags t WHERE t.label = '# competition announcement';`,
	)
	if err != nil {
		return Tag{}, err
	}

	var tag Tag
	for rows.Next() {
		err = rows.Scan(&tag.Id, &tag.Label, &tag.Color)
		if err != nil {
			return Tag{}, err
		}
	}

	return tag, nil
}

func (a *AnnouncementState) GetTags(db *pgxpool.Pool) error {
	rows, err := db.Query(
		context.Background(),
		`SELECT t.tag_id, t.label, t.color FROM tags t JOIN announcement_tags at ON t.tag_id = at.tag_id WHERE at.announcement_id = $1;`,
		a.Id,
	)
	if err != nil {
		return err
	}

	a.Tags = make([]Tag, 0)

	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.Id, &tag.Label, &tag.Color)
		if err != nil {
			return err
		}

		a.Tags = append(a.Tags, tag)
	}

	return nil
}

func (a *AnnouncementState) GetEmojiCounters(db *pgxpool.Pool) error {
	rows, err := db.Query(
		context.Background(),
		`SELECT ar.announcement_reaction_id, ar.emoji, ar.by FROM announcement_reaction ar WHERE ar.announcement_id = $1 AND ar."set" = TRUE;`,
		a.Id,
	)
	if err != nil {
		return err
	}

	a.EmojiCounters = make([]EmojiCounter, 0)

	for rows.Next() {
		var emojiCounter EmojiCounter
		err = rows.Scan(&emojiCounter.Id, &emojiCounter.Emoji, &emojiCounter.By)
		if err != nil {
			return err
		}

		a.EmojiCounters = append(a.EmojiCounters, emojiCounter)
	}

	return nil
}

func (a *AnnouncementState) RemoveAllTags(db *pgxpool.Pool, tx pgx.Tx) ([]int, error) {
	rows, err := tx.Query(
		context.Background(),
		`SELECT tag_id FROM announcement_tags WHERE announcement_id = $1;`,
		a.Id,
	)
	if err != nil {
		return []int{}, err
	}

	tag_ids := make([]int, 0)
	for rows.Next() {
		var tag_id int
		err = rows.Scan(&tag_id)
		if err != nil {
			return []int{}, err
		}

		tag_ids = append(tag_ids, tag_id)
	}

	_, err = tx.Exec(
		context.Background(),
		`DELETE FROM announcement_tags WHERE announcement_id = $1;`,
		a.Id,
	)
	return tag_ids, err
}

func (a *AnnouncementState) AddTags(tx pgx.Tx, tag_ids []int) error {
	for _, tag := range a.Tags {
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO announcement_tags (announcement_id, tag_id) VALUES ($1, $2);`,
			a.Id,
			tag.Id,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AnnouncementState) IsRead(db *pgxpool.Pool) error {
	rows, err := db.Query(
		context.Background(),
		`SELECT read FROM announcement_read WHERE announcement_id = $1 AND user_id = $2;`,
		a.Id,
		a.AuthorId,
	)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&a.Read)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AnnouncementState) MarkRead(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE announcement_read SET read = TRUE, read_timestamp = CURRENT_TIMESTAMP WHERE announcement_id = $1 AND user_id = $2;`,
		a.Id,
		a.AuthorId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *AnnouncementState) MakeAnnouncementUnreadForEveryone(tx pgx.Tx) (string, string) {
	users, _, logMessage, returnMessage := GetUsersFromDB(tx, "_")

	if logMessage != "" || returnMessage != "" {
		return logMessage, returnMessage
	}

	for _, user := range users {
		err := user.CreateAnnouncementReadConnection(tx, a.Id)
		if err != nil {
			return "ERR user.MakeAnnouncement in AnnouncementState.MakeAnnouncementUnreadForEveryone: " + err.Error(), "Failed to make announcement unread for user with id: " + strconv.Itoa(
				user.Id,
			)
		}
	}

	return "", ""
}

func (a *AnnouncementState) Create(db *pgxpool.Pool, envMap map[string]string) (string, string) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return "ERR db.Begin in AnnouncementState.Create: " + err.Error(), "Failed to start transaction."
	}

	row := tx.QueryRow(
		context.Background(),
		`INSERT INTO announcements (title, content, author_id, created_at) VALUES ($1,$2,$3,CURRENT_TIMESTAMP) RETURNING announcement_id;`,
		a.Title,
		a.Content,
		a.AuthorId,
	)
	row.Scan(&a.Id)

	for _, tag := range a.Tags {
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO announcement_tags (announcement_id, tag_id) VALUES ($1,$2);`,
			a.Id,
			tag.Id,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return "ERR tx.Exec INSERT INTO announcement_tags in AnnouncementState.Create: " + err.Error(), "Failed to insert announcement tag connections into database."
		}
	}

	logMessage, returnMessage := a.MakeAnnouncementUnreadForEveryone(tx)
	if logMessage != "" || returnMessage != "" {
		return logMessage, returnMessage
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return "ERR tx.commit in AnnouncementState.Create: " + err.Error(), "Failed to finish transaction."
	}

	return "", ""
}
