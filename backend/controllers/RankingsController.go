package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
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
	Times []string `json:"solves"`
}

func ProcessRankingsRows(rows pgx.Rows, single bool) ([]RankingsEntry, error) {
	return []RankingsEntry{}, nil
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
		if _type != "Single" && _type != "Average" {
			log.Println("ERR invalid type in GetRankings (" + _type + "): " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Invalid type (neither single nor average).")
			return
		}
		single := _type == "Single"

		region := c.Query("region")
		regionSplit := strings.Split(region, "|")
		if len(regionSplit) != 2 || (regionSplit[1] != "World" && regionSplit[1] != "Continent" && regionSplit[1] != "Country") {
			log.Println("ERR region in GetRankings (" + region + "): " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Invalid region/region format.")
			return
		}
		regionType := regionSplit[0]
		regionPrecise := regionSplit[1]

		var rankings []RankingsEntry

		if regionType == "World" {
			rows, err := db.Query(context.Background(), `SELECT u.name, c.iso2, c.name, r.competition_id, comp.name FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.user_id JOIN competitions comp ON comp.competition_id = r.competition_id WHERE r.event_id = $1;`, eid)
			if err != nil {
				log.Println("ERR db.Query (World) in GetRankings (" + region + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query rankings entries from database.")
				return
			}

			rankings, err = ProcessRankingsRows(rows, single)
			if err != nil {
				log.Println("ERR ProcessRankingsRows in GetRankings (" + region + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to process rankings entries from database.")
				return
			}
		} else {
			regionTypeColumn := "cont.name"
			if regionType == "Country" { regionTypeColumn = "c.name" }
			rows, err := db.Query(context.Background(), `SELECT u.name, c.iso2, c.name, r.competition_id, comp.name FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.user_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN continents cont ON cont.continent_id = c.continent_id WHERE r.event_id = $1 AND ` + regionTypeColumn + ` = $2;`, eid, regionPrecise)
			if err != nil {
				log.Println("ERR db.Query (" + regionType + ") in GetRankings (" + region + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query rankings entries from database.")
				return
			}

			rankings, err = ProcessRankingsRows(rows, single)
			if err != nil {
				log.Println("ERR ProcessRankingsRows (" + regionType + ") in GetRankings (" + region + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to process rankings entries from database.")
				return
			}
		}

		c.IndentedJSON(http.StatusOK, rankings)
	}
}