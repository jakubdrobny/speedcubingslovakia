package controllers

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

func GetManageRolesUsers(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		manageRolesUsers := make([]models.ManageRolesUser, 0)

		rows, err := db.Query(context.Background(), `SELECT u.user_id, u.name, u.isadmin FROM users u;`)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		uid := c.MustGet("uid").(int)

		for rows.Next() {
			var manageRolesUser models.ManageRolesUser
			err = rows.Scan(&manageRolesUser.Id, &manageRolesUser.Name, &manageRolesUser.Isadmin)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}
			if uid != manageRolesUser.Id {
				manageRolesUsers = append(manageRolesUsers, manageRolesUser)
			}
		}

		c.IndentedJSON(http.StatusOK, manageRolesUsers)
	}
}

func PutManageRolesUsers(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		var manageRolesUsers []models.ManageRolesUser

		if err := c.BindJSON(&manageRolesUsers); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}

		tx, err := db.Begin(context.Background())
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			tx.Rollback(context.Background())
			return;
		}
		
		for _, manageRolesUser := range manageRolesUsers {
			_, err = tx.Exec(context.Background(), `UPDATE users SET isadmin = $1 WHERE user_id = $2;`, manageRolesUser.Isadmin, manageRolesUser.Id)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err);
				tx.Rollback(context.Background())
				return
			}
		}

		tx.Commit(context.Background())

		c.IndentedJSON(http.StatusCreated, manageRolesUsers)
	}
}


func PostLogIn(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func (c *gin.Context) {
		reqBodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		code := string(reqBodyBytes)
		authInfo, err := models.GetAuthInfo(code, envMap)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		user, err := models.GetUserInfoFromWCA(&authInfo, envMap)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		exists, err := user.Exists(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		if exists {
			err = user.Update(db)
		} else {
			err = user.Insert(db)
		}

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		authInfo.AvatarUrl = user.AvatarUrl
		authInfo.WcaId = user.WcaId
		authInfo.AccessToken, err = utils.CreateToken(user.Id, envMap["JWT_SECRET_KEY"])
		authInfo.IsAdmin = user.IsAdmin
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		
		c.IndentedJSON(http.StatusOK, authInfo)
	}
}