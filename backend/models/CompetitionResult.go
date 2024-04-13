package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CompetitionResult struct {
	Username string `json:"username"`
	CountryName string `json:"country_name"`
	CountryIso2 string `json:"country_iso2"`
	Single string `json:"single"`
	Average string `json:"average"`
	Times []string `json:"times"`
}

func GetResultsFromCompetitionByEventName(db *pgxpool.Pool, cid string, eid int) ([]CompetitionResult, error) {
	rows, err := db.Query(context.Background(), `SELECT u.name, c.name, c.iso2, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, rs.visible FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN events e ON e.event_id = r.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE r.competition_id = $1 AND r.event_id = $2;`, cid, eid)
	if err != nil { return []CompetitionResult{}, err }

	competitionResults := make([]CompetitionResult, 0)
	for rows.Next() {
		var competitionResult CompetitionResult
		var resultEntry ResultEntry
		
		err = rows.Scan(&competitionResult.Username, &competitionResult.CountryName, &competitionResult.CountryIso2, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format, &resultEntry.Status.Visible)
		if err != nil { return []CompetitionResult{}, err }

		if !resultEntry.Competed() || !resultEntry.Status.Visible { continue; }

		competitionResult.Single = resultEntry.SingleFormatted()
		
		avg, err := resultEntry.AverageFormatted()
		if err != nil { return []CompetitionResult{}, err }
		competitionResult.Average = avg

		formattedTimes, err := resultEntry.GetFormattedTimes()
		if err != nil { return []CompetitionResult{}, err }
		competitionResult.Times = formattedTimes

		competitionResults = append(competitionResults, competitionResult)
	}

	return competitionResults, nil
}