package controllers

import (
	"context"
	"math/rand"
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
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		
		var resultEntries []models.ResultEntry

		if competitionId == "_" && userName == "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re WHERE re.event_id = $1;`, eventId)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else if competitionId == "_" && userName != "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN users u ON u.user_id = re.user_id WHERE re.event_id = $1 AND u.name = $2;`, eventId, userName)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else if competitionId != "_" && userName == "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re WHERE re.event_id = $1 AND re.competition_id = $2;`, eventId, competitionId)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN users u ON u.user_id = re.user_id WHERE re.event_id = $1 AND re.competition_id = $2 AND u.name = $3;`, eventId, competitionId, userName)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
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
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}

		resultEntry, err := models.GetResultEntryById(db, body.ResultId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}

		statusId := 3
		if !body.Verdict { statusId = 2 }
		resultStatus, err := models.GetResultsStatus(db, statusId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}

		resultEntry.Status = resultStatus
		err = resultEntry.Update(db, true)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return
		}

		c.IndentedJSON(http.StatusCreated, "")
	}
}



func GetResultsByIdAndEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		eventId, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		competitionId := c.Param("cid")
		userId := c.MustGet("uid").(int)

		user, err := models.GetUserById(db, userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		event, err := models.GetEventById(db, eventId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		competition, err := models.GetCompetitionByIdObject(db, competitionId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return	
		}

		resultEntry, err := models.GetResultEntry(db, userId, competitionId, eventId)

		if err != nil {
			if err.Error() != "not found" {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			} else {
				approvedResultsStatus, err := models.GetResultsStatus(db, 3)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return	
				}

				resultEntry = models.ResultEntry{
					Id: rand.Int(),
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
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
			}
		} else {
			currentStatus, err := models.GetResultsStatus(db, resultEntry.Status.Id)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
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
				c.IndentedJSON(http.StatusInternalServerError, err)
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
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}
		
		err = resultEntry.Update(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusCreated, resultEntry)
	}
}