package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tag struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
	Color string `json:"color"`
}

func GetAvailableTags(db *pgxpool.Pool) ([]Tag, error) {
	rows, err := db.Query(context.Background(), "SELECT t.tag_id, t.label, t.color FROM tags t;")
	if err != nil {
		return []Tag{}, err
	}

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.Id, &tag.Label, &tag.Color)
		if err != nil {
			return []Tag{}, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
