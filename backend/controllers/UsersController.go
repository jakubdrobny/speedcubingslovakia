package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

func GetManageUsers(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer utils.PrintStack(&err)

		ctx := c.Request.Context()

		manageUsers, err := models.ViewManageUsers(ctx, db)
		if err != nil {
			err = fmt.Errorf("%w: when calling models.ViewManageUsers", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to query users.")
			return
		}

		c.IndentedJSON(http.StatusOK, manageUsers)
	}
}

func ManageUserRole(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer utils.PrintStack(&err)

		ctx := c.Request.Context()

		var manageUser models.ManageUser
		if err := c.ShouldBindJSON(&manageUser); err != nil {
			err = fmt.Errorf("%w: when parsing request body", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse request body.")
			return
		}

		if err = manageUser.UpdateRole(ctx, db); err != nil {
			err = fmt.Errorf("%w: when updating user role", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to update user role.")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully update user role.")
	}
}

func PostLogIn(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()

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

		exists, err := user.Exists(ctx, db)
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
			err = user.Insert(ctx, db)

			go func() {
				if err := user.SendNewUserMailAsync(ctx, db, envMap); err != nil {
					utils.PrintStack(&err)
				}
			}()
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

func GetSearchUsers(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer utils.PrintStack(&err)

		query := c.Query("query")
		if len(query) < 2 {
			c.IndentedJSON(http.StatusOK, []models.ManageUser{})
			return
		}

		ctx := c.Request.Context()
		users, err := models.SearchUsers(ctx, db, query)
		if err != nil {
			err = fmt.Errorf("%w: when calling models.SearchUsers", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to search for users.")
			return
		}

		c.IndentedJSON(http.StatusOK, users)
	}
}

func GetUserMapData(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := os.ReadFile("CountriesGeo.json")
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

func FindDuplicateUser(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer utils.PrintStack(&err)

		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid user ID provided.")
			return
		}

		ctx := c.Request.Context()
		duplicate, found, err := models.FindFuzzyDuplicateUser(ctx, db, userID)
		if err != nil {
			err = fmt.Errorf("%w: when calling models.FindFuzzyDuplicateUser", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to query for duplicate user.")
			return
		}

		if !found {
			c.IndentedJSON(http.StatusNotFound, "Duplicate user not found.")
			return
		}

		c.IndentedJSON(http.StatusOK, duplicate)
	}
}

type MergeUsersRequest struct {
	OldUserID int `json:"old_user_id" binding:"required"`
	NewUserID int `json:"new_user_id" binding:"required"`
}

func MergeUsers(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer utils.PrintStack(&err)

		var req MergeUsersRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("%w: when parsing request body", err)
			c.IndentedJSON(http.StatusBadRequest, "Invalid request body.")
			return
		}

		ctx := c.Request.Context()
		if err := models.MergeUsers(ctx, db, req.OldUserID, req.NewUserID); err != nil {
			err = fmt.Errorf("%w: when merging users", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to merge users.")
			return
		}

		c.IndentedJSON(http.StatusOK, "Users merged successfully.")
	}
}
