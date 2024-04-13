package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jakubdrobny/speedcubingslovakia/backend/middlewares"
	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"

	"math/rand"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	envMap, err := godotenv.Read(".env.development")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load enviromental variables from file: %v\n", err)
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), envMap["DB_URL"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type"},
        ExposeHeaders: []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

	api_v1 := router.Group("/api")
	
	results := api_v1.Group("/results", middlewares.AuthMiddleWare(db, envMap))
	{
		results.GET("/edit/:uname/:cid/:eid", middlewares.AdminMiddleWare(), GetResultsQuery(db))
		results.GET("/compete/:cid/:eid", GetResultsByIdAndEvent(db))
		results.POST("/save", PostResults(db))
		results.POST("/save-validation", middlewares.AdminMiddleWare(), PostResultsValidation(db))
	}

	events := api_v1.Group("/events")
	{
		events.GET("/", GetEvents(db))
	}

	competitions := api_v1.Group("/competitions")
	{
		competitions.GET("/filter/:filter", GetFilteredCompetitions(db))
		competitions.GET("/id/:id", GetCompetitionById(db))
		competitions.POST("/", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), PostCompetition(db))
		competitions.PUT("/", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), PutCompetition(db))
		competitions.GET("/results/:cid/:eid", getResultsFromCompetition(db))
	}

	users := api_v1.Group("/users")
	{
		users.GET("/manage-roles", middlewares.AuthMiddleWare(db, envMap), GetManageRolesUsers(db))
		users.PUT("/manage-roles", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), PutManageRolesUsers(db))
	}

	router.GET("/api/auth/admin", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), func(c *gin.Context) { c.IndentedJSON(http.StatusAccepted, "authorized")});

	router.POST("/api/login", postLogIn(db, envMap))

	router.Run("localhost:8080")
}

func getResultsFromCompetition(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		cid := c.Param("cid")
		eid, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		competitionResults, err := models.GetResultsFromCompetitionByEventName(db, cid, eid)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusAccepted, competitionResults)
	}
}

func postLogIn(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
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

func GetResultsQuery(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		userName := c.Param("uname")
		competitionId := c.Param("cid")
		eventId, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		
		var resultEntries []models.ResultEntry

		if competitionId == "_" && userName == "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re WHERE re.event_id = $1;`, eventId)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else if competitionId == "_" && userName != "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN users u ON u.user_id = re.user_id WHERE re.event_id = $1 AND u.name = $2;`, eventId, userName)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else if competitionId != "_" && userName == "_" {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re WHERE re.event_id = $1 AND re.competition_id = $2;`, eventId, competitionId)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		} else {
			rows, err := db.Query(context.Background(), `SELECT re.result_id FROM results re JOIN users u ON u.user_id = re.user_id WHERE re.event_id = $1 AND re.competition_id = $2 AND u.name = $3;`, eventId, competitionId, userName)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}

			for rows.Next() {
				var resultEntryId int
				err = rows.Scan(&resultEntryId)
				if err != nil { 
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}

				resultEntry, err := models.GetResultEntryById(db, resultEntryId)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
				resultEntries = append(resultEntries, resultEntry)
			}
		}

		c.IndentedJSON(http.StatusOK, resultEntries)
	}
}

func GetResultsByIdAndEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		eventId, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		competitionId := c.Param("cid")
		userId := c.MustGet("uid").(int)

		user, err := models.GetUserById(db, userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		event, err := models.GetEventById(db, eventId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		competition, err := models.GetCompetitionByIdObject(db, competitionId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return	
		}

		resultEntry, err := models.GetResultEntry(db, userId, competitionId, eventId)

		if err != nil {
			if err.Error() != "not found" {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			} else {
				approvedResultsStatus, err := models.GetResultsStatus(db, 3)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return	
				}

				resultEntry = models.ResultEntry{
					Id: rand.Int(),
					Userid: userId,
					Username: user.Name,
					Competitionid: competitionId,
					Competitionname: competition.Name,
					Eventid: event.Id,
					Eventname: event.Displayname,
					Iconcode: event.Iconcode,
					Format: event.Format,
					Solve1: "DNS",
					Solve2: "DNS",
					Solve3: "DNS",
					Solve4: "DNS",
					Solve5: "DNS",
					Comment: "",
					Status: approvedResultsStatus,
				}

				err = resultEntry.Insert(db)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
			}
		} else {
			currentStatus, err := models.GetResultsStatus(db, resultEntry.Status.Id)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}
			
			resultEntry.Status = currentStatus
			resultEntry.Eventname = event.Displayname
			resultEntry.Competitionname = competition.Name
			resultEntry.Username = user.Name
			resultEntry.Iconcode = event.Iconcode
			resultEntry.Format = event.Format

			err = resultEntry.Update(db)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}
		}


		c.IndentedJSON(http.StatusOK, resultEntry)
	}
}

func PostResults(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		var resultEntry models.ResultEntry
		var err error

		if err = c.BindJSON(&resultEntry); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}
		
		err = resultEntry.Update(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusCreated, resultEntry)
	}
}

func PostResultsValidation(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		type ReqBody struct {
			ResultId int `json:"resultId"`
			Verdict bool `json:"verdict"`
		}
		var body ReqBody

		if err := c.BindJSON(&body); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}

		resultEntry, err := models.GetResultEntryById(db, body.ResultId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}

		statusId := 3
		if !body.Verdict { statusId = 2 }
		resultStatus, err := models.GetResultsStatus(db, statusId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}

		resultEntry.Status = resultStatus
		err = resultEntry.Update(db, true)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return
		}

		c.IndentedJSON(http.StatusCreated, "")
	}
}

func GetEvents(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		events, err := models.GetAvailableEvents(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			c.IndentedJSON(http.StatusOK, events);
		}
	}
}

func GetFilteredCompetitions(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		filter := c.Param("filter")
		
		result := make([]models.CompetitionData, 0);
		competitions, err := models.GetAllCompetitions(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return
		}

		now := time.Now()
		if filter == "Past" {
			for _, competition := range competitions {
				if competition.Enddate.Before(now) {
					result = append(result, competition)
				}
			}
		} else if filter == "Current" {
			for _, competition := range competitions {
				if competition.Startdate.Before(now) && now.Before(competition.Enddate) {
					result = append(result, competition)
				}
			}
		} else if filter == "Future" {
			for _, competition := range competitions {
				if now.Before(competition.Startdate) {
					result = append(result, competition)
				}
			}
		}

		c.IndentedJSON(http.StatusOK, result);
	}
}

func GetCompetitionById(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")	
		
		rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c WHERE c.competition_id = $1;`, id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}

		var competition models.CompetitionData
		found := false

		for rows.Next() {
			err := rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)	
				return;
			}
			found = true
		}

		if !found {
			c.IndentedJSON(http.StatusInternalServerError, "Competition not found.")	
			return;
		}

		err = competition.GetEvents(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}

		err = competition.GetScrambles(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}
		
		c.IndentedJSON(http.StatusOK, competition)
	}
}

func PostCompetition(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var competition models.CompetitionData

		if err := c.BindJSON(&competition); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "what bro")
			return
		}

		competition.RecomputeCompetitionId()

		tx, err := db.Begin(context.Background())
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			tx.Rollback(context.Background())
			return
		}

		_, err = tx.Exec(context.Background(), `INSERT INTO competitions (competition_id, name, startdate, enddate) VALUES ($1,$2,$3,$4);`, competition.Id, competition.Name, competition.Startdate, competition.Enddate)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			tx.Rollback(context.Background())
			return
		}

		for _, event := range competition.Events {
			_, err := tx.Exec(context.Background(), `INSERT INTO competition_events (competition_id, event_id) VALUES ($1,$2);`, competition.Id, event.Id)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				tx.Rollback(context.Background())
				return
			}
		}

		for _, scrambleSet := range competition.Scrambles {
			for scrambleIdx, scramble := range scrambleSet.Scrambles {
				_, err := tx.Exec(context.Background(), `INSERT INTO scrambles (scramble, event_id, competition_id, "order") VALUES ($1,$2,$3,$4);`, scramble, scrambleSet.Event.Id, competition.Id, scrambleIdx + 1)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					tx.Rollback(context.Background())
					return
				}
			}
		}

		err = tx.Commit(context.Background())
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return	
		}

		c.IndentedJSON(http.StatusCreated, competition)
	}
}

func PutCompetition(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var competition models.CompetitionData

		if err := c.BindJSON(&competition); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		_, err := db.Exec(context.Background(), `UPDATE competitions SET name = $1, startdate = $2, enddate = $3, timestamp = CURRENT_TIMESTAMP WHERE competition_id = $4;`, competition.Name, competition.Startdate, competition.Enddate, competition.Id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		err = models.UpdateCompetitionEvents(&competition, db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusCreated, competition)
	}
}

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