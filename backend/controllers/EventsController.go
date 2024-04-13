package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func GetEvents(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		events, err := models.GetAvailableEvents(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			c.IndentedJSON(http.StatusOK, events);
		}
	}
}