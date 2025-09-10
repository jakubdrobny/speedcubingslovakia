package tablespostgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DSN struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func (d DSN) encode() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Password, d.User, d.Password, d.DBName)
}

func (d DSN) DB() *sql.DB {
	pool, err := sql.Open("pgx", d.encode())
	if err != nil {
		panic(err)
	}
	pool.SetConnMaxLifetime(time.Minute * 10)
	pool.SetMaxIdleConns(10)
	pool.SetMaxOpenConns(40)

	err = pool.Ping()
	if err != nil {
		panic(err)
	}

	return pool
}

type DbExecutor interface {
	QueryRowContext(ctx context.Context, queryString string, args ...any) (*sql.Row, error)
	QueryContext(ctx context.Context, queryString string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, queryString string, args ...any) (sql.Result, error)
}

func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
