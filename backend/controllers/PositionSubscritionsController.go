package controllers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

// if userId = 0, returns subscriptions for all users
func PositionSubscriptionFromDB(
	db interfaces.DB,
	userId int,
) ([]models.WCACompAnnouncementsPositionSubscriptions, error) {
	queryString := `SELECT wca_competitions_announcements_position_subscription_id as id, latitude_degrees, longitude_degrees, radius, user_id FROM wca_competitions_announcements_position_subscriptions`
	args := []any{}
	if userId != 0 {
		queryString += " WHERE user_id = $1"
		args = append(args, userId)
	}
	queryString += ";"

	rows, err := db.Query(
		context.Background(),
		queryString,
		args...,
	)
	if err != nil {
		slog.Error(
			"ERR db.Query(wca_competitions_announcements_position_subscriptions) in PositionSubscriptionFromDB.",
			"error",
			err,
			"user_id",
			userId,
		)
		return []models.WCACompAnnouncementsPositionSubscriptions{}, err
	}

	subscriptions := make([]models.WCACompAnnouncementsPositionSubscriptions, 0)
	for rows.Next() {
		sub := models.WCACompAnnouncementsPositionSubscriptions{
			New:  false,
			Open: false,
		}
		err = rows.Scan(
			&sub.Id,
			&sub.LatitudeDegrees,
			&sub.LongitudeDegrees,
			&sub.Radius,
			&sub.UserId,
		)
		if err != nil {
			slog.Error(
				"ERR rows.Scan(position_subscription=id,lat_deg,long_deg,radius) in GetWCACompAnnouncementsPositionSubscriptions",
				"error",
				err,
			)
			return []models.WCACompAnnouncementsPositionSubscriptions{}, err
		}

		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func GetWCACompAnnouncementsPositionSubscriptions(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")

		slog.Info("Querying position subscription data...", "user_id", uid)

		subscriptions, err := PositionSubscriptionFromDB(db, uid)
		if err != nil {
			slog.Error(
				"ERR PositionSubscriptionFromDB in GetWCACompAnnouncementsPositionSubscriptions.",
				"error",
				err,
				"user_id",
				uid,
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to query position subscriptions from db.",
			)
			return
		}

		slog.Info("Successfully queried position subscription data.", "user_id", uid)
		c.IndentedJSON(http.StatusOK, subscriptions)
	}
}

func UpdateWCAAnnouncementsPositionSubscriptions(db interfaces.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")

		slog.Info("Updating position subscription...", "user_id", uid)
		var subscription models.WCACompAnnouncementsPositionSubscriptions

		if err := c.ShouldBindJSON(&subscription); err != nil {
			slog.Error(
				"Failed to parse body in UpdateWCAAnnouncementsPositionSubscriptions.",
				"error",
				err,
				"user_id",
				uid,
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to parse position subscription data.",
			)
			return
		}

		if subscription.HasOutOfRangeCoords() {
			slog.Error(
				"Someone is trying to put a marker outside the earth :DD",
				"subscription",
				subscription,
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Invalid marker position. Please be in the central earth on the map :D",
			)
			return
		}

		subscription.UserId = uid

		tx, err := db.Begin(context.Background())
		if err != nil {
			slog.Error(
				"ERR db.Begin in UpdateWCAAnnouncementsPositionSubscriptions.",
				"error",
				err,
				"subscription",
				subscription,
			)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to start db transaction.")
			return
		}
		defer tx.Rollback(context.Background())

		exists, err := subscription.Exists(tx)
		if err != nil {
			slog.Error(
				"ERR subscription.Exists in UpdateWCAAnnouncementsPositionSubscriptions.",
				"error",
				err,
				"subscription",
				subscription,
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to check if subscription already exists.",
			)
			return
		}

		if exists {
			err := subscription.Update(tx)
			if err != nil {
				slog.Error(
					"ERR subscription.Update in UpdateWCAAnnouncementsPositionSubscriptions.",
					"error",
					err,
					"subscription",
					subscription,
				)
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed to update position subscription into db.",
				)
				return
			}
		} else {
			// no update happened => we need to insert
			subscriptionId, err := subscription.Insert(tx)
			if err != nil {
				slog.Error(
					"ERR subscription.Insert in UpdateWCAAnnouncementsPositionSubscriptions.",
					"error",
					err,
					"subscription",
					subscription,
				)
				c.IndentedJSON(
					http.StatusInternalServerError,
					"Failed to insert position subscription into db.",
				)
				return
			}

			subscription.Id = subscriptionId
		}

		subscription.New = false
		subscription.Open = false

		err = tx.Commit(context.Background())
		if err != nil {
			slog.Error(
				"ERR tx.commit in UpdateWCAAnnouncementsPositionSubscriptions.",
				"error",
				err,
				"subscription",
				subscription,
			)
			return
		}

		slog.Info("Updated position subscription successfully.", "subscription", subscription)
		c.IndentedJSON(http.StatusOK, subscription)
	}
}

func DeleteWCAAnnouncementsPositionSubscriptions(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")

		slog.Info("Deleting position subscription...", "user_id", uid)
		var subscription models.WCACompAnnouncementsPositionSubscriptions

		if err := c.ShouldBindJSON(&subscription); err != nil {
			slog.Error(
				"ERR c.ShouldBindJson(&subscription) in DeleteWCAAnnouncementsPositionSubscriptions.",
				"error",
				err,
				"subscription",
				subscription,
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to parse position subscription data.",
			)
			return
		}

		subscription.UserId = uid

		_, err := subscription.Delete(db)
		if err != nil {
			slog.Error(
				"ERR db.Exec(delete position subscription) in DeleteWCAAnnouncementsPositionSubscriptions.",
				"error",
				err,
				"subscription",
				subscription,
			)
			c.IndentedJSON(
				http.StatusInternalServerError,
				"Failed to delete position subscription from db.",
			)
			return
		}

		slog.Info("Deleted position subscription successfully.", "subscription", subscription)
		c.IndentedJSON(http.StatusOK, "Subscriptions inserted/updated correctly.")
	}
}
