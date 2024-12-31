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

func MarkAuthorization(
	c *gin.Context,
	db *pgxpool.Pool,
	envMap map[string]string,
	markHeaders bool,
) bool {
	authDetails, err := models.GetAuthDetailsFromHeader(c, envMap["JWT_SECRET_KEY"])
	if err != nil {
		if markHeaders {
			c.IndentedJSON(http.StatusUnauthorized, err)
			c.Abort()
		}
		return false
	}

	user, err := models.GetUserById(db, authDetails.UserId)
	if err != nil {
		if markHeaders {
			c.IndentedJSON(http.StatusInternalServerError, err)
			c.Abort()
		}
		return false
	}

	c.Set("uid", user.Id)
	c.Set("isadmin", user.IsAdmin)

	return true
}

func AuthMiddleWare(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		MarkAuthorization(c, db, envMap, true)

		c.Next()
	}
}
