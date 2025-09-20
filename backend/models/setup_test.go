package models_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDb *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("%w: when running postgres container", err))
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatal(fmt.Errorf("%w: when terminating postgres container", err))
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(fmt.Errorf("%w: when getting container connection string", err))
	}

	testDb, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(fmt.Errorf("%w: when connecting to test database", err))
	}

	if err := runMigrations(connStr); err != nil {
		log.Fatal(fmt.Errorf("%w: when running migrations", err))
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func runMigrations(dbURL string) error {
	migrationsPath, err := filepath.Abs("../../database/migrations")
	if err != nil {
		return fmt.Errorf("%w: when getting absolute path for migrations", err)
	}

	sourceURL := "file://" + migrationsPath

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return fmt.Errorf("%w: when creating new migrate instance", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%w: when applying up migrations", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
