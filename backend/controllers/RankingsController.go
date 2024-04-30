package controllers

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
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
	Username string `json:"username"`
	CountryISO2 string `json:"country_iso2"`
	CountryName string `json:"country_name"`
	Result string `json:"result"`
	CompetitionId string `json:"competitionId"`
	CompetitionName string `json:"competitionName"`
	Times []string `json:"times"`
}

func MergeNonUniqueRankings(rankings []RankingsEntry) ([]RankingsEntry) {
	result := make([]RankingsEntry, 0)
	best := make(map[string]RankingsEntry)

	for _, rankingsEntry := range rankings {
		entry, ok := best[rankingsEntry.Username]
		if !ok || utils.ParseSolveToMilliseconds(entry.Result) > utils.ParseSolveToMilliseconds(rankingsEntry.Result) {
			best[rankingsEntry.Username] = rankingsEntry
		}
	}

	for _, v := range best { result = append(result, v) }

	sort.Slice(result, func (i int, j int) bool { return utils.ParseSolveToMilliseconds(result[i].Result) < utils.ParseSolveToMilliseconds(result[j].Result)})

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

		if regionType == "World" {
			rows, err := db.Query(context.Background(), `SELECT u.name, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN events e ON e.event_id = r.event_id WHERE r.event_id = $1;`, eid)
			if err != nil {
				log.Println("ERR db.Query (World) in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query rankings entries from database.")
				return
			}

			for rows.Next() {
				var rankingsEntry RankingsEntry
				var resultsEntry models.ResultEntry
				err := rows.Scan(&rankingsEntry.Username, &rankingsEntry.CountryISO2, &rankingsEntry.CountryName, &rankingsEntry.CompetitionId, &rankingsEntry.CompetitionName, &resultsEntry.Solve1, &resultsEntry.Solve2, &resultsEntry.Solve3, &resultsEntry.Solve4, &resultsEntry.Solve5, &resultsEntry.Format)
				if err != nil {
					log.Println("ERR scanning rows in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed to query rows from database.")
					return
				}

				if single {
					rankingsEntry.Result = resultsEntry.SingleFormatted()
					if utils.ParseSolveToMilliseconds(rankingsEntry.Result) >= constants.VERY_SLOW { continue; }
					rankingsEntry.Times = make([]string, 0)
				} else {
					resultFormatted, err := resultsEntry.AverageFormatted()
					if err != nil {
						log.Println("ERR AverageFormatted in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
						c.IndentedJSON(http.StatusInternalServerError, "Failed to calculate average in rankings entry.")
						return
					}
					rankingsEntry.Result = resultFormatted
					if utils.ParseSolveToMilliseconds(rankingsEntry.Result) >= constants.VERY_SLOW { continue; }
					rankingsEntry.Times, _ = resultsEntry.GetFormattedTimes()
				}
				rankings = append(rankings, rankingsEntry)
			}
		} else {
			regionTypeColumn := "cont.name"
			if regionType == "Country" { regionTypeColumn = "c.name" }
			rows, err := db.Query(context.Background(), `SELECT u.name, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN continents cont ON cont.continent_id = c.continent_id JOIN events e ON r.event_id = e.event_id WHERE r.event_id = $1 AND ` + regionTypeColumn + ` = $2;`, eid, regionPrecise)
			if err != nil {
				log.Println("ERR db.Query (" + regionType + ") in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query rankings entries from database.")
				return
			}

			for rows.Next() {
				var rankingsEntry RankingsEntry
				var resultsEntry models.ResultEntry
				err := rows.Scan(&rankingsEntry.Username, &rankingsEntry.CountryISO2, &rankingsEntry.CountryName, &rankingsEntry.CompetitionId, &rankingsEntry.CompetitionName, &resultsEntry.Solve1, &resultsEntry.Solve2, &resultsEntry.Solve3, &resultsEntry.Solve4, &resultsEntry.Solve5, &resultsEntry.Format)
				if err != nil {
					log.Println("ERR scanning rows in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed to query rows from database.")
					return
				}

				if single {
					rankingsEntry.Result = resultsEntry.SingleFormatted()
					if utils.ParseSolveToMilliseconds(rankingsEntry.Result) >= constants.VERY_SLOW { continue; }
					rankingsEntry.Times = make([]string, 0)
				} else {
					resultFormatted, err := resultsEntry.AverageFormatted()
					if err != nil {
						log.Println("ERR AverageFormatted in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
						c.IndentedJSON(http.StatusInternalServerError, "Failed to calculate average in rankings entry.")
						return
					}
					rankingsEntry.Result = resultFormatted
					if utils.ParseSolveToMilliseconds(rankingsEntry.Result) >= constants.VERY_SLOW { continue; }
					rankingsEntry.Times = resultsEntry.GetSolves()
				}
				rankings = append(rankings, rankingsEntry)
			}
		}

		rankings = MergeNonUniqueRankings(rankings)

		c.IndentedJSON(http.StatusOK, rankings)
	}
}