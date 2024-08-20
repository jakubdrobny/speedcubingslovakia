package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
)

func GetAnnouncementById(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")	
		
		rows, err := db.Query(context.Background(), `SELECT a.announcement_id, a.title, a.content, u.wcaid, u.name FROM announcements a JOIN users u ON u.user_id = a.author_id WHERE a.announcement_id = $1;`, id)
		if err != nil {
			log.Println("ERR db.Query in GetAnnouncementById: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed querying announcement by id.")
			return;
		}

		var announcement models.AnnouncementState
		found := false

		for rows.Next() {
			err := rows.Scan(&announcement.Id, &announcement.Title, &announcement.Content, &announcement.AuthorWcaId, &announcement.AuthorUsername)
			if err != nil {
				log.Println("ERR scanning announcement data in GetAnnouncementById: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Failed parsing announcement from database.")
				return;
			}
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

		c.IndentedJSON(http.StatusCreated, announcement)
	}
}
