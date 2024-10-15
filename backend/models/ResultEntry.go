package models

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/email"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type ResultEntry struct {
	Id              int           `json:"id"`
	Userid          int           `json:"userid"`
	Username        string        `json:"username"`
	WcaId           string        `json:"wcaid"`
	Competitionid   string        `json:"competitionid"`
	Competitionname string        `json:"competitionname"`
	Eventid         int           `json:"eventid"`
	Eventname       string        `json:"eventname"`
	Iconcode        string        `json:"iconcode"`
	Format          string        `json:"format"`
	Solve1          string        `json:"solve1"`
	Solve2          string        `json:"solve2"`
	Solve3          string        `json:"solve3"`
	Solve4          string        `json:"solve4"`
	Solve5          string        `json:"solve5"`
	Comment         string        `json:"comment"`
	Status          ResultsStatus `json:"status"`
	BadFormat       bool          `json:"badFormat"`
	Scrambles       []string      `json:"-"`
	Email           string        `json:"-"`
}

func (r *ResultEntry) Insert(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`INSERT INTO results (competition_id, user_id, event_id, solve1, solve2, solve3, solve4, solve5, comment, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		r.Competitionid,
		r.Userid,
		r.Eventid,
		r.Solve1,
		r.Solve2,
		r.Solve3,
		r.Solve4,
		r.Solve5,
		r.Comment,
		r.Status.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ResultEntry) CheckFormats(isfmc bool) {
	r.BadFormat = false

	if isfmc {
		return
	}

	if !utils.CheckFormat(r.Solve1) {
		r.Solve1 = "DNF"
		r.BadFormat = true
	}
	if !utils.CheckFormat(r.Solve2) {
		r.Solve2 = "DNF"
		r.BadFormat = true
	}
	if !utils.CheckFormat(r.Solve3) {
		r.Solve3 = "DNF"
		r.BadFormat = true
	}
	if !utils.CheckFormat(r.Solve4) {
		r.Solve4 = "DNF"
		r.BadFormat = true
	}
	if !utils.CheckFormat(r.Solve5) {
		r.Solve5 = "DNF"
		r.BadFormat = true
	}
}

func (r *ResultEntry) Validate(db *pgxpool.Pool, isfmc bool, scrambles []string) error {
	var err error
	if r.IsSuspicous(isfmc, scrambles) {
		r.Status, err = GetResultsStatus(db, 1) // waitingForApproval
		if err != nil {
			return err
		}
	} else {
		if !r.IsMBLD() {
			r.CheckFormats(isfmc)
		}
		r.Status, err = GetResultsStatus(db, 3) // approved
		if err != nil {
			return err
		}
	}

	return nil
}

func IsValidTimePeriod(db *pgxpool.Pool, competitionId string) (bool, error) {
	competition, err := GetCompetitionByIdObject(db, competitionId)
	if err != nil {
		return false, err
	}

	return competition.Startdate.Before(time.Now()) && time.Now().Before(competition.Enddate), nil
}

func (r *ResultEntry) LoadId(db *pgxpool.Pool) error {
	rows, err := db.Query(
		context.Background(),
		`SELECT result_id FROM results WHERE user_id = $1 AND competition_id = $2 AND event_id = $3;`,
		r.Userid,
		r.Competitionid,
		r.Eventid,
	)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&r.Id)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (r *ResultEntry) ValidateMultiEntry(entry string) string {
	if entry == "DNS" || entry == "0/0 00:00:00" {
		return "DNS"
	}

	re, _ := regexp.Compile("[0-9]{1,2}[/][0-9]{1,2}[ ][0-1][0-9]:[0-5][0-9]:[0-5][0-9]")
	if !re.MatchString(entry) {
		r.BadFormat = true
		return "DNF"
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
	var err error

	if r.Id == 0 {
		err = r.LoadId(db)
		if err != nil {
			return err
		}
	}

	if isfmc {
		r.Scrambles, err = utils.GetScramblesByResultEntryId(db, r.Eventid, r.Competitionid)
		if err != nil {
			return err
		}
	} else {
		r.Scrambles = make([]string, 5)
	}

	if r.Iconcode == "333mbf" {
		r.Solve1 = r.ValidateMultiEntry(r.Solve1)
		r.Solve2 = r.ValidateMultiEntry(r.Solve2)
		r.Solve3 = r.ValidateMultiEntry(r.Solve3)
		r.Solve4 = r.ValidateMultiEntry(r.Solve4)
		r.Solve5 = r.ValidateMultiEntry(r.Solve5)
	}

	if len(valid) == 0 || (len(valid) > 0 && !valid[0]) {
		err := r.Validate(db, isfmc, r.Scrambles)
		if err != nil {
			return err
		}
	}

	ok, err := IsValidTimePeriod(db, r.Competitionid)
	if err != nil {
		return err
	}

	if ok {
		_, err := db.Exec(
			context.Background(),
			`UPDATE results SET solve1 = $1, solve2 = $2, solve3 = $3, solve4 = $4, solve5 = $5, comment = $6, status_id = $7, timestamp = CURRENT_TIMESTAMP WHERE user_id = $8 AND competition_id = $9 AND event_id = $10;`,
			r.Solve1,
			r.Solve2,
			r.Solve3,
			r.Solve4,
			r.Solve5,
			r.Comment,
			r.Status.Id,
			r.Userid,
			r.Competitionid,
			r.Eventid,
		)
		if err != nil {
			return err
		}
	}

	return nil
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
	// Not necessary since using new GetWorldRecords function
	//if (len(r.Iconcode) >= 10 && r.Iconcode[:10] == "unofficial") || r.Iconcode == "333ft" {
	//return false
	//}

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
	if err != nil {
		return false
	}

	if r.IsMBLD() {
		return curSingle < constants.VERY_SLOW && float64(recSingle-curSingle) > 1e-9
	}

	return (recSingle < constants.VERY_SLOW && float64(recSingle-curSingle) > 1e-9) ||
		(recAverage < constants.VERY_SLOW && float64(recAverage-curAverage) > 1e-9)
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

	return values
}

func (r *ResultEntry) GetSolves(isfmc bool, scrambles []string) []string {
	values := make([]string, 0)

	values = append(values, utils.GetSolve(r.Solve1, isfmc, scrambles[0]))
	values = append(values, utils.GetSolve(r.Solve2, isfmc, scrambles[1]))
	values = append(values, utils.GetSolve(r.Solve3, isfmc, scrambles[2]))
	values = append(values, utils.GetSolve(r.Solve4, isfmc, scrambles[3]))
	values = append(values, utils.GetSolve(r.Solve5, isfmc, scrambles[4]))

	return values
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
				if !r.Competed() {
					return constants.DNS
				}
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
	return r.Solve1 != "DNS" || r.Solve2 != "DNS" || r.Solve3 != "DNS" || r.Solve4 != "DNS" ||
		r.Solve5 != "DNS"
}

func GetResultEntry(
	db *pgxpool.Pool,
	competitorId int,
	competitionId string,
	eventId int,
) (ResultEntry, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT re.result_id, re.competition_id, re.user_id, re.event_id, re.solve1, re.solve2, re.solve3, re.solve4, re.solve5, re.comment, re.status_id FROM results re WHERE re.user_id = $1 AND re.competition_id = $2 AND re.event_id = $3;`,
		competitorId,
		competitionId,
		eventId,
	)
	if err != nil {
		return ResultEntry{}, err
	}

	var resultEntry ResultEntry
	found := false
	for rows.Next() {
		err = rows.Scan(
			&resultEntry.Id,
			&resultEntry.Competitionid,
			&resultEntry.Userid,
			&resultEntry.Eventid,
			&resultEntry.Solve1,
			&resultEntry.Solve2,
			&resultEntry.Solve3,
			&resultEntry.Solve4,
			&resultEntry.Solve5,
			&resultEntry.Comment,
			&resultEntry.Status.Id,
		)
		if err != nil {
			return ResultEntry{}, err
		}
		found = true
	}

	if !found {
		return ResultEntry{}, fmt.Errorf("not found")
	}

	return resultEntry, nil
}

func GetResultEntryById(db *pgxpool.Pool, resultId int) (ResultEntry, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT re.result_id, re.competition_id, re.user_id, re.event_id, re.solve1, re.solve2, re.solve3, re.solve4, re.solve5, re.comment, re.status_id, c.name, e.displayname, rs.approvalfinished, rs.approved, rs.visible, rs.displayname, u.name, e.format, e.iconcode FROM results re JOIN competitions c ON c.competition_id = re.competition_id JOIN events e ON e.event_id = re.event_id JOIN results_status rs ON results_status_id = re.status_id JOIN users u ON u.user_id = re.user_id WHERE re.result_id = $1;`,
		resultId,
	)
	if err != nil {
		return ResultEntry{}, err
	}

	var resultEntry ResultEntry
	found := false
	for rows.Next() {
		err = rows.Scan(
			&resultEntry.Id,
			&resultEntry.Competitionid,
			&resultEntry.Userid,
			&resultEntry.Eventid,
			&resultEntry.Solve1,
			&resultEntry.Solve2,
			&resultEntry.Solve3,
			&resultEntry.Solve4,
			&resultEntry.Solve5,
			&resultEntry.Comment,
			&resultEntry.Status.Id,
			&resultEntry.Competitionname,
			&resultEntry.Eventname,
			&resultEntry.Status.ApprovalFinished,
			&resultEntry.Status.Approved,
			&resultEntry.Status.Visible,
			&resultEntry.Status.Displayname,
			&resultEntry.Username,
			&resultEntry.Format,
			&resultEntry.Iconcode,
		)
		if err != nil {
			return ResultEntry{}, err
		}
		found = true
	}

	if !found {
		return ResultEntry{}, err
	}

	return resultEntry, nil
}

func (r *ResultEntry) FormatMultiSingle(single int) string {
	if utils.ParseMultiToMilliseconds(r.Solve1) == single {
		return r.Solve1
	}
	if utils.ParseMultiToMilliseconds(r.Solve2) == single {
		return r.Solve2
	}
	if utils.ParseMultiToMilliseconds(r.Solve3) == single {
		return r.Solve3
	}
	if utils.ParseMultiToMilliseconds(r.Solve4) == single {
		return r.Solve4
	}
	if utils.ParseMultiToMilliseconds(r.Solve5) == single {
		return r.Solve5
	}

	return "DNS"
}

func (r *ResultEntry) SingleFormatted(isfmc bool, scrambles []string) string {
	single := r.Single(isfmc, scrambles)

	if r.Iconcode == "333mbf" {
		if single == constants.DNF {
			return "DNF"
		}

		return r.FormatMultiSingle(single)
	}

	return utils.FormatTime(single, isfmc)
}

func (r *ResultEntry) AverageFormatted(isfmc bool, scrambles []string) (string, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil {
		return "", err
	}

	return utils.FormatTime(r.Average(noOfSolves, isfmc, scrambles), isfmc), nil
}

func (r *ResultEntry) GetFormattedTimes(isfmc bool, scrambles []string) ([]string, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil {
		return []string{}, err
	}

	solves := r.GetSolves(isfmc, scrambles)
	solves = solves[:noOfSolves]
	if noOfSolves <= 3 {
		if r.Iconcode == "333mbf" {
			return utils.FormatMultiTimes(solves), nil
		}
		return solves, nil
	}

	type SolveTuple struct {
		FormattedTime     string
		TimeInMiliseconds int
		Index             int
	}
	sortedSolves := make([]SolveTuple, 0)

	for idx, val := range solves {
		sortedSolves = append(
			sortedSolves,
			SolveTuple{val, utils.ParseSolveToMilliseconds(val, false, ""), idx},
		)
	}

	sort.Slice(
		sortedSolves,
		func(i int, j int) bool { return sortedSolves[i].TimeInMiliseconds < sortedSolves[j].TimeInMiliseconds },
	)
	solves[sortedSolves[0].Index] = "(" + solves[sortedSolves[0].Index] + ")"
	solves[sortedSolves[len(sortedSolves)-1].Index] = "(" + solves[sortedSolves[len(sortedSolves)-1].Index] + ")"

	return solves, nil
}

func GetFormattedTimes(times []string, format string, scrambles []string) ([]string, error) {
	resultEntry := ResultEntry{
		Format: format,
		Solve1: times[0],
		Solve2: times[1],
		Solve3: times[2],
		Solve4: times[3],
		Solve5: times[4],
	}
	return resultEntry.GetFormattedTimes(resultEntry.IsFMC(), scrambles)
}

func (r *ResultEntry) IsFMC() bool {
	return r.Iconcode == "333fm"
}

func (r *ResultEntry) GetSolveIdx(s string) int {
	if r.Solve1 == s {
		return 0
	}
	if r.Solve2 == s {
		return 1
	}
	if r.Solve3 == s {
		return 2
	}
	if r.Solve4 == s {
		return 3
	}
	if r.Solve5 == s {
		return 4
	}

	return -1
}

func (r *ResultEntry) IsAverageOfX() bool {
	return len(r.Format) > 0 && r.Format[0] == 'a'
}

// r.Format must be set
func (r *ResultEntry) ShowPossibleAverages() (bool, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil {
		return false, err
	}

	if noOfSolves == 1 || !r.IsAverageOfX() {
		return false, nil
	}

	ok := true
	for idx, solve := range []string{r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5} {
		if idx < noOfSolves-1 {
			ok = ok && solve != "DNS"
		} else {
			ok = ok && solve == "DNS"
			break
		}
	}

	return ok, nil
}

func (r *ResultEntry) GetNthSolve(noOfSolves int) string {
	if noOfSolves == 1 {
		return r.Solve1
	}
	if noOfSolves == 2 {
		return r.Solve2
	}
	if noOfSolves == 3 {
		return r.Solve3
	}
	if noOfSolves == 4 {
		return r.Solve4
	}
	return r.Solve5
}

func (r *ResultEntry) SetNthSolve(noOfSolves int, newSolveValue string) {
	if noOfSolves == 1 {
		r.Solve1 = newSolveValue
	}
	if noOfSolves == 2 {
		r.Solve2 = newSolveValue
	}
	if noOfSolves == 3 {
		r.Solve3 = newSolveValue
	}
	if noOfSolves == 4 {
		r.Solve4 = newSolveValue
	}
	r.Solve5 = newSolveValue
}

// load scrambles first into resultEntry.Scrambles
func (r *ResultEntry) GetBPA() (string, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil {
		return "", err
	}
	if ok, _ := r.ShowPossibleAverages(); !ok {
		return "", fmt.Errorf("did not finish first %d solves", noOfSolves)
	}

	oldSolveN := r.GetNthSolve(noOfSolves)
	single := r.SingleFormatted(r.IsFMC(), r.Scrambles)
	r.SetNthSolve(noOfSolves, single)

	average, err := r.AverageFormatted(r.IsFMC(), r.Scrambles)
	if err != nil {
		return "", err
	}

	r.SetNthSolve(noOfSolves, oldSolveN)

	return average, nil
}

func (r *ResultEntry) GetWPA() (string, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil {
		return "", err
	}
	if ok, _ := r.ShowPossibleAverages(); !ok {
		return "", fmt.Errorf("did not finish first %d solves", noOfSolves)
	}

	oldSolveN := r.GetNthSolve(noOfSolves)
	r.SetNthSolve(noOfSolves, "DNF")

	average, err := r.AverageFormatted(r.IsFMC(), r.Scrambles)
	if err != nil {
		return "", err
	}

	r.SetNthSolve(noOfSolves, oldSolveN)

	return average, nil
}

func (r *ResultEntry) FinishedCompeting() (bool, error) {
	noOfSolves, err := utils.GetNoOfSolves(r.Format)
	if err != nil {
		return false, err
	}

	for solveNo := 1; solveNo <= noOfSolves; solveNo++ {
		if r.GetNthSolve(solveNo) == "DNS" {
			return false, nil
		}
	}

	return true, nil
}

func (r *ResultEntry) GetCompetitionPlace(db *pgxpool.Pool) (string, error) {
	results, err := GetResultsFromCompetitionByEventName(db, r.Competitionid, r.Eventid)
	if err != nil {
		return "", err
	}

	for _, result := range results.Results {
		if result.Username == r.Username {
			return utils.PlaceFromDotToEnglish(result.Place), nil
		}
	}

	return "LAST", nil
}

func (r *ResultEntry) SuspicousChangeInResults(previouslySavedTimes []string, noOfSolves int) bool {
	for idx, newTime := range []string{r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5} {
		if idx >= noOfSolves {
			continue
		}

		oldTime := previouslySavedTimes[idx]
		if oldTime != "DNS" && oldTime != newTime {
			return true
		}
	}

	return false
}

func (r *ResultEntry) GetSuspicousChangeTimesHTML(
	previouslySavedTimes []string,
	noOfSolves int,
	oldTimesFormatted, newTimesFormatted []string,
) string {
	prevItems, currItems := make([]string, noOfSolves), make([]string, noOfSolves)

	for idx, newTime := range []string{r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5} {
		if idx >= noOfSolves {
			continue
		}

		oldTime := previouslySavedTimes[idx]
		color := "black"
		if oldTime != "DNS" && oldTime != newTime {
			color = "red"
		}

		prevItems[idx] = "<td style=\"text-align: center; border: 1px solid black; color:" + color + ";\">" + oldTimesFormatted[idx] + "</td>"
		currItems[idx] = "<td style=\"text-align: center; border: 1px solid black; color:" + color + ";\">" + newTimesFormatted[idx] + "</td>"
	}

	return "<table style=\"border: 1px solid black;\"><tr><th style=\"border: 1px solid black;\">Previous times:</th>" + strings.Join(
		prevItems,
		"",
	) + "</tr><tr><th style=\"border: 1px solid black;\">Current times:</th>" + strings.Join(
		currItems,
		"",
	) + "</tr></table>"
}

func (r *ResultEntry) SendSuspicousMail(
	c *gin.Context,
	db *pgxpool.Pool,
	envMap map[string]string,
	previouslySavedTimes []string,
) {
	select {
	case <-c.Request.Context().Done():
		if r.Iconcode == "333fm" {
			log.Println("Change in FMC results. Not sending an email.")
			return
		}

		noOfSolves, err := utils.GetNoOfSolves(r.Format)
		if err != nil {
			log.Println(
				"ERR utils.GetScramGetNoOfSolvesblesByResultEntryId in r.SendSuspicousMail: " + err.Error(),
			)
			return
		}

		suspicousChangeInResults := r.SuspicousChangeInResults(previouslySavedTimes, noOfSolves)
		suspicousResult := !r.Status.ApprovalFinished

		if suspicousResult || suspicousChangeInResults {
			log.Println("Sending email...")

			scrambles, err := utils.GetScramblesByResultEntryId(db, r.Eventid, r.Competitionid)
			if err != nil {
				log.Println(
					"ERR utils.GetScramblesByResultEntryId in r.SendSuspicousMail: " + err.Error(),
				)
				return
			}

			average, err := r.AverageFormatted(r.IsFMC(), scrambles)
			if err != nil {
				log.Println("ERR r.AverageFormatted in r.SendSuspicousMail: " + err.Error())
				return
			}

			newTimesFormatted, err := r.GetFormattedTimes(r.IsFMC(), scrambles)
			if err != nil {
				log.Println("ERR r.GetFormattedTimes in r.SendSuspicousMail: " + err.Error())
				return
			}

			oldTimesFormatted, err := GetFormattedTimes(previouslySavedTimes, r.Format, scrambles)
			if err != nil {
				log.Println("ERR GetFormattedTimes in r.SendSuspicousMail: " + err.Error())
				return
			}

			adminToken, err := utils.CreateToken(
				1,
				envMap["JWT_SECRET_KEY"],
				60*24,
			) // admin token for a day
			if err != nil {
				log.Println("ERR utils.CreateToken in r.SendSuspicousMail: " + err.Error())
				return
			}

			mailSubject := "Suspicous"
			if suspicousResult {
				mailSubject += " result"
			}
			if suspicousChangeInResults {
				if suspicousResult {
					mailSubject += " and"
				}
				mailSubject += " change in results"
			}
			mailSubject += " detected !!!"

			backendEnv := os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV")
			if backendEnv == "development" {
				mailSubject = "DEVELOPMENT: " + mailSubject
			}

			r.Email, err = GetEmailByWCAID(db, r.WcaId)
			if err != nil {
				log.Println("ERR GetEmailByWCAID in r.SendSuspicousMail: " + err.Error())
				return
			}

			content := "<html>" +
				"<head>" +
				"<style>" +
				`.mui-joy-btn { font-size: 0.875rem; box-sizing: border-box; border-radius: 6px; border: none; background-color: transparent; display: inline-flex; align-items: center; justify-content: center; position: relative; text-decoration: none; font-weight: 600; }
					.mui-joy-btn-soft-success { color: #0a470a; background-color: #e3fbe3; }
					.mui-joy-btn-soft-danger { color: #7d1212; background-color: #fce4e4; }` +
				"</style></head><body>" +
				"<b>Username:</b> <a href=\"" + envMap["WEBSITE_HOME"] + "/profile/" + r.WcaId + "\">" + r.Username + "</a><br>" +
				"<b>Email:</b> " + r.Email + "<br>" +
				"<b>Competition:</b> <a href=\"" + envMap["WEBSITE_HOME"] + "/competition/" + r.Competitionid + "\">" + r.Competitionname + "</a><br>" +
				"<b>Event:</b> " + r.Eventname + "<br>" +
				"<b>Single:</b> " + r.SingleFormatted(r.IsFMC(), scrambles) + "<br>" +
				"<b>Average:</b> " + average + "<br>"

			if suspicousChangeInResults {
				suspicousChangeTimesHTML := r.GetSuspicousChangeTimesHTML(
					previouslySavedTimes,
					noOfSolves,
					oldTimesFormatted,
					newTimesFormatted,
				)
				content += suspicousChangeTimesHTML
			} else {
				content += "<b>Times:</b> " + strings.Join(newTimesFormatted, ", ") + "<br>"
			}
			content +=
				"<b>Comment:</b> " + r.Comment + "<br>" +
					"<a class=\"mui-joy-btn mui-joy-btn-soft-danger\" style=\"padding:10px;\" " +
					"href=\"" + envMap["MAIL_VALIDATE_URL"] + "?resultId=" + strconv.Itoa(r.Id) + "&verdict=false&atoken=" + adminToken + "\">Deny</a>&nbsp;" +
					"<a class=\"mui-joy-btn mui-joy-btn-soft-success\" style=\"padding:10px;\" " +
					"href=\"" + envMap["MAIL_VALIDATE_URL"] + "?resultId=" + strconv.Itoa(r.Id) + "&verdict=true&atoken=" + adminToken + "\">Allow</a><br>" +
					"<span style=\"font-size: 0.5rem\">Token for validating these results will expire in 24 hours.</span>" +
					"</body></html>"

			err = email.SendMail(
				envMap["MAIL_USERNAME"],
				envMap["MAIL_USERNAME"],
				mailSubject,
				content,
				envMap,
			)
			if err != nil {
				log.Println("ERR email.SendMail in r.SendSuspicousMail: " + err.Error())
				return
			}

			log.Println("Successfully sent mail about suspicous result.")
		}
	}

}

func (r *ResultEntry) GetPreviouslySavedTimes(db *pgxpool.Pool) ([]string, error) {
	times := make([]string, 5)

	err := db.QueryRow(context.Background(), `SELECT solve1, solve2, solve3, solve4, solve5 FROM results r WHERE r.result_id = $1;`, r.Id).
		Scan(&times[0], &times[1], &times[2], &times[3], &times[4])
	if err != nil {
		return []string{}, err
	}

	return times, nil
}
