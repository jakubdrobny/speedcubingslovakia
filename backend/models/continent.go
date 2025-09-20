package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

type Continent struct {
	Id         string
	Name       string
	RecordName string
}

func (c *Continent) Get(ctx context.Context, db interfaces.DB, name string) error {
	err := db.QueryRow(ctx, `SELECT continent_id, name, recordName FROM continents WHERE name = $1`, name).Scan(&c.Id, &c.Name, &c.RecordName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("continent with name=%s not found", name)
		}
		return fmt.Errorf("%w: when querying continent by name=%s", err, name)
	}

	return nil
}

func (c Continent) Insert(ctx context.Context, db interfaces.DB) error {
	_, err := db.Exec(ctx, `INSERT INTO continents (continent_id, name, recordName) VALUES ($1, $2, $3)`, c.Id, c.Name, c.RecordName)
	if err != nil {
		return fmt.Errorf("%w: when executing insert continent statement for continent=%+v", err, c)
	}

	return nil
}
