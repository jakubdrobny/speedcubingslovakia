package models

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

type WCACompAnnouncementsPositionSubscriptions struct {
	Id               int     `json:"id"`
	UserId           int     `json:"-"`
	LatitudeDegrees  float64 `json:"lat"`
	LongitudeDegrees float64 `json:"long"`
	Radius           float64 `json:"radius"`
	New              bool    `json:"new"`
	Open             bool    `json:"open"`
}

// returns the insert id and error
func (s *WCACompAnnouncementsPositionSubscriptions) Insert(
	db interfaces.DB,
) (int, error) {
	var id int

	err := db.QueryRow(
		context.Background(),
		`INSERT INTO wca_competitions_announcements_position_subscriptions (radius, latitude_degrees, longitude_degrees, user_id) VALUES ($1,$2,$3,$4) RETURNING wca_competitions_announcements_position_subscription_id;`,
		s.Radius,
		s.LatitudeDegrees,
		s.LongitudeDegrees,
		s.UserId,
	).Scan(&id)
	if err != nil {
		slog.Error(
			"ERR db.QueryRow(INSERT wca_competitions_announcements_position_subscriptions).",
			"error",
			err,
			"subscription",
			s,
		)
	}

	return id, err
}

func (s *WCACompAnnouncementsPositionSubscriptions) Update(
	db interfaces.DB,
) (pgconn.CommandTag, error) {
	tag, err := db.Exec(
		context.Background(),
		`UPDATE wca_competitions_announcements_position_subscriptions SET radius = $1 WHERE wca_competitions_announcements_position_subscription_id = $2 AND user_id = $3;`,
		s.Radius,
		s.Id,
		s.UserId,
	)
	if err != nil {
		slog.Error(
			"ERR db.Exec(UPDATE wca_competitions_announcements_position_subscriptions).",
			"error",
			err,
			"subscription",
			s,
		)
	}

	return tag, err
}

func (s *WCACompAnnouncementsPositionSubscriptions) Delete(
	db interfaces.DB,
) (pgconn.CommandTag, error) {
	tag, err := db.Exec(
		context.Background(),
		`DELETE FROM wca_competitions_announcements_position_subscriptions WHERE user_id = $1 AND wca_competitions_announcements_position_subscription_id = $2`,
		s.UserId,
		s.Id,
	)
	if err != nil {
		slog.Error(
			"ERR db.Exec(DELETE wca_competitions_announcements_position_subscriptions).",
			"error",
			err,
			"subscription",
			s,
		)
	}

	return tag, err
}
