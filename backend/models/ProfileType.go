package models

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type ProfileTypeBasicsRegion struct {
	Name string `json:"name"`
	Iso2 string `json:"iso2"`
}

type ProfileTypeBasics struct {
	Name string `json:"name"`
	Imageurl string `json:"imageurl"`
	Region ProfileTypeBasicsRegion `json:"region"`
	Wcaid string `json:"wcaid"`
	Sex string `json:"sex"`
	NoOfCompetitions int `json:"noOfCompetitions"`
	CompletedSolves int `json:"completedSolves"`
}

type PersonalBestEntry struct {
	NR string `json:"nr"`
	CR string `json:"cr"`
	WR string `json:"wr"`
	Value string `json:"value"`
}

type ProfileTypePersonalBests struct {
	EventId int `json:"eventid"`
	EventName string `json:"eventName"`
	EventIconCode string `json:"eventIconcode"`
	Average PersonalBestEntry `json:"average"`
	Single PersonalBestEntry `json:"single"`
}

type MedalCollection struct {
	Gold int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
}

type RecordCollection struct {
	NR int `json:"nr"`
	CR int `json:"cr"`
	WR int `json:"wr"`
}

type ProfileTypeResultHistoryEntry struct {
	CompetitionId string `json:"competitionId"`
	CompetitionName string `json:"competitionName"`
	Place string `json:"place"`
	Single string `json:"single"`
	Average string `json:"average"`
	Solves []string `json:"solves"`
}

type ProfileTypeResultHistory struct {
	EventId int `json:"eventId"`
	EventName string `json:"eventName"`
	EventIconCode string `json:"eventIconcode"`
	EventFormat string `json:"eventFormat"`
	History []ProfileTypeResultHistoryEntry `json:"history"`
}

type ProfileType struct {
	Basics ProfileTypeBasics `json:"basics"`
	PersonalBests []ProfileTypePersonalBests `json:"personalBests"`
	MedalCollection MedalCollection `json:"medalCollection"`
	RecordCollection RecordCollection `json:"recordCollection"`
	ResultsHistory []ProfileTypeResultHistory `json:"resultsHistory"`
}

func GetNoOfCompetitions(db *pgxpool.Pool, uid int) (int, error) {
	rows, err := db.Query(context.Background(), `SELECT COUNT(*) FROM (SELECT r.competition_id FROM results r WHERE r.user_id = $1 AND ((r.solve1 NOT LIKE 'DNS' AND r.solve1 NOT LIKE 'DNF') OR (r.solve2 NOT LIKE 'DNS' AND r.solve2 NOT LIKE 'DNF') OR (r.solve3 NOT LIKE 'DNS' AND r.solve3 NOT LIKE 'DNF') OR (r.solve4 NOT LIKE 'DNS' AND r.solve4 NOT LIKE 'DNF') OR (r.solve5 NOT LIKE 'DNS' AND r.solve5 NOT LIKE 'DNF')) GROUP BY r.competition_id);`, uid);
	if err != nil { return 0, err }

	var noOfCompetitions int
	for rows.Next() {
		err = rows.Scan(&noOfCompetitions)
		if err != nil { return 0, err }
	}

	return noOfCompetitions, nil
}

func GetCompletedSolves(db *pgxpool.Pool, uid int) (int, error) {
	rows, err := db.Query(context.Background(), `SELECT solve1, solve2, solve3, solve4, solve5 FROM results r WHERE r.user_id = $1;`, uid);
	if err != nil { return 0, err }

	completedSolves := 0
	for rows.Next() {
		solves := make([]string, 5)
		err = rows.Scan(&solves[0], &solves[1], &solves[2], &solves[3], &solves[4])
		if err != nil { return 0, err }

		for _, solve := range solves { if solve != "DNF" && solve != "DNS" { completedSolves++ } }
	}

	return completedSolves, nil
}

func (p *ProfileType) LoadBasics(db *pgxpool.Pool, uid int) (error) {
	rows, err := db.Query(context.Background(), `SELECT u.name, u.avatarurl, c.name, c.iso2, (CASE WHEN u.wcaid LIKE '' THEN u.name ELSE u.wcaid END) AS wcaid, u.sex FROM users u JOIN countries c ON c.country_id = u.country_id WHERE u.user_id = $1;`, uid);
	if err != nil { return err }

	for rows.Next() {
		err = rows.Scan(&p.Basics.Name, &p.Basics.Imageurl, &p.Basics.Region.Name, &p.Basics.Region.Iso2, &p.Basics.Wcaid, &p.Basics.Sex)
		if err != nil { return err }
		if p.Basics.Sex == "m" {p.Basics.Sex = "Male"}
		if p.Basics.Sex == "f" {p.Basics.Sex = "Female"}
		if p.Basics.Sex == "o" {p.Basics.Sex = "?"}

		p.Basics.NoOfCompetitions, err = GetNoOfCompetitions(db, uid)
		if err != nil { return err }

		p.Basics.CompletedSolves, err = GetCompletedSolves(db, uid)
		if err != nil { return err }
	}

	return nil
}

func LoadBestAverage(db *pgxpool.Pool, user User, eid int) (string, error) {
	rows, err := db.Query(context.Background(), `SELECT r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN events e ON e.event_id = r.event_id WHERE r.user_id = $1 AND r.event_id = $2;`, user.Id, eid);
	if err != nil { return "", err }

	average := constants.DNS

	isfmc := false
	for rows.Next() {
		var resultEntry ResultEntry
		err = rows.Scan(&resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format, &resultEntry.Iconcode, &resultEntry.Eventid, &resultEntry.Competitionid)
		if err != nil { return "", err }

		isfmc = resultEntry.IsFMC()

		scrambles, err := utils.GetScramblesByResultEntryId(db, resultEntry.Eventid, resultEntry.Competitionid)
		if err != nil { return "", err }

		averageCandidate, err := resultEntry.AverageFormatted(isfmc, scrambles)
		if err != nil { return "", err }
		utils.CompareSolves(&average, averageCandidate, false, "")
	}

	return utils.FormatTime(average, false), nil
}

func (p *ProfileTypePersonalBests) LoadAverage(db *pgxpool.Pool, user User) (error) {
	average, err := LoadBestAverage(db, user, p.EventId)
	if err != nil { return err }

	if utils.ParseSolveToMilliseconds(average, false, "") >= constants.VERY_SLOW { return err }

	nrRank, err := LoadNRRank(db, user, average, 1, p.EventId)
	if err != nil { return err }

	crRank, err := LoadCRRank(db, user, average, 1, p.EventId)
	if err != nil { return err }

	wrRank, err := LoadWRRank(db, average, 1, p.EventId)
	if err != nil { return err }

	p.Average.Value = average
	p.Average.NR = nrRank
	p.Average.CR = crRank
	p.Average.WR = wrRank

	return nil
}

func LoadBestSingle(db *pgxpool.Pool, user User, eid int) (string, error) {
	rows, err := db.Query(context.Background(), `SELECT r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN events e ON e.event_id = r.event_id WHERE r.user_id = $1 AND r.event_id = $2;`, user.Id, eid);
	if err != nil { return "", err }

	single := constants.DNS
	formattedSingle := "DNS"

	for rows.Next() {
		var resultEntry ResultEntry
		err = rows.Scan(&resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format, &resultEntry.Iconcode, &resultEntry.Eventid, &resultEntry.Competitionid)
		if err != nil { return "", err }

		scrambles, err := utils.GetScramblesByResultEntryId(db, resultEntry.Eventid, resultEntry.Competitionid)
		if err != nil { return "", err }

		utils.CompareSolves(&single, resultEntry.SingleFormatted(resultEntry.IsFMC(), scrambles), false, "")
		if single == resultEntry.Single(resultEntry.IsFMC(), scrambles) {
			formattedSingle = resultEntry.SingleFormatted(resultEntry.IsFMC(), scrambles)
		}
	}

	return formattedSingle, nil
}

func LoadRankFromRows(rows pgx.Rows, result string, average int, db *pgxpool.Pool) (string, error) {
	results := make(map[int]int)

	isfmc := false
	for rows.Next() {
		var resultEntry ResultEntry
		err := rows.Scan(&resultEntry.Userid, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format, &resultEntry.Iconcode, &resultEntry.Eventid, &resultEntry.Competitionid)
		if err != nil { return "", err }

		val, ok := results[resultEntry.Userid]; 
		if !ok { val = constants.DNS }

		isfmc = resultEntry.IsFMC()

		scrambles, err := utils.GetScramblesByResultEntryId(db, resultEntry.Eventid, resultEntry.Competitionid)
		if err != nil { return "", err }

		if average == 0 {
			utils.CompareSolves(&val, resultEntry.SingleFormatted(isfmc, scrambles), false, "")
		} else { 
			tmpAverageFormatted, err := resultEntry.AverageFormatted(isfmc, scrambles)
			if err != nil { return "", err }
			utils.CompareSolves(&val, tmpAverageFormatted, false, "")
		}
		results[resultEntry.Userid] = val
	}

	resultsArr := make([]int, 0)
	for _, _result := range results {
		if _result < constants.VERY_SLOW {
			resultsArr = append(resultsArr, _result)
		}
	}

	sort.Slice(resultsArr, func (i int, j int) bool { return resultsArr[i] < resultsArr[j] })
	
	resultInMili := utils.ParseSolveToMilliseconds(result, false, "")
	rank := 1
	for ; rank <= len(resultsArr) && resultsArr[rank - 1] < resultInMili; rank++ {}

	return fmt.Sprint(rank), nil
}

func LoadNRRank(db *pgxpool.Pool, user User, result string, average int, eid int) (string, error) {
	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN events e ON e.event_id = r.event_id JOIN users u ON r.user_id = u.user_id JOIN countries c ON c.country_id = u.country_id WHERE u.country_id = $1 AND r.event_id = $2;`, user.CountryId, eid);
	if err != nil { return "", err }
	
	return LoadRankFromRows(rows, result, average, db)
}

func LoadCRRank(db *pgxpool.Pool, user User, result string, average int, eid int) (string, error) {
	err := user.LoadContinent(db)
	if err != nil { return "", err }

	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN events e ON e.event_id = r.event_id JOIN users u ON r.user_id = u.user_id JOIN countries c ON c.country_id = u.country_id JOIN continents con ON con.continent_id = c.continent_id WHERE c.continent_id = $1 AND r.event_id = $2;`, user.ContinentId, eid);
	if err != nil { return "", err }

	return LoadRankFromRows(rows, result, average, db)
}

func LoadWRRank(db *pgxpool.Pool, result string, average int, eid int) (string, error) {
	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN events e ON e.event_id = r.event_id JOIN users u ON r.user_id = u.user_id WHERE r.event_id = $1;`, eid);
	if err != nil { return "", err }

	return LoadRankFromRows(rows, result, average, db)
}

func (p *ProfileTypePersonalBests) LoadSingle(db *pgxpool.Pool, user User) (error) {
	single, err := LoadBestSingle(db, user, p.EventId)
	if err != nil { return err }

	if utils.ParseSolveToMilliseconds(single, false, "") >= constants.VERY_SLOW { return err }

	nrRank, err := LoadNRRank(db, user, single, 0, p.EventId)
	if err != nil { return err }

	crRank, err := LoadCRRank(db, user, single, 0, p.EventId)
	if err != nil { return err }

	wrRank, err := LoadWRRank(db, single, 0, p.EventId)
	if err != nil { return err }

	p.Single.Value = single
	p.Single.NR = nrRank
	p.Single.CR = crRank
	p.Single.WR = wrRank

	return nil
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
 
func (p *ProfileType) LoadPersonalBests(db *pgxpool.Pool, user User) (error) {
	rows, err := db.Query(context.Background(), `SELECT e.fulldisplayname, e.iconcode, e.event_id FROM results r JOIN events e ON e.event_id = r.event_id WHERE r.user_id = $1 GROUP BY e.fulldisplayname, e.iconcode, e.event_id ORDER BY e.event_id;`, user.Id);
	if err != nil { return err }

	p.PersonalBests = make([]ProfileTypePersonalBests, 0)
	for rows.Next() {
		var pbEntry ProfileTypePersonalBests
		err = rows.Scan(&pbEntry.EventName, &pbEntry.EventIconCode, &pbEntry.EventId)
		if err != nil { return err }
		
		p.PersonalBests = append(p.PersonalBests, pbEntry)
	}
	
	newPersonalBests := make([]ProfileTypePersonalBests, 0)
	for idx := range p.PersonalBests {
		ismbld := p.PersonalBests[idx].EventIconCode == "333mbf"
		if !ismbld {
			err = p.PersonalBests[idx].LoadAverage(db, user)
			if err != nil { return err }
		}
		
		err = p.PersonalBests[idx].LoadSingle(db, user)
		if err != nil { return err }

		if utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Single.Value, false, "") >= constants.VERY_SLOW && (ismbld || utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Average.Value, false, "") >= constants.VERY_SLOW) {
			continue			
		}

		if utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Single.Value, false, "") >= constants.VERY_SLOW {
			p.PersonalBests[idx].ClearSingle()
		} else if ismbld || utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Average.Value, false, "") >= constants.VERY_SLOW {
			p.PersonalBests[idx].ClearAverage()
		}

		newPersonalBests = append(newPersonalBests, p.PersonalBests[idx])
	}

	p.PersonalBests = newPersonalBests

	return nil
}

func ComputePlacement(db *pgxpool.Pool, uname string, cid string, eid int) (string, error) {
	competitionResults, err := GetResultsFromCompetitionByEventName(db, cid, eid)
	if err != nil { return "", err }

	placement := 1
	for ; placement <= len(competitionResults) && competitionResults[placement - 1].Username != uname; placement++ {}

	return fmt.Sprint(placement), nil
}

func CreateEventHistoryForUser(db *pgxpool.Pool, user User, event CompetitionEvent, p *ProfileType) (ProfileTypeResultHistory, error) {
	var history ProfileTypeResultHistory

	history.EventId = event.Id
	history.EventName = event.Fulldisplayname
	history.EventIconCode = event.Iconcode
	history.EventFormat = event.Format

	rows, err := db.Query(context.Background(), `SELECT r.competition_id, c.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, u.name, e.format, e.iconcode, r.event_id FROM results r JOIN competitions c ON c.competition_id = r.competition_id JOIN users u ON u.user_id = r.user_id JOIN events e ON e.event_id = r.event_id WHERE r.user_id = $1 AND r.event_id = $2 ORDER BY c.enddate DESC;`, user.Id, event.Id)
	if err != nil { return ProfileTypeResultHistory{}, err }

	history.History = make([]ProfileTypeResultHistoryEntry, 0)
	for rows.Next() {
		var resultEntry ResultEntry
		err = rows.Scan(&resultEntry.Competitionid, &resultEntry.Competitionname, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Username, &resultEntry.Format, &resultEntry.Iconcode, &resultEntry.Eventid)
		if err != nil { return ProfileTypeResultHistory{}, err }

		scrambles, err := utils.GetScramblesByResultEntryId(db, resultEntry.Eventid, resultEntry.Competitionid)
		if err != nil { return ProfileTypeResultHistory{}, err }

		ismbld := resultEntry.Iconcode == "333mbf"

		var historyEntry ProfileTypeResultHistoryEntry
		historyEntry.CompetitionId = resultEntry.Competitionid
		historyEntry.CompetitionName = resultEntry.Competitionname
		historyEntry.Single = resultEntry.SingleFormatted(resultEntry.IsFMC(), scrambles)
		if historyEntry.Single == "DNS" { continue }
		if !ismbld {
			historyEntry.Average, err = resultEntry.AverageFormatted(resultEntry.IsFMC(), scrambles)
			if err != nil { return ProfileTypeResultHistory{}, err }
		}
		historyEntry.Solves, err = resultEntry.GetFormattedTimes(resultEntry.IsFMC(), scrambles)
		if err != nil { return ProfileTypeResultHistory{}, err }
		historyEntry.Place, err = ComputePlacement(db, resultEntry.Username, resultEntry.Competitionid, event.Id)
		if err != nil { return ProfileTypeResultHistory{}, err }
		
		canIncreaseMedalCount := (event.Format[0] == 'b' && utils.ParseSolveToMilliseconds(historyEntry.Single, false, "") < constants.VERY_SLOW) || ((!ismbld && event.Format[0] != 'b' && utils.ParseSolveToMilliseconds(historyEntry.Average, false, "") < constants.VERY_SLOW))
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

		history.History = append(history.History, historyEntry)
	}

	return history, nil
}

func (p *ProfileType) LoadHistory(db *pgxpool.Pool, user User) (error) {
	events, err := GetAvailableEvents(db)
	if err != nil { return err }

	p.ResultsHistory = make([]ProfileTypeResultHistory, 0)
	for _, event := range events {
		eventHistory, err := CreateEventHistoryForUser(db, user, event, p)
		if err != nil { return err }
		if len(eventHistory.History) > 0 {
			p.ResultsHistory = append(p.ResultsHistory, eventHistory)
		}
	}

	return nil
}

type RecordersEntry struct {
	time time.Time
	single int
	singleRecorders []int
	average int
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

func UpdateRecordersEntry(oldRecordersEntry *RecordersEntry, newRecordersEntry RecordersEntry, uid int, records *int, ismbld bool) {
	if newRecordersEntry.single <= oldRecordersEntry.single { IsRecorder2(newRecordersEntry.singleRecorders, uid, records) }

	if newRecordersEntry.single < oldRecordersEntry.single {
		oldRecordersEntry.single = newRecordersEntry.single
		oldRecordersEntry.singleRecorders = newRecordersEntry.singleRecorders
	} else if newRecordersEntry.single == oldRecordersEntry.single {
		oldRecordersEntry.singleRecorders = append(oldRecordersEntry.singleRecorders, newRecordersEntry.singleRecorders...)
	}
	
	if !ismbld {
		if newRecordersEntry.average <= oldRecordersEntry.average { IsRecorder2(newRecordersEntry.averageRecorders, uid, records) }

		if newRecordersEntry.average < oldRecordersEntry.average {
			oldRecordersEntry.average = newRecordersEntry.average
			oldRecordersEntry.averageRecorders = newRecordersEntry.averageRecorders
		} else if newRecordersEntry.average == oldRecordersEntry.average {
			oldRecordersEntry.averageRecorders = append(oldRecordersEntry.averageRecorders, newRecordersEntry.averageRecorders...)
		}
	}
}

func CountRecordsInEventFromRows(rows pgx.Rows, uid int, db *pgxpool.Pool) (int, error) {
	recorders := make(map[time.Time]RecordersEntry)

	ismbld := false
	for rows.Next() {
		var resultEntry ResultEntry
		var date time.Time
		err := rows.Scan(&resultEntry.Userid, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &date, &resultEntry.Format, &resultEntry.Iconcode, &resultEntry.Eventid, &resultEntry.Competitionid)
		if err != nil { return 0, err }

		scrambles, err := utils.GetScramblesByResultEntryId(db, resultEntry.Eventid, resultEntry.Competitionid)
		if err != nil { return 0, err }
		resultEntry.Scrambles = scrambles

		single := resultEntry.SingleFormatted(resultEntry.IsFMC(), resultEntry.Scrambles)
		singleMili := utils.ParseSolveToMilliseconds(single, false, "")
		
		ismbld = resultEntry.Iconcode == "333mbf"

		var averageMili int
		if !ismbld {
			average, err := resultEntry.AverageFormatted(resultEntry.IsFMC(), resultEntry.Scrambles)
			if err != nil { return 0, err }
			averageMili = utils.ParseSolveToMilliseconds(average, false, "")
		}
		
		recordersEntry, ok := recorders[date]
		if !ok {
			recorders[date] = RecordersEntry{date, constants.DNS, make([]int, 0), constants.DNS, make([]int ,0)}
			recordersEntry = recorders[date]
		}

		if singleMili < constants.VERY_SLOW {
			if singleMili < recordersEntry.single {
				recordersEntry.single = singleMili
				recordersEntry.singleRecorders = make([]int, 0)
			}
			if singleMili <= recordersEntry.single { recordersEntry.singleRecorders = append(recordersEntry.singleRecorders, resultEntry.Userid) }
		}

		if !ismbld {
			if averageMili < constants.VERY_SLOW {
				if averageMili < recordersEntry.average {
					recordersEntry.average = averageMili
					recordersEntry.averageRecorders = make([]int, 0)
				}
				if averageMili <= recordersEntry.average { recordersEntry.averageRecorders = append(recordersEntry.averageRecorders, resultEntry.Userid) }
			}
		}

		recorders[date] = recordersEntry
	}
	
	recordersArr := make([]RecordersEntry, 0)
	for _, v := range recorders { recordersArr = append(recordersArr, v) }
	sort.Slice(recordersArr, func (i int, j int) bool { return recordersArr[i].time.Before(recordersArr[j].time) })

	if len(recorders) == 0 { return 0, nil }
	
	records := 0

	recordersEntry := recordersArr[0]
	IsRecorder(&recordersEntry, uid, &records)
	
	for idx := 1; idx < len(recordersArr); idx++ {
		UpdateRecordersEntry(&recordersEntry, recordersArr[idx], uid, &records, ismbld)
	}

	return records, nil
}

func CountWRs(db *pgxpool.Pool, user User, eid int) (int, error) {
	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, c.enddate, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN competitions c ON c.competition_id = r.competition_id JOIN users u ON u.user_id = r.user_id JOIN events e ON e.event_id = r.event_id WHERE r.event_id = $1;`, eid);
	if err != nil { return 0, err }

	return CountRecordsInEventFromRows(rows, user.Id, db)
}

func CountCRs(db *pgxpool.Pool, user User, eid int) (int, error) {
	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, c.enddate, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN competitions c ON c.competition_id = r.competition_id JOIN users u ON u.user_id = r.user_id JOIN events e ON e.event_id = r.event_id JOIN countries ON countries.country_id = u.country_id WHERE r.event_id = $1 AND countries.continent_id = $2;`, eid, user.ContinentId);
	if err != nil { return 0, err }

	return CountRecordsInEventFromRows(rows, user.Id, db)
}

func CountNRsInEvent(db *pgxpool.Pool, user User, eid int) (int, error) {
	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, c.enddate, e.format, e.iconcode, r.event_id, r.competition_id FROM results r JOIN competitions c ON c.competition_id = r.competition_id JOIN users u ON u.user_id = r.user_id JOIN events e ON e.event_id = r.event_id WHERE r.event_id = $1 AND u.country_id = $2;`, eid, user.CountryId);
	if err != nil { return 0, err }

	return CountRecordsInEventFromRows(rows, user.Id, db)
}

func (p *ProfileType) LoadRecordCollection(db *pgxpool.Pool, user User) (error) {
	for _, history := range p.ResultsHistory {
		wr, err := CountWRs(db, user, history.EventId)
		if err != nil { return err }
		p.RecordCollection.WR += wr

		cr, err := CountCRs(db, user, history.EventId)
		if err != nil { return err }
		p.RecordCollection.CR += cr

		nr, err := CountNRsInEvent(db, user, history.EventId)
		if err != nil { return err }
		p.RecordCollection.NR += nr
	}

	return nil
}

func (p *ProfileType) Load(db *pgxpool.Pool, uid int) (error) {
	err := p.LoadBasics(db, uid)
	if err != nil { return err }

	user, err := GetUserById(db, uid)
	if err != nil { return err }
	err = user.LoadContinent(db)
	if err != nil { return err }

	err = p.LoadPersonalBests(db, user)
	if err != nil { return err }

	err = p.LoadHistory(db, user)
	if err != nil { return err }

	err = p.LoadRecordCollection(db, user)
	if err != nil { return err }

	return nil
}
