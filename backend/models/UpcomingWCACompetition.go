package models

import "time"

type UpcomingWCACompetition struct {
	Id              int                `json:"id"`
	Name            string             `json:"name"`
	Startdate       time.Time          `json:"startdate"`
	Enddate         time.Time          `json:"enddate"`
	Registered      int                `json:"registered"`
	CompetitorLimit int                `json:"competitor_limit"`
	VenueAddress    string             `json:"venue_address"`
	Url             string             `json:"url"`
	Events          []CompetitionEvent `json:"events"`
}
