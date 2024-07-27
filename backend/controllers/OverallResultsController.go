package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type RegionSelectGroup struct {
	GroupName string `json:"groupName"`;
	GroupMembers []string `json:"groupMembers"`;
}

func GetRegionsGrouped(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		regionSelectGroups := make([]RegionSelectGroup, 0)
		regionSelectGroups = append(regionSelectGroups, RegionSelectGroup{"World", []string{"World"}})

		continents, err := utils.GetContinents(db)
		if err != nil {
			log.Println("ERR GetContinents in GetRegionsGrouped: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying continents from database.")
			return
		}
		regionSelectGroups = append(regionSelectGroups, RegionSelectGroup{"Continent", continents})

		countries, err := utils.GetCountries(db)
		if err != nil {
			log.Println("ERR GetCountries in GetRegionsGrouped: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying countries from database.")
			return
		}
		regionSelectGroups = append(regionSelectGroups, RegionSelectGroup{"Country", countries})

		c.IndentedJSON(http.StatusOK, regionSelectGroups)
	}
}

type RankingsEntry struct {
	Place string `json:"place"`
	Username string `json:"username"`
	WcaId string `json:"wca_id"`
	CountryISO2 string `json:"country_iso2"`
	CountryName string `json:"country_name"`
	Result string `json:"result"`
	CompetitionId string `json:"competitionId"`
	CompetitionName string `json:"competitionName"`
	Times []string `json:"times"`
}

type RecordsItem struct {
	EventName string `json:"eventname"`
	Iconcode string `json:"iconcode"`
	Entries []RecordsItemEntry `json:"entries"`
}

type RecordsItemEntry struct {
	Type string `json:"type"` // Single or Average
	Username string `json:"username"`
	WcaId string `json:"wcaId"`
	Result string `json:"result"`
	CountryIso2 string `json:"countryIso2"`
	CountryName string `json:"countryName"`
	CompetitionName string `json:"competitionName"`
	CompetitionId string `json:"competitionId"`
	CompetitionEndDate time.Time `json:"-"`
	Solves []string `json:"solves"`
}

func AddPlacementToRankings(rankings []RankingsEntry) {
	if len(rankings) == 0 { return }
	
	oldIdx := 0

	for idx := range rankings {
		if idx == 0 {
			rankings[0].Place = "1."
		} else {
			if utils.ParseSolveToMilliseconds(rankings[oldIdx].Result, false, "") != utils.ParseSolveToMilliseconds(rankings[idx].Result, false, "") {
				rankings[idx].Place = fmt.Sprintf("%d.", idx + 1)
				oldIdx = idx
			}
		}
	}
}

func MergeNonUniqueRankings(rankings []RankingsEntry, isfmc bool) ([]RankingsEntry) {
	result := make([]RankingsEntry, 0)
	best := make(map[string]RankingsEntry)

	for _, rankingsEntry := range rankings {
		entry, ok := best[rankingsEntry.Username]
		if !ok || utils.ParseSolveToMilliseconds(entry.Result, false, "") > utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") {
			best[rankingsEntry.Username] = rankingsEntry
		}
	}

	for _, v := range best { result = append(result, v) }

	sort.Slice(result, func (i int, j int) bool {
		val1, val2 := utils.ParseSolveToMilliseconds(result[i].Result, false, ""), utils.ParseSolveToMilliseconds(result[j].Result, false, "")
		if val1 == val2 { return result[i].Username < result[j].Username }
		return val1 < val2
	})

	return result
}

func GetRankings(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		eid, err := strconv.Atoi(c.Query("eid"))
		if err != nil {
			log.Println("ERR strconv(eid) in GetRankings: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing eventId.")
			return
		}

		_type := c.Query("type")
		if _type != "single" && _type != "average" {
			log.Println("ERR invalid type in GetRankings (" + _type + "): invalid type, should be single/average.")
			c.IndentedJSON(http.StatusInternalServerError, "Invalid type (neither single nor average).")
			return
		}
		single := _type == "single"

		regionType := c.Query("regionGroup")
		regionPrecise := c.Query("region")

		rankings := make([]RankingsEntry, 0)

		isfmc := false

		var rows pgx.Rows

		if regionType == "World" {
			rows, err = db.Query(context.Background(), `SELECT u.name, u.wcaid, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, rs.visible FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN events e ON e.event_id = r.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE r.event_id = $1 AND rs.visible IS TRUE;`, eid)
			if err != nil {
				log.Println("ERR db.Query (World) in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query rankings entries from database.")
				return
			}
		} else {
			regionTypeColumn := "cont.name"
			if regionType == "Country" { regionTypeColumn = "c.name" }
			rows, err = db.Query(context.Background(), `SELECT u.name, u.wcaid, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, rs.visible FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN continents cont ON cont.continent_id = c.continent_id JOIN events e ON r.event_id = e.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE r.event_id = $1 AND ` + regionTypeColumn + ` = $2 AND rs.visible IS TRUE;`, eid, regionPrecise)
			if err != nil {
				log.Println("ERR db.Query (" + regionType + ") in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query rankings entries from database.")
				return
			}
		}

		for rows.Next() {
			var rankingsEntry RankingsEntry
			var resultsEntry models.ResultEntry
			err := rows.Scan(&rankingsEntry.Username, &rankingsEntry.WcaId, &rankingsEntry.CountryISO2, &rankingsEntry.CountryName, &rankingsEntry.CompetitionId, &rankingsEntry.CompetitionName, &resultsEntry.Solve1, &resultsEntry.Solve2, &resultsEntry.Solve3, &resultsEntry.Solve4, &resultsEntry.Solve5, &resultsEntry.Format, &resultsEntry.Iconcode, &resultsEntry.Eventid, &resultsEntry.Status.Visible)
			if err != nil {
				log.Println("ERR scanning rows in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query rows from database.")
				return
			}

			if rankingsEntry.WcaId == "" { rankingsEntry.WcaId = rankingsEntry.Username }
			isfmc = utils.IsFMC(resultsEntry.Iconcode)
			scrambles, err := utils.GetScramblesByResultEntryId(db, resultsEntry.Eventid, rankingsEntry.CompetitionId)
			if err != nil {
				log.Println("ERR GetScramblesByResultEntryId in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to load scrambles.")
				return
			}

			if single {
				rankingsEntry.Result = resultsEntry.SingleFormatted(isfmc, scrambles)
				if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") >= constants.VERY_SLOW { continue; }
				rankingsEntry.Times = make([]string, 0)
			} else if resultsEntry.Iconcode != "333mbf" {
				resultFormatted, err := resultsEntry.AverageFormatted(isfmc, scrambles)
				if err != nil {
					log.Println("ERR AverageFormatted in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed to calculate average in rankings entry.")
					return
				}
				rankingsEntry.Result = resultFormatted
				if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") >= constants.VERY_SLOW { continue; }
				rankingsEntry.Times, _ = resultsEntry.GetFormattedTimes(isfmc, scrambles)
			}
			rankings = append(rankings, rankingsEntry)
		}

		rankings = MergeNonUniqueRankings(rankings, isfmc)
		AddPlacementToRankings(rankings)

		c.IndentedJSON(http.StatusOK, rankings)
	}
}

const (ALL_EVENT = -1)

func GetRecords(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		eid, err := strconv.Atoi(c.Query("eid"))
		if err != nil {
			log.Println("ERR strconv(eid) in GetRankings: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing eventId.")
			return
		}

		regionType := c.Query("regionGroup")
		regionPrecise := c.Query("region")

		recordItems := make([]RecordsItem, 0)

		isfmc := false

		var rows pgx.Rows
		
		if regionType == "World" {
			queryString := `SELECT u.name, u.wcaid, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, rs.visible, comp.enddate FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN events e ON e.event_id = r.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE rs.visible IS TRUE`;
			eidQueryPart := ` AND r.event_id = $1;`
			if eid != ALL_EVENT { queryString += eidQueryPart
			} else { queryString += `;`; }

			rows, err = db.Query(context.Background(), queryString, eid)
			if err != nil {
				log.Println("ERR db.Query (World) in GetRecords (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query records entries from database.")
				return
			}
		} else {
			regionTypeColumn := "cont.name"
			if regionType == "Country" { regionTypeColumn = "c.name" }

			queryString := `SELECT u.name, u.wcaid, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, rs.visible, comp.enddate FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN continents cont ON cont.continent_id = c.continent_id JOIN events e ON r.event_id = e.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE ` + regionTypeColumn + ` = $1 AND rs.visible IS TRUE `
			eidQueryPart := ` AND r.event_id = $2;`
			if eid != ALL_EVENT { queryString += eidQueryPart
			} else { queryString += `;`; }

			rows, err = db.Query(context.Background(), queryString, regionPrecise, eid)
			if err != nil {
				log.Println("ERR db.Query (" + regionType + ") in GetRecords (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query records entries from database.")
				return
			}
		}

		singleEntries := make(map[int][]RecordsItemEntry)
		averageEntries := make(map[int][]RecordsItemEntry)

		for rows.Next() {
			var rankingsEntry RankingsEntry
			var resultsEntry models.ResultEntry
			var competitionEndDate time.Time
			err := rows.Scan(&rankingsEntry.Username, &rankingsEntry.WcaId, &rankingsEntry.CountryISO2, &rankingsEntry.CountryName, &rankingsEntry.CompetitionId, &rankingsEntry.CompetitionName, &resultsEntry.Solve1, &resultsEntry.Solve2, &resultsEntry.Solve3, &resultsEntry.Solve4, &resultsEntry.Solve5, &resultsEntry.Format, &resultsEntry.Iconcode, &resultsEntry.Eventid, &resultsEntry.Status.Visible, &competitionEndDate)
			if err != nil {
				log.Println("ERR scanning rows in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query rows from database.")
				return
			}

			if rankingsEntry.WcaId == "" { rankingsEntry.WcaId = rankingsEntry.Username }
			isfmc = utils.IsFMC(resultsEntry.Iconcode)
			scrambles, err := utils.GetScramblesByResultEntryId(db, resultsEntry.Eventid, rankingsEntry.CompetitionId)
			if err != nil {
				log.Println("ERR GetScramblesByResultEntryId in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to load scrambles.")
				return
			}

			recordsItemEntrySingle := RecordsItemEntry{
				Type: "Single",
				Username: rankingsEntry.Username,
				WcaId: rankingsEntry.WcaId,
				CountryIso2: rankingsEntry.CountryISO2,
				CountryName: rankingsEntry.CountryName,
				CompetitionName: rankingsEntry.CompetitionName,
				CompetitionId: rankingsEntry.CompetitionId,
				CompetitionEndDate: competitionEndDate,
				Solves: []string{},
			}
			recordsItemEntryAverage := recordsItemEntrySingle
			recordsItemEntryAverage.Type = "Average"

			recordsItemEntrySingle.Result = resultsEntry.SingleFormatted(isfmc, scrambles)
			recordsItemEntryAverage.Result = "DNS"
			if resultsEntry.Iconcode != "333mbf" {
				resultFormatted, err := resultsEntry.AverageFormatted(isfmc, scrambles)
				if err != nil {
					log.Println("ERR AverageFormatted in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed to calculate average in rankings entry.")
					return
				}
				recordsItemEntryAverage.Result  = resultFormatted
				if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") >= constants.VERY_SLOW { continue; }
				recordsItemEntryAverage.Solves, _ = resultsEntry.GetFormattedTimes(isfmc, scrambles)
			}
			
			
		}

		c.IndentedJSON(http.StatusOK, recordItems)
	}
}

