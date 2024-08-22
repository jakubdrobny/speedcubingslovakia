package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/middlewares"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func GetAnnouncementById(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")
		uid, exists := c.Get("uid")
		if exists { uid = uid.(int) }

		var rows pgx.Rows
		var err error
		if !exists {
			rows, err = db.Query(context.Background(), `SELECT a.announcement_id, a.title, a.content, u.wcaid, u.name FROM announcements a JOIN users u ON u.user_id = a.author_id WHERE a.announcement_id = $1;`, id)
		} else {
			rows, err = db.Query(context.Background(), `SELECT a.announcement_id, a.title, a.content, u.wcaid, u.name, ar.read FROM announcements a JOIN users u ON u.user_id = a.author_id JOIN announcement_read ar ON ar.announcement_id = a.announcement_id WHERE a.announcement_id = $1 AND ar.user_id = $2;`, id, uid)
		}
		if err != nil {
			log.Println("ERR db.Query in GetAnnouncementById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying announcement by id.")
			return;
		}

		var announcement models.AnnouncementState
		found := false

		for rows.Next() {
			err := rows.Scan(&announcement.Id, &announcement.Title, &announcement.Content, &announcement.AuthorWcaId, &announcement.AuthorUsername, &announcement.Read)
			if err != nil {
				log.Println("ERR scanning announcement data in GetAnnouncementById: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed parsing announcement from database.")
				return;
			}

			if !exists { announcement.Read = true }
			found = true
		}

		if !found {
			log.Println("ERR announcement with id: ", id, " not found in GetAnnouncementById.")
			c.IndentedJSON(http.StatusInternalServerError, "Announcement not found.")	
			return;
		}

		err = announcement.GetTags(db)
		if err != nil {
			log.Println("ERR GetTags in GetAnnouncementById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get announcement tags.")
			return;
		}

		c.IndentedJSON(http.StatusOK, announcement)
	}
}

func GetAnnouncements(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func (c *gin.Context) {
		uidExists := middlewares.MarkAuthorization(c, db, envMap, false)

		uid, _ := c.Get("uid")
		if uidExists { uid = uid.(int) }
		fmt.Println(uid, uidExists)

		var rows pgx.Rows
		var err error 
		if !uidExists {
			rows, err = db.Query(context.Background(), `SELECT a.announcement_id, a.title, a.content, u.wcaid, u.name FROM announcements a JOIN users u ON u.user_id = a.author_id;`)
		} else {
			rows, err = db.Query(context.Background(), `SELECT a.announcement_id, a.title, a.content, u.wcaid, u.name, ar.read FROM announcements a JOIN users u ON u.user_id = a.author_id JOIN announcement_read ar ON ar.announcement_id = a.announcement_id WHERE ar.user_id = $1;`, uid)
		}

		if err != nil {
			log.Println("ERR db.Query in GetAnnouncements: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying announcements.")
			return;
		}

		announcements := make([]models.AnnouncementState, 0)

		for rows.Next() {
			var announcement models.AnnouncementState

			if !uidExists {
				err = rows.Scan(&announcement.Id, &announcement.Title, &announcement.Content, &announcement.AuthorWcaId, &announcement.AuthorUsername)
				announcement.Read = true
			} else {
				err = rows.Scan(&announcement.Id, &announcement.Title, &announcement.Content, &announcement.AuthorWcaId, &announcement.AuthorUsername, &announcement.Read)
			}

			if err != nil {
				log.Println("ERR scanning announcement data in GetAnnouncements: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed parsing announcement from database.")
				return;
			}

			err = announcement.GetTags(db)
			if err != nil {
				log.Println("ERR GetTags in GetAnnouncements: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed to get announcement tags.")
				return;
			}

			announcements = append(announcements, announcement)
		}

		c.IndentedJSON(http.StatusOK, announcements)
	}
}

func PutAnnouncement(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var announcement models.AnnouncementState
		uid := c.MustGet("uid").(int)

		if err := c.BindJSON(&announcement); err != nil {
			log.Println("ERR BindJSON(&announcement) in PutAnnouncement: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse announcement data.")
			return
		}

		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Println("ERR db.begin in PutAnnouncement: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to start transaction.")
			tx.Rollback(context.Background())
			return
		}

		_, err = tx.Exec(context.Background(), `UPDATE announcements SET title = $1, content = $2, author_id = $3, timestamp = CURRENT_TIMESTAMP WHERE announcement_id = $4;`, announcement.Title, announcement.Content, uid, announcement.Id)
		if err != nil {
			log.Println("ERR tx.Exec UPDATE announcements in PutAnnouncement: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to update announcement info in database.")
			tx.Rollback(context.Background())
			return
		}

		err = models.UpdateAnnouncementTags(&announcement, db, tx, envMap)
		if err != nil {
			log.Println("ERR UpdateAnnouncementTags in PutAnnouncement: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to update announcement tag connections in database.")
			tx.Rollback(context.Background())
			return
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Println("ERR tx.commit in in PutAnnouncement: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to finish transaction.")
			return	
		}

		c.IndentedJSON(http.StatusCreated, announcement)
	}
}

func PostAnnouncement(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var announcement models.AnnouncementState

		if err := c.BindJSON(&announcement); err != nil {
			log.Println("ERR BindJSON(&announcement) in PostAnnouncement: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse announcement data.")
			return
		}

		announcement.AuthorId = c.MustGet("uid").(int)

		errLog, errOut := announcement.Create(db, envMap)
		if errLog != "" && errOut != "" {
			log.Println(errLog)
			c.IndentedJSON(http.StatusInternalServerError, errOut)
			return
		}

		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Println("ERR tx.Begin in PostAnnouncement: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to start transaction.")
			tx.Rollback(context.Background())
			return
		}

		logMessage, returnMessage := announcement.MakeAnnouncementUnreadForEveryone(tx)
		if logMessage != "" {
			log.Println(logMessage)
			c.IndentedJSON(http.StatusInternalServerError, returnMessage)
			tx.Rollback(context.Background())
			return
		}
		announcement.Read = false

		c.IndentedJSON(http.StatusCreated, announcement)
	}
}

func ReadAnnouncement(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var announcement models.AnnouncementState

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println("ERR strconv.Atoi in ReadAnnouncement of id (" + strconv.Itoa(id) + "): " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to parse announcement id to string.")
			return
		}

		announcement.Id = id
		announcement.AuthorId = c.MustGet("uid").(int)

		err = announcement.IsRead(db)
		if err != nil {
			log.Println("ERR announcement.IsRead in ReadAnnouncement (" + strconv.Itoa(announcement.Id) + "): " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to check if announcement is read.")
			return
		}

		if !announcement.Read {
			err = announcement.MarkRead(db)
			log.Println("ERR announcement.MarkRead in ReadAnnouncement (" + strconv.Itoa(announcement.Id) + "): " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to make announcement read.")
			return
		}

		c.IndentedJSON(http.StatusOK, announcement)
	}
}