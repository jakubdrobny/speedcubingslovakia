package interfaces

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Query(ctx context.Context, queryString string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, queryString string, arguments ...any) pgx.Row
	Exec(
		ctx context.Context,
		queryString string,
		arguments ...interface{},
	) (pgconn.CommandTag, error)
}
