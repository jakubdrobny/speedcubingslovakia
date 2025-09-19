package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

func GetSubscriptionStats(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer utils.PrintStack(&err)

		ctx := c.Request.Context()
		stats, err := models.GetSubscriptionStats(ctx, db)
		if err != nil {
			err = fmt.Errorf("%w: when calling models.GetSubscriptionStats", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to query subscription stats.")
			return
		}

		c.IndentedJSON(http.StatusOK, stats)
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
