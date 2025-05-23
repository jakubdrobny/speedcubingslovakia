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

func GetResultsQuery(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("uname")
		competitionName := c.Param("cname")
		resultsStatusDisplayName := c.Param("rsname")
		eventId, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			log.Println("ERR in strconv(eventId) in GetResultsQuery: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse eventId.")
			return
		}

		resultEntries := make([]models.ResultEntry, 0)

		if competitionName == "_" && userName == "_" {
			rows, err := db.Query(
				context.Background(),
				`SELECT re.result_id FROM results re JOIN results_status rs ON rs.results_status_id = re.status_id WHERE re.event_id = $1 AND rs.displayname LIKE $2;`,
				eventId,
				resultsStatusDisplayName,
			)
			if err != nil {
				log.Println(
					"ERR db.Query in GetResultsQuery (competitionName not set and userName not set): " + err.Error(),
				)
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed querying result entry from database.",
				)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil {
					log.Println(
						"ERR scanning resultEntryId in GetResultsQuery (competitionId not set and userName not set): " + err.Error(),
					)
					c.IndentedJSON(
						http.StatusInternalServerError,
						"Failed querying result entry from database.",
					)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					log.Println(
						"ERR GetResultEntryById in GetResultsQuery (competitionId not set and userName not set): " + err.Error(),
					)
					c.IndentedJSON(
						http.StatusInternalServerError,
						"Failed getting result entry from database.",
					)
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else if competitionName == "_" && userName != "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN users u ON u.user_id = re.user_id JOIN results_status rs ON rs.results_status_id = re.status_id WHERE re.event_id = $1 AND UPPER(u.name) LIKE UPPER('%' || $2 || '%') AND rs.displayname LIKE $3;`, eventId, userName, resultsStatusDisplayName)
			if err != nil {
				log.Println("ERR db.Query in GetResultsQuery (competitionName not set and userName set): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil {
					log.Println("ERR scanning resultEntryId in GetResultsQuery (competitionName not set and userName set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					log.Println("ERR GetResultEntryById in GetResultsQuery (competitionName not set and userName set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else if competitionName != "_" && userName == "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN competitions c ON c.competition_id = re.competition_id JOIN results_status rs ON rs.results_status_id = re.status_id WHERE re.event_id = $1 AND UPPER(c.name) LIKE UPPER('%' || $2 || '%') AND rs.displayname LIKE $3;`, eventId, competitionName, resultsStatusDisplayName)
			if err != nil {
				log.Println("ERR db.Query in GetResultsQuery (competitionName set and userName not set): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil {
					log.Println("ERR scanning resultEntryId in GetResultsQuery (competitionName set and userName not set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					log.Println("ERR GetResultEntryById in GetResultsQuery (competitionName set and userName not set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN users u ON u.user_id = re.user_id JOIN competitions c ON c.competition_id = re.competition_id JOIN results_status rs ON rs.results_status_id = re.status_id WHERE re.event_id = $1 AND UPPER(c.name) LIKE UPPER('%' || $2 || '%') AND UPPER(u.name) LIKE UPPER('%' || $3 || '%') AND rs.displayname = $4;`, eventId, competitionName, userName, resultsStatusDisplayName)
			if err != nil {
				log.Println("ERR db.Query in GetResultsQuery (competitionName set and userName set): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil {
					log.Println("ERR scanning resultEntryId in GetResultsQuery (competitionName set and userName set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					log.Println("ERR GetResultEntryById in GetResultsQuery (competitionName set and userName set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		}

		c.IndentedJSON(http.StatusOK, resultEntries)
	}
}

type ValidateResultsBody struct {
	ResultId int  `json:"resultId"`
	Verdict  bool `json:"verdict"`
}

func ValidateResults(db *pgxpool.Pool, body ValidateResultsBody, isadmin bool) (string, string) {
	resultEntry, err := models.GetResultEntryById(db, body.ResultId)
	if err != nil {
		return "ERR GetResultEntryById in PostResultsValidation: " + err.Error(), "Failed getting result entry from database."
	}

	statusId := 3
	if !body.Verdict {
		statusId = 2
	}
	resultStatus, err := models.GetResultsStatus(db, statusId)
	if err != nil {
		return "ERR GetResultsStatus in PostResultsValidation: " + err.Error(), "Failed getting result status in database."
	}

	resultEntry.Status = resultStatus
	err = resultEntry.Update(db, isadmin, resultEntry.IsFMC(), true)
	if err != nil {
		return "ERR resultEntry.Update in PostResultsValidation: " + err.Error(), "Failed updating result entry in database."
	}

	return "", ""
}

func GetResultsValidation(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		resultId, err := strconv.Atoi(c.DefaultQuery("resultId", "0"))
		if err != nil {
			log.Println("ERR strconv.Atoi in GetResutsValidation: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse result id.")
			return
		}

		verdict, err := strconv.ParseBool(c.DefaultQuery("verdict", "false"))
		if err != nil {
			log.Println("ERR strconv.ParseBool in GetResutsValidation: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse verdict.")
			return
		}

		body := ValidateResultsBody{ResultId: resultId, Verdict: verdict}

		isadmin := c.MustGet("isadmin").(bool)
		logMsg, retMsg := ValidateResults(db, body, isadmin)
		if logMsg != "" || retMsg != "" {
			log.Println(logMsg)
			c.IndentedJSON(http.StatusInternalServerError, retMsg)
			return
		}

		retMsg = "Result APPROVED."
		if !verdict {
			retMsg = "Result DENIED."
		}

		c.IndentedJSON(http.StatusCreated, retMsg)
	}
}

func PostResultsValidation(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body ValidateResultsBody

		if err := c.BindJSON(&body); err != nil {
			log.Println("ERR BindJSON in PostResultsValidation: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing data.")
			return
		}

		isadmin := c.MustGet("isadmin").(bool)
		logMsg, retMsg := ValidateResults(db, body, isadmin)
		if logMsg != "" || retMsg != "" {
			log.Println(logMsg)
			c.IndentedJSON(http.StatusInternalServerError, retMsg)
			return
		}

		c.IndentedJSON(http.StatusCreated, "")
	}
}

func GetResultsByIdAndEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventId, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			log.Println("ERR strconv.eid in GetResultsByIdAndEvent: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing eventId.")
			return
		}

		competitionId := c.Param("cid")
		userId := c.MustGet("uid").(int)

		user, err := models.GetUserById(db, userId)
		if err != nil {
			log.Println("ERR GetUserById in GetResultsByIdAndEvent: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed getting user information from database.",
			)
			return
		}

		event, err := models.GetEventById(db, eventId)
		if err != nil {
			log.Println("ERR GetEventById in GetResultsByIdAndEvent: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed getting event information from database.",
			)
			return
		}

		competition, err := models.GetCompetitionByIdObject(db, competitionId)
		if err != nil {
			log.Println("ERR GetCompetitionByIdObject in GetResultsByIdAndEvent: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed getting competition information from database.",
			)
			return
		}

		resultEntry, err := models.GetResultEntry(db, userId, competitionId, eventId)

		if err != nil {
			if err.Error() != "not found" {
				log.Println("ERR GetResultEntry in GetResultsByIdAndEvent: " + err.Error())
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed getting result entry from database.",
				)
				return
			} else {
				approvedResultsStatus, err := models.GetResultsStatus(db, 3)
				if err != nil {
					log.Println("ERR GetResultsStatus.approved in GetResultsByIdAndEvent: " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result status in database.")
					return
				}

				resultEntry = models.ResultEntry{
					Userid:          userId,
					Username:        user.Name,
					WcaId:           user.WcaId,
					Competitionid:   competitionId,
					Competitionname: competition.Name,
					Eventid:         event.Id,
					Eventname:       event.Displayname,
					Iconcode:        event.Iconcode,
					Format:          event.Format,
					Solve1:          "DNS",
					Solve2:          "DNS",
					Solve3:          "DNS",
					Solve4:          "DNS",
					Solve5:          "DNS",
					Comment:         "",
					Status:          approvedResultsStatus,
				}

				err = resultEntry.Insert(db)
				if err != nil {
					log.Println("ERR resultEntry.Insert in GetResultsByIdAndEvent: " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed inserting results into database.")
					return
				}
			}
		} else {
			currentStatus, err := models.GetResultsStatus(db, resultEntry.Status.Id)
			if err != nil {
				log.Println("ERR GetResultsStatus.resultEntry.Status.Id in GetResultsByIdAndEvent: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed getting result status in database.")
				return
			}

			resultEntry.Status = currentStatus
			resultEntry.Eventname = event.Displayname
			resultEntry.Competitionname = competition.Name
			resultEntry.Username = user.Name
			resultEntry.WcaId = user.WcaId
			resultEntry.Iconcode = event.Iconcode
			resultEntry.Format = event.Format

			err = resultEntry.Update(db, false, resultEntry.IsFMC())
			if err != nil {
				log.Println("ERR resultEntry.Update in GetResultsByIdAndEvent: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed updating results in database.")
				return
			}
		}

		c.IndentedJSON(http.StatusOK, resultEntry)
	}
}

func PostResults(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resultEntry models.ResultEntry
		var err error

		if err = c.BindJSON(&resultEntry); err != nil {
			log.Println("ERR BindJSON in PostResults: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing data.")
			return
		}

		isadmin := c.MustGet("isadmin").(bool)
		uid := c.MustGet("uid").(int)
		if !isadmin && uid != resultEntry.Userid {
			c.IndentedJSON(http.StatusCreated, "Nope")
			return
		}

		if resultEntry.Id == 0 {
			err = resultEntry.LoadId(db)
			if err != nil {
				log.Println("ERR resultEntry.LoadId in PostResults: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry id.")
				return
			}
		}

		previousTimes, err := resultEntry.GetPreviouslySavedTimes(db)
		if err != nil {
			log.Println("ERR resultEntry.GetPreviouslySavedTimes in PostResults: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed querying result entry in database.",
			)
			return
		}

		err = resultEntry.Update(db, isadmin, resultEntry.IsFMC())
		if err != nil {
			log.Println("ERR resultEntry.Update in PostResults: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed updating results in database.")
			return
		}

		go resultEntry.SendSuspicousMail(c, db, envMap, previousTimes)

		c.IndentedJSON(http.StatusCreated, resultEntry)
	}
}

func GetProfileResults(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		uid, err := models.GetUserByWCAID(db, id)
		if err != nil {
			log.Println("ERR in GetProfileResults in GetUserByWCAID: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Finding user by WCA ID in database failed.",
			)
			return
		}

		if uid == 0 {
			uid, err = models.GetUserByName(db, id)
			if err != nil {
				log.Println("ERR in GetProfileResults in GetUserByName: " + err.Error())
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Finding user by name in database failed.",
				)
				return
			}
		}

		var profileResults models.ProfileType
		err = profileResults.Load(db, uid)
		if err != nil {
			log.Println("ERR in GetProfileResults in ProfileType.Load: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Retrieving profile results failed.")
			return
		}

		c.IndentedJSON(http.StatusOK, profileResults)
	}
}

type RegionSelectGroup struct {
	GroupName    string   `json:"groupName"`
	GroupMembers []string `json:"groupMembers"`
}

func GetRegionsGrouped(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		regionSelectGroups := make([]RegionSelectGroup, 0)
		regionSelectGroups = append(
			regionSelectGroups,
			RegionSelectGroup{"World", []string{"World"}},
		)

		continents, err := utils.GetContinents(db)
		if err != nil {
			log.Println("ERR GetContinents in GetRegionsGrouped: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed querying continents from database.",
			)
			return
		}
		regionSelectGroups = append(regionSelectGroups, RegionSelectGroup{"Continent", continents})

		countries, err := models.GetCountries(db)
		if err != nil {
			log.Println("ERR GetCountries in GetRegionsGrouped: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed querying countries from database.",
			)
			return
		}
		regionSelectGroups = append(
			regionSelectGroups,
			RegionSelectGroup{
				"Country",
				utils.Map(countries, func(c models.Country) string { return c.Name }),
			},
		)

		c.IndentedJSON(http.StatusOK, regionSelectGroups)
	}
}

type RankingsEntry struct {
	Place           string   `json:"place"`
	Username        string   `json:"username"`
	WcaId           string   `json:"wca_id"`
	CountryISO2     string   `json:"country_iso2"`
	CountryName     string   `json:"country_name"`
	Result          string   `json:"result"`
	CompetitionId   string   `json:"competitionId"`
	CompetitionName string   `json:"competitionName"`
	Times           []string `json:"times"`
}

type RecordsItem struct {
	EventName string             `json:"eventname"`
	Iconcode  string             `json:"iconcode"`
	Entries   []RecordsItemEntry `json:"entries"`
}

type RecordsItemEntry struct {
	Type               string    `json:"type"` // Single or Average
	Username           string    `json:"username"`
	WcaId              string    `json:"wcaId"`
	Result             string    `json:"result"`
	CountryIso2        string    `json:"countryIso2"`
	CountryName        string    `json:"countryName"`
	CompetitionName    string    `json:"competitionName"`
	CompetitionId      string    `json:"competitionId"`
	Solves             []string  `json:"solves"`
	CompetitionEndDate time.Time `json:"-"`
	EventName          string    `json:"-"`
	IconCode           string    `json:"-"`
}

func AddPlacementToRankings(rankings []RankingsEntry) {
	if len(rankings) == 0 {
		return
	}

	oldIdx := 0

	for idx := range rankings {
		if idx == 0 {
			rankings[0].Place = "1."
		} else {
			if utils.ParseSolveToMilliseconds(rankings[oldIdx].Result, false, "") != utils.ParseSolveToMilliseconds(rankings[idx].Result, false, "") {
				rankings[idx].Place = fmt.Sprintf("%d.", idx+1)
				oldIdx = idx
			}
		}
	}
}

func MergeNonUniqueRankings(rankings []RankingsEntry, isfmc bool) []RankingsEntry {
	result := make([]RankingsEntry, 0)
	best := make(map[string]RankingsEntry)

	for _, rankingsEntry := range rankings {
		entry, ok := best[rankingsEntry.Username]
		if !ok ||
			utils.ParseSolveToMilliseconds(
				entry.Result,
				false,
				"",
			) > utils.ParseSolveToMilliseconds(
				rankingsEntry.Result,
				false,
				"",
			) {
			best[rankingsEntry.Username] = rankingsEntry
		}
	}

	for _, v := range best {
		result = append(result, v)
	}

	return result
}

func GetRankings(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		eid, err := strconv.Atoi(c.Query("eid"))
		if err != nil {
			log.Println("ERR strconv(eid) in GetRankings: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing eventId.")
			return
		}

		_type := c.Query("type")
		if _type != "single" && _type != "average" {
			log.Println(
				"ERR invalid type in GetRankings (" + _type + "): invalid type, should be single/average.",
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Invalid type (neither single nor average).",
			)
			return
		}
		single := _type == "single"

		regionType := c.Query("regionGroup")
		regionPrecise := c.Query("region")
		queryType := c.Query("queryType")
		numOfEntries, err := strconv.Atoi(c.Query("numOfEntries"))
		if err != nil {
			log.Println("ERR strconv(numOfEntries) in GetRankings: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing numOfEntries.")
			return
		}

		persons := queryType == "Persons"
		if numOfEntries != 100 && numOfEntries != 1000 {
			log.Println(
				"ERR invalid no. of entries (" + string(
					numOfEntries,
				) + ") in query in GetRankings. Possible values: 100, 1000",
			)
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing eventId.")
			return
		}

		if !persons && queryType != "Results" {
			log.Println(
				"ERR invalid query type (" + string(
					queryType,
				) + ") in query in GetRankings. Possible values: Persons, Results",
			)
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing eventId.")
			return
		}

		rankings := make([]RankingsEntry, 0)

		if eid == -1 {
			competitionResults, err := models.GetOverallResults(db, "", regionType, regionPrecise)
			if err != nil {
				log.Println("ERR models.GetOverallResults in GetRankings: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed getting overall rankings.")
				return
			}

			for _, result := range competitionResults {
				var rankingsEntry RankingsEntry
				rankingsEntry.Place = result.Place
				rankingsEntry.Username = result.Username
				rankingsEntry.WcaId = result.WcaId
				rankingsEntry.CountryName = result.CountryName
				rankingsEntry.CountryISO2 = result.CountryIso2
				rankingsEntry.Result = result.Score

				rankings = append(rankings, rankingsEntry)
			}
		} else {
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
				if regionType == "Country" {
					regionTypeColumn = "c.name"
				}
				rows, err = db.Query(context.Background(), `SELECT u.name, u.wcaid, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, rs.visible FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN continents cont ON cont.continent_id = c.continent_id JOIN events e ON r.event_id = e.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE r.event_id = $1 AND `+regionTypeColumn+` = $2 AND rs.visible IS TRUE;`, eid, regionPrecise)
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

				if rankingsEntry.WcaId == "" {
					rankingsEntry.WcaId = rankingsEntry.Username
				}
				isfmc = utils.IsFMC(resultsEntry.Iconcode)
				scrambles, err := utils.GetScramblesByResultEntryId(db, resultsEntry.Eventid, rankingsEntry.CompetitionId)
				if err != nil {
					log.Println("ERR GetScramblesByResultEntryId in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed to load scrambles.")
					return
				}

				ismbld := resultsEntry.Iconcode == "333mbf"

				if single {
					if persons {
						rankingsEntry.Result = resultsEntry.SingleFormatted(isfmc, scrambles)
						if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") >= constants.VERY_SLOW {
							continue
						}
						rankingsEntry.Times = make([]string, 0)

						rankings = append(rankings, rankingsEntry)
					} else {
						rankingsEntry.Times = make([]string, 0)

						noOfSolves, _ := utils.GetNoOfSolves(resultsEntry.Format)

						for idx, solve := range []string{resultsEntry.Solve1, resultsEntry.Solve2, resultsEntry.Solve3, resultsEntry.Solve4, resultsEntry.Solve5} {
							if idx >= noOfSolves {
								break
							}

							result := utils.ParseSolveToMilliseconds(solve, isfmc, scrambles[idx])
							if ismbld {
								rankingsEntry.Result = solve
							} else {
								rankingsEntry.Result = utils.FormatTime(result, isfmc)
							}
							if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") < constants.VERY_SLOW {
								rankings = append(rankings, rankingsEntry)
							}
						}
					}
				} else if !ismbld && resultsEntry.Format != "bo1" {
					resultFormatted, err := resultsEntry.AverageFormatted(isfmc, scrambles)
					if err != nil {
						log.Println("ERR AverageFormatted in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error())
						c.IndentedJSON(http.StatusInternalServerError, "Failed to calculate average in rankings entry.")
						return
					}
					rankingsEntry.Result = resultFormatted
					if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") >= constants.VERY_SLOW {
						continue
					}
					rankingsEntry.Times, _ = resultsEntry.GetFormattedTimes(isfmc, scrambles)

					rankings = append(rankings, rankingsEntry)
				}
			}

			if persons {
				rankings = MergeNonUniqueRankings(rankings, isfmc)
			}
			sort.Slice(rankings, func(i int, j int) bool {
				val1, val2 := utils.ParseSolveToMilliseconds(rankings[i].Result, false, ""), utils.ParseSolveToMilliseconds(rankings[j].Result, false, "")
				if val1 == val2 {
					return rankings[i].Username < rankings[j].Username
				}
				return val1 < val2
			})
			AddPlacementToRankings(rankings)
		}

		if len(rankings) > numOfEntries {
			rankings = rankings[:numOfEntries]
		}

		c.IndentedJSON(http.StatusOK, rankings)
	}
}

const (
	ALL_EVENT = -1
)

func GetRecords(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			queryString := `SELECT u.name, u.wcaid, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, rs.visible, comp.enddate, e.fulldisplayname FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN events e ON e.event_id = r.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE rs.visible IS TRUE`
			eidQueryPart := ` AND r.event_id = $1;`
			if eid != ALL_EVENT {
				rows, err = db.Query(context.Background(), queryString+eidQueryPart, eid)
			} else {
				rows, err = db.Query(context.Background(), queryString+`;`)
			}
			if err != nil {
				log.Println(
					"ERR db.Query (World) in GetRecords (" + regionType + "+" + regionPrecise + "): " + err.Error(),
				)
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed to query records entries from database.",
				)
				return
			}
		} else {
			regionTypeColumn := "cont.name"
			if regionType == "Country" {
				regionTypeColumn = "c.name"
			}

			queryString := `SELECT u.name, u.wcaid, c.iso2, c.name, r.competition_id, comp.name, r.solve1, r.solve2, r.solve3, r.solve4, r.solve5, e.format, e.iconcode, r.event_id, rs.visible, comp.enddate, e.fulldisplayname FROM results r JOIN users u ON u.user_id = r.user_id JOIN countries c ON c.country_id = u.country_id JOIN competitions comp ON comp.competition_id = r.competition_id JOIN continents cont ON cont.continent_id = c.continent_id JOIN events e ON r.event_id = e.event_id JOIN results_status rs ON rs.results_status_id = r.status_id WHERE ` + regionTypeColumn + ` = $1 AND rs.visible IS TRUE `
			eidQueryPart := ` AND r.event_id = $2;`
			if eid != ALL_EVENT {
				rows, err = db.Query(context.Background(), queryString+eidQueryPart, regionPrecise, eid)
			} else {
				rows, err = db.Query(context.Background(), queryString+`;`, regionPrecise)
			}
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
			err := rows.Scan(
				&rankingsEntry.Username,
				&rankingsEntry.WcaId,
				&rankingsEntry.CountryISO2,
				&rankingsEntry.CountryName,
				&rankingsEntry.CompetitionId,
				&rankingsEntry.CompetitionName,
				&resultsEntry.Solve1,
				&resultsEntry.Solve2,
				&resultsEntry.Solve3,
				&resultsEntry.Solve4,
				&resultsEntry.Solve5,
				&resultsEntry.Format,
				&resultsEntry.Iconcode,
				&resultsEntry.Eventid,
				&resultsEntry.Status.Visible,
				&competitionEndDate,
				&resultsEntry.Eventname,
			)
			if err != nil {
				log.Println(
					"ERR scanning rows in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error(),
				)
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed to query rows from database.",
				)
				return
			}

			if rankingsEntry.WcaId == "" {
				rankingsEntry.WcaId = rankingsEntry.Username
			}
			isfmc = utils.IsFMC(resultsEntry.Iconcode)
			scrambles, err := utils.GetScramblesByResultEntryId(
				db,
				resultsEntry.Eventid,
				rankingsEntry.CompetitionId,
			)
			if err != nil {
				log.Println(
					"ERR GetScramblesByResultEntryId in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error(),
				)
				c.IndentedJSON(http.StatusInternalServerError, "Failed to load scrambles.")
				return
			}

			recordsItemEntrySingle := RecordsItemEntry{
				Type:               "Single",
				Username:           rankingsEntry.Username,
				WcaId:              rankingsEntry.WcaId,
				CountryIso2:        rankingsEntry.CountryISO2,
				CountryName:        rankingsEntry.CountryName,
				CompetitionName:    rankingsEntry.CompetitionName,
				CompetitionId:      rankingsEntry.CompetitionId,
				Solves:             []string{},
				CompetitionEndDate: competitionEndDate,
				EventName:          resultsEntry.Eventname,
				IconCode:           resultsEntry.Iconcode,
			}

			recordsItemEntryAverage := recordsItemEntrySingle
			recordsItemEntryAverage.Type = "Average"
			if val, err := utils.GetNoOfSolves(resultsEntry.Format); err == nil && val < 5 {
				recordsItemEntryAverage.Type = "Mean"
			}

			recordsItemEntrySingle.Result = resultsEntry.SingleFormatted(isfmc, scrambles)
			if utils.ParseSolveToMilliseconds(
				recordsItemEntrySingle.Result,
				false,
				"",
			) >= constants.VERY_SLOW {
				continue
			}
			singleEntry, ok := singleEntries[resultsEntry.Eventid]
			if !ok ||
				utils.ParseSolveToMilliseconds(
					singleEntry[0].Result,
					false,
					"",
				) >= utils.ParseSolveToMilliseconds(
					recordsItemEntrySingle.Result,
					false,
					"",
				) {
				if !ok ||
					utils.ParseSolveToMilliseconds(
						singleEntry[0].Result,
						false,
						"",
					) > utils.ParseSolveToMilliseconds(
						recordsItemEntrySingle.Result,
						false,
						"",
					) {
					singleEntries[resultsEntry.Eventid] = []RecordsItemEntry{}
				}
				singleEntries[resultsEntry.Eventid] = append(
					singleEntries[resultsEntry.Eventid],
					recordsItemEntrySingle,
				)
			}

			recordsItemEntryAverage.Result = "DNS"

			if resultsEntry.Iconcode != "333mbf" && resultsEntry.Format != "bo1" {
				resultFormatted, err := resultsEntry.AverageFormatted(isfmc, scrambles)
				if err != nil {
					log.Println(
						"ERR AverageFormatted in GetRankings (" + regionType + "+" + regionPrecise + "): " + err.Error(),
					)
					c.IndentedJSON(
						http.StatusInternalServerError,
						"Failed to calculate average in rankings entry.",
					)
					return
				}
				recordsItemEntryAverage.Result = resultFormatted
				if utils.ParseSolveToMilliseconds(
					recordsItemEntryAverage.Result,
					false,
					"",
				) >= constants.VERY_SLOW {
					continue
				}
				recordsItemEntryAverage.Solves, _ = resultsEntry.GetFormattedTimes(isfmc, scrambles)

				averageEntry, ok := averageEntries[resultsEntry.Eventid]
				if !ok ||
					utils.ParseSolveToMilliseconds(
						averageEntry[0].Result,
						false,
						"",
					) >= utils.ParseSolveToMilliseconds(
						recordsItemEntryAverage.Result,
						false,
						"",
					) {
					if !ok ||
						utils.ParseSolveToMilliseconds(
							averageEntry[0].Result,
							false,
							"",
						) > utils.ParseSolveToMilliseconds(
							recordsItemEntryAverage.Result,
							false,
							"",
						) {
						averageEntries[resultsEntry.Eventid] = []RecordsItemEntry{}
					}
					averageEntries[resultsEntry.Eventid] = append(
						averageEntries[resultsEntry.Eventid],
						recordsItemEntryAverage,
					)
				}
			}
		}

		eventIDs := make([]int, len(singleEntries))
		idx := 0
		for key := range singleEntries {
			eventIDs[idx] = key
			idx++
		}
		sort.Slice(eventIDs, func(i int, j int) bool { return eventIDs[i] < eventIDs[j] })

		for _, eventID := range eventIDs {
			if len(singleEntries[eventID]) == 0 {
				continue
			}

			recordItem := RecordsItem{
				EventName: singleEntries[eventID][0].EventName,
				Iconcode:  singleEntries[eventID][0].IconCode,
				Entries:   []RecordsItemEntry{},
			}

			sort.Slice(singleEntries[eventID], func(i int, j int) bool {
				item1 := singleEntries[eventID][i]
				item2 := singleEntries[eventID][j]
				if item1.CompetitionEndDate.Equal(item1.CompetitionEndDate) {
					return item1.Username < item2.Username
				}
				return item1.CompetitionEndDate.Before(item2.CompetitionEndDate)
			})
			recordItem.Entries = append(recordItem.Entries, singleEntries[eventID]...)
			sort.Slice(averageEntries[eventID], func(i int, j int) bool {
				item1 := averageEntries[eventID][i]
				item2 := averageEntries[eventID][j]
				if item1.CompetitionEndDate.Equal(item1.CompetitionEndDate) {
					return item1.Username < item2.Username
				}
				return item1.CompetitionEndDate.Before(item2.CompetitionEndDate)
			})
			recordItem.Entries = append(recordItem.Entries, averageEntries[eventID]...)

			recordItems = append(recordItems, recordItem)
		}

		c.IndentedJSON(http.StatusOK, recordItems)
	}
}

type AverageInfo struct {
	Single              string   `json:"single"`
	Average             string   `json:"average"`
	Times               []string `json:"times"`
	Bpa                 string   `json:"bpa"`
	Wpa                 string   `json:"wpa"`
	ShowPossibleAverage bool     `json:"showPossibleAverage"`
	FinishedCompeting   bool     `json:"finishedCompeting"`
	Place               string   `json:"place"`
	SingleRecord        string   `json:"singleRecord"`
	SingleRecordColor   string   `json:"singleRecordColor"`
	AverageRecord       string   `json:"averageRecord"`
	AverageRecordColor  string   `json:"averageRecordColor"`
}

func GetAverageInfo(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resultEntry models.ResultEntry
		var err error

		if err := c.BindJSON(&resultEntry); err != nil {
			log.Println("ERR BindJSON(&resultEntry) in GetAverageInfo: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse result entry.")
			return
		}

		var averageInfo AverageInfo

		resultEntry.Scrambles, err = utils.GetScramblesByResultEntryId(
			db,
			resultEntry.Eventid,
			resultEntry.Competitionid,
		)
		if err != nil {
			log.Println("ERR utils.GetScramblesByResultEntryId in GetAverageInfo: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to get scrambles for result entry.",
			)
			return
		}
		averageInfo.Single = resultEntry.SingleFormatted(resultEntry.IsFMC(), resultEntry.Scrambles)

		avg, err := resultEntry.AverageFormatted(resultEntry.IsFMC(), resultEntry.Scrambles)
		if err != nil {
			log.Println("ERR resultEntry.AverageFormatted in GetAverageInfo: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to get average for result entry.",
			)
			return
		}
		averageInfo.Average = avg

		formattedTimes, err := resultEntry.GetFormattedTimes(
			resultEntry.IsFMC(),
			resultEntry.Scrambles,
		)
		if err != nil {
			log.Println("ERR resultEntry.GetFormattedTimes in GetAverageInfo: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get times for result entry.")
			return
		}
		averageInfo.Times = formattedTimes

		ok, err := resultEntry.ShowPossibleAverages()
		if err != nil {
			log.Println("ERR resultEntry.ShowPossibleAverages in GetAverageInfo: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to check if should calculate BPA/WPA.",
			)
			return
		}

		if ok {
			averageInfo.ShowPossibleAverage = true

			averageInfo.Bpa, err = resultEntry.GetBPA()
			if err != nil {
				log.Println("ERR resultEntry.GetBPA in GetAverageInfo: " + err.Error())
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed to get BPA for result entry.",
				)
				return
			}

			averageInfo.Wpa, err = resultEntry.GetWPA()
			if err != nil {
				log.Println("ERR resultEntry.GetWPA in GetAverageInfo: " + err.Error())
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed to get WPA for result entry.",
				)
				return
			}
		}

		averageInfo.FinishedCompeting, err = resultEntry.FinishedCompeting()
		if err != nil {
			log.Println("ERR resultEntry.FinishedCompeting in GetAverageInfo: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to check if you finished competing.",
			)
			return
		}

		if averageInfo.FinishedCompeting {
			averageInfo.Place, err = resultEntry.GetCompetitionPlace(db)
			if err != nil {
				log.Println("ERR resultEntry.GetCompetitionPlace in GetAverageInfo: " + err.Error())
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed to get competition place for result entry.",
				)
				return
			}
		}

		c.IndentedJSON(http.StatusOK, averageInfo)
	}
}

func GetAverageInfoRecords(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Body struct {
			ResultEntry models.ResultEntry `json:"resultEntry"`
			AverageInfo AverageInfo        `json:"averageInfo"`
		}
		var body Body
		if err := c.BindJSON(&body); err != nil {
			log.Println("ERR BindJSON(&body) in GetAverageInfo: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse request body.")
			return
		}

		resultEntry := body.ResultEntry
		averageInfo := body.AverageInfo

		var err error

		averageInfo.FinishedCompeting, err = resultEntry.FinishedCompeting()
		if err != nil {
			log.Println(
				"ERR resultEntry.FinishedCompeting in GetAverageInfoRecords: " + err.Error(),
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to check if you finished competing.",
			)
			return
		}

		if averageInfo.FinishedCompeting {
			uid := c.MustGet("uid").(int)

			user, err := models.GetUserById(db, uid)
			if err != nil {
				log.Println("ERR models.GetUserById in GetAverageInfo: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to get user info.")
				return
			}
			err = user.LoadContinent(db)
			if err != nil {
				log.Println("ERR user.LoadContinent in GetAverageInfo: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to get user continent.")
				return
			}

			var profileType models.ProfileType
			_, err = profileType.LoadPersonalBests(db, &user, resultEntry.Eventid)
			if err != nil {
				log.Println("ERR profileType.LoadPersonalBests in GetAverageInfo: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to get user continent.")
				return
			}

			for idx := range profileType.PersonalBests {
				personalBestEntry := profileType.PersonalBests[idx]
				if personalBestEntry.Event.Id != resultEntry.Eventid {
					continue
				}

				if averageInfo.Single == personalBestEntry.Single.Value {
					if personalBestEntry.Single.WR == "1" {
						averageInfo.SingleRecord = "WR"
						averageInfo.SingleRecordColor = constants.WR_COLOR
					} else if personalBestEntry.Single.CR == "1" {
						averageInfo.SingleRecord = "CR"
						averageInfo.SingleRecordColor = constants.CR_COLOR
					} else if personalBestEntry.Single.NR == "1" {
						averageInfo.SingleRecord = "NR"
						averageInfo.SingleRecordColor = constants.NR_COLOR
					} else {
						averageInfo.SingleRecord = "PB"
						averageInfo.SingleRecordColor = constants.PR_COLOR
					}
				}

				if averageInfo.Average == personalBestEntry.Average.Value {
					if personalBestEntry.Average.WR == "1" {
						averageInfo.AverageRecord = "WR"
						averageInfo.AverageRecordColor = constants.WR_COLOR
					} else if personalBestEntry.Average.CR == "1" {
						averageInfo.AverageRecord = "CR"
						averageInfo.AverageRecordColor = constants.CR_COLOR
					} else if personalBestEntry.Average.NR == "1" {
						averageInfo.AverageRecord = "NR"
						averageInfo.AverageRecordColor = constants.NR_COLOR
					} else {
						averageInfo.AverageRecord = "PB"
						averageInfo.AverageRecordColor = constants.PR_COLOR
					}
				}
			}
		}

		c.IndentedJSON(http.StatusOK, averageInfo)
	}
}
