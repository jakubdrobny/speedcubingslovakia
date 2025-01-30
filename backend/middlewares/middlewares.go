package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
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
	db interfaces.DB,
	envMap map[string]string,
) {
	c.Set("authorized", false)

	authDetails, err := models.GetAuthDetailsFromHeader(c, envMap["JWT_SECRET_KEY"])
	if err != nil {
		// we can only get here, if the status code should be unauthorized
		c.Set("unauthorization_reason", err)
		return
	}

	user, err := models.GetUserById(db, authDetails.UserId)
	if err != nil {
		c.Set("user_id_error", err)
		return
	}

	c.Set("authorized", true)
	c.Set("uid", user.Id)
	c.Set("isadmin", user.IsAdmin)
}

func Authorization(db interfaces.DB, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		MarkAuthorization(c, db, envMap)

		c.Next()
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authorized := c.GetBool("authorized"); !authorized {
			unauthorizationReason, unauthorizationReasonExists := c.Get("unauthorization_reason")
			if unauthorizationReasonExists {
				slog.LogAttrs(
					context.Background(),
					slog.LevelWarn,
					"ERR models.GetAuthDetailsFromHeader in AuthMiddleWare in MarkAuthorization",
					slog.Any("unauthorization_reason", unauthorizationReason),
				)
				c.IndentedJSON(http.StatusUnauthorized, unauthorizationReason)
				c.Abort()
				return
			}

			userIdError, userIdErrorExists := c.Get("user_id_error")
			if userIdErrorExists {
				slog.LogAttrs(
					context.Background(),
					slog.LevelWarn,
					"ERR models.GetUserById in MarkAuthorization in AuthMiddleware",
					slog.Any("error",
						userIdError),
				)
				c.IndentedJSON(http.StatusInternalServerError, userIdError)
				c.Abort()
				return
			}

			slog.LogAttrs(context.Background(), slog.LevelError, "Should never get here?")
			c.IndentedJSON(http.StatusUnauthorized, "Unauthorized.")
			c.Abort()
			return
		}

		c.Next()
	}
}
