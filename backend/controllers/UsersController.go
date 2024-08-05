package controllers

import (
	"context"
	"io"
	"log"
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
			log.Println("ERR db.Query SELECT from users in GetManageRolesUsers: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to query users from database.")
			return
		}

		uid := c.MustGet("uid").(int)

		for rows.Next() {
			var manageRolesUser models.ManageRolesUser
			err = rows.Scan(&manageRolesUser.Id, &manageRolesUser.Name, &manageRolesUser.Isadmin)
			if err != nil {
				log.Println("ERR scanning ManageRolesUser in GetManageRolesUsers: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to query users from database.")
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
			log.Println("ERR BindJSON(&manageRolesUsers) in PutManageRolesUsers: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse incoming data.")
			return;
		}

		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Println("ERR db.Begin in GetManageRolesUsers: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to start transaction.")
			tx.Rollback(context.Background())
			return;
		}
		
		for _, manageRolesUser := range manageRolesUsers {
			_, err = tx.Exec(context.Background(), `UPDATE users SET isadmin = $1 WHERE user_id = $2;`, manageRolesUser.Isadmin, manageRolesUser.Id)
			if err != nil {
				log.Println("ERR tx.Exec UPDATE users in GetManageRolesUsers: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to update user roles in database.")
				tx.Rollback(context.Background())
				return
			}
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Println("ERR tx.Commit in GetManageRolesUsers: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to commit transaction.")
			return
		}

		c.IndentedJSON(http.StatusCreated, manageRolesUsers)
	}
}


func PostLogIn(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func (c *gin.Context) {
		reqBodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("ERR io.ReadAll in PostLogIn: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse incoming data.")
			return
		}

		code := string(reqBodyBytes)
		authInfo, err := models.GetAuthInfo(code, envMap)
		if err != nil {
			log.Println("ERR GetAuthInfo in PostLogIn: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed getting auth info for WCA.")
			return
		}

		user, err := models.GetUserInfoFromWCA(&authInfo, envMap)
		if err != nil {
			log.Println("ERR GetUserInfoFromWCA in PostLogIn: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed getting user info from WCA.")
			return
		}

		exists, err := user.Exists(db)
		if err != nil {
			log.Println("ERR user.Exists in PostLogIn: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed getting user info from database.")
			return
		}

		if exists {
			err = user.Update(db)
		} else {
			err = user.Insert(db)
		}

		if err != nil {
			log.Println("ERR (", exists, ")user.Update or (", !exists, ")user.Insert in PostLogIn: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed updating/insert user data into database.")
			return
		}

		authInfo.AvatarUrl = user.AvatarUrl
		authInfo.WcaId = user.WcaId
		if user.WcaId == "" { authInfo.WcaId = user.Name }
		authInfo.AccessToken, err = utils.CreateToken(user.Id, envMap["JWT_SECRET_KEY"], authInfo.ExpiresIn)
		authInfo.IsAdmin = user.IsAdmin
		if err != nil {
			log.Println("ERR CreateToken in PostLogIn: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed creating token.")
			return
		}
		
		c.IndentedJSON(http.StatusOK, authInfo)
	}
}

func GetSearchUsers(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		query := c.Query("query")

		searchUsers := make([]models.SearchUser, 0)
		if query == "_" {
			rows, err := db.Query(context.Background(), `SELECT u.name, (CASE WHEN u.wcaid LIKE '' THEN u.name ELSE u.wcaid END) AS wcaid FROM users u ORDER BY u.name;`)
			if err != nil {
				log.Println("ERR db.Query all in GetSearchUsers (" + query + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying all users.")
				return
			}

			for rows.Next() {
				var searchUser models.SearchUser
				err = rows.Scan(&searchUser.Username, &searchUser.WCAID)
				if err != nil {
					log.Println("ERR scanning all in GetSearchUsers with query (" + query + "): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying all users.")
					return
				}

				searchUsers = append(searchUsers, searchUser)
			}
		} else {
			rows, err := db.Query(context.Background(), `SELECT u.name, (CASE WHEN u.wcaid LIKE '' THEN u.name ELSE u.wcaid END) AS wcaid FROM users u WHERE u.wcaid LIKE $1 ORDER BY u.name;`, query)
			if err != nil {
				log.Println("ERR db.Query wcaid in GetSearchUsers (" + query + "): " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed querying users by WCAID.")
				return
			}

			for rows.Next() {
				var searchUser models.SearchUser
				err = rows.Scan(&searchUser.Username, &searchUser.WCAID)
				if err != nil {
					log.Println("ERR scanning wcaid in GetSearchUsers (" + query + "): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying users by WCAID.")
					return
				}

				searchUsers = append(searchUsers, searchUser)
			}

			if len(searchUsers) == 0 {
				rows, err := db.Query(context.Background(), `SELECT u.name, (CASE WHEN u.wcaid LIKE '' THEN u.name ELSE u.wcaid END) AS wcaid FROM users u WHERE LOWER(u.name) LIKE LOWER('%' || $1 || '%') ORDER BY u.name;`, query)
				if err != nil {
					log.Println("ERR db.Query name in GetSearchUsers (" + query + "): " + err.Error())
					c.IndentedJSON(http.StatusInternalServerError, "Failed querying users by name.")
					return
				}

				for rows.Next() {
					var searchUser models.SearchUser
					err = rows.Scan(&searchUser.Username, &searchUser.WCAID)
					if err != nil {
						log.Println("ERR scanning name in GetSearchUsers (" + query + "): " + err.Error())
						c.IndentedJSON(http.StatusInternalServerError, "Failed querying users by name.")
						return
					}

					searchUsers = append(searchUsers, searchUser)
				}
			}
		}

		c.IndentedJSON(http.StatusOK, searchUsers)
	}
}