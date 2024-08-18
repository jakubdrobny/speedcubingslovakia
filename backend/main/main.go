package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/jakubdrobny/speedcubingslovakia/backend/controllers"
	"github.com/jakubdrobny/speedcubingslovakia/backend/middlewares"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, fmt.Sprintf("Too many requests. Try again in %ds.", int(math.Round(time.Until(info.ResetTime).Seconds()))))
}

func main() {
	envMap, err := godotenv.Read(fmt.Sprintf(".env.%s", os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV")))
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

	cronScheduler := cron.New()
	cronScheduler.AddFunc("@every 7d", func () { controllers.AddNewWeeklyCompetition(db, envMap) })
	cronScheduler.Start()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type"},
        ExposeHeaders: []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))
	
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 100,
	})

	rateLimitMiddleWare := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc: keyFunc,
	})

	router.Use(rateLimitMiddleWare)

	api_v1 := router.Group("/api")
	
	results := api_v1.Group("/results")
	{
		results.GET("/edit/:uname/:cname/:eid/:rsname", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), controllers.GetResultsQuery(db))
		results.GET("/compete/:cid/:eid", middlewares.AuthMiddleWare(db, envMap), controllers.GetResultsByIdAndEvent(db))
		results.POST("/save", middlewares.AuthMiddleWare(db, envMap), controllers.PostResults(db))
		results.POST("/save-validation", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), controllers.PostResultsValidation(db))
		results.GET("/rankings", controllers.GetRankings(db))
		results.GET("/records", controllers.GetRecords(db))
		results.GET("/regions/grouped", controllers.GetRegionsGrouped(db))
		results.GET("/profile/:id", controllers.GetProfileResults(db))
	}

	events := api_v1.Group("/events")
	{
		events.GET("/", controllers.GetEvents(db))
	}

	resultsStatuses := api_v1.Group("/resultsStatuses")
	{
		resultsStatuses.GET("/", controllers.GetResultsStatuses(db))
	}

	competitions := api_v1.Group("/competitions")
	{
		competitions.GET("/filter/:filter", controllers.GetFilteredCompetitions(db))
		competitions.GET("/id/:id", controllers.GetCompetitionById(db))
		competitions.POST("/", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), controllers.PostCompetition(db, envMap))
		competitions.PUT("/", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), controllers.PutCompetition(db, envMap))
		competitions.GET("/results/:cid/:eid", controllers.GetResultsFromCompetition(db))
	}

	users := api_v1.Group("/users")
	{
		users.GET("/manage-roles", middlewares.AuthMiddleWare(db, envMap), controllers.GetManageRolesUsers(db))
		users.PUT("/manage-roles", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), controllers.PutManageRolesUsers(db))
		users.POST("/login", controllers.PostLogIn(db, envMap))
		users.GET("/search", controllers.GetSearchUsers(db))
		users.GET("/auth/admin", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), func(c *gin.Context) { c.IndentedJSON(http.StatusAccepted, "authorized")});
	}

	router.Run("localhost:8000")
}