package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/jakubdrobny/speedcubingslovakia/backend/controllers"
)

func main() {
	envMap, err := godotenv.Read(
		fmt.Sprintf("../.env.%s", os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV")),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load environmental variables from file: %v\n", err)
		return
	}

	db, err := pgxpool.New(context.Background(), envMap["DB_URL"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	err = controllers.DeletePastWCACompetitions(db)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Something went wrong during deleting past WCA competitions: %v\n",
			err,
		)
		return
	}
}
