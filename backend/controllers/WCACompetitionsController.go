package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/email"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

func GetWCACompetitionRegions(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		countries, err := models.GetCountries(db)
		if err != nil {
			log.Println("ERR GetCountries in GetRegionsGrouped: " + err.Error())
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed querying countries from database.",
			)
			return
		}
		countryGroup :=
			RegionSelectGroup{
				"Country",
				utils.Map(countries, func(c models.Country) string { return c.Name }),
			}

		if usIdx := slices.Index(countryGroup.GroupMembers, "United States"); usIdx != -1 {
			countryGroup.GroupMembers = utils.RemoveFromSliceBad(
				countryGroup.GroupMembers,
				usIdx,
			)
			for _, stateName := range constants.US_STATE_NAMES {
				countryGroup.GroupMembers = append(
					countryGroup.GroupMembers,
					fmt.Sprintf("United States, %v", stateName),
				)
			}
			slices.Sort(countryGroup.GroupMembers)
		}

		c.IndentedJSON(http.StatusOK, []RegionSelectGroup{countryGroup})
	}
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

func constructContent(
	notifEntry map[string][]models.UpcomingWCACompetition,
	username string,
	events []models.CompetitionEvent,
) string {
	content := "<!DOCTYPE html><html lang=\"en-US\"><head><title>New WCA competitions announced</title></head><body>"

	country_ids := make([]string, len(notifEntry))
	for country_id := range notifEntry {
		country_ids = append(country_ids, country_id)
	}
	sort.Slice(country_ids, func(i, j int) bool {
		return country_ids[i] < country_ids[j]
	})

	content += fmt.Sprintf("Hi %s!<br/><br/>", username)
	content += "there have been new WCA competitions announced in regions you have subscribed to. Here :<br/><br/>"

	content += "<table style=\"border-collapse: collapse;\">"
	for _, country_id := range country_ids {
		comps := notifEntry[country_id]
		sort.Slice(comps, func(i, j int) bool {
			return comps[i].Startdate.Before(comps[j].Startdate)
		})

		if len(comps) == 0 {
			continue
		}

		borderBottomStyle := " style=\"border-bottom: 1px solid black;\""
		content += fmt.Sprintf(
			"<tr%s><td style=\"vertical-align:middle;\"><img style=\"vertical-align: middle;\" title=\"Flag of %s\" alt=\"flag of %s\" src=\"https://flagpedia.net/data/flags/h20/%s.png\"/><h1 style=\"vertical-align: middle; display: inline-block; margin: 0; padding-left: 10px;\">%s</h1></td></tr>",
			borderBottomStyle,
			comps[0].CountryName,
			comps[0].CountryName,
			strings.ToLower(comps[0].CountryIso2),
			comps[0].CountryName,
		)

		for _, comp := range comps {
			content += fmt.Sprintf(
				"<tr><td style=\"padding-left: 10px\"><h2 style=\"margin: 0\">%s</h2></td></tr>",
				comp.Name,
			)

			content += fmt.Sprintf(
				"<tr><td style=\"padding-left: 20px\"><b>Place:</b> <span style=\"font-weight: normal;\">%s</span></td></tr>",
				comp.VenueAddress,
			)
			content += fmt.Sprintf(
				"<tr><td style=\"padding-left: 20px\"><b>Date:</b> <span style=\"font-weight: normal;\">%s</span></td></tr>",
				comp.DateFormatted(),
			)
			content += fmt.Sprintf(
				"<tr><td style=\"padding-left: 20px\"><b>Competitor limit:</b> <span style=\"font-weight: normal;\">%d</span></td></tr>",
				comp.CompetitorLimit,
			)
			content += fmt.Sprintf(
				"<tr><td style=\"padding-left: 20px\"><b>Registration opens:</b> <span style=\"font-weight: normal;\">%s</span></td></tr>",
				comp.RegistrationOpen.UTC().Format("02 Jan 2006 15:04:05 MST"),
			)
			content += fmt.Sprintf(
				"<tr><td style=\"padding-left: 20px\"><b>Events:</b> <span style=\"font-weight: normal;\">%s</span></td></tr>",
				strings.Join(comp.GetEventNamesFromCompetitionEvents(events), ", "),
			)
			content += fmt.Sprintf(
				"<tr><td style=\"font-weight: normal; padding-left: 20px\">For more info click <a href=\"%s\"><b>here</b></a>.</td></tr>",
				comp.Url,
			)
		}
		content += "<tr><td>&nbsp;</td></tr>"
	}
	content += "</table><br/>"

	content += fmt.Sprintf(
		"Thank you for subscribing to our competition announcement newsletter.<br/><br/>If you want to prepare for WCA competitions and compete with your friends, don't forget to compete in Online Weekly Competitions at our <a href=\"https://speedcubingslovakia.sk/competitions\"><b>website</b></a>.<br/><br/>",
	)
	content += fmt.Sprintf("Have a great day.<br/><br/><i>Jakub</i></body></html>")

	return content
}

// notifications if user_id -> country_id -> [] of comps
func SendCompAnnouncementSubscriptions(
	db *pgxpool.Pool,
	envMap map[string]string,
	notifications map[int]map[string][]models.UpcomingWCACompetition,
) error {
	log.Println("Sending email notifications to WCA competitions announcements subscribers...")

	events, err := models.GetAvailableEvents(db)
	if err != nil {
		log.Println(
			"ERR models.GetAvailableEvents in SendCompAnnouncementSubscriptions: " + err.Error(),
		)
		return err
	}

	for userId, notifEntry := range notifications {
		user, err := models.GetUserById(db, userId)
		if err != nil {
			log.Println(
				"ERR models.GetUserById in SendCompAnnouncementSubscriptions: " + err.Error(),
			)
			return err
		}
		log.Println("Sending email notification to user: " + user.Name)

		from := envMap["MAIL_USERNAME"]
		to := user.Email
		subject := "New WCA competitions announced"
		if os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV") == "development" {
			subject = "DEVELOPMENT: " + subject
		}
		content := constructContent(notifEntry, user.Name, events)

		err = email.SendMail(from, to, subject, content, envMap)
		if err != nil {
			log.Println("ERR email.SendMail in SendCompAnnouncementSubscriptions: " + err.Error())
			return err
		}

		log.Println("Email sent successfully.")
	}

	log.Println("All emails sent successfully.")

	return nil
}

func MakeCompAnnouncementContent(
	comp models.UpcomingWCACompetition,
	events []models.CompetitionEvent,
) string {
	timeLoc, _ := time.LoadLocation("Europe/Bratislava")
	content := fmt.Sprintf(
		"Hello everyone,\n\nnew WCA competition in Slovakia has just been announced:\n\n**Name:** %s<br>**Place:** %s<br>**Date:** %s<br>**Events:** %s\n\n**Registration** starts on **%s**. Mark it in your calendars so you don't miss it.\n\nFor more info check out the [competition website](%s).\n\nHope to see you there.\n\nSpeedcubing Slovakia",
		comp.Name,
		comp.DateFormatted(),
		comp.VenueAddress,
		strings.Join(comp.GetEventNamesFromCompetitionEvents(events), ", "),
		comp.RegistrationOpen.UTC().
			In(timeLoc).
			Format("2 Jan 2006 15:04:05"),
		comp.Url,
	)

	return content
}

// make announcements for newly announced WCA competitions in Slovakia
func MakeCompAnnouncementAnnouncements(
	db *pgxpool.Pool,
	envMap map[string]string,
	comps []models.UpcomingWCACompetition,
) error {
	competitions := "competition"
	if len(comps) != 1 {
		competitions += "s"
	}
	log.Printf(
		"Creating announcements for %d newly announced %s in Slovakia\n",
		len(comps),
		competitions,
	)

	sort.Slice(comps, func(i, j int) bool {
		return comps[i].Startdate.Before(comps[j].Startdate)
	})

	events, err := models.GetAvailableEvents(db)
	if err != nil {
		log.Println(
			"ERR models.GetAvailableEvents in SendCompAnnouncementSubscriptions: " + err.Error(),
		)
		return err
	}

	compAnnouncementTag, err := models.GetCompetitionAnnouncementTag(db)
	if err != nil {
		log.Println(
			"ERR models.GetCompetitionAnnouncementTag in MakeCompAnnouncementContent: " + err.Error(),
		)
		return err
	}

	for _, comp := range comps {
		announcement := models.AnnouncementState{
			Title:    comp.Name + " announced",
			Content:  MakeCompAnnouncementContent(comp, events),
			AuthorId: 1,
			Tags:     []models.Tag{compAnnouncementTag},
		}

		logMsg, retMsg := announcement.Create(db, envMap)
		if logMsg != "" || retMsg != "" {
			log.Println("Failed to create announcement, checkout the error log below.")
			log.Println(logMsg)
			return errors.New(retMsg)
		}

		log.Printf("%s announcement was made successfully.\n", comp.Name)
	}

	log.Println("Competitions announced successfully.")

	return nil
}

func CheckUpcomingWCACompetitions(db *pgxpool.Pool, envMap map[string]string) error {
	log.Println("Querying countries...")

	countries, err := models.GetCountries(db)
	if err != nil {
		log.Println("ERR models.GetCountries in CheckUpcomingWCACompetitions: " + err.Error())
		return err
	}

	log.Println("Starting db transaction...")
	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("ERR db.Begin in CheckUpcomingWCACompetitions: " + err.Error())
		return err
	}
	defer tx.Rollback(context.Background())

	log.Println("Checking if already announced comps are loaded in db...")
	upcomingCompsFromDB, err := GetSavedUpcomingWCACompetitions(db, "_")
	if err != nil {
		log.Println(
			"ERR GetSavedUpcomingWCACompetitions in CheckUpcomingWCACompetitions: " + err.Error(),
		)
		return err
	}

	notifySubscribers := len(upcomingCompsFromDB) > 0
	notifications := make(map[int]map[string][]models.UpcomingWCACompetition)
	newlyAnnouncedSlovakComps := make([]models.UpcomingWCACompetition, 0)

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
				log.Printf(
					"==========\nFailed to unmarshal this body: %v\n==========\n",
					string(body),
				)
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
					CountryName:      country.Name,
					CountryIso2:      country.Iso2,
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
					// for now don't log anything if competition exists so i can easily check which comps were not in the db after each run
					//log.Println(
					//"Competition " + upcomingWCACompetition.Name + " is already in the database.",
					//)
				} else {
					log.Println("Competition " + upcomingWCACompetition.Name + " saved successfully.")

					if !notifySubscribers {
						continue
					}

					log.Println("Querying subscribers...")
					rows, err := tx.Query(context.Background(), `SELECT user_id FROM wca_competitions_announcements_subscriptions WHERE country_id = $1;`, country.Id)
					if err != nil {
						log.Println("ERR tx.Query(subscriptions) for " + country.Id + " in CheckUpcomingWCACompetitions: " + err.Error())
						return err
					}

					// find subscribers and add notification to them
					for rows.Next() {
						var currentUserId int
						err = rows.Scan(&currentUserId)
						if err != nil {
							log.Println("ERR rows.Scan(user_id) in CheckUpcomingWCACompetitions: " + err.Error())
							return err
						}

						if _, ok := notifications[currentUserId]; !ok {
							notifications[currentUserId] = make(map[string][]models.UpcomingWCACompetition)
						}
						if _, ok := notifications[currentUserId][country.Id]; !ok {
							notifications[currentUserId][country.Id] = make([]models.UpcomingWCACompetition, 0)
						}

						notifications[currentUserId][country.Id] = append(notifications[currentUserId][country.Id], upcomingWCACompetition)
					}

					// if comp added is slovak, add it to the list of comps need to make an announcement for
					if country.Iso2 == "SK" {
						newlyAnnouncedSlovakComps = append(newlyAnnouncedSlovakComps, upcomingWCACompetition)
					}
				}
			}

			page += 1
		}
	}

	defer func() {
		if notifySubscribers {
			SendCompAnnouncementSubscriptions(db, envMap, notifications)
			MakeCompAnnouncementAnnouncements(db, envMap, newlyAnnouncedSlovakComps)
		}
	}()

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
			`SELECT c.country_id, c.name, CASE WHEN s.user_id IS NULL THEN false ELSE true END AS subscribed FROM countries c LEFT JOIN wca_competitions_announcements_subscriptions s ON s.country_id = c.country_id AND user_id = $1;`,
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
			a, b := rows.Values()
			log.Printf("values: %v, error; %v\n", a, b)
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
