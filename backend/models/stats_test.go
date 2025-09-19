package models_test

// import (
// 	"context"
// 	"testing"
//
// 	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
// 	"github.com/stretchr/testify/require"
// )
//
// func TestGetSubscriptionStats(t *testing.T) {
// 	ctx := t.Context()
// 	country, _, err := models.TestInsertCountry(ctx, testDb)
// 	require.NoError(t, err)
//
// 	user := models.NewTestUser(country.Id)
// 	err = user.Insert(ctx, testDb)
// 	require.NoError(t, err)
//
// 	_, err := testDb.Exec(context.Background(),
// 		"INSERT INTO wca_competitions_announcements_subscriptions (country_id, user_id) VALUES ($1, $2)",
// 		countryID, userCountryOnly.Id,
// 	)
// 	require.NoError(t, err)
//
// 	// User 2: Only subscribes to a position.
// 	userPositionOnly := testInsertUser(t, newTestUser("PositionSubUser"))
// 	_, err = testDb.Exec(context.Background(),
// 		"INSERT INTO wca_competitions_announcements_position_subscriptions (user_id) VALUES ($1)",
// 		userPositionOnly.Id,
// 	)
// 	require.NoError(t, err)
//
// 	// User 3: Subscribes to both a country and a position to test distinct counts.
// 	userBoth := testInsertUser(t, newTestUser("BothSubUser"))
// 	_, err = testDb.Exec(context.Background(),
// 		"INSERT INTO wca_competitions_announcements_subscriptions (country_id, user_id) VALUES ($1, $2)",
// 		countryID, userBoth.Id,
// 	)
// 	require.NoError(t, err)
// 	_, err = testDb.Exec(context.Background(),
// 		"INSERT INTO wca_competitions_announcements_position_subscriptions (user_id) VALUES ($1)",
// 		userBoth.Id,
// 	)
// 	require.NoError(t, err)
//
// 	// User 4: Has no subscriptions and should not be counted.
//
// 	// ACT: Call the function we want to test.
// 	stats, err := models.GetSubscriptionStats(context.Background(), testDb)
// 	require.NoError(t, err, "GetSubscriptionStats should not return an error")
// 	require.Equal(t, 2, stats.CountrySubscriptions, "should be 2 unique users with country subscriptions")
// 	require.Equal(t, 2, stats.PositionSubscriptions, "should be 2 unique users with position subscriptions")
// 	require.Equal(t, 3, stats.TotalUniqueUsers, "should be 3 total unique users with any subscription")
// }
//
// // TestGetUserSubscriptionDetails verifies the detailed subscription list.
// func TestGetUserSubscriptionDetails(t *testing.T) {
// 	// ARRANGE: Use a similar setup to the aggregate test.
// 	countryID, countryName, countryISO2 := newTestCountry()
// 	testInsertCountry(t, countryID, countryName, countryISO2)
//
// 	userCountryOnly := testInsertUser(t, newTestUser("DetailCountryUser"))
// 	_, err := testDb.Exec(context.Background(),
// 		"INSERT INTO wca_competitions_announcements_subscriptions (country_id, user_id) VALUES ($1, $2)",
// 		countryID, userCountryOnly.Id,
// 	)
// 	require.NoError(t, err)
//
// 	userPositionOnly := testInsertUser(t, newTestUser("DetailPositionUser"))
// 	_, err = testDb.Exec(context.Background(),
// 		"INSERT INTO wca_competitions_announcements_position_subscriptions (user_id) VALUES ($1)",
// 		userPositionOnly.Id,
// 	)
// 	require.NoError(t, err)
// 	// Add a second position subscription for the same user to test COUNT(*)
// 	_, err = testDb.Exec(context.Background(),
// 		"INSERT INTO wca_competitions_announcements_position_subscriptions (user_id, radius) VALUES ($1, 100)", // Different radius to avoid unique constraint
// 		userPositionOnly.Id,
// 	)
// 	require.NoError(t, err)
//
// 	userNoSubs := testInsertUser(t, newTestUser("DetailNoSubUser"))
//
// 	// ACT: Call the function we want to test.
// 	details, err := models.GetUserSubscriptionDetails(context.Background(), testDb)
//
// 	// ASSERT: Check the length and content of the returned slice.
// 	require.NoError(t, err, "GetUserSubscriptionDetails should not return an error")
//
// 	// The function should only return users who have subscriptions.
// 	// We expect to see userCountryOnly and userPositionOnly. userNoSubs should be excluded.
// 	require.Len(t, details, 2, "Should only return users with subscriptions")
//
// 	// For easier assertions, convert the slice to a map keyed by user ID.
// 	detailsMap := make(map[int]models.UserSubscriptionDetail)
// 	for _, d := range details {
// 		detailsMap[d.ID] = d
// 	}
//
// 	// Verify details for the user with only country subscriptions.
// 	countryUserDetail, ok := detailsMap[userCountryOnly.Id]
// 	require.True(t, ok, "User with country subscription should be in the results")
// 	require.Equal(t, 1, countryUserDetail.CountrySubCount, "Country sub count should be 1")
// 	require.Equal(t, 0, countryUserDetail.PositionSubCount, "Position sub count should be 0")
//
// 	// Verify details for the user with only position subscriptions.
// 	positionUserDetail, ok := detailsMap[userPositionOnly.Id]
// 	require.True(t, ok, "User with position subscription should be in the results")
// 	require.Equal(t, 0, positionUserDetail.CountrySubCount, "Country sub count should be 0")
// 	require.Equal(t, 2, positionUserDetail.PositionSubCount, "Position sub count should be 2")
//
// 	// Verify the user with no subscriptions is not in the map.
// 	_, ok = detailsMap[userNoSubs.Id]
// 	require.False(t, ok, "User with no subscriptions should not be in the results")
// }
