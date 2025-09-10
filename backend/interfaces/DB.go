package interfaces

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Query(ctx context.Context, queryString string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, queryString string, args ...any) pgx.Row
	Exec(ctx context.Context, queryString string, args ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}
