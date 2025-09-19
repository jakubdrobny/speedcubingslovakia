package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jakubdrobny/speedcubingslovakia/backend/controllers"
	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSubscriptionStats(t *testing.T) {
	subscriptionStats := models.SubscriptionStats{
		PositionSubscriptions: 5,
		CountrySubscriptions:  10,
		TotalUniqueUsers:      12,
	}
	subscriptionStatsJson, err := json.Marshal(subscriptionStats)
	require.NoError(t, err)

	callbackOk := func(ctx context.Context, db interfaces.DB) (models.SubscriptionStats, error) {
		return subscriptionStats, nil
	}
	callbackFailed := func(ctx context.Context, db interfaces.DB) (models.SubscriptionStats, error) {
		return models.SubscriptionStats{}, errors.New("error")
	}

	tests := []struct {
		name                 string
		callback             func(ctx context.Context, db interfaces.DB) (models.SubscriptionStats, error)
		expectedResponseCode int
		expectedMsg          string
	}{
		{
			name:                 "successfull",
			callback:             callbackOk,
			expectedResponseCode: http.StatusOK,
			expectedMsg:          string(subscriptionStatsJson),
		},
		{
			name:                 "callback failed",
			callback:             callbackFailed,
			expectedResponseCode: http.StatusInternalServerError,
			expectedMsg:          "Failed to query subscription stats.",
		},
	}

	for _, testcase := range tests {
		req := httptest.NewRequest("GET", "/subscriptions", nil)
		rr := httptest.NewRecorder()

		handler := controllers.GetSubscriptionStats(nil, testcase.callback)
		handler(rr, req)

		assert.Equal(t, testcase.expectedResponseCode, rr.Code)
		assert.Contains(t, rr.Body.String(), testcase.expectedMsg)
	}
}
