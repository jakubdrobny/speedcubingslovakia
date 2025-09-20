package models_test

import (
	"math/rand/v2"
	"testing"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/stretchr/testify/require"
)

func TestWCACompAnnouncementPositionSubscription(t *testing.T) {
	ctx := t.Context()

	t.Run("get + insert", func(t *testing.T) {
		sub := models.WCACompAnnouncementsPositionSubscription{}
		err := sub.Get(ctx, testDb, -1)
		require.Error(t, err)

		sub2, _, _, _, err := models.TestInsertWCACompAnnouncementsPositionSubscription(ctx, testDb)
		require.NoError(t, err)

		err = sub.Get(ctx, testDb, sub2.Id)
		require.NoError(t, err)
		require.Equal(t, sub, sub2)
	})

	t.Run("updateRadius", func(t *testing.T) {
		sub, _, _, _, err := models.TestInsertWCACompAnnouncementsPositionSubscription(ctx, testDb)
		require.NoError(t, err)

		sub2 := sub
		for sub2.Radius == sub.Radius {
			sub2.Radius = rand.IntN(20000)
		}
		err = sub2.UpdateRadius(ctx, testDb)
		require.NoError(t, err)

		var sub3 models.WCACompAnnouncementsPositionSubscription
		err = sub3.Get(ctx, testDb, sub2.Id)
		require.NoError(t, err)
		require.Equal(t, sub2.Radius, sub3.Radius)
		require.NotEqual(t, sub.Radius, sub3.Radius)
	})

	t.Run("exists", func(t *testing.T) {
		var sub models.WCACompAnnouncementsPositionSubscription
		ok, err := sub.Exists(ctx, testDb)
		require.NoError(t, err)
		require.False(t, ok)

		sub2, _, _, _, err := models.TestInsertWCACompAnnouncementsPositionSubscription(ctx, testDb)
		require.NoError(t, err)

		ok, err = sub2.Exists(ctx, testDb)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("delete", func(t *testing.T) {
		sub, _, _, _, err := models.TestInsertWCACompAnnouncementsPositionSubscription(ctx, testDb)
		require.NoError(t, err)

		err = sub.Delete(ctx, testDb)
		require.NoError(t, err)

		ok, err := sub.Exists(ctx, testDb)
		require.NoError(t, err)
		require.False(t, ok)
	})
}
