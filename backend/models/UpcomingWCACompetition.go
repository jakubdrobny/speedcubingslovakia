package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type UpcomingWCACompetition struct {
	Id                string             `json:"id"`
	Name              string             `json:"name"`
	Startdate         time.Time          `json:"startdate"`
	Enddate           time.Time          `json:"enddate"`
	Registered        int                `json:"registered"`
	RegistrationOpen  time.Time          `json:"registrationOpen"`
	RegistrationClose time.Time          `json:"registrationClose"`
	LatitudeDegrees   float64            `json:"latitudeDegrees"`
	LongitudeDegrees  float64            `json:"longitudeDegrees"`
	CompetitorLimit   int                `json:"competitorLimit"`
	VenueAddress      string             `json:"venueAddress"`
	Url               string             `json:"url"`
	Events            []CompetitionEvent `json:"events"`
	CountryId         string             `json:"-"`
	CountryName       string             `json:"-"`
	CountryIso2       string             `json:"-"`
	State             string             `json:"-"`
	City              string             `json:"-"`
}

type GetWCACompetitionsResponse struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	Startdate         string    `json:"start_date"`
	Enddate           string    `json:"end_date"`
	RegistrationOpen  time.Time `json:"registration_open"`
	RegistrationClose time.Time `json:"registration_close"`
	LatitudeDegrees   float64   `json:"latitude_degrees"`
	LongitudeDegrees  float64   `json:"longitude_degrees"`
	CompetitorLimit   int       `json:"competitor_limit"`
	Url               string    `json:"url"`
	CountryIso2       string    `json:"country_iso2"`
	VenueAddress      string    `json:"venue_address"`
	City              string    `json:"city"`
	EventIds          []string  `json:"event_ids"`
}

type UpcomingWCACompetitionRegistration struct {
	Id int `json:"id"`
}

func (c *UpcomingWCACompetition) GetRegistered(db pgx.Tx) error {
	url := fmt.Sprintf(
		"https://www.worldcubeassociation.org/api/v0/competitions/%s/registrations",
		c.Id,
	)
	body, err := utils.GetRequest(url)
	if err != nil {
		log.Println(
			"ERR utils.GetRequest(" + url + ") in UpcomingWCACompetition.GetRegistered: " + err.Error(),
		)
		return err
	}

	regs := []UpcomingWCACompetitionRegistration{}
	err = json.Unmarshal(body, &regs)
	if err != nil && string(body) != "" {
		log.Println(
			"ERR json.Unmarshal(" + string(
				body,
			) + ") in UpcomingWCACompetition.GetRegistered: " + err.Error(),
		)
		return err
	}

	c.Registered = len(regs)

	return nil
}

func (c *UpcomingWCACompetition) SaveEvents(db pgx.Tx) error {
	for _, event := range c.Events {
		_, err := db.Exec(
			context.Background(),
			`INSERT INTO upcoming_wca_competition_events (upcoming_wca_competition_id, event_id) SELECT $1 as upcoming_wca_competition_id, event_id FROM events e WHERE e.iconcode = $2 ON CONFLICT (upcoming_wca_competition_id, event_id) DO NOTHING;`,
			c.Id,
			event.Iconcode,
		)
		if err != nil {
			log.Println(
				"ERR db.Exec(insert into upcoming_wca_competition_events) in UpcomingWCACompetition.Save: " + err.Error(),
			)
			return err
		}
	}

	return nil
}

func (c *UpcomingWCACompetition) DeleteEvents(db pgx.Tx) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM upcoming_wca_competition_events WHERE upcoming_wca_competition_id = $1;`,
		c.Id,
	)
	if err != nil {
		log.Println(
			"ERR db.Exec(delete upcoming wca comp events) in UpcomingWCACompetition.DeleteEvents: " + err.Error(),
		)
		return err
	}

	return nil
}

func (c *UpcomingWCACompetition) UpdateEvents(db pgx.Tx) error {
	err := c.DeleteEvents(db)
	if err != nil {
		log.Println(
			"ERR UpcomingWCACompetition.DeleteEvents in UpcomingWCACompetition.UpdateEvents: " + err.Error(),
		)
		return err
	}

	err = c.SaveEvents(db)
	if err != nil {
		log.Println(
			"ERR UpcomingWCACompetition.Save in UpcomingWCACompetition.UpdateEvents: " + err.Error(),
		)
		return err
	}

	return nil
}

func (c *UpcomingWCACompetition) Save(db pgx.Tx) (pgconn.CommandTag, error) {
	res, err := db.Exec(
		context.Background(),
		`INSERT INTO upcoming_wca_competitions (upcoming_wca_competition_id, name, startdate, enddate, registered, competitor_limit, venue_address, url, country_id, registration_open, registration_close, latitude_degrees, longitude_degrees, state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) ON CONFLICT (upcoming_wca_competition_id) DO NOTHING;`,
		c.Id,
		c.Name,
		c.Startdate,
		c.Enddate,
		c.Registered,
		c.CompetitorLimit,
		c.VenueAddress,
		c.Url,
		c.CountryId,
		c.RegistrationOpen,
		c.RegistrationClose,
		c.LatitudeDegrees,
		c.LongitudeDegrees,
		c.State,
	)
	if err != nil {
		log.Println(
			"ERR db.Exec(insert into upcoming_wca_competitions) in UpcomingWCACompetition.Save: " + err.Error(),
		)
		return pgconn.CommandTag{}, err
	}

	if res.RowsAffected() != 0 {
		err = c.SaveEvents(db)
		if err != nil {
			log.Println(
				"ERR UpcomingWCACompetition.SaveEvents in UpcomingWCACompetition.Save: " + err.Error(),
			)
			return pgconn.CommandTag{}, err
		}
	} else {
		_, err := db.Exec(context.Background(), `UPDATE upcoming_wca_competitions SET name = $1, startdate = $2, enddate = $3, registered = $4, competitor_limit = $5, venue_address = $6, url = $7, country_id = $8, registration_open = $9, registration_close = $10, latitude_degrees = $11, longitude_degrees = $12, state = $13 WHERE upcoming_wca_competition_id = $14;`, c.Name, c.Startdate, c.Enddate, c.Registered, c.CompetitorLimit, c.VenueAddress, c.Url, c.CountryId, c.RegistrationOpen, c.RegistrationClose, c.LatitudeDegrees, c.LongitudeDegrees, c.State, c.Id)
		if err != nil {
			log.Println("ERR db.Exec(update upcoming_wca_competitions) in UpcomingWCACompetition.Save: " + err.Error())
			return pgconn.CommandTag{}, err
		}

		err = c.UpdateEvents(db)
		if err != nil {
			log.Println("ERR UpcomingWCACompetition.UpdateEvents in UpcomingWCACompetition.Save: " + err.Error())
			return pgconn.CommandTag{}, err
		}
	}

	return res, nil
}

func (c *UpcomingWCACompetition) DateFormatted() string {
	layout := "02 Jan 2006"

	startdateFormatted := c.Startdate.Format(layout)
	enddateFormatted := c.Enddate.Format(layout)

	dateFormatted := startdateFormatted
	if enddateFormatted != startdateFormatted {
		dateFormatted += " - " + enddateFormatted
	}

	return dateFormatted
}

func (c *UpcomingWCACompetition) GetEventNamesFromCompetitionEvents(
	events []CompetitionEvent,
) []string {
	return utils.Map(c.Events, func(ue CompetitionEvent) string {
		idx := slices.IndexFunc(events, func(e CompetitionEvent) bool {
			return e.Iconcode == ue.Iconcode
		})

		if idx == -1 {
			return ""
		}

		return events[idx].Displayname
	})
}

// if city is in format {city_name}, {state_name} works
// otherwise puts ""
// ONLY LOAD US STATES
func (c *UpcomingWCACompetition) LoadState() {
	if c.CountryIso2 != "US" {
		c.State = ""
		return
	}

	citySplitByCommaAndSpace := strings.Split(c.City, ", ")
	if len(citySplitByCommaAndSpace) != 2 {
		c.State = ""
	} else {
		state := citySplitByCommaAndSpace[1]
		if idx := slices.Index(constants.US_STATE_NAMES, state); idx != -1 {
			c.State = state
		} else {
			c.State = "Territories"
		}
	}
}
