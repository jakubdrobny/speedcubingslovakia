package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

type WCACompAnnouncementsPositionSubscription struct {
	Id               int     `json:"id"`
	UserId           int     `json:"-"`
	LatitudeDegrees  float64 `json:"lat"`
	LongitudeDegrees float64 `json:"long"`
	Radius           int     `json:"radius"`
	New              bool    `json:"new"`
	Open             bool    `json:"open"`
}

func (s *WCACompAnnouncementsPositionSubscription) Get(ctx context.Context, db interfaces.DB, id int) error {
	err := db.QueryRow(ctx,
		`SELECT ps.wca_competitions_announcements_position_subscription_id, ps.user_id, ps.latitude_degrees, ps.longitude_degrees, ps.radius
		FROM wca_competitions_announcements_position_subscriptions ps
		WHERE wca_competitions_announcements_position_subscription_id = $1
		`, id).Scan(&s.Id, &s.UserId, &s.LatitudeDegrees, &s.LongitudeDegrees, &s.Radius)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("wca comp announcements position subscription with id=%d not found", id)
		}
		return fmt.Errorf("%w: when querying wca comp announcements position subscription with id=%d", err, id)
	}

	return nil
}

func (s *WCACompAnnouncementsPositionSubscription) Insert(
	ctx context.Context, db interfaces.DB,
) error {
	err := db.QueryRow(
		ctx,
		`INSERT INTO wca_competitions_announcements_position_subscriptions (radius, latitude_degrees, longitude_degrees, user_id)
			VALUES ($1,$2,$3,$4)
			RETURNING wca_competitions_announcements_position_subscription_id;
		`,
		s.Radius,
		s.LatitudeDegrees,
		s.LongitudeDegrees,
		s.UserId,
	).Scan(&s.Id)
	if err != nil {
		return fmt.Errorf("%w: when insert wca comp announcements position subscription=%+v", err, *s)
	}

	return nil
}

func (s WCACompAnnouncementsPositionSubscription) UpdateRadius(
	ctx context.Context, db interfaces.DB,
) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE wca_competitions_announcements_position_subscriptions
			SET radius = $1
			WHERE wca_competitions_announcements_position_subscription_id = $2 AND user_id = $3;
		`,
		s.Radius,
		s.Id,
		s.UserId,
	)
	if err != nil {
		return fmt.Errorf("%w: when updating wca comp announcements position subscription=%+v", err, s)
	}

	return nil
}

func (s WCACompAnnouncementsPositionSubscription) Exists(ctx context.Context, db interfaces.DB) (bool, error) {
	var exists bool
	err := db.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1 
			FROM wca_competitions_announcements_position_subscriptions 
			WHERE latitude_degrees = $1 AND longitude_degrees = $2 AND radius = $3 AND user_id = $4
		);`,
		s.LatitudeDegrees,
		s.LongitudeDegrees,
		s.Radius,
		s.UserId,
	).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%w: when checking existance of wca comp announcements position subscription=%+v", err, s)
	}

	return exists, nil
}

func (s WCACompAnnouncementsPositionSubscription) Delete(
	ctx context.Context, db interfaces.DB,
) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM wca_competitions_announcements_position_subscriptions
			WHERE user_id = $1 AND wca_competitions_announcements_position_subscription_id = $2`,
		s.UserId,
		s.Id,
	)
	if err != nil {
		return fmt.Errorf("%w: when deleting wca comp announcements position subscription=%+v", err, s)
	}

	return nil
}

func (s WCACompAnnouncementsPositionSubscription) HasOutOfRangeCoords() bool {
	return s.LatitudeDegrees+180 < -constants.EPS || s.LatitudeDegrees-180 > constants.EPS ||
		s.LongitudeDegrees+180 < -constants.EPS ||
		s.LongitudeDegrees-180 > constants.EPS
}
