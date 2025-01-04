package controllers

import (
	"context"
	"encoding/json"
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
	return func(c *gin.Context) {
		filter := c.Param("filter")

		result := make([]models.CompetitionData, 0)
		competitions, err := models.GetAllCompetitions(db)
		if err != nil {
			log.Println("ERR GetAllCompetitions in GetFilteredCompetitions: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to query all competitions in database.",
			)
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

		sort.Slice(result, func(i int, j int) bool {
			if filter == "Past" {
				return result[i].Enddate.After(result[j].Enddate)
			}
			return result[i].Enddate.Before(result[j].Enddate)
		})

		c.IndentedJSON(http.StatusOK, result)
	}
}

func GetCompetitionById(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		rows, err := db.Query(
			context.Background(),
			`SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c WHERE c.competition_id = $1;`,
			id,
		)
		if err != nil {
			log.Println("ERR db.Query in GetCompetitionById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying competition by id.")
			return
		}

		var competition models.CompetitionData
		found := false

		for rows.Next() {
			err := rows.Scan(
				&competition.Id,
				&competition.Name,
				&competition.Startdate,
				&competition.Enddate,
			)
			if err != nil {
				log.Println("ERR scanning competition data in GetCompetitionById: " + err.Error())
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed parsing competition from database.",
				)
				return
			}
			found = true
		}

		if !found {
			log.Println("ERR competition with id: ", id, " not found in GetCompetitionById.")
			c.IndentedJSON(http.StatusInternalServerError, "Competition not found.")
			return
		}

		err = competition.GetEvents(db)
		if err != nil {
			log.Println("ERR GetEvents in GetCompetitionById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get competition events.")
			return
		}

		err = competition.GetScrambles(db)
		if err != nil {
			log.Println("ERR GetScrambles in GetCompetitionById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get competition scrambles.")
			return
		}

		c.IndentedJSON(http.StatusOK, competition)
	}
}

func CreateCompetition(
	db *pgxpool.Pool,
	competition models.CompetitionData,
	envMap map[string]string,
) (string, string) {
	competition.RecomputeCompetitionId()
	err := competition.GenerateScrambles(envMap)
	if err != nil {
		return "ERR GenerateScrambles in PostCompetition: " + err.Error(), "Failed to generate scrambles."
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return "ERR db.Begin in PostCompetition: " + err.Error(), "Failed to start transaction."
	}

	_, err = tx.Exec(
		context.Background(),
		`INSERT INTO competitions (competition_id, name, startdate, enddate) VALUES ($1,$2,$3,$4);`,
		competition.Id,
		competition.Name,
		competition.Startdate,
		competition.Enddate,
	)
	if err != nil {
		tx.Rollback(context.Background())
		return "ERR tx.Exec INSERT INTO competitions in PostCompetition: " + err.Error(), "Failed inserting competition into database."
	}

	for _, event := range competition.Events {
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO competition_events (competition_id, event_id) VALUES ($1,$2);`,
			competition.Id,
			event.Id,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return "ERR tx.Exec INSERT INTO competition_events in PostCompetition: " + err.Error(), "Failed to insert competition events connections into database."
		}
	}

	for _, scrambleSet := range competition.Scrambles {
		for scrambleIdx, scramble := range scrambleSet.Scrambles {
			_, err := tx.Exec(
				context.Background(),
				`INSERT INTO scrambles (scramble, event_id, competition_id, "order", img) VALUES ($1,$2,$3,$4,$5);`,
				scramble.Scramble,
				scrambleSet.Event.Id,
				competition.Id,
				scrambleIdx+1,
				scramble.Img,
			)
			if err != nil {
				tx.Rollback(context.Background())
				return "ERR tx.Exec INSERT INTO scrambles in PostCompetition: " + err.Error(), "Failed to insert scrambles into database."
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return "ERR tx.commit in PostCompetition: " + err.Error(), "Failed to finish transaction."
	}

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

		_, err = tx.Exec(
			context.Background(),
			`UPDATE competitions SET name = $1, startdate = $2, enddate = $3, timestamp = CURRENT_TIMESTAMP WHERE competition_id = $4;`,
			competition.Name,
			competition.Startdate,
			competition.Enddate,
			competition.Id,
		)
		if err != nil {
			log.Println("ERR tx.Exec UPDATE competitions in PutCompetition: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to update competition info in database.",
			)
			tx.Rollback(context.Background())
			return
		}

		err = models.UpdateCompetitionEvents(&competition, db, tx, envMap)
		if err != nil {
			log.Println("ERR UpdateCompetitionEvents in PutCompetition: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to update competition event connections in database.",
			)
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
	return func(c *gin.Context) {
		cid := c.Param("cid")
		eid, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			log.Println("ERR strconv(eventId) in GetResultsFromCompetition: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse eventId.")
			return
		}

		competitionResults, err := models.GetResultsFromCompetitionByEventName(db, cid, eid)
		if err != nil {
			log.Println(
				"ERR GetResultsFromCompetitionByEventName in GetResultsFromCompetition: " + err.Error(),
			)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get competition results.")
			return
		}

		c.IndentedJSON(http.StatusAccepted, competitionResults)
	}
}

func GetNewWeeklyCompetitionInfo(db *pgxpool.Pool) (models.CompetitionData, error) {
	var competition models.CompetitionData

	rows, err := db.Query(
		context.Background(),
		`SELECT c.name, c.enddate FROM competitions c WHERE c.competition_id LIKE ('WeeklyCompetition%') ORDER BY c.enddate DESC;`,
	)
	if err != nil {
		return models.CompetitionData{}, err
	}

	competition.Name = "Weekly Competition 1"
	competition.Startdate = utils.NextMonday()
	for rows.Next() {
		var latest models.CompetitionData
		err = rows.Scan(&latest.Name, &latest.Enddate)
		if err != nil {
			return models.CompetitionData{}, err
		}

		nameSplit := strings.Split(latest.Name, " ")
		log.Println(nameSplit)
		if len(nameSplit) != 3 {
			return models.CompetitionData{}, fmt.Errorf(
				"Invalid last competition name format: " + latest.Name + ". Should be Weekly Competition {number}",
			)
		}

		newCompNum, err := strconv.Atoi(nameSplit[2])
		if err != nil {
			return models.CompetitionData{}, err
		}

		competition.Name = "Weekly Competition " + fmt.Sprint(newCompNum+1)
		competition.Startdate = latest.Enddate

		rows.Close()
		break
	}

	competition.Enddate = competition.Startdate.AddDate(0, 0, 7)

	events, err := models.GetAvailableEvents(db)
	if err != nil {
		return models.CompetitionData{}, err
	}
	competition.Events = events

	return competition, nil
}

func AddNewWeeklyCompetition(db *pgxpool.Pool, envMap map[string]string) {
	competition, err := GetNewWeeklyCompetitionInfo(db)
	if err != nil {
		log.Println(
			"ERR failed GetNewWeeklyCompetitionInfo in AddNewWeeklyCompetition: " + err.Error(),
		)
		return
	}

	log.Println(competition)

	errLog, errOut := CreateCompetition(db, competition, envMap)
	if errLog != "" && errOut != "" {
		log.Println(errLog)
		log.Println("ERR_OUT: " + errOut)
		return
	}

	log.Println("Competition successfully created !!!")
	log.Println(competition)
}

func GetUpcomingWCACompetitions(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		regionPrecise := c.Query("regionPrecise")
		region, err := models.GetCountryByName(db, regionPrecise)
		if err != nil {
			log.Println("ERR models.GetCountryByName in GetUpcomingWCACompetitions: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to get country information from name.",
			)
			return
		}

		if region.Id == "" {
			c.IndentedJSON(http.StatusInternalServerError, "Country does not exists.")
			return
		}

		upcomingCompetitions, err := GetSavedUpcomingWCACompetitions(db, region.Id)
		if err != nil {
			log.Println(
				"ERR GetSavedUpcomingWCACompetitions in GetUpcomingWCACompetitions: " + err.Error(),
			)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to load competitions.")
			return
		}

		c.IndentedJSON(http.StatusOK, upcomingCompetitions)
	}
}

func GetUpcomingWCACompetitionEvents(
	db *pgxpool.Pool,
	cid string,
) ([]models.CompetitionEvent, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT e.iconcode FROM upcoming_wca_competition_events uwce JOIN events e ON uwce.event_id = e.event_id WHERE uwce.upcoming_wca_competition_id = $1;`,
		cid,
	)
	if err != nil {
		log.Println(
			"ERR db.Query(upcoming_wca_competition_events) in GetUpcomingWCACompetitionEvents: " + err.Error(),
		)
		return []models.CompetitionEvent{}, err
	}

	events := make([]models.CompetitionEvent, 0)
	for rows.Next() {
		var event models.CompetitionEvent
		err = rows.Scan(&event.Iconcode)
		if err != nil {
			log.Println(
				"ERR rows.scan(event.iconcode) in GetUpcomingWCACompetitionEvents: " + err.Error(),
			)
			return []models.CompetitionEvent{}, err
		}

		events = append(events, event)
	}

	return events, nil
}

// set countryId = '_' to query for all competitions
func GetSavedUpcomingWCACompetitions(
	db *pgxpool.Pool,
	countryId string,
) ([]models.UpcomingWCACompetition, error) {
	queryString := `SELECT upcoming_wca_competition_id as id, name, startdate, enddate, registered, competitor_limit, venue_address, url, registration_open FROM upcoming_wca_competitions`
	args := make([]interface{}, 0)
	if countryId != "_" {
		queryString += " WHERE country_id = $1"
		args = append(args, countryId)
	}
	queryString += " ORDER BY enddate;"

	rows, err := db.Query(
		context.Background(),
		queryString,
		args...,
	)
	if err != nil {
		log.Println(
			"ERR db.Query(upcoming_wca_competitions) in GetSavedUpcomingWCACompetitions: " + err.Error(),
		)
		return []models.UpcomingWCACompetition{}, err
	}

	upcoming_comps := make([]models.UpcomingWCACompetition, 0)
	for rows.Next() {
		var upcoming_comp models.UpcomingWCACompetition
		err = rows.Scan(
			&upcoming_comp.Id,
			&upcoming_comp.Name,
			&upcoming_comp.Startdate,
			&upcoming_comp.Enddate,
			&upcoming_comp.Registered,
			&upcoming_comp.CompetitorLimit,
			&upcoming_comp.VenueAddress,
			&upcoming_comp.Url,
			&upcoming_comp.RegistrationOpen,
		)
		if err != nil {
			log.Println(
				"ERR rows.scan(upcoming_comp) in GetSavedUpcomingWCACompetitions: " + err.Error(),
			)
			return []models.UpcomingWCACompetition{}, err
		}

		events, err := GetUpcomingWCACompetitionEvents(db, upcoming_comp.Id)
		if err != nil {
			log.Println(
				"ERR GetUpcomingWCACompetitionEvents in GetSavedUpcomingWCACompetitions: " + err.Error(),
			)
			return []models.UpcomingWCACompetition{}, err
		}

		upcoming_comp.Events = events
		upcoming_comps = append(upcoming_comps, upcoming_comp)
	}

	return upcoming_comps, nil
}

func CheckUpcomingWCACompetitions(db *pgxpool.Pool) error {
	countries, err := models.GetCountries(db)
	if err != nil {
		log.Println("ERR models.GetCountries in CheckUpcomingWCACompetitions: " + err.Error())
		return err
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("ERR db.Begin in CheckUpcomingWCACompetitions: " + err.Error())
		return err
	}
	defer tx.Rollback(context.Background())

	//upcomingCompsFromDB, err := GetSavedUpcomingWCACompetitions(db, "_")
	//if err != nil {
	//log.Println(
	//"ERR GetSavedUpcomingWCACompetitions in CheckUpcomingWCACompetitions: " + err.Error(),
	//)
	//return err
	//}

	//notifySubscribers := len(upcomingCompsFromDB) > 0
	//notifications := make(map[int]map[string][]models.UpcomingWCACompetition)

	for _, country := range countries {
		page := 1
		can := true
		for can {
			url := fmt.Sprintf(
				"https://www.worldcubeassociation.org/api/v0/competitions?country_iso2=%s&page=%d&sort=-end_date",
				country.Iso2,
				page,
			)
			body, err := utils.GetRequest(url)
			if err != nil {
				log.Println(
					"ERR utils.GetRequest(url=" + url + ") in CheckUpcomingWCACompetitions: " + err.Error(),
				)
				return err
			}

			var respComps []models.GetWCACompetitionsResponse
			err = json.Unmarshal(body, &respComps)
			if err != nil {
				log.Println("ERR json.Unmarshal in CheckUpcomingWCACompetitions: " + err.Error())
				return err
			}

			if len(respComps) < 25 {
				can = false
			}

			for _, respComp := range respComps {
				const layout = "2006-01-02"
				enddate, _ := time.Parse(layout, respComp.Enddate)
				if enddate.Before(time.Now().Round(0)) {
					can = false
					break
				}

				startdate, _ := time.Parse(layout, respComp.Startdate)
				upcomingWCACompetition := models.UpcomingWCACompetition{
					Id:              respComp.Id,
					Name:            respComp.Name,
					Startdate:       startdate,
					Enddate:         enddate,
					CompetitorLimit: respComp.CompetitorLimit,
					VenueAddress:    respComp.VenueAddress + ", " + respComp.City + ", " + country.Name,
					Url:             respComp.Url,
					Events: utils.Map(
						respComp.EventIds,
						func(iconcode string) models.CompetitionEvent { return models.CompetitionEvent{Iconcode: iconcode} },
					),
					CountryId:        country.Id,
					RegistrationOpen: respComp.RegistrationOpen,
				}

				err = upcomingWCACompetition.GetRegistered(tx)
				if err != nil {
					log.Println(
						"ERR upcomingWCACompetition.GetRegistered in CheckUpcomingWCACompetitions: " + err.Error(),
					)
					return err
				}

				res, err := upcomingWCACompetition.Save(tx)
				if err != nil {
					log.Println(
						"ERR upcomingWCACompetition.Save in CheckUpcomingWCACompetitions: " + err.Error(),
					)
					return err
				}

				if res.RowsAffected() == 0 {
					log.Println(
						"Competition " + upcomingWCACompetition.Name + " is already in the database.",
					)
				} else {
					log.Println("Competition " + upcomingWCACompetition.Name + " saved successfully.")
				}
			}

			page += 1
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Println("ERR tx.Commit in CheckUpcomingWCACompetitions: " + err.Error())
		return err
	}

	return nil
}

func DeletePastWCACompetitions(db *pgxpool.Pool) error {
	res, err := db.Exec(
		context.Background(),
		`DELETE FROM upcoming_wca_competitions WHERE date_trunc('day', now()) > enddate;`,
	)
	if err != nil {
		log.Println(
			"ERR db.Exec(delete upcoming_wca_comps) in DeletePastWCACompetitions: " + err.Error(),
		)
		return err
	}

	fmt.Printf(
		"Successfully deleted %d upcoming WCA competitions from db, because they are not upcoming anymore :DD\n",
		res.RowsAffected(),
	)

	return nil
}
