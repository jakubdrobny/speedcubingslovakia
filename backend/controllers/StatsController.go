package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func GetAdminStats(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var adminStatsCollection models.AdminStatsCollection

		c.IndentedJSON(http.StatusOK, adminStatsCollection)
	}
}
