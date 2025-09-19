package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

func GetSubscriptionStats(db interfaces.DB, getSubscriptionStats func(context.Context, interfaces.DB) (models.SubscriptionStats, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		var err error
		defer utils.PrintStack(&err)

		stats, err := getSubscriptionStats(ctx, db)
		if err != nil {
			err = fmt.Errorf("%w: when calling models.GetSubscriptionStats", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to query subscription stats."))
			return
		}

		responseJson, err := json.Marshal(stats)
		if err != nil {
			err = fmt.Errorf("%w: when marshalling response=%+v", err, stats)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to serialize response."))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(responseJson)
	}
}

func GetUserSubscriptionDetails(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer utils.PrintStack(&err)

		ctx := c.Request.Context()
		details, err := models.GetUserSubscriptionDetails(ctx, db)
		if err != nil {
			err = fmt.Errorf("%w: when calling models.GetUserSubscriptionDetails", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to query subscription details.")
			return
		}

		c.IndentedJSON(http.StatusOK, details)
	}
}
