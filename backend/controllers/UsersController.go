package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

func GetManageUsers(db interfaces.DbExecutor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer utils.PrintStack(err)

		ctx := c.Request.Context()

		manageUsers, err := models.ViewManageUsers(ctx, db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Failed to query users.")
			return
		}

		c.IndentedJSON(http.StatusOK, manageUsers)
	}
}

func PostLogIn(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed getting user info from database.",
			)
			return
		}

		if exists {
			err = user.Update(db)
		} else {
			err = user.Insert(db)
		}

		if err != nil {
			log.Println(
				"ERR (",
				exists,
				")user.Update or (",
				!exists,
				")user.Insert in PostLogIn: "+err.Error(),
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed updating/insert user data into database.",
			)
			return
		}

		authInfo.AvatarUrl = user.AvatarUrl
		authInfo.WcaId = user.WcaId
		if user.WcaId == "" {
			authInfo.WcaId = user.Name
		}
		authInfo.AccessToken, err = utils.CreateToken(
			user.Id,
			envMap["JWT_SECRET_KEY"],
			authInfo.ExpiresIn,
		)
		authInfo.IsAdmin = user.IsAdmin
		authInfo.Username = user.Name
		if err != nil {
			log.Println("ERR CreateToken in PostLogIn: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed creating token.")
			return
		}

		c.IndentedJSON(http.StatusOK, authInfo)
	}
}

func GetSearchUsers(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")

		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Println("ERR db.begin in GetSearchUsers: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to start transaction.")
			tx.Rollback(context.Background())
			return
		}

		searchUsers, statusCode, logMessage, returnMessage := models.GetUsersFromDB(tx, query)
		if statusCode == http.StatusInternalServerError {
			log.Println(logMessage)
			c.IndentedJSON(statusCode, returnMessage)
			return
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Println("ERR tx.commit in GetSearchUsers: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to finish transaction.")
			return
		}

		c.IndentedJSON(statusCode, searchUsers)
	}
}

func GetUserMapData(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := os.ReadFile("main/CountriesGeo.json")
		if err != nil {
			log.Println("ERR os.ReadFile in GetUserMapData: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to load map data.")
			return
		}

		var featureCollection models.FeatureCollection
		err = json.Unmarshal(buf, &featureCollection)
		if err != nil {
			log.Println("ERR json.Unmarshal in GetUserMapData: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse map data.")
			return
		}

		usersByCountry, logMsg, retMsg, err := models.GetUsersByCountryWithKinchScore(db)
		if err != nil {
			log.Println(logMsg)
			c.IndentedJSON(http.StatusInternalServerError, retMsg)
			return
		}

		for idx := range featureCollection.Features {
			countryIso2 := featureCollection.Features[idx].Properties.CountryIso2
			if _, ok := usersByCountry[countryIso2]; !ok {
				usersByCountry[countryIso2] = make([]models.MapDataUser, 0)
			}
			featureCollection.Features[idx].Properties.Users = usersByCountry[countryIso2]
		}

		c.IndentedJSON(http.StatusOK, featureCollection)
	}
}
