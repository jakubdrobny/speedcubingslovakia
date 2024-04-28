package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jakubdrobny/speedcubingslovakia/backend/controllers"
	"github.com/jakubdrobny/speedcubingslovakia/backend/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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
		results.GET("/edit/:cid/:uname/:eid", middlewares.AdminMiddleWare(), controllers.GetResultsQuery(db))
		results.GET("/compete/:cid/:eid", controllers.GetResultsByIdAndEvent(db))
		results.POST("/save", controllers.PostResults(db))
		results.POST("/save-validation", middlewares.AdminMiddleWare(), controllers.PostResultsValidation(db))
	}

	events := api_v1.Group("/events")
	{
		events.GET("/", controllers.GetEvents(db))
	}

	competitions := api_v1.Group("/competitions")
	{
		competitions.GET("/filter/:filter", controllers.GetFilteredCompetitions(db))
		competitions.GET("/id/:id", controllers.GetCompetitionById(db))
		competitions.POST("/", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), controllers.PostCompetition(db))
		competitions.PUT("/", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), controllers.PutCompetition(db))
		competitions.GET("/results/:cid/:eid", controllers.GetResultsFromCompetition(db))
	}

	users := api_v1.Group("/users")
	{
		users.GET("/manage-roles", middlewares.AuthMiddleWare(db, envMap), controllers.GetManageRolesUsers(db))
		users.PUT("/manage-roles", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), controllers.PutManageRolesUsers(db))
		users.POST("/login", controllers.PostLogIn(db, envMap))
		users.GET("/auth/admin", middlewares.AuthMiddleWare(db, envMap), middlewares.AdminMiddleWare(), func(c *gin.Context) { c.IndentedJSON(http.StatusAccepted, "authorized")});
	}

	rankings := api_v1.Group("/rankings")
	{
		rankings.GET("/profile/:id", controllers.GetProfileResults(db))
	}

	router.Run("localhost:8000")
}