package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func GetTags(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		statuses, err := models.GetAvailableTags(db)
		if err != nil {
			log.Println("ERR GetAvailableTags in GetTags: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying tags from database.")
		} else {
			c.IndentedJSON(http.StatusOK, statuses)
		}
	}
}
