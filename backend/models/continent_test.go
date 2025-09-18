package models_test

import (
	"testing"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/stretchr/testify/require"
)

func TestContinent(t *testing.T) {
	ctx := t.Context()

	t.Run("get + insert", func(t *testing.T) {
		c := models.Continent{}
		err := c.Get(ctx, testDb, "invalid")
		require.Error(t, err)

		c2, err := models.TestInsertContinent(ctx, testDb)
		require.NoError(t, err)

		err = c.Get(ctx, testDb, c2.Name)
		require.NoError(t, err)
		require.Equal(t, c, c2)
	})
}
