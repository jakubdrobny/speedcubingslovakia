package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func GetResultsStatuses(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		statuses, err := models.GetAvailableResultsStatuses(db)
		if err != nil {
			log.Println("ERR GetAvailableResultsStatuses in GetResultsStatuses: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying statuses from database.")
		} else {
			c.IndentedJSON(http.StatusOK, statuses)
		}
	}
}
