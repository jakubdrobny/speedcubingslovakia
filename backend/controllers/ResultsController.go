package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func GetResultsQuery(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		userName := c.Param("uname")
		competitionId := c.Param("cid")
		eventId, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			log.Println("ERR in strconv(eventId) in GetResultsQuery: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse eventId.")
			return
		}
		
		resultEntries := make([]models.ResultEntry, 0)

		if competitionId == "_" && userName == "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re WHERE re.event_id = $1;`, eventId)
			if err != nil {
				log.Println("ERR db.Query in GetResultsQuery (competitionId not set and userId not set): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					log.Println("ERR scanning resultEntryId in GetResultsQuery (competitionId not set and userId not set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					log.Println("ERR GetResultEntryById in GetResultsQuery (competitionId not set and userId not set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else if competitionId == "_" && userName != "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN users u ON u.user_id = re.user_id WHERE re.event_id = $1 AND UPPER(u.name) LIKE UPPER('%' || $2 || '%');`, eventId, userName)
			if err != nil {
				log.Println("ERR db.Query in GetResultsQuery (competitionId not set and userId set): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					log.Println("ERR scanning resultEntryId in GetResultsQuery (competitionId not set and userId set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					log.Println("ERR GetResultEntryById in GetResultsQuery (competitionId not set and userId set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else if competitionId != "_" && userName == "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re WHERE re.event_id = $1 AND UPPER(re.competition_id) LIKE UPPER('%' || $2 || '%');`, eventId, competitionId)
			if err != nil {
				log.Println("ERR db.Query in GetResultsQuery (competitionId set and userId not set): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					log.Println("ERR scanning resultEntryId in GetResultsQuery (competitionId set and userId not set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					log.Println("ERR GetResultEntryById in GetResultsQuery (competitionId set and userId not set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN users u ON u.user_id = re.user_id WHERE re.event_id = $1 AND UPPER(re.competition_id) LIKE UPPER('%' || $2 || '%') AND UPPER(u.name) LIKE UPPER('%' || $3 || '%');`, eventId, competitionId, userName)
			if err != nil {
				log.Println("ERR db.Query in GetResultsQuery (competitionId set and userId set): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					log.Println("ERR scanning resultEntryId in GetResultsQuery (competitionId set and userId set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying result entry from database.")
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					log.Println("ERR GetResultEntryById in GetResultsQuery (competitionId set and userId set): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		}

		c.IndentedJSON(http.StatusOK, resultEntries)
	}
}

func PostResultsValidation(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		type ReqBody struct {
			ResultId int `json:"resultId"`
			Verdict bool `json:"verdict"`
		}
		var body ReqBody

		if err := c.BindJSON(&body); err != nil {
			log.Println("ERR BindJSON in PostResultsValidation: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing data.")
			return;
		}

		resultEntry, err := models.GetResultEntryById(db, body.ResultId)
		if err != nil {
			log.Println("ERR GetResultEntryById in PostResultsValidation: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
			return;
		}

		statusId := 3
		if !body.Verdict { statusId = 2 }
		resultStatus, err := models.GetResultsStatus(db, statusId)
		if err != nil {
			log.Println("ERR GetResultsStatus in PostResultsValidation: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed getting result status in database.")
			return;
		}

		resultEntry.Status = resultStatus
		err = resultEntry.Update(db, true)
		if err != nil {
			log.Println("ERR resultEntry.Update in PostResultsValidation: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed updating result entry in database.")
			return
		}

		c.IndentedJSON(http.StatusCreated, "")
	}
}



func GetResultsByIdAndEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
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
			c.IndentedJSON(http.StatusInternalServerError, "Failed getting user information from database.")
			return
		}

		event, err := models.GetEventById(db, eventId)
		if err != nil {
			log.Println("ERR GetEventById in GetResultsByIdAndEvent: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed getting event information from database.")
			return
		}

		competition, err := models.GetCompetitionByIdObject(db, competitionId)
		if err != nil {
			log.Println("ERR GetCompetitionByIdObject in GetResultsByIdAndEvent: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed getting competition information from database.")
			return	
		}

		resultEntry, err := models.GetResultEntry(db, userId, competitionId, eventId)

		if err != nil {
			if err.Error() != "not found" {
				log.Println("ERR GetResultEntry in GetResultsByIdAndEvent: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed getting result entry from database.")
				return
			} else {
				approvedResultsStatus, err := models.GetResultsStatus(db, 3)
				if err != nil {
					log.Println("ERR GetResultsStatus.approved in GetResultsByIdAndEvent: " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed getting result status in database.")
					return	
				}

				resultEntry = models.ResultEntry{
					Userid: userId,
					Username: user.Name,
					Competitionid: competitionId,
					Competitionname: competition.Name,
					Eventid: event.Id,
					Eventname: event.Displayname,
					Iconcode: event.Iconcode,
					Format: event.Format,
					Solve1: "DNS",
					Solve2: "DNS",
					Solve3: "DNS",
					Solve4: "DNS",
					Solve5: "DNS",
					Comment: "",
					Status: approvedResultsStatus,
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
			resultEntry.Iconcode = event.Iconcode
			resultEntry.Format = event.Format

			err = resultEntry.Update(db)
			if err != nil {
				log.Println("ERR resultEntry.Update in GetResultsByIdAndEvent: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed updating results in database.")
				return
			}
		}

		c.IndentedJSON(http.StatusOK, resultEntry)
	}
}

func PostResults(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		var resultEntry models.ResultEntry
		var err error

		if err = c.BindJSON(&resultEntry); err != nil {
			log.Println("ERR BindJSON in PostResults: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed parsing data.")
			return;
		}

		err = resultEntry.Update(db)
		if err != nil {
			log.Println("ERR resultEntry.Update in PostResults: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed updating results in database.")
			return
		}

		c.IndentedJSON(http.StatusCreated, resultEntry)
	}
}

func GetProfileResults(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")

		uid, err := models.GetUserByWCAID(db, id)
		if err != nil {
			log.Println("ERR in GetProfileResults in GetUserByWCAID: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Finding user by WCA ID in database failed.")
			return
		}

		if uid == 0 {
			uid, err = models.GetUserByName(db, id)
			if err != nil {
				log.Println("ERR in GetProfileResults in GetUserByName: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Finding user by name in database failed.")
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