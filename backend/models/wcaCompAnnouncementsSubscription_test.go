package models_test

import (
	"testing"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/stretchr/testify/require"
)

func TestWCACompAnnouncementsSubscription(t *testing.T) {
	ctx := t.Context()

	t.Run("get + insert", func(t *testing.T) {
		sub := models.WCACompAnnouncementsSubscription{}
		err := sub.Get(ctx, testDb, -1)
		require.Error(t, err)

		sub2, _, _, _, err := models.TestInsertWCACompAnnouncementsSubscription(ctx, testDb)
		require.NoError(t, err)

		err = sub.Get(ctx, testDb, sub2.Id)
		require.NoError(t, err)
		require.Equal(t, sub2, sub)
	})
}
