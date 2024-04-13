package models

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type ResultEntry struct {
	Id int `json:"id"`
	Userid int `json:"userid"`
	Username string `json:"username"`
	Competitionid string `json:"competitionid"`
	Competitionname string `json:"competitionname"`
	Eventid int `json:"eventid"`
	Eventname string `json:"eventname"`
	Iconcode string `json:"iconcode"`
	Format string `json:"format"`
	Solve1 string `json:"solve1"`
	Solve2 string `json:"solve2"`
	Solve3 string `json:"solve3"`
	Solve4 string `json:"solve4"`
	Solve5 string `json:"solve5"`
	Comment string `json:"comment"`
	Status ResultsStatus `json:"status"`
}

func (r *ResultEntry) Insert(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `INSERT INTO results (competition_id, user_id, event_id, solve1, solve2, solve3, solve4, solve5, comment, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`, r.Competitionid, 1, r.Eventid, r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5, r.Comment, r.Status.Id)
	if err != nil {
		return err
	}

	return nil;
}

func (r *ResultEntry) Validate(db *pgxpool.Pool) (error) {
	var err error
	if (r.IsSuspicous()) {
		r.Status, err = GetResultsStatus(db, 1) // waitingForApproval
		if err != nil { return err }
	} else {
		r.Status, err = GetResultsStatus(db, 3) // approved
		if err != nil { return err }
	}

	competition, err := GetCompetitionByIdObject(db, r.Competitionid)
	if err != nil { return err }

	if time.Now().Before(competition.Startdate) || time.Now().After(competition.Enddate) {
		return fmt.Errorf("competition has not started yet or has already finished")
	}

	return nil
}

func (r *ResultEntry) Update(db *pgxpool.Pool, valid ...bool) error {
	if len(valid) == 0 || (len(valid) > 0 && !valid[0]) {
		err := r.Validate(db)
		if err != nil { return err }
	}
	
	_, err := db.Exec(context.Background(), `UPDATE results SET solve1 = $1, solve2 = $2, solve3 = $3, solve4 = $4, solve5 = $5, comment = $6, status_id = $7, timestamp = CURRENT_TIMESTAMP WHERE user_id = $8 AND competition_id = $9 AND event_id = $10;`, r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5, r.Comment, r.Status.Id, r.Userid, r.Competitionid, r.Eventid)
	if err != nil { return err }

	return nil;
}


func (r *ResultEntry) Single() int {
	res := math.MaxInt

	utils.CompareSolves(&res, r.Solve1)
	utils.CompareSolves(&res, r.Solve2)
	utils.CompareSolves(&res, r.Solve3)
	utils.CompareSolves(&res, r.Solve4)
	utils.CompareSolves(&res, r.Solve5)

	return res
}


func (r *ResultEntry) IsSuspicous() bool {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil {
		return false
	}

	curSingle, curAverage := r.Single(), r.Average(noOfSolves)

	recSingle, recAverage, err := utils.GetWorldRecords(r.Iconcode)
	if err != nil { return false }

	return float64(recSingle - curSingle) > 1e-9 || float64(recAverage - curAverage) > 1e-9;
}

func (r *ResultEntry) GetSolvesInMiliseconds() []int {
	values := make([]int, 0)

	values = append(values, utils.ParseSolveToMilliseconds(r.Solve1))
	values = append(values, utils.ParseSolveToMilliseconds(r.Solve2))
	values = append(values, utils.ParseSolveToMilliseconds(r.Solve3))
	values = append(values, utils.ParseSolveToMilliseconds(r.Solve4))
	values = append(values, utils.ParseSolveToMilliseconds(r.Solve5))

	return values;
}

func (r *ResultEntry) GetSolves() []string {
	values := make([]string, 0)

	values = append(values, r.Solve1)
	values = append(values, r.Solve2)
	values = append(values, r.Solve3)
	values = append(values, r.Solve4)
	values = append(values, r.Solve5)

	return values;
}

func (r *ResultEntry) Average(noOfSolves int) int {
	solves := r.GetSolvesInMiliseconds()
	sort.Ints(solves)

	sum := 0
	cntBad := 0

	for idx, solve := range solves {
		if idx >= noOfSolves {
			break
		}

		if solve >= constants.VERY_SLOW {
			cntBad++
			if (noOfSolves == 5 && cntBad > 1) || (noOfSolves == 3 && cntBad > 0) {
				if !r.Competed() { return constants.DNS }
				return constants.DNF
			}
		}

		if noOfSolves == 3 || (noOfSolves == 5 && idx > 0 && idx < 4) {
			sum += solve
		}
	}

	return sum / 3
}

func (r *ResultEntry) Competed() bool {
	return r.Solve1 != "DNS" || r.Solve2 != "DNS" || r.Solve3 != "DNS" || r.Solve4 != "DNS" || r.Solve5 != "DNS"
}

func GetResultEntry(db *pgxpool.Pool, competitorId int, competitionId string, eventId int) (ResultEntry, error) {
	rows, err := db.Query(context.Background(), `SELECT re.result_id, re.competition_id, re.user_id, re.event_id, re.solve1, re.solve2, re.solve3, re.solve4, re.solve5, re.comment, re.status_id FROM results re WHERE re.user_id = $1 AND re.competition_id = $2 AND re.event_id = $3;`, competitorId, competitionId, eventId)
	if err != nil { return ResultEntry{}, err }

	var resultEntry ResultEntry
	found := false
	for rows.Next() {
		err = rows.Scan(&resultEntry.Id, &resultEntry.Competitionid, &resultEntry.Userid, &resultEntry.Eventid, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Comment, &resultEntry.Status.Id)
		if err != nil { return ResultEntry{}, err }
		found = true
	}

	if !found { return ResultEntry{}, fmt.Errorf("not found") }

	return resultEntry, nil
}

func GetResultEntryById(db *pgxpool.Pool, resultId int) (ResultEntry, error) {
	rows, err := db.Query(context.Background(), `SELECT re.result_id, re.competition_id, re.user_id, re.event_id, re.solve1, re.solve2, re.solve3, re.solve4, re.solve5, re.comment, re.status_id, c.name, e.displayname, rs.approvalfinished, rs.approved, rs.visible, rs.displayname, u.name, e.format, e.iconcode FROM results re JOIN competitions c ON c.competition_id = re.competition_id JOIN events e ON e.event_id = re.event_id JOIN results_status rs ON results_status_id = re.status_id JOIN users u ON u.user_id = re.user_id WHERE re.result_id = $1;`, resultId)
	if err != nil { return ResultEntry{}, err }

	var resultEntry ResultEntry
	found := false
	for rows.Next() {
		err = rows.Scan(&resultEntry.Id, &resultEntry.Competitionid, &resultEntry.Userid, &resultEntry.Eventid, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Comment, &resultEntry.Status.Id, &resultEntry.Competitionname, &resultEntry.Eventname, &resultEntry.Status.ApprovalFinished, &resultEntry.Status.Approved, &resultEntry.Status.Visible, &resultEntry.Status.Displayname, &resultEntry.Username, &resultEntry.Format, &resultEntry.Iconcode)
		if err != nil { return ResultEntry{}, err }
		found = true
	}

	if !found { return ResultEntry{}, err }

	return resultEntry, nil
}

func (r *ResultEntry) SingleFormatted() string {
	return utils.FormatTime(r.Single())
}

func (r *ResultEntry) AverageFormatted() (string, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil { return "", err }

	return utils.FormatTime(r.Average(noOfSolves)), nil
}

func (r *ResultEntry) GetFormattedTimes() ([]string, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil { return []string{}, err }

	solves := r.GetSolves()
	if noOfSolves == 3 { return solves, nil }

	type SolveTuple struct {
		FormattedTime string
		TimeInMiliseconds int
		Index int
	}
	sortedSolves := make([]SolveTuple, 0)

	for idx, val := range solves {
		sortedSolves = append(sortedSolves, SolveTuple{val, utils.ParseSolveToMilliseconds(val), idx})
	}

	sort.Slice(sortedSolves, func (i int, j int) bool { return sortedSolves[i].TimeInMiliseconds < sortedSolves[j].TimeInMiliseconds })
	solves[sortedSolves[0].Index] = "(" + solves[sortedSolves[0].Index] + ")"
	solves[sortedSolves[len(sortedSolves) - 1].Index] = "(" + solves[sortedSolves[len(sortedSolves) - 1].Index] + ")"

	return solves, nil
}

