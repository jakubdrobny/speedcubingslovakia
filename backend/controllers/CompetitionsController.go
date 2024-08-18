package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)


func GetFilteredCompetitions(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		filter := c.Param("filter")
		
		result := make([]models.CompetitionData, 0);
		competitions, err := models.GetAllCompetitions(db)
		if err != nil {
			log.Println("ERR GetAllCompetitions in GetFilteredCompetitions: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to query all competitions in database.");
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

		sort.Slice(result, func (i int, j int) bool {
			if filter == "Past" { return result[i].Enddate.After(result[j].Enddate) }
			return result[i].Enddate.Before(result[j].Enddate)
		})

		c.IndentedJSON(http.StatusOK, result);
	}
}

func GetCompetitionById(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")	
		
		rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c WHERE c.competition_id = $1;`, id)
		if err != nil {
			log.Println("ERR db.Query in GetCompetitionById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying competition by id.")
			return;
		}

		var competition models.CompetitionData
		found := false

		for rows.Next() {
			err := rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
			if err != nil {
				log.Println("ERR scanning competition data in GetCompetitionById: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed parsing competition from database.")
				return;
			}
			found = true
		}

		if !found {
			log.Println("ERR competition with id: ", id, " not found in GetCompetitionById.")
			c.IndentedJSON(http.StatusInternalServerError, "Competition not found.")	
			return;
		}

		err = competition.GetEvents(db)
		if err != nil {
			log.Println("ERR GetEvents in GetCompetitionById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get competition events.")
			return;
		}

		err = competition.GetScrambles(db)
		if err != nil {
			log.Println("ERR GetScrambles in GetCompetitionById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get competition scrambles.")
			return;
		}

		c.IndentedJSON(http.StatusOK, competition)
	}
}

func CreateCompetition(db *pgxpool.Pool, competition models.CompetitionData, envMap map[string]string) (string, string) {
	competition.RecomputeCompetitionId()
	err := competition.GenerateScrambles(envMap)
	if err != nil { return "ERR GenerateScrambles in PostCompetition: " + err.Error(), "Failed to generate scrambles." }

	tx, err := db.Begin(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return "ERR db.Begin in PostCompetition: " + err.Error(), "Failed to start transaction."
	}

	_, err = tx.Exec(context.Background(), `INSERT INTO competitions (competition_id, name, startdate, enddate) VALUES ($1,$2,$3,$4);`, competition.Id, competition.Name, competition.Startdate, competition.Enddate)
	if err != nil {
		tx.Rollback(context.Background())
		return "ERR tx.Exec INSERT INTO competitions in PostCompetition: " + err.Error(), "Failed inserting competition into database."
	}

	for _, event := range competition.Events {
		_, err := tx.Exec(context.Background(), `INSERT INTO competition_events (competition_id, event_id) VALUES ($1,$2);`, competition.Id, event.Id)
		if err != nil {
			tx.Rollback(context.Background())
			return "ERR tx.Exec INSERT INTO competition_events in PostCompetition: " + err.Error(), "Failed to insert competition events connections into database."
		}
	}

	for _, scrambleSet := range competition.Scrambles {
		for scrambleIdx, scramble := range scrambleSet.Scrambles {
			_, err := tx.Exec(context.Background(), `INSERT INTO scrambles (scramble, event_id, competition_id, "order", img) VALUES ($1,$2,$3,$4,$5);`, scramble.Scramble, scrambleSet.Event.Id, competition.Id, scrambleIdx + 1, scramble.Img)
			if err != nil {
				tx.Rollback(context.Background())
				return "ERR tx.Exec INSERT INTO scrambles in PostCompetition: " + err.Error(), "Failed to insert scrambles into database."
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil { return "ERR tx.commit in PostCompetition: " + err.Error(), "Failed to finish transaction." }

	return "", ""
}

func PostCompetition(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var competition models.CompetitionData

		if err := c.BindJSON(&competition); err != nil {
			log.Println("ERR BindJSON(&competition) in PostCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse competition data.")
			return
		}

		errLog, errOut := CreateCompetition(db, competition, envMap)
		if errLog != "" && errOut != "" {
			log.Println(errLog)
			c.IndentedJSON(http.StatusInternalServerError, errOut)
			return
		}

		c.IndentedJSON(http.StatusCreated, competition)
	}
}

func PutCompetition(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var competition models.CompetitionData

		if err := c.BindJSON(&competition); err != nil {
			log.Println("ERR BindJSON(&competition) in PutCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse competition data.")
			return
		}

		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Println("ERR db.begin in PutCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to start transaction.")
			tx.Rollback(context.Background())
			return
		}

		_, err = tx.Exec(context.Background(), `UPDATE competitions SET name = $1, startdate = $2, enddate = $3, timestamp = CURRENT_TIMESTAMP WHERE competition_id = $4;`, competition.Name, competition.Startdate, competition.Enddate, competition.Id)
		if err != nil {
			log.Println("ERR tx.Exec UPDATE competitions in PutCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to update competition info in database.")
			tx.Rollback(context.Background())
			return
		}

		err = models.UpdateCompetitionEvents(&competition, db, tx, envMap)
		if err != nil {
			log.Println("ERR UpdateCompetitionEvents in PutCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to update competition event connections in database.")
			tx.Rollback(context.Background())
			return
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Println("ERR tx.commit in in PutCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to finish transaction.")
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
			log.Println("ERR strconv(eventId) in GetResultsFromCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse eventId.")
			return
		}
		
		competitionResults, err := models.GetResultsFromCompetitionByEventName(db, cid, eid)
		if err != nil {
			log.Println("ERR GetResultsFromCompetitionByEventName in GetResultsFromCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get competition results.")
			return
		}

		c.IndentedJSON(http.StatusAccepted, competitionResults)
	}
}

func GetNewWeeklyCompetitionInfo(db *pgxpool.Pool) (models.CompetitionData, error) {
	var competition models.CompetitionData

	rows, err := db.Query(context.Background(), `SELECT c.name, c.enddate FROM competitions c WHERE c.competition_id LIKE ('WeeklyCompetition%') ORDER BY c.competition_id DESC;`)
	if err != nil { return models.CompetitionData{}, err }

	competition.Name = "Weekly Competition 1"
	competition.Startdate = utils.NextMonday()
	for rows.Next() {
		var latest models.CompetitionData
		err = rows.Scan(&latest.Name, &latest.Enddate)
		if err != nil { return models.CompetitionData{}, err }

		nameSplit := strings.Split(latest.Name, " ")
		if len(nameSplit) != 3 { return models.CompetitionData{}, fmt.Errorf("Invalid last competition name format: " + latest.Name + ". Should be Weekly Competition {number}")}

		newCompNum, err := strconv.Atoi(nameSplit[2])
		if err != nil { return models.CompetitionData{}, err }

		competition.Name = "Weekly Competition " + fmt.Sprint(newCompNum + 1)
		competition.Startdate = latest.Enddate
		
		rows.Close()
		break
	}

	competition.Enddate = competition.Startdate.AddDate(0, 0, 7)

	events, err := models.GetAvailableEvents(db)
	if err != nil { return models.CompetitionData{}, err }
	competition.Events = events

	return competition, nil
}

func AddNewWeeklyCompetition(db *pgxpool.Pool, envMap map[string]string) {
	competition, err := GetNewWeeklyCompetitionInfo(db)
	if err != nil {
		log.Println("ERR failed GetNewWeeklyCompetitionInfo in AddNewWeeklyCompetition: " + err.Error())
		return
	}

	errLog, errOut := CreateCompetition(db, competition, envMap)
	if errLog != "" && errOut != "" {
		log.Println(errLog)
		log.Println("ERR_OUT: " + errOut)
		return
	}

	log.Println("Competition successfully created !!!")
	log.Println(competition)
}