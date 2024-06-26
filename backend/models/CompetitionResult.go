package models

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type CompetitionResult struct {
	Username string `json:"username"`
	WcaId string `json:"wca_id"`
	CountryName string `json:"country_name"`
	CountryIso2 string `json:"country_iso2"`
	Single string `json:"single"`
	Average string `json:"average"`
	Times []string `json:"times"`
	Score string `json:"score"`
}

type BestEntry struct {
	Single int
	Average int
}

func GetNewBest(pBest BestEntry, resultEntry ResultEntry, noOfSolves int, scrambles []string) BestEntry {
	single, average := resultEntry.Single(resultEntry.IsFMC(), scrambles), resultEntry.Average(noOfSolves, resultEntry.IsFMC(), scrambles)
	if single < pBest.Single { pBest.Single = single }
	if average < pBest.Average { pBest.Average = average }

	return pBest
}

func ComputeBests(bests map[int]BestEntry, rows []KinchQueryRow) error {
	for _, row := range rows {
		resultEntry := row.ResultEntry
		
		if !resultEntry.Competed() || !resultEntry.Status.Visible { continue; }

		eid := resultEntry.Eventid
		
		noOfSolves, err := utils.GetNoOfSolves(resultEntry.Format)
		if err != nil { return err }
		bests[eid] = GetNewBest(bests[eid], resultEntry, noOfSolves, resultEntry.Scrambles)
	}

	return nil
}

type KinchQueryRow struct {
	CompetitionResult CompetitionResult
	ResultEntry ResultEntry
}

func GetScores(rows []KinchQueryRow, bests map[int]BestEntry, noOfEvents int, db *pgxpool.Pool) ([]CompetitionResult, error) {
	cums := make(map[int]float64)
	res := make(map[int]CompetitionResult)

	for _, row := range rows {
		competitionResult := row.CompetitionResult
		resultEntry := row.ResultEntry

		if !resultEntry.Competed() || !resultEntry.Status.Visible { continue; }

		noOfSolves, err := utils.GetNoOfSolves(resultEntry.Format)
		if err != nil { return []CompetitionResult{}, nil }

		// KINCH RANKS - 4bld, 5bld, mbld sa berie single, 3bld a fmc lepsi z single,average a ostatne average
		single := resultEntry.Single(resultEntry.IsFMC(), resultEntry.Scrambles)
		singleContrib := float64(bests[resultEntry.Eventid].Single) / float64(single)
		ismbld := resultEntry.Iconcode == "333mbf"
		if ismbld { singleContrib = float64(single) / float64(bests[resultEntry.Eventid].Single) }
		if single >= constants.VERY_SLOW { singleContrib = 0. }

		average := resultEntry.Average(noOfSolves, resultEntry.IsFMC(), resultEntry.Scrambles)
		averageContrib := float64(bests[resultEntry.Eventid].Average) / float64(average)
		if average >= constants.VERY_SLOW { averageContrib = 0. }

		var finalContrib float64 = averageContrib
		if resultEntry.Iconcode == "444bf" || resultEntry.Eventname == "555bf" || ismbld {
			finalContrib = singleContrib
		} else if resultEntry.Eventname == "333bf" || resultEntry.Iconcode == "333fm" {
			finalContrib = math.Max(finalContrib, singleContrib)
		}

		cums[resultEntry.Userid] += finalContrib * 100
		res[resultEntry.Userid] = competitionResult
	}

	competitionResults := make([]CompetitionResult, 0)
	for uid, cum := range cums {
		competitionResult := res[uid]
		competitionResult.Score = fmt.Sprintf("%.2f", cum / float64(noOfEvents))
		competitionResults = append(competitionResults, competitionResult)
	}

	sort.Slice(competitionResults, func (i int, j int) bool {
		a, err := strconv.ParseFloat(competitionResults[i].Score, 64)
		if err != nil { return true }
		b, err := strconv.ParseFloat(competitionResults[j].Score, 64)
		if err != nil { return true }

		return a - b > 1e-9;
	})

	return competitionResults, nil
}

func GetKinchQueryRows(rawRows pgx.Rows, db *pgxpool.Pool) ([]KinchQueryRow, error) {
	rows := make([]KinchQueryRow, 0)

	for rawRows.Next() {
		var competitionResult CompetitionResult
		var resultEntry ResultEntry
		
		err := rawRows.Scan(&resultEntry.Userid, &competitionResult.WcaId, &competitionResult.Username, &competitionResult.CountryName, &competitionResult.CountryIso2, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format, &resultEntry.Status.Visible, &resultEntry.Eventid, &resultEntry.Iconcode, &resultEntry.Eventid, &resultEntry.Competitionid)
		if err != nil { return []KinchQueryRow{}, err }

		if competitionResult.WcaId == "" { competitionResult.WcaId = competitionResult.Username }

		scrambles, err := utils.GetScramblesByResultEntryId(db, resultEntry.Eventid, resultEntry.Competitionid)
		if err != nil { return []KinchQueryRow{}, err }
		resultEntry.Scrambles = scrambles

		rows = append(rows, KinchQueryRow{competitionResult, resultEntry})
	}

	return rows, nil
}

func GetOverallResults(db *pgxpool.Pool, cid string) ([]CompetitionResult, error) {
	rawRows, err := db.Query(context.Background(), `SELECT u.user_id, u.wcaid, u.name, c.name, c.iso2, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, rs.visible, e.event_id, e.iconcode, r.event_id, r.competition_id FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN events e ON e.event_id = r.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE r.competition_id = $1;`, cid)
	if err != nil { return []CompetitionResult{}, err }
	rows, err := GetKinchQueryRows(rawRows, db)
	if err != nil { return []CompetitionResult{}, err }

	competition, err := GetCompetitionByIdObject(db, cid)
	if err != nil { return []CompetitionResult{}, err }
	err = competition.GetEvents(db)
	if err != nil { return []CompetitionResult{}, err }
	
	bests := make(map[int]BestEntry)
	for _, ev := range competition.Events { bests[ev.Id] = BestEntry{constants.DNS, constants.DNS} }
	noOfEvents := len(competition.Events) - 1

	err = ComputeBests(bests, rows)
	if err != nil { return []CompetitionResult{}, err }

	competitionResults, err := GetScores(rows, bests, noOfEvents, db)
	if err != nil { return []CompetitionResult{}, err }

	return competitionResults, nil
}

func GetResultsFromCompetitionByEventName(db *pgxpool.Pool, cid string, eid int) ([]CompetitionResult, error) {
	if (eid == -1) {
		competitionResults, err := GetOverallResults(db, cid)
		if err != nil { return []CompetitionResult{}, err}
		return competitionResults, nil
	}
	
	rows, err := db.Query(context.Background(), `SELECT u.name, u.wcaid, c.name, c.iso2, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, rs.visible, e.iconcode, r.event_id, r.competition_id FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN events e ON e.event_id = r.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE r.competition_id = $1 AND r.event_id = $2;`, cid, eid)
	if err != nil { return []CompetitionResult{}, err }
	
	competitionResults := make([]CompetitionResult, 0)
	format := ""
	
	for rows.Next() {
		var competitionResult CompetitionResult
		var resultEntry ResultEntry
		
		err = rows.Scan(&competitionResult.Username, &competitionResult.WcaId, &competitionResult.CountryName, &competitionResult.CountryIso2, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format, &resultEntry.Status.Visible, &resultEntry.Iconcode, &resultEntry.Eventid, &resultEntry.Competitionid)
		if err != nil { return []CompetitionResult{}, err }
		
		if competitionResult.WcaId == "" { competitionResult.WcaId = competitionResult.Username }
		
		if !resultEntry.Competed() || !resultEntry.Status.Visible { continue; }
		
		scrambles, err := utils.GetScramblesByResultEntryId(db, resultEntry.Eventid, resultEntry.Competitionid)
		if err != nil { return []CompetitionResult{}, err }
		
		competitionResult.Single = resultEntry.SingleFormatted(resultEntry.IsFMC(), scrambles)
		
		avg, err := resultEntry.AverageFormatted(resultEntry.IsFMC(), scrambles)
		if err != nil { return []CompetitionResult{}, err }
		competitionResult.Average = avg

		formattedTimes, err := resultEntry.GetFormattedTimes(resultEntry.IsFMC(), scrambles)
		if err != nil { return []CompetitionResult{}, err }
		competitionResult.Times = formattedTimes

		competitionResults = append(competitionResults, competitionResult)
		format = resultEntry.Format
	}

	if len(format) > 0 {
		sort.Slice(competitionResults, func (i int, j int) bool {
			if format[0] == 'b' { return utils.ParseSolveToMilliseconds(competitionResults[i].Single, false, "") < utils.ParseSolveToMilliseconds(competitionResults[j].Single, false, "")}
			return utils.ParseSolveToMilliseconds(competitionResults[i].Average, false, "") < utils.ParseSolveToMilliseconds(competitionResults[j].Average, false, "")
		})
	}

	return competitionResults, nil
}