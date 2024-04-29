package models

import (
	"context"
	"fmt"

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
	Gold string `json:"gold"`
	Silver string `json:"silver"`
	Bronze string `json:"bronze"`
}

type RecordCollection struct {
	NR string `json:"nr"`
	CR string `json:"cr"`
	WR string `json:"wr"`
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
	EventName string `json:"eventName"`
	EventIconCode string `json:"eventIconcode"`
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
	rows, err := db.Query(context.Background(), `SELECT u.name, u.avatarurl, c.country_id, c.iso2, u.wcaid, u.sex FROM users u JOIN countries c ON c.country_id = u.country_id WHERE u.user_id = $1;`, uid);
	if err != nil { return err }

	for rows.Next() {
		err = rows.Scan(&p.Basics.Name, &p.Basics.Imageurl, &p.Basics.Region.Name, &p.Basics.Region.Iso2, &p.Basics.Wcaid, &p.Basics.Sex)
		if err != nil { return err }
		if p.Basics.Sex == "m" {p.Basics.Sex = "Male"}
		if p.Basics.Sex == "f" {p.Basics.Sex = "Female"}

		p.Basics.NoOfCompetitions, err = GetNoOfCompetitions(db, uid)
		if err != nil { return err }

		p.Basics.CompletedSolves, err = GetCompletedSolves(db, uid)
		if err != nil { return err }
	}

	return nil
}

func LoadBestAverage(db *pgxpool.Pool, user User, eid int) (string, error) {
	rows, err := db.Query(context.Background(), `SELECT r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format FROM results r JOIN events e ON e.event_id = r.event_id WHERE r.user_id = $1 AND r.event_id = $2;`, user.Id, eid);
	if err != nil { return "", err }

	average := constants.DNS

	for rows.Next() {
		var resultEntry ResultEntry
		err = rows.Scan(&resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format)
		if err != nil { return "", err }

		averageCandidate, err := resultEntry.AverageFormatted()
		if err != nil { return "", err }
		utils.CompareSolves(&average, averageCandidate)
	}

	return utils.FormatTime(average), nil
}

func (p *ProfileTypePersonalBests) LoadAverage(db *pgxpool.Pool, user User) (error) {
	average, err := LoadBestAverage(db, user, p.EventId)
	if err != nil { return err }

	if utils.ParseSolveToMilliseconds(average) >= constants.VERY_SLOW { return err }

	nrRank, err := LoadNRRank(db, user, average, 1, p.EventId)
	if err != nil { return err }

	crRank, err := LoadCRRank(db, user, average, 1, p.EventId)
	if err != nil { return err }

	wrRank, err := LoadWRRank(db, user, average, 1, p.EventId)
	if err != nil { return err }

	p.Average.Value = average
	p.Average.NR = nrRank
	p.Average.CR = crRank
	p.Average.WR = wrRank

	return nil
}

func LoadBestSingle(db *pgxpool.Pool, user User, eid int) (string, error) {
	rows, err := db.Query(context.Background(), `SELECT r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format FROM results r JOIN events e ON e.event_id = r.event_id WHERE r.user_id = $1 AND r.event_id = $2;`, user.Id, eid);
	if err != nil { return "", err }

	single := constants.DNS

	for rows.Next() {
		var resultEntry ResultEntry
		err = rows.Scan(&resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format)
		if err != nil { return "", err }

		utils.CompareSolves(&single, resultEntry.SingleFormatted())
	}

	return utils.FormatTime(single), nil
}

func LoadRankFromRows(rows pgx.Rows, result string, average int) (string, error) {
	results := make(map[int]int)

	for rows.Next() {
		var resultEntry ResultEntry
		err := rows.Scan(&resultEntry.Userid, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Format)
		if err != nil { return "", err }

		val, ok := results[resultEntry.Userid]; 
		if !ok { val = constants.DNS }
		if average == 0 {
			utils.CompareSolves(&val, resultEntry.SingleFormatted())
		} else { 
			tmpAverageFormatted, err := resultEntry.AverageFormatted()
			if err != nil { return "", err }
			utils.CompareSolves(&val, tmpAverageFormatted) 
		}
		results[resultEntry.Userid] = val
	}

	resultsArr := make([]int, 0)
	for _, _result := range results {
		if _result < constants.VERY_SLOW {
			resultsArr = append(resultsArr, _result)
		}
	}

	resultInMili := utils.ParseSolveToMilliseconds(result)
	rank := 1
	for ; rank <= len(resultsArr) && resultsArr[rank - 1] < resultInMili; rank++ {
		fmt.Println(result, resultInMili, resultsArr[rank - 1])
	}

	return fmt.Sprint(rank), nil
}

func LoadNRRank(db *pgxpool.Pool, user User, result string, average int, eid int) (string, error) {
	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format FROM results r JOIN events e ON e.event_id = r.event_id JOIN users u ON r.user_id = u.user_id JOIN countries c ON c.country_id = u.country_id WHERE u.country_id = $1 AND r.user_id = $2 AND r.event_id = $3;`, user.CountryId, user.Id, eid);
	if err != nil { return "", err }
	
	return LoadRankFromRows(rows, result, average)
}

func LoadCRRank(db *pgxpool.Pool, user User, result string, average int, eid int) (string, error) {
	err := user.LoadContinent(db)
	if err != nil { return "", err }

	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format FROM results r JOIN events e ON e.event_id = r.event_id JOIN users u ON r.user_id = u.user_id JOIN countries c ON c.country_id = u.country_id JOIN continents con ON con.continent_id = c.continent_id WHERE c.continent_id = $1 AND r.user_id = $2 AND r.event_id = $3;`, user.ContinentId, user.Id, eid);
	if err != nil { return "", err }

	return LoadRankFromRows(rows, result, average)
}

func LoadWRRank(db *pgxpool.Pool, user User, result string, average int, eid int) (string, error) {
	rows, err := db.Query(context.Background(), `SELECT r.user_id, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format FROM results r JOIN events e ON e.event_id = r.event_id JOIN users u ON r.user_id = u.user_id WHERE r.user_id = $1 AND r.event_id = $2;`, user.Id, eid);
	if err != nil { return "", err }

	return LoadRankFromRows(rows, result, average)
}

func (p *ProfileTypePersonalBests) LoadSingle(db *pgxpool.Pool, user User) (error) {
	single, err := LoadBestSingle(db, user, p.EventId)
	if err != nil { return err }

	if utils.ParseSolveToMilliseconds(single) >= constants.VERY_SLOW { return err }

	nrRank, err := LoadNRRank(db, user, single, 0, p.EventId)
	if err != nil { return err }

	crRank, err := LoadCRRank(db, user, single, 0, p.EventId)
	if err != nil { return err }

	wrRank, err := LoadWRRank(db, user, single, 0, p.EventId)
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
 
func (p *ProfileType) LoadPersonBests(db *pgxpool.Pool, user User) (error) {
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
		err = p.PersonalBests[idx].LoadAverage(db, user)
		if err != nil { return err }
		
		err = p.PersonalBests[idx].LoadSingle(db, user)
		if err != nil { return err }

		if utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Single.Value) >= constants.VERY_SLOW && utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Average.Value) >= constants.VERY_SLOW {
			continue			
		}

		if utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Single.Value) >= constants.VERY_SLOW {
			p.PersonalBests[idx].ClearSingle()
		} else if utils.ParseSolveToMilliseconds(p.PersonalBests[idx].Average.Value) >= constants.VERY_SLOW {
			p.PersonalBests[idx].ClearAverage()
		}

		newPersonalBests = append(newPersonalBests, p.PersonalBests[idx])
	}

	p.PersonalBests = newPersonalBests

	return nil
}

func (p *ProfileType) Load(db *pgxpool.Pool, uid int) (error) {
	err := p.LoadBasics(db, uid)
	if err != nil { return err }

	user, err := GetUserById(db, uid)
	if err != nil { return err }

	err = p.LoadPersonBests(db, user)
	if err != nil { return err }

	return nil
}