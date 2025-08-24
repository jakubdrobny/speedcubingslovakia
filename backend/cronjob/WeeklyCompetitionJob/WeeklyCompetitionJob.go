package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jakubdrobny/speedcubingslovakia/backend/controllers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	envMap, err := godotenv.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load environmental variables from file: %v\n", err)
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), envMap["DB_URL"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	controllers.AddNewWeeklyCompetition(db, envMap)
}
