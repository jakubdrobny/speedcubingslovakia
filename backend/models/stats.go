package models

import (
	"context"
	"fmt"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

type SubscriptionStats struct {
	PositionSubscriptions int `json:"position_subscriptions"`
	CountrySubscriptions  int `json:"country_subscriptions"`
	TotalUniqueUsers      int `json:"total_unique_users"`
}

type UserSubscriptionDetail struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	WCAID            string `json:"wca_id"`
	CountryName      string `json:"country_name"`
	CountryISO2      string `json:"country_iso2"`
	CountrySubCount  int    `json:"country_sub_count"`
	PositionSubCount int    `json:"position_sub_count"`
}

func GetSubscriptionStats(ctx context.Context, db interfaces.DB) (SubscriptionStats, error) {
	var stats SubscriptionStats
	query := `
		SELECT
			(SELECT COUNT(DISTINCT user_id) FROM wca_competitions_announcements_position_subscriptions) AS position_subscriptions,
			(SELECT COUNT(DISTINCT user_id) FROM wca_competitions_announcements_subscriptions) AS country_subscriptions,
			(SELECT COUNT(DISTINCT user_id) FROM (
				SELECT user_id FROM wca_competitions_announcements_position_subscriptions
				UNION
				SELECT user_id FROM wca_competitions_announcements_subscriptions
			) AS all_subscribers) AS total_unique_users;
	`
	err := db.QueryRow(ctx, query).Scan(
		&stats.PositionSubscriptions,
		&stats.CountrySubscriptions,
		&stats.TotalUniqueUsers,
	)
	if err != nil {
		return SubscriptionStats{}, fmt.Errorf("%w: when querying subscription stats", err)
	}

	return stats, nil
}

func GetUserSubscriptionDetails(ctx context.Context, db interfaces.DB) ([]UserSubscriptionDetail, error) {
	query := `
		SELECT
			u.user_id,
			u.name,
			u.wcaid,
			c.name AS country_name,
			c.iso2 AS country_iso2,
			COALESCE(country_subs.count, 0) AS country_sub_count,
			COALESCE(position_subs.count, 0) AS position_sub_count
		FROM
			users u
		JOIN
			countries c ON u.country_id = c.country_id
		LEFT JOIN (
			SELECT user_id, COUNT(*) as count
			FROM wca_competitions_announcements_subscriptions
			GROUP BY user_id
		) AS country_subs ON u.user_id = country_subs.user_id
		LEFT JOIN (
			SELECT user_id, COUNT(*) as count
			FROM wca_competitions_announcements_position_subscriptions
			GROUP BY user_id
		) AS position_subs ON u.user_id = position_subs.user_id
		WHERE COALESCE(country_subs.count, 0) > 0 OR COALESCE(position_subs.count, 0) > 0
		ORDER BY
			u.name;
	`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return []UserSubscriptionDetail{}, fmt.Errorf("%w: when querying user subscription details", err)
	}
	defer rows.Close()

	details := []UserSubscriptionDetail{}
	for rows.Next() {
		var detail UserSubscriptionDetail
		if err := rows.Scan(
			&detail.Id,
			&detail.Name,
			&detail.WCAID,
			&detail.CountryName,
			&detail.CountryISO2,
			&detail.CountrySubCount,
			&detail.PositionSubCount,
		); err != nil {
			return []UserSubscriptionDetail{}, fmt.Errorf("%w: when scanning user subscription detail row", err)
		}
		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		return []UserSubscriptionDetail{}, fmt.Errorf("%w: when scanning user subscription details rows", err)
	}

	return details, nil
}
