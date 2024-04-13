package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func AdminMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		isadmin := c.MustGet("isadmin").(bool)
		if !isadmin {
			c.IndentedJSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthMiddleWare(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authDetails, err := models.GetAuthDetailsFromHeader(c, envMap["JWT_SECRET_KEY"])
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, err)
			c.Abort()
			return
		}
		
		user, err := models.GetUserById(db, authDetails.UserId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.Set("uid", user.Id)
		c.Set("isadmin", user.IsAdmin)

		c.Next()
	}
}