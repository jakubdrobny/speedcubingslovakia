package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)


func GetFilteredCompetitions(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		filter := c.Param("filter")
		
		result := make([]models.CompetitionData, 0);
		competitions, err := models.GetAllCompetitions(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return
		}

		now := time.Now()
		if filter == "Past" {
			for _, competition := range competitions {
				if competition.Enddate.Before(now) {
					result = append(result, competition)
				}
			}
		} else if filter == "Current" {
			for _, competition := range competitions {
				if competition.Startdate.Before(now) && now.Before(competition.Enddate) {
					result = append(result, competition)
				}
			}
		} else if filter == "Future" {
			for _, competition := range competitions {
				if now.Before(competition.Startdate) {
					result = append(result, competition)
				}
			}
		}

		c.IndentedJSON(http.StatusOK, result);
	}
}

func GetCompetitionById(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")	
		
		rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c WHERE c.competition_id = $1;`, id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}

		var competition models.CompetitionData
		found := false

		for rows.Next() {
			err := rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)	
				return;
			}
			found = true
		}

		if !found {
			c.IndentedJSON(http.StatusInternalServerError, "Competition not found.")	
			return;
		}

		err = competition.GetEvents(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}

		err = competition.GetScrambles(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}
		
		c.IndentedJSON(http.StatusOK, competition)
	}
}

func PostCompetition(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var competition models.CompetitionData

		if err := c.BindJSON(&competition); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "what bro")
			return
		}

		competition.RecomputeCompetitionId()

		tx, err := db.Begin(context.Background())
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			tx.Rollback(context.Background())
			return
		}

		_, err = tx.Exec(context.Background(), `INSERT INTO competitions (competition_id, name, startdate, enddate) VALUES ($1,$2,$3,$4);`, competition.Id, competition.Name, competition.Startdate, competition.Enddate)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			tx.Rollback(context.Background())
			return
		}

		for _, event := range competition.Events {
			_, err := tx.Exec(context.Background(), `INSERT INTO competition_events (competition_id, event_id) VALUES ($1,$2);`, competition.Id, event.Id)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				tx.Rollback(context.Background())
				return
			}
		}

		for _, scrambleSet := range competition.Scrambles {
			for scrambleIdx, scramble := range scrambleSet.Scrambles {
				_, err := tx.Exec(context.Background(), `INSERT INTO scrambles (scramble, event_id, competition_id, "order") VALUES ($1,$2,$3,$4);`, scramble, scrambleSet.Event.Id, competition.Id, scrambleIdx + 1)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					tx.Rollback(context.Background())
					return
				}
			}
		}

		err = tx.Commit(context.Background())
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return	
		}

		c.IndentedJSON(http.StatusCreated, competition)
	}
}

func PutCompetition(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var competition models.CompetitionData

		if err := c.BindJSON(&competition); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		_, err := db.Exec(context.Background(), `UPDATE competitions SET name = $1, startdate = $2, enddate = $3, timestamp = CURRENT_TIMESTAMP WHERE competition_id = $4;`, competition.Name, competition.Startdate, competition.Enddate, competition.Id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		err = models.UpdateCompetitionEvents(&competition, db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusCreated, competition)
	}
}

func GetResultsFromCompetition(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		cid := c.Param("cid")
		eid, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		competitionResults, err := models.GetResultsFromCompetitionByEventName(db, cid, eid)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusAccepted, competitionResults)
	}
}