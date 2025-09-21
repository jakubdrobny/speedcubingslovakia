package models

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type ProfileTypeBasicsRegion struct {
	Name string `json:"name"`
	Iso2 string `json:"iso2"`
}

type ProfileTypeBasics struct {
	Name             string                  `json:"name"`
	Imageurl         string                  `json:"imageurl"`
	Region           ProfileTypeBasicsRegion `json:"region"`
	Wcaid            string                  `json:"wcaid"`
	Sex              string                  `json:"sex"`
	NoOfCompetitions int                     `json:"noOfCompetitions"`
	CompletedSolves  int                     `json:"completedSolves"`
}

type PersonalBestEntry struct {
	NR    string `json:"nr"`
	CR    string `json:"cr"`
	WR    string `json:"wr"`
	Value string `json:"value"`
}

type ProfileTypePersonalBests struct {
	EventId       int               `json:"eventid"`
	EventName     string            `json:"eventName"`
	EventIconCode string            `json:"eventIconcode"`
	Average       PersonalBestEntry `json:"average"`
	Single        PersonalBestEntry `json:"single"`
	Event         CompetitionEvent  `json:"-"`
}

type MedalCollection struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
}

type RecordCollection struct {
	NR int `json:"nr"`
	CR int `json:"cr"`
	WR int `json:"wr"`
}

type ProfileTypeResultHistoryEntry struct {
	CompetitionId      string    `json:"competitionId"`
	CompetitionName    string    `json:"competitionName"`
	CompetitionEnddate time.Time `json:"-"`
	Place              string    `json:"place"`
	Single             string    `json:"single"`
	Average            string    `json:"average"`
	Solves             []string  `json:"solves"`
	SingleRecord       string    `json:"singleRecord"`
	AverageRecord      string    `json:"averageRecord"`
	SingleRecordColor  string    `json:"singleRecordColor"`
	AverageRecordColor string    `json:"averageRecordColor"`
}

type ProfileTypeResultHistory struct {
	EventId       int                             `json:"eventId"`
	EventName     string                          `json:"eventName"`
	EventIconCode string                          `json:"eventIconcode"`
	EventFormat   string                          `json:"eventFormat"`
	History       []ProfileTypeResultHistoryEntry `json:"history"`
}

type ProfileType struct {
	Basics           ProfileTypeBasics          `json:"basics"`
	PersonalBests    []ProfileTypePersonalBests `json:"personalBests"`
	MedalCollection  MedalCollection            `json:"medalCollection"`
	RecordCollection RecordCollection           `json:"recordCollection"`
	ResultsHistory   []ProfileTypeResultHistory `json:"resultsHistory"`
}

func GetNoOfCompetitions(db *pgxpool.Pool, uid int) (int, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT COUNT(*) FROM (SELECT r.competition_id FROM results r WHERE r.user_id = $1 AND ((r.solve1 NOT LIKE 'DNS' AND r.solve1 NOT LIKE 'DNF') OR (r.solve2 NOT LIKE 'DNS' AND r.solve2 NOT LIKE 'DNF') OR (r.solve3 NOT LIKE 'DNS' AND r.solve3 NOT LIKE 'DNF') OR (r.solve4 NOT LIKE 'DNS' AND r.solve4 NOT LIKE 'DNF') OR (r.solve5 NOT LIKE 'DNS' AND r.solve5 NOT LIKE 'DNF')) GROUP BY r.competition_id);`,
		uid,
	)
	if err != nil {
		return 0, err
	}

	var noOfCompetitions int
	for rows.Next() {
		err = rows.Scan(&noOfCompetitions)
		if err != nil {
			return 0, err
		}
	}

	return noOfCompetitions, nil
}

func GetCompletedSolves(db *pgxpool.Pool, uid int) (int, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT solve1, solve2, solve3, solve4, solve5 FROM results r WHERE r.user_id = $1;`,
		uid,
	)
	if err != nil {
		return 0, err
	}

	completedSolves := 0
	for rows.Next() {
		solves := make([]string, 5)
		err = rows.Scan(&solves[0], &solves[1], &solves[2], &solves[3], &solves[4])
		if err != nil {
			return 0, err
		}

		for _, solve := range solves {
			if solve != "DNF" && solve != "DNS" {
				completedSolves++
			}
		}
	}

	return completedSolves, nil
}

func (p *ProfileType) LoadBasics(db *pgxpool.Pool, uid int) error {
	rows, err := db.Query(
		context.Background(),
		`SELECT u.name, u.avatarurl, c.name, c.iso2, (CASE WHEN u.wcaid LIKE '' THEN u.name ELSE u.wcaid END) AS wcaid, u.sex FROM users u JOIN countries c ON c.country_id = u.country_id WHERE u.user_id = $1;`,
		uid,
	)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(
			&p.Basics.Name,
			&p.Basics.Imageurl,
			&p.Basics.Region.Name,
			&p.Basics.Region.Iso2,
			&p.Basics.Wcaid,
			&p.Basics.Sex,
		)
		if err != nil {
			return err
		}
		if p.Basics.Sex == "m" {
			p.Basics.Sex = "Male"
		}
		if p.Basics.Sex == "f" {
			p.Basics.Sex = "Female"
		}
		if p.Basics.Sex == "o" {
			p.Basics.Sex = "?"
		}

		p.Basics.NoOfCompetitions, err = GetNoOfCompetitions(db, uid)
		if err != nil {
			return err
		}

		p.Basics.CompletedSolves, err = GetCompletedSolves(db, uid)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadBestSingleAndAverage(
	db *pgxpool.Pool,
	resultEntries *[]ResultEntry,
	ismbld bool,
) (string, string, error) {
	single := constants.DNS
	average := constants.DNS
	formattedSingle := "DNS"
	var err error

	for _, resultEntry := range *resultEntries {
		isfmc := resultEntry.IsFMC()
		scrambles := make([]string, 5)
		if isfmc {
			scrambles, err = utils.GetScramblesByResultEntryId(
				db,
				resultEntry.Eventid,
				resultEntry.Competitionid,
			)
			if err != nil {
				return "", "", err
			}
		}

		utils.CompareSolves(&single, resultEntry.SingleFormatted(isfmc, scrambles), false, "")
		if single == resultEntry.Single(resultEntry.IsFMC(), scrambles) {
			formattedSingle = resultEntry.SingleFormatted(resultEntry.IsFMC(), scrambles)
		}

		if !ismbld {
			averageCandidate, err := resultEntry.AverageFormatted(isfmc, scrambles)
			if err != nil {
				return "", "", err
			}
			utils.CompareSolves(&average, averageCandidate, false, "")
		}
	}

	return formattedSingle, utils.FormatTime(average, false), nil
}

type ResultsSingleAverageEntry struct {
	Single      int
	Average     int
	ContinentId string
	CountryId   string
}

func ProcessResultEntryToLoadRank(
	eventResultsRow *EventResultsRow,
	db *pgxpool.Pool,
	user *User,
	results map[int]ResultsSingleAverageEntry,
	ismbld bool,
) error {
	resultEntry := eventResultsRow.ResultEntry
	val, ok := results[resultEntry.Userid]
	if !ok {
		val.Single = constants.DNS
		val.Average = constants.DNS
	}

	var err error

	isfmc := resultEntry.IsFMC()
	scrambles := make([]string, 5)
	if isfmc {
		scrambles, err = utils.GetScramblesByResultEntryId(
			db,
			resultEntry.Eventid,
			resultEntry.Competitionid,
		)
		if err != nil {
			return err
		}
	}

	utils.CompareSolves(&val.Single, resultEntry.SingleFormatted(isfmc, scrambles), false, "")
	if !ismbld {
		tmpAverageFormatted, err := resultEntry.AverageFormatted(isfmc, scrambles)
		if err != nil {
			return err
		}
		utils.CompareSolves(&val.Average, tmpAverageFormatted, false, "")
	}
	val.ContinentId = eventResultsRow.ContinentId
	val.CountryId = eventResultsRow.CountryId
	results[resultEntry.Userid] = val

	return nil
}

type RankResultsEntry struct {
	Value       int
	ContinentId int
	CountryId   int
}

func GetRankFromResults(
	results map[int]ResultsSingleAverageEntry,
	single string,
	average string,
	user *User,
) (PersonalBestRanks, error) {
	type ResultsArrEntry struct {
		Value       int
		ContinentId string
		CountryId   string
	}

	resultsArrSingle := make([]ResultsArrEntry, 0)
	resultsArrAverage := make([]ResultsArrEntry, 0)
	for _, _result := range results {
		if _result.Single < constants.VERY_SLOW {
			resultsArrSingle = append(resultsArrSingle, ResultsArrEntry{
				Value:       _result.Single,
				ContinentId: _result.ContinentId,
				CountryId:   _result.CountryId,
			})
		}
		if _result.Average < constants.VERY_SLOW {
			resultsArrAverage = append(resultsArrAverage, ResultsArrEntry{
				Value:       _result.Average,
				ContinentId: _result.ContinentId,
				CountryId:   _result.CountryId,
			})
		}
	}

	sort.Slice(
		resultsArrSingle,
		func(i int, j int) bool { return resultsArrSingle[i].Value < resultsArrSingle[j].Value },
	)
	sort.Slice(
		resultsArrAverage,
		func(i int, j int) bool { return resultsArrAverage[i].Value < resultsArrAverage[j].Value },
	)

	singleResultInMili := utils.ParseSolveToMilliseconds(single, false, "")
	averageResultInMili := utils.ParseSolveToMilliseconds(average, false, "")

	personalBestRanks := PersonalBestRanks{
		Single: Ranks{
			NR: "1",
			CR: "1",
			WR: "1",
		},
		Average: Ranks{
			NR: "1",
			CR: "1",
			WR: "1",
		},
	}

	for nrIdx, crIdx, wrIdx, curIdx, nrPos, crPos, wrPos := -1, -1, 0, 1, 1, 1, 1; curIdx < len(resultsArrSingle) && resultsArrSingle[curIdx].Value <= singleResultInMili; curIdx++ {
		if nrIdx == -1 && resultsArrSingle[curIdx-1].CountryId == user.CountryId {
			nrIdx = curIdx - 1
		}
		if crIdx == -1 && resultsArrSingle[curIdx-1].ContinentId == user.ContinentId {
			crIdx = curIdx - 1
		}

		if nrIdx != -1 && resultsArrSingle[nrIdx].Value < resultsArrSingle[curIdx].Value &&
			resultsArrSingle[curIdx].CountryId == user.CountryId {
			personalBestRanks.Single.NR = fmt.Sprint(nrPos + 1)
			nrIdx = curIdx
			nrPos++
		}
		if crIdx != -1 && resultsArrSingle[crIdx].Value < resultsArrSingle[curIdx].Value &&
			resultsArrSingle[curIdx].ContinentId == user.ContinentId {
			personalBestRanks.Single.CR = fmt.Sprint(crPos + 1)
			crIdx = curIdx
			crPos++
		}
		if resultsArrSingle[wrIdx].Value < resultsArrSingle[curIdx].Value {
			personalBestRanks.Single.WR = fmt.Sprint(wrPos + 1)
			wrIdx = curIdx
			wrPos++
		}
	}

	for nrIdx, crIdx, wrIdx, curIdx, nrPos, crPos, wrPos := -1, -1, 0, 1, 1, 1, 1; curIdx < len(resultsArrAverage) && resultsArrAverage[curIdx].Value <= averageResultInMili; curIdx++ {
		if nrIdx == -1 && resultsArrAverage[curIdx-1].CountryId == user.CountryId {
			nrIdx = curIdx - 1
		}
		if crIdx == -1 && resultsArrAverage[curIdx-1].ContinentId == user.ContinentId {
			crIdx = curIdx - 1
		}

		if nrIdx != -1 && resultsArrAverage[nrIdx].Value < resultsArrAverage[curIdx].Value &&
			resultsArrAverage[curIdx].CountryId == user.CountryId {
			personalBestRanks.Average.NR = fmt.Sprint(nrPos + 1)
			nrIdx = curIdx
			nrPos++
		}
		if crIdx != -1 && resultsArrAverage[crIdx].Value < resultsArrAverage[curIdx].Value &&
			resultsArrAverage[curIdx].ContinentId == user.ContinentId {
			personalBestRanks.Average.CR = fmt.Sprint(crPos + 1)
			crIdx = curIdx
			crPos++
		}
		if resultsArrAverage[wrIdx].Value < resultsArrAverage[curIdx].Value {
			personalBestRanks.Average.WR = fmt.Sprint(wrPos + 1)
			wrIdx = curIdx
			wrPos++
		}
	}

	return personalBestRanks, nil
}

func LoadRank(
	single string,
	average string,
	eid int,
	rows *[]EventResultsRow,
	ismbld bool,
	db *pgxpool.Pool,
	user *User,
) (PersonalBestRanks, error) {
	results := make(map[int]ResultsSingleAverageEntry)

	for _, row := range *rows {
		err := ProcessResultEntryToLoadRank(&row, db, user, results, ismbld)
		if err != nil {
			return PersonalBestRanks{}, err
		}
	}

	return GetRankFromResults(results, single, average, user)
}

type EventResultsRow struct {
	ResultEntry ResultEntry
	Date        time.Time
	ContinentId string
	CountryId   string
	CountryName string
	CountryIso2 string
}

func LoadEventRows(db *pgxpool.Pool, eid int) ([]EventResultsRow, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, c.enddate, e.format, e.iconcode, r.event_id, r.competition_id, countries.continent_id, u.country_id, c.name, u.name, u.wcaid, countries.country_id, countries.iso2, rs.visible FROM results r JOIN competitions c ON c.competition_id = r.competition_id JOIN users u ON u.user_id = r.user_id JOIN events e ON e.event_id = r.event_id JOIN results_status rs ON rs.results_status_id = r.status_id JOIN countries countries ON countries.country_id = u.country_id WHERE rs.visible IS TRUE AND r.event_id = $1 ORDER BY c.enddate DESC;`,
		eid,
	)
	if err != nil {
		return []EventResultsRow{}, err
	}

	eventResultsRows := make([]EventResultsRow, 0)
	for rows.Next() {
		var eventResultsRow EventResultsRow
		err := rows.Scan(
			&eventResultsRow.ResultEntry.Userid,
			&eventResultsRow.ResultEntry.Solve1,
			&eventResultsRow.ResultEntry.Solve2,
			&eventResultsRow.ResultEntry.Solve3,
			&eventResultsRow.ResultEntry.Solve4,
			&eventResultsRow.ResultEntry.Solve5,
			&eventResultsRow.Date,
			&eventResultsRow.ResultEntry.Format,
			&eventResultsRow.ResultEntry.Iconcode,
			&eventResultsRow.ResultEntry.Eventid,
			&eventResultsRow.ResultEntry.Competitionid,
			&eventResultsRow.ContinentId,
			&eventResultsRow.CountryId,
			&eventResultsRow.ResultEntry.Competitionname,
			&eventResultsRow.ResultEntry.Username,
			&eventResultsRow.ResultEntry.WcaId,
			&eventResultsRow.CountryName,
			&eventResultsRow.CountryIso2,
			&eventResultsRow.ResultEntry.Status.Visible,
		)
		if err != nil {
			return []EventResultsRow{}, err
		}

		eventResultsRows = append(eventResultsRows, eventResultsRow)
	}

	return eventResultsRows, nil
}

type Ranks struct {
	NR string
	CR string
	WR string
}

type PersonalBestRanks struct {
	Single  Ranks
	Average Ranks
}

func (p *ProfileTypePersonalBests) LoadSingleAndAverage(
	db *pgxpool.Pool,
	user *User,
	resultEntries *[]ResultEntry,
	ismbld bool,
) ([]EventResultsRow, error) {
	single, average, err := LoadBestSingleAndAverage(db, resultEntries, ismbld)
	if err != nil {
		return []EventResultsRow{}, err
	}

	if utils.ParseSolveToMilliseconds(single, false, "") >= constants.VERY_SLOW {
		return []EventResultsRow{}, err
	}

	eventResultsRows, err := LoadEventRows(db, p.EventId)
	if err != nil {
		return []EventResultsRow{}, err
	}

	ranks, err := LoadRank(single, average, p.EventId, &eventResultsRows, ismbld, db, user)
	if err != nil {
		return []EventResultsRow{}, err
	}

	p.Single.Value = single
	p.Single.NR = ranks.Single.NR
	p.Single.CR = ranks.Single.CR
	p.Single.WR = ranks.Single.WR

	if !ismbld && utils.ParseSolveToMilliseconds(average, false, "") < constants.VERY_SLOW {
		p.Average.Value = average
		p.Average.NR = ranks.Average.NR
		p.Average.CR = ranks.Average.CR
		p.Average.WR = ranks.Average.WR
	}

	return eventResultsRows, nil
}

func (p *ProfileTypePersonalBests) ClearSingle() {
	p.Single.Value = ""
	p.Single.NR = ""
	p.Single.CR = ""
	p.Single.WR = ""
}

func (p *ProfileTypePersonalBests) ClearAverage() {
	p.Average.Value = ""
	p.Average.NR = ""
	p.Average.CR = ""
	p.Average.WR = ""
}

func GetPersonalResultEntriesInEvent(db *pgxpool.Pool, uid int, eid int) ([]ResultEntry, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN events e ON e.event_id = r.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE r.user_id = $1 AND r.event_id = $2 AND rs.visible IS TRUE;`,
		uid,
		eid,
	)
	if err != nil {
		return []ResultEntry{}, err
	}

	resultEntries := make([]ResultEntry, 0)

	for rows.Next() {
		var resultEntry ResultEntry
		err = rows.Scan(
			&resultEntry.Solve1,
			&resultEntry.Solve2,
			&resultEntry.Solve3,
			&resultEntry.Solve4,
			&resultEntry.Solve5,
			&resultEntry.Format,
			&resultEntry.Iconcode,
			&resultEntry.Eventid,
			&resultEntry.Competitionid,
		)
		if err != nil {
			return []ResultEntry{}, err
		}

		resultEntries = append(resultEntries, resultEntry)
	}

	return resultEntries, nil
}

// eid = 0 - all events, eid > 0 - only event with that id
func (p *ProfileType) LoadPersonalBests(
	db *pgxpool.Pool,
	user *User,
	eid int,
) (map[int][]EventResultsRow, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT e.fulldisplayname, e.iconcode, e.event_id, e.format, e.displayname FROM results r JOIN events e ON e.event_id = r.event_id WHERE r.user_id = $1 GROUP BY e.fulldisplayname, e.iconcode, e.event_id ORDER BY e.event_id;`,
		user.Id,
	)
	if err != nil {
		return map[int][]EventResultsRow{}, err
	}

	p.PersonalBests = make([]ProfileTypePersonalBests, 0)
	for rows.Next() {
		var pbEntry ProfileTypePersonalBests
		var eventFormat string
		err = rows.Scan(
			&pbEntry.EventName,
			&pbEntry.EventIconCode,
			&pbEntry.EventId,
			&eventFormat,
			&pbEntry.Event.Displayname,
		)
		if err != nil {
			return map[int][]EventResultsRow{}, err
		}

		pbEntry.Event = CompetitionEvent{
			Id:              pbEntry.EventId,
			Fulldisplayname: pbEntry.EventName,
			Displayname:     pbEntry.Event.Displayname,
			Format:          eventFormat,
			Iconcode:        pbEntry.EventIconCode,
		}

		if eid == 0 || pbEntry.EventId == eid {
			p.PersonalBests = append(p.PersonalBests, pbEntry)
		}
	}

	newPersonalBests := make([]ProfileTypePersonalBests, 0)
	eventsResultRows := make(map[int][]EventResultsRow)

	for idx := range p.PersonalBests {
		checkAverage := p.PersonalBests[idx].EventIconCode == "333mbf" ||
			p.PersonalBests[idx].Event.Format == "bo1"
		eid := p.PersonalBests[idx].EventId

		personalResultEntries, err := GetPersonalResultEntriesInEvent(db, user.Id, eid)
		if err != nil {
			return map[int][]EventResultsRow{}, err
		}

		eventResultRows, err := p.PersonalBests[idx].LoadSingleAndAverage(
			db,
			user,
			&personalResultEntries,
			checkAverage,
		)
		if err != nil {
			return map[int][]EventResultsRow{}, err
		}

		if utils.ParseSolveToMilliseconds(
			p.PersonalBests[idx].Single.Value,
			false,
			"",
		) >= constants.VERY_SLOW &&
			(checkAverage || utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Average.Value, false, "") >= constants.VERY_SLOW) {
			continue
		}

		if utils.ParseSolveToMilliseconds(
			p.PersonalBests[idx].Single.Value,
			false,
			"",
		) >= constants.VERY_SLOW {
			p.PersonalBests[idx].ClearSingle()
		} else if checkAverage || p.PersonalBests[idx].Event.Format == "bo1" || utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Average.Value, false, "") >= constants.VERY_SLOW {
			p.PersonalBests[idx].ClearAverage()
		}

		newPersonalBests = append(newPersonalBests, p.PersonalBests[idx])
		eventsResultRows[eid] = eventResultRows
	}

	p.PersonalBests = newPersonalBests

	return eventsResultRows, nil
}

func GetResultsFromCompetitionFromRows(
	rows []EventResultsRow,
	db *pgxpool.Pool,
) ([]CompetitionResult, error) {
	competitionResults := make([]CompetitionResult, 0)
	format := ""
	var err error

	for _, row := range rows {
		resultEntry := row.ResultEntry

		competitionResult := CompetitionResult{
			Username:    resultEntry.Username,
			WcaId:       resultEntry.WcaId,
			CountryName: row.CountryName,
			CountryIso2: row.CountryIso2,
			UserId:      resultEntry.Userid,
		}

		if competitionResult.WcaId == "" {
			competitionResult.WcaId = competitionResult.Username
		}

		if !resultEntry.Competed() || !resultEntry.Status.Visible {
			continue
		}

		scrambles := make([]string, 5)
		if resultEntry.IsFMC() {
			scrambles, err = utils.GetScramblesByResultEntryId(
				db,
				resultEntry.Eventid,
				resultEntry.Competitionid,
			)
			if err != nil {
				return []CompetitionResult{}, err
			}
		}

		competitionResult.Single = resultEntry.SingleFormatted(resultEntry.IsFMC(), scrambles)

		avg, err := resultEntry.AverageFormatted(resultEntry.IsFMC(), scrambles)
		if err != nil {
			return []CompetitionResult{}, err
		}
		competitionResult.Average = avg

		formattedTimes, err := resultEntry.GetFormattedTimes(resultEntry.IsFMC(), scrambles)
		if err != nil {
			return []CompetitionResult{}, err
		}
		competitionResult.Times = formattedTimes

		competitionResults = append(competitionResults, competitionResult)
		format = resultEntry.Format
	}

	if len(format) > 0 {
		sort.Slice(competitionResults, func(i int, j int) bool {
			name1, name2 := competitionResults[i].Username, competitionResults[j].Username
			d := CompareCompetitionResults(competitionResults[i], competitionResults[j], format)
			if d == 0 {
				return name1 < name2
			}
			return d > 0
		})

		AddPlacementToCompetitionResults(competitionResults, format)
	}

	return competitionResults, nil
}

func ComputePlacementForCompetition(
	rows *[]EventResultsRow,
	firstRowIdx int,
	lastRowIdx int,
	uid int,
	eventFormat string,
	db *pgxpool.Pool,
) (string, error) {
	competitionResults, err := GetResultsFromCompetitionFromRows(
		(*rows)[firstRowIdx:lastRowIdx+1],
		db,
	)
	if err != nil {
		return "", err
	}

	placement := 1
	for idx := 0; idx+1 < len(competitionResults) && competitionResults[idx].UserId != uid; idx++ {
		if CompareCompetitionResults(
			competitionResults[idx],
			competitionResults[idx+1],
			eventFormat,
		) > 0 {
			placement = idx + 2
		}
	}

	return fmt.Sprint(placement), nil
}

func AddRecordsToHistory(
	history *ProfileTypeResultHistory,
	recorders Recorders,
	uid int,
	p *ProfileType,
) {
	singleSoFar := "DNS"
	averageSoFar := "DNS"
	checkAverage := history.EventIconCode == "333mbf" || history.EventFormat == "bo1"

	for historyIdx := len(history.History) - 1; historyIdx >= 0; historyIdx-- {
		historyEntry := history.History[historyIdx]
		competitionEndDate := historyEntry.CompetitionEnddate

		currentSingle := utils.ParseSolveToMilliseconds(historyEntry.Single, false, "")
		if currentSingle <= utils.ParseSolveToMilliseconds(singleSoFar, false, "") &&
			currentSingle < constants.VERY_SLOW {
			history.History[historyIdx].SingleRecord = "PR"
			singleSoFar = historyEntry.Single
		}

		currentAverage := utils.ParseSolveToMilliseconds(historyEntry.Average, false, "")
		if !checkAverage &&
			currentAverage <= utils.ParseSolveToMilliseconds(averageSoFar, false, "") &&
			currentAverage < constants.VERY_SLOW {
			history.History[historyIdx].AverageRecord = "PR"
			averageSoFar = historyEntry.Average
		}

		recordersForComp, ok := recorders.NR[competitionEndDate]
		if ok {
			if recordersForComp.single == currentSingle {
				history.History[historyIdx].SingleRecord = "NR"
			}
			if !checkAverage && recordersForComp.average == currentAverage {
				history.History[historyIdx].AverageRecord = "NR"
			}
		}

		recordersForComp, ok = recorders.CR[competitionEndDate]
		if ok {
			if recordersForComp.single == currentSingle {
				history.History[historyIdx].SingleRecord = "CR"
			}
			if !checkAverage && recordersForComp.average == currentAverage {
				history.History[historyIdx].AverageRecord = "CR"
			}
		}

		recordersForComp, ok = recorders.WR[competitionEndDate]
		if ok {
			if recordersForComp.single == currentSingle {
				history.History[historyIdx].SingleRecord = "WR"
			}
			if !checkAverage && recordersForComp.average == currentAverage {
				history.History[historyIdx].AverageRecord = "WR"
			}
		}

		switch history.History[historyIdx].SingleRecord {
		case "WR":
			history.History[historyIdx].SingleRecordColor = constants.WR_COLOR
			p.RecordCollection.WR++
		case "CR":
			history.History[historyIdx].SingleRecordColor = constants.CR_COLOR
			p.RecordCollection.CR++
		case "NR":
			history.History[historyIdx].SingleRecordColor = constants.NR_COLOR
			p.RecordCollection.NR++
		case "PR":
			history.History[historyIdx].SingleRecord = ""
			history.History[historyIdx].SingleRecordColor = constants.PR_COLOR
		}

		if !checkAverage {
			switch history.History[historyIdx].AverageRecord {
			case "WR":
				history.History[historyIdx].AverageRecordColor = constants.WR_COLOR
				p.RecordCollection.WR++
			case "CR":
				history.History[historyIdx].AverageRecordColor = constants.CR_COLOR
				p.RecordCollection.CR++
			case "NR":
				history.History[historyIdx].AverageRecordColor = constants.NR_COLOR
				p.RecordCollection.NR++
			case "PR":
				history.History[historyIdx].AverageRecord = ""
				history.History[historyIdx].AverageRecordColor = constants.PR_COLOR
			}
		}
	}
}

func (p *ProfileType) CreateEventHistoryForUser(
	db *pgxpool.Pool,
	user *User,
	event CompetitionEvent,
	rows []EventResultsRow,
	recorders Recorders,
) error {
	var history ProfileTypeResultHistory

	history.EventId = event.Id
	history.EventName = event.Fulldisplayname
	history.EventIconCode = event.Iconcode
	history.EventFormat = event.Format

	history.History = make([]ProfileTypeResultHistoryEntry, 0)

	lastCompStartingIdx := 0
	hasUser := false
	var historyEntry ProfileTypeResultHistoryEntry
	var err error

	for curIdx, row := range rows {
		resultEntry := row.ResultEntry

		hasAverage := resultEntry.Iconcode == "333mbf" || resultEntry.Format == "bo1"
		hasUser = hasUser || resultEntry.Userid == user.Id
		if resultEntry.Userid == user.Id {
			isfmc := resultEntry.IsFMC()
			scrambles := make([]string, 5)
			if isfmc {
				scrambles, err = utils.GetScramblesByResultEntryId(
					db,
					resultEntry.Eventid,
					resultEntry.Competitionid,
				)
				if err != nil {
					return err
				}
			}

			historyEntry.CompetitionId = resultEntry.Competitionid
			historyEntry.CompetitionName = resultEntry.Competitionname
			historyEntry.Single = resultEntry.SingleFormatted(resultEntry.IsFMC(), scrambles)
			if historyEntry.Single == "DNS" {
				hasUser = false
				continue
			}
			if !hasAverage {
				historyEntry.Average, err = resultEntry.AverageFormatted(
					resultEntry.IsFMC(),
					scrambles,
				)
				if err != nil {
					return err
				}
			}
			historyEntry.Solves, err = resultEntry.GetFormattedTimes(resultEntry.IsFMC(), scrambles)
			if err != nil {
				return err
			}
		}

		if curIdx == len(rows)-1 ||
			resultEntry.Competitionid != rows[curIdx+1].ResultEntry.Competitionid {
			if hasUser {
				historyEntry.Place, err = ComputePlacementForCompetition(
					&rows,
					lastCompStartingIdx,
					curIdx,
					user.Id,
					event.Format,
					db,
				)
				if err != nil {
					return err
				}

				canIncreaseMedalCount := (event.Format[0] == 'b' && utils.ParseSolveToMilliseconds(historyEntry.Single, false, "") < constants.VERY_SLOW) ||
					(!hasAverage && event.Format[0] != 'b' && utils.ParseSolveToMilliseconds(historyEntry.Average, false, "") < constants.VERY_SLOW)
				if canIncreaseMedalCount {
					switch historyEntry.Place {
					case "1":
						p.MedalCollection.Gold++
					case "2":
						p.MedalCollection.Silver++
					case "3":
						p.MedalCollection.Bronze++
					}
				}
				historyEntry.CompetitionEnddate = row.Date

				history.History = append(history.History, historyEntry)
			}

			lastCompStartingIdx = curIdx + 1
			hasUser = false
			historyEntry = ProfileTypeResultHistoryEntry{}
		}
	}

	if len(history.History) > 0 {
		AddRecordsToHistory(&history, recorders, user.Id, p)
		p.ResultsHistory = append(p.ResultsHistory, history)
	}

	return nil
}

func (p *ProfileType) LoadHistory(
	db *pgxpool.Pool,
	user *User,
	recorders map[int]Recorders,
	rows map[int][]EventResultsRow,
) error {
	p.ResultsHistory = make([]ProfileTypeResultHistory, 0)
	for _, personalBestEntry := range p.PersonalBests {
		err := p.CreateEventHistoryForUser(
			db,
			user,
			personalBestEntry.Event,
			rows[personalBestEntry.EventId],
			recorders[personalBestEntry.EventId],
		)
		if err != nil {
			return err
		}
	}

	return nil
}

type RecordersEntry struct {
	time             time.Time
	single           int
	singleRecorders  []int
	average          int
	averageRecorders []int
}

func IsRecorder(recordersEntry *RecordersEntry, uid int, records *int) {
	for _, id := range recordersEntry.singleRecorders {
		if id == uid {
			*records++
			break
		}
	}

	for _, id := range recordersEntry.averageRecorders {
		if id == uid {
			*records++
			break
		}
	}
}

func IsRecorder2(recorders []int, uid int, records *int) {
	for _, id := range recorders {
		if id == uid {
			*records++
			break
		}
	}
}

func IsRecorder3(recorders []int, uid int) bool {
	for _, id := range recorders {
		if id == uid {
			return true
		}
	}

	return false
}

func UpdateRecordersEntry(
	oldRecordersEntry *RecordersEntry,
	newRecordersEntry RecordersEntry,
	uid int,
	records *int,
	ismbld bool,
) {
	if newRecordersEntry.single <= oldRecordersEntry.single {
		IsRecorder2(newRecordersEntry.singleRecorders, uid, records)
	}

	if newRecordersEntry.single < oldRecordersEntry.single {
		oldRecordersEntry.single = newRecordersEntry.single
		oldRecordersEntry.singleRecorders = newRecordersEntry.singleRecorders
	} else if newRecordersEntry.single == oldRecordersEntry.single {
		oldRecordersEntry.singleRecorders = append(oldRecordersEntry.singleRecorders, newRecordersEntry.singleRecorders...)
	}

	if !ismbld {
		if newRecordersEntry.average <= oldRecordersEntry.average {
			IsRecorder2(newRecordersEntry.averageRecorders, uid, records)
		}

		if newRecordersEntry.average < oldRecordersEntry.average {
			oldRecordersEntry.average = newRecordersEntry.average
			oldRecordersEntry.averageRecorders = newRecordersEntry.averageRecorders
		} else if newRecordersEntry.average == oldRecordersEntry.average {
			oldRecordersEntry.averageRecorders = append(oldRecordersEntry.averageRecorders, newRecordersEntry.averageRecorders...)
		}
	}
}

func UpdateRecordersByDate(
	recorders map[time.Time]RecordersEntry,
	date time.Time,
	singleMili int,
	averageMili int,
	ismbld bool,
	uid int,
) {
	recordersEntry, ok := recorders[date]
	if !ok {
		recorders[date] = RecordersEntry{
			date,
			constants.DNS,
			make([]int, 0),
			constants.DNS,
			make([]int, 0),
		}
		recordersEntry = recorders[date]
	}

	if singleMili < constants.VERY_SLOW {
		if singleMili < recordersEntry.single {
			recordersEntry.single = singleMili
			recordersEntry.singleRecorders = make([]int, 0)
		}
		if singleMili <= recordersEntry.single {
			recordersEntry.singleRecorders = append(recordersEntry.singleRecorders, uid)
		}
	}

	if !ismbld {
		if averageMili < constants.VERY_SLOW {
			if averageMili < recordersEntry.average {
				recordersEntry.average = averageMili
				recordersEntry.averageRecorders = make([]int, 0)
			}
			if averageMili <= recordersEntry.average {
				recordersEntry.averageRecorders = append(recordersEntry.averageRecorders, uid)
			}
		}
	}

	recorders[date] = recordersEntry
}

func ProcessRecorders(recorders map[time.Time]RecordersEntry, uid int, ismbld bool) int {
	recordersArr := make([]RecordersEntry, 0)
	for _, v := range recorders {
		recordersArr = append(recordersArr, v)
	}
	sort.Slice(
		recordersArr,
		func(i int, j int) bool { return recordersArr[i].time.Before(recordersArr[j].time) },
	)

	if len(recorders) == 0 {
		return 0
	}

	records := 0

	recordersEntry := recordersArr[0]
	IsRecorder(&recordersEntry, uid, &records)

	for idx := 1; idx < len(recordersArr); idx++ {
		UpdateRecordersEntry(&recordersEntry, recordersArr[idx], uid, &records, ismbld)
	}

	return records
}

type Recorders struct {
	WR map[time.Time]RecordersEntry
	CR map[time.Time]RecordersEntry
	NR map[time.Time]RecordersEntry
}

func (p *ProfileType) CountRecordsInEventFromRows(
	eventResultsRows *[]EventResultsRow,
	user *User,
	db *pgxpool.Pool,
) (Recorders, error) {
	recordersWR := make(map[time.Time]RecordersEntry)
	recordersCR := make(map[time.Time]RecordersEntry)
	recordersNR := make(map[time.Time]RecordersEntry)
	// uid := user.Id
	contId := user.ContinentId
	countryId := user.CountryId
	var lastDateWR, lastDateCR, lastDateNR time.Time

	var err error
	checkAverage := false
	for eventResultRowIdx := len(*eventResultsRows) - 1; eventResultRowIdx >= 0; eventResultRowIdx-- {
		eventResultRow := (*eventResultsRows)[eventResultRowIdx]

		resultEntry := eventResultRow.ResultEntry
		date := eventResultRow.Date
		currentContId := eventResultRow.ContinentId
		currentCountryId := eventResultRow.CountryId

		isfmc := resultEntry.IsFMC()
		scrambles := make([]string, 5)
		if isfmc {
			scrambles, err = utils.GetScramblesByResultEntryId(
				db,
				resultEntry.Eventid,
				resultEntry.Competitionid,
			)
			if err != nil {
				return Recorders{}, err
			}
		}
		resultEntry.Scrambles = scrambles

		single := resultEntry.SingleFormatted(resultEntry.IsFMC(), resultEntry.Scrambles)
		singleMili := utils.ParseSolveToMilliseconds(single, false, "")

		checkAverage = resultEntry.Iconcode == "333mbf" || resultEntry.Format == "bo1"

		var averageMili int
		if !checkAverage {
			average, err := resultEntry.AverageFormatted(resultEntry.IsFMC(), resultEntry.Scrambles)
			if err != nil {
				return Recorders{}, err
			}
			averageMili = utils.ParseSolveToMilliseconds(average, false, "")
		}

		if !lastDateWR.IsZero() {
			recordersWR[date] = recordersWR[lastDateWR]
		}
		UpdateRecordersByDate(
			recordersWR,
			date,
			singleMili,
			averageMili,
			checkAverage,
			resultEntry.Userid,
		)
		lastDateWR = date

		if currentContId == contId {
			if !lastDateCR.IsZero() {
				recordersCR[date] = recordersCR[lastDateCR]
			}
			UpdateRecordersByDate(
				recordersCR,
				date,
				singleMili,
				averageMili,
				checkAverage,
				resultEntry.Userid,
			)
			lastDateCR = date
		}

		if currentCountryId == countryId {
			if !lastDateNR.IsZero() {
				recordersNR[date] = recordersNR[lastDateNR]
			}
			UpdateRecordersByDate(
				recordersNR,
				date,
				singleMili,
				averageMili,
				checkAverage,
				resultEntry.Userid,
			)
			lastDateNR = date
		}
	}

	// p.RecordCollection.WR += ProcessRecorders(recordersWR, uid, ismbld)
	// p.RecordCollection.CR += ProcessRecorders(recordersCR, uid, ismbld)
	// p.RecordCollection.NR += ProcessRecorders(recordersNR, uid, ismbld)

	return Recorders{WR: recordersWR, CR: recordersCR, NR: recordersNR}, nil
}

func (p *ProfileType) LoadRecordCollection(
	db *pgxpool.Pool,
	user *User,
	rows map[int][]EventResultsRow,
) (map[int]Recorders, error) {
	recorders := make(map[int]Recorders, len(p.PersonalBests))

	for _, personalBestEntry := range p.PersonalBests {
		eid := personalBestEntry.EventId
		currentRows, ok := rows[eid]
		if !ok {
			return map[int]Recorders{}, nil
		}

		currentRecorders, err := p.CountRecordsInEventFromRows(&currentRows, user, db)
		if err != nil {
			return map[int]Recorders{}, err
		}
		recorders[eid] = currentRecorders
	}

	p.RecordCollection.NR -= p.RecordCollection.CR
	p.RecordCollection.CR -= p.RecordCollection.WR

	return recorders, nil
}

func (p *ProfileType) Load(db *pgxpool.Pool, uid int) error {
	err := p.LoadBasics(db, uid)
	if err != nil {
		return err
	}

	user, err := GetUserById(db, uid)
	if err != nil {
		return err
	}
	err = user.LoadContinent(db)
	if err != nil {
		return err
	}

	rows, err := p.LoadPersonalBests(db, &user, 0)
	if err != nil {
		return err
	}

	recorders, err := p.LoadRecordCollection(db, &user, rows)
	if err != nil {
		return err
	}

	err = p.LoadHistory(db, &user, recorders, rows)
	if err != nil {
		return err
	}

	return nil
}
