package models_test

import (
	"testing"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/stretchr/testify/require"
)

func TestGetSubscriptionStats(t *testing.T) {
	ctx := t.Context()

	// subscribes to country only
	_, _, _, _, err := models.TestInsertWCACompAnnouncementsSubscription(ctx, testDb)
	require.NoError(t, err)

	// subscribes to position only
	_, _, _, _, err = models.TestInsertWCACompAnnouncementsPositionSubscription(ctx, testDb)
	require.NoError(t, err)

	// subscribes to both country and position
	user3, _, _, err := models.TestInsertUser(ctx, testDb)
	require.NoError(t, err)

	sub3_1 := models.NewTestWCACompAnnouncementsSubscription(user3.Id, user3.CountryId)
	err = sub3_1.Insert(ctx, testDb)
	require.NoError(t, err)

	sub3_2 := models.NewTestWCACompAnnouncementsPositionSubscription(user3.Id)
	err = sub3_2.Insert(ctx, testDb)
	require.NoError(t, err)

	// no subscriptions
	_, _, _, err = models.TestInsertUser(ctx, testDb)
	require.NoError(t, err)

	stats, err := models.GetSubscriptionStats(ctx, testDb)
	require.NoError(t, err)
	require.GreaterOrEqual(t, 2, stats.CountrySubscriptions, "should be at least 2 unique users with country subscriptions")
	require.GreaterOrEqual(t, 2, stats.PositionSubscriptions, "should be at least 2 unique users with position subscriptions")
	require.GreaterOrEqual(t, 3, stats.TotalUniqueUsers, "should be at least 3 total unique users with any subscription")
}

func TestGetUserSubscriptionDetails(t *testing.T) {
	ctx := t.Context()

	_, user1, _, _, err := models.TestInsertWCACompAnnouncementsSubscription(ctx, testDb)
	require.NoError(t, err)

	_, user2, _, _, err := models.TestInsertWCACompAnnouncementsPositionSubscription(ctx, testDb)
	require.NoError(t, err)

	sub2_2 := models.NewTestWCACompAnnouncementsPositionSubscription(user2.Id)
	err = sub2_2.Insert(ctx, testDb)
	require.NoError(t, err)

	user3, _, _, err := models.TestInsertUser(ctx, testDb)
	require.NoError(t, err)

	details, err := models.GetUserSubscriptionDetails(ctx, testDb)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(details), 2, "should only return users with subscriptions")

	detailsMap := make(map[int]models.UserSubscriptionDetail)
	for _, d := range details {
		detailsMap[d.Id] = d
	}

	countryUserDetail, ok := detailsMap[user1.Id]
	require.True(t, ok)
	require.Equal(t, 1, countryUserDetail.CountrySubCount)
	require.Equal(t, 0, countryUserDetail.PositionSubCount)

	positionUserDetail, ok := detailsMap[user2.Id]
	require.True(t, ok)
	require.Equal(t, 0, positionUserDetail.CountrySubCount)
	require.Equal(t, 2, positionUserDetail.PositionSubCount)

	_, ok = detailsMap[user3.Id]
	require.False(t, ok)
}
