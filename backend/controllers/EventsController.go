package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func GetEvents(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := models.GetAvailableEvents(db)
		if err != nil {
			log.Println("ERR GetAvailableEvents in GetEvents: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying events from database.")
		} else {
			c.IndentedJSON(http.StatusOK, events)
		}
	}
}
