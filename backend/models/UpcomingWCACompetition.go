package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type UpcomingWCACompetition struct {
	Id               string             `json:"id"`
	Name             string             `json:"name"`
	Startdate        time.Time          `json:"startdate"`
	Enddate          time.Time          `json:"enddate"`
	Registered       int                `json:"registered"`
	RegistrationOpen time.Time          `json:"registrationOpen"`
	CompetitorLimit  int                `json:"competitorLimit"`
	VenueAddress     string             `json:"venueAddress"`
	Url              string             `json:"url"`
	Events           []CompetitionEvent `json:"events"`
	CountryId        string             `json:"-"`
}

type GetWCACompetitionsResponse struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	Startdate        string    `json:"start_date"`
	Enddate          string    `json:"end_date"`
	RegistrationOpen time.Time `json:"registration_open"`
	CompetitorLimit  int       `json:"competitor_limit"`
	Url              string    `json:"url"`
	CountryIso2      string    `json:"country_iso2"`
	VenueAddress     string    `json:"venue_address"`
	City             string    `json:"city"`
	EventIds         []string  `json:"event_ids"`
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
		return err
	}

	var regs []UpcomingWCACompetitionRegistration
	err = json.Unmarshal(body, &regs)
	if err != nil {
		return err
	}

	c.Registered = len(regs)

	return nil
}

func (c *UpcomingWCACompetition) Save(db pgx.Tx) error {
	res, err := db.Exec(
		context.Background(),
		`INSERT INTO upcoming_wca_competitions (upcoming_wca_competition_id, name, startdate, enddate, registered, competitor_limit, venue_address, url, country_id, registration_open) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (upcoming_wca_competition_id) DO NOTHING;`,
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
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() != 0 {
		for _, event := range c.Events {
			_, err := db.Exec(
				context.Background(),
				`INSERT INTO upcoming_wca_competition_events (upcoming_wca_competition_id, event_id) SELECT $1 as upcoming_wca_competition_id, event_id FROM events e WHERE e.iconcode = $2 ON CONFLICT (upcoming_wca_competition_id, event_id) DO NOTHING;`,
				c.Id,
				event.Iconcode,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
