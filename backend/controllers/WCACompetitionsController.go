package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

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

func GetWCACompAnnouncementSubscriptions(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")

		rows, err := db.Query(
			context.Background(),
			`SELECT c.country_id, c.name, CASE WHEN s.user_id IS NULL THEN false ELSE true END AS subscribed FROM countries c FULL JOIN wca_competitions_announcements_subscriptions s ON s.country_id = c.country_id AND user_id = $1;`,
			uid,
		)
		if err != nil {
			log.Println(
				"ERR db.Query(announcements_subscriptions) in GetWCACompAnnouncementSubscriptions: " + err.Error(),
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to query subscription data from db.",
			)
			return
		}

		subscriptions := make([]models.WCACompAnnouncementSubscriptions, 0)
		for rows.Next() {
			var sub models.WCACompAnnouncementSubscriptions
			err = rows.Scan(&sub.CountryId, &sub.CountryName, &sub.Subscribed)
			if err != nil {
				log.Println(
					"ERR rows.Scan(subscription=countryId,countryName,subscribed) in GetWCACompAnnouncementSubscriptions: " + err.Error(),
				)
				c.IndentedJSON(http.StatusInternalServerError, "Failed to parse subscription data.")
				return
			}

			subscriptions = append(subscriptions, sub)
		}

		c.IndentedJSON(http.StatusOK, subscriptions)
	}
}

func UpdateWCAAnnouncementSubscriptions(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body models.UpdateWCAAnnouncementSubscriptionsRequestBody

		if err := c.BindJSON(&body); err != nil {
			log.Println(
				"ERR c.BindJson(Update wca comp sub) in UpdateWCAAnnouncementSubscriptionsRequestBody: " + err.Error(),
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to parse subscription update data.",
			)
			return
		}

		uid := c.GetInt("uid")
		queryString := ``
		fmt.Println(uid, body.CountryName)
		args := make([]interface{}, 0)
		if !body.Subscribed {
			queryString = `DELETE FROM wca_competitions_announcements_subscriptions WHERE user_id = $1 AND country_id = (SELECT c.country_id FROM countries c WHERE c.name = $2);`
		} else {
			queryString = `INSERT INTO wca_competitions_announcements_subscriptions (user_id, country_id) SELECT $1 as user_id, c.country_id as country_id FROM countries c WHERE c.name = $2 ON CONFLICT (country_id, user_id) DO NOTHING;`
		}
		args = append(args, uid, body.CountryName)

		_, err := db.Exec(context.Background(), queryString, args...)
		if err != nil {
			log.Println(
				"ERR db.Exec(update wca comp announcement sub) in UpdateWCAAnnouncementSubscriptionsRequestBody: " + err.Error(),
			)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to update subscription.")
			return
		}

		c.IndentedJSON(http.StatusOK, body)
	}
}
