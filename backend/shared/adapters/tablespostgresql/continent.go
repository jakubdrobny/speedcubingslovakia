package tablespostgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/repository"
)

func (t *Continent) Get(ctx context.Context, x DbExecutor, id int) error {
	err := x.QueryRowContext(ctx, `
		SELECT c.name, c.recordName
		FROM continents c
		WHERE c.id = ?
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
	}
}
