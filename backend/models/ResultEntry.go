package models

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	Scrambles []string `json:"-"`
}

func (r *ResultEntry) Insert(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `INSERT INTO results (competition_id, user_id, event_id, solve1, solve2, solve3, solve4, solve5, comment, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`, r.Competitionid, r.Userid, r.Eventid, r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5, r.Comment, r.Status.Id)
	if err != nil {
		return err
	}

	return nil;
}

func (r *ResultEntry) Validate(db *pgxpool.Pool, isfmc bool, scrambles []string) (error) {
	var err error
	if (r.IsSuspicous(isfmc, scrambles)) {
		r.Status, err = GetResultsStatus(db, 1) // waitingForApproval
		if err != nil { return err }
	} else {
		r.Status, err = GetResultsStatus(db, 3) // approved
		if err != nil { return err }
	}

	return nil
}

func IsValidTimePeriod(db *pgxpool.Pool, competitionId string) (bool, error) {
	competition, err := GetCompetitionByIdObject(db, competitionId)
	if err != nil { return false, err }

	return competition.Startdate.Before(time.Now()) && time.Now().Before(competition.Enddate), nil
}

func (r *ResultEntry) LoadId(db *pgxpool.Pool) error {
	rows, err := db.Query(context.Background(), `SELECT result_id FROM results WHERE user_id = $1 AND competition_id = $2 AND event_id = $3;`, r.Userid, r.Competitionid, r.Eventid)
	if err != nil { return err }

	for rows.Next() {
		err = rows.Scan(&r.Id)
		if err != nil { return nil }
	}

	return nil
}

func ValidateMultiEntry(entry string) string {
	if entry == "DNS" || entry == "0/0 00:00:00" { return "DNS" }

	r, _ := regexp.Compile("[0-9]{1,2}[/][0-9]{1,2}[ ][0-1][0-9]:[0-5][0-9]:[0-5][0-9]")
	if !r.MatchString(entry) {
		return "DNS"
	}

	cubes := strings.Split(entry, " ")[0]
	solved, _ := strconv.Atoi(strings.Split(cubes, "/")[0])
	attempted, _ := strconv.Atoi(strings.Split(cubes, "/")[1])

	if solved > attempted || attempted > constants.MBLD_MAX_CUBES_PER_ATTEMPT {
		return "DNS"
	}

	return entry
}

func (r *ResultEntry) Update(db *pgxpool.Pool, isfmc bool, valid ...bool) error {
	err := r.LoadId(db)
	if err != nil { return err }

	if (isfmc) {
		r.Scrambles, err = utils.GetScramblesByResultEntryId(db, r.Eventid, r.Competitionid)
		if err != nil { return err }
	} else {
		r.Scrambles = make([]string, 5)
	}

	if r.Iconcode == "333mbf" {
		r.Solve1 = ValidateMultiEntry(r.Solve1)
		r.Solve2 = ValidateMultiEntry(r.Solve2)
		r.Solve3 = ValidateMultiEntry(r.Solve3)
		r.Solve4 = ValidateMultiEntry(r.Solve4)
		r.Solve5 = ValidateMultiEntry(r.Solve5)
	}

	if len(valid) == 0 || (len(valid) > 0 && !valid[0]) {
		err := r.Validate(db, isfmc, r.Scrambles)
		if err != nil { return err }
	}

	ok, err := IsValidTimePeriod(db, r.Competitionid)
	if err != nil { return err }

	if ok {
		_, err := db.Exec(context.Background(), `UPDATE results SET solve1 = $1, solve2 = $2, solve3 = $3, solve4 = $4, solve5 = $5, comment = $6, status_id = $7, timestamp = CURRENT_TIMESTAMP WHERE user_id = $8 AND competition_id = $9 AND event_id = $10;`, r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5, r.Comment, r.Status.Id, r.Userid, r.Competitionid, r.Eventid)
		if err != nil { return err }
	}

	return nil;
}


func (r *ResultEntry) Single(isfmc bool, scrambles []string) int {
	res := math.MaxInt

	utils.CompareSolves(&res, r.Solve1, isfmc, scrambles[0])
	utils.CompareSolves(&res, r.Solve2, isfmc, scrambles[1])
	utils.CompareSolves(&res, r.Solve3, isfmc, scrambles[2])
	utils.CompareSolves(&res, r.Solve4, isfmc, scrambles[3])
	utils.CompareSolves(&res, r.Solve5, isfmc, scrambles[4])

	return res
}


func (r *ResultEntry) IsSuspicous(isfmc bool, scrambles []string) bool {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil {
		return false
	}

	curSingle, curAverage := r.Single(isfmc, scrambles), r.Average(noOfSolves, isfmc, scrambles)
	if isfmc {
		curSingle = utils.ParseSolveToMilliseconds(utils.FormatTime(curSingle, true), false, "")
		curAverage = utils.ParseSolveToMilliseconds(utils.FormatTime(curAverage, true), false, "")
	}

	recSingle, recAverage, err := utils.GetWorldRecords(r.Iconcode)
	if err != nil { return false }

	if r.IsMBLD() {
		return curSingle < constants.VERY_SLOW && float64(recSingle - curSingle) > 1e-9
	}

	return float64(recSingle - curSingle) > 1e-9 || float64(recAverage - curAverage) > 1e-9;
}

func (r *ResultEntry) IsMBLD() bool {
	return r.Iconcode == "333mbf"
}

func (r *ResultEntry) GetSolvesInMiliseconds(isfmc bool, scrambles []string) []int {
	values := make([]int, 0)

	values = append(values, utils.ParseSolveToMilliseconds(r.Solve1, isfmc, scrambles[0]))
	values = append(values, utils.ParseSolveToMilliseconds(r.Solve2, isfmc, scrambles[1]))
	values = append(values, utils.ParseSolveToMilliseconds(r.Solve3, isfmc, scrambles[2]))
	values = append(values, utils.ParseSolveToMilliseconds(r.Solve4, isfmc, scrambles[3]))
	values = append(values, utils.ParseSolveToMilliseconds(r.Solve5, isfmc, scrambles[4]))

	return values;
}

func (r *ResultEntry) GetSolves(isfmc bool, scrambles []string) []string {
	values := make([]string, 0)

	values = append(values, utils.GetSolve(r.Solve1, isfmc, scrambles[0]))
	values = append(values, utils.GetSolve(r.Solve2, isfmc, scrambles[1]))
	values = append(values, utils.GetSolve(r.Solve3, isfmc, scrambles[2]))
	values = append(values, utils.GetSolve(r.Solve4, isfmc, scrambles[3]))
	values = append(values, utils.GetSolve(r.Solve5, isfmc, scrambles[4]))

	return values;
}

func (r *ResultEntry) Average(noOfSolves int, isfmc bool, scrambles []string) int {
	solves := r.GetSolvesInMiliseconds(isfmc, scrambles)
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

func (r *ResultEntry) FormatMultiSingle(single int) string {
	if utils.ParseMultiToMilliseconds(r.Solve1) == single { return r.Solve1 }
	if utils.ParseMultiToMilliseconds(r.Solve2) == single { return r.Solve2 }
	if utils.ParseMultiToMilliseconds(r.Solve3) == single { return r.Solve3 }
	if utils.ParseMultiToMilliseconds(r.Solve4) == single { return r.Solve4 }
	if utils.ParseMultiToMilliseconds(r.Solve5) == single { return r.Solve5 }

	return "DNS"
}

func (r *ResultEntry) SingleFormatted(isfmc bool, scrambles []string) string {
	single := r.Single(isfmc, scrambles)

	if r.Iconcode == "333mbf" {
		if single == constants.DNF { return "DNF" }
		
		return r.FormatMultiSingle(single)
	}

	return utils.FormatTime(single, isfmc)
}

func (r *ResultEntry) AverageFormatted(isfmc bool, scrambles []string) (string, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil { return "", err }

	return utils.FormatTime(r.Average(noOfSolves, isfmc, scrambles), isfmc), nil
}

func (r *ResultEntry) GetFormattedTimes(isfmc bool, scrambles []string) ([]string, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil { return []string{}, err }

	solves := r.GetSolves(isfmc, scrambles)
	solves = solves[:noOfSolves]
	if noOfSolves == 3 {
		if r.Iconcode == "333mbf" {
			return utils.FormatMultiTimes(solves), nil
		}
		return solves, nil
	}

	type SolveTuple struct {
		FormattedTime string
		TimeInMiliseconds int
		Index int
	}
	sortedSolves := make([]SolveTuple, 0)

	for idx, val := range solves {
		sortedSolves = append(sortedSolves, SolveTuple{val, utils.ParseSolveToMilliseconds(val, false, ""), idx})
	}

	sort.Slice(sortedSolves, func (i int, j int) bool { return sortedSolves[i].TimeInMiliseconds < sortedSolves[j].TimeInMiliseconds })
	solves[sortedSolves[0].Index] = "(" + solves[sortedSolves[0].Index] + ")"
	solves[sortedSolves[len(sortedSolves) - 1].Index] = "(" + solves[sortedSolves[len(sortedSolves) - 1].Index] + ")"

	return solves, nil
}

func (r *ResultEntry) IsFMC() bool {
	return r.Iconcode == "333fm"
}

func (r *ResultEntry) GetSolveIdx(s string) int {
	if r.Solve1 == s { return 0 }
	if r.Solve2 == s { return 1 }
	if r.Solve3 == s { return 2 }
	if r.Solve4 == s { return 3 }
	if r.Solve5 == s { return 4 }

	return -1
}