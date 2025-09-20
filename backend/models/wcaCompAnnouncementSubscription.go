package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

type WCACompAnnouncementsSubscription struct {
	Id        int
	UserId    int
	CountryId string
	State     string
}

func (s *WCACompAnnouncementsSubscription) Get(ctx context.Context, db interfaces.DB, id int) error {
	err := db.QueryRow(ctx, `
		SELECT s.wca_competitions_announcements_subscription_id, s.user_id, s.country_id, s.state
		FROM wca_competitions_announcements_subscriptions s
		WHERE s.wca_competitions_announcements_subscription_id = $1
	`, id).Scan(&s.Id, &s.UserId, &s.CountryId, &s.State)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("wca comp announcements subscription with id=%d not found", id)
		}
		return fmt.Errorf("%w: when querying wca comp announcements subscription with id=%d", err, id)
	}

	return nil
}

func (s *WCACompAnnouncementsSubscription) Insert(ctx context.Context, db interfaces.DB) error {
	err := db.QueryRow(ctx, `
		INSERT INTO wca_competitions_announcements_subscriptions (user_id, country_id, state)
		VALUES ($1, $2, $3) RETURNING wca_competitions_announcements_subscription_id
	`, s.UserId, s.CountryId, s.State).Scan(&s.Id)
	if err != nil {
		return fmt.Errorf("%w: when inserting wca comp announcements subscription=%+v", err, s)
	}

	return nil
}
