package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testSqlDb interfaces.DB

func loadMigrationsFileNames() ([]string, error) {
	migrationsDirPath := "../../database/migrations"
	entries, err := os.ReadDir(migrationsDirPath)
	if err != nil {
		return []string{}, fmt.Errorf("%w: when listing migrations directory with path %s", err, migrationsDirPath)
	}

	filenames := []string{}
	for _, entry := range entries {
		filename := entry.Name()
		if strings.HasSuffix(filename, ".up.sql") {
			filenames = append(filenames, filename)
		}
	}

	return filenames, nil
}

func TestMain(m *testing.M) {
	var err error
	defer utils.PrintStack(err)

	ctx := context.Background()

	dbName := "speedcubingslovakiadb_local"
	dbUser := "admin"
	dbPassword := "password"

	migrationsScripts, err := loadMigrationsFileNames()
	if err != nil {
		err = fmt.Errorf("%w: when calling loadMigrationsFileNames", err)
		return
	}

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithOrderedInitScripts(migrationsScripts...),
	)
	if err != nil {
		err = fmt.Errorf("%w: when running postgres testcontainer", err)
		return
	}

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			utils.PrintStack(fmt.Errorf("%w: when terminating postgres testcontainer", err))
		}
	}()

	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		err = fmt.Errorf("%w: when getting postgres testcontainer connection string", err)
		return
	}

	testSqlDb, err = sql.Open("pgx", connectionString)
	if err != nil {
		err = fmt.Errorf("%w: when opening pgx connection with connectionString: %v", err, connectionString)
		return
	}

	err = testSqlDb.PingContext(ctx)
	if err != nil {
		err = fmt.Errorf("%w: when testing database connection", err)
		return
	}

	fmt.Println("Successfully connected to the postgres testcontainer database.")

	code := m.Run()
	os.Exit(code)
}
