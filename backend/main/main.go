package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/jakubdrobny/speedcubingslovakia/backend/controllers"
	"github.com/jakubdrobny/speedcubingslovakia/backend/logging"
	"github.com/jakubdrobny/speedcubingslovakia/backend/middlewares"
)

func main() {
	logger := logging.CustomLogger()

	slog.SetDefault(logger)

	envMap, err := godotenv.Read(
		fmt.Sprintf(".env.%s", os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV")),
	)
	if err != nil {
		slog.Error("unable to load environmental variables from file", "error", err)
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), envMap["DB_URL"])
	if err != nil {
		slog.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:3000", "http://localhost:3000", "http://frontend:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(logging.GinLoggerMiddleware(logger), logging.GinRecoveryMiddleware(logger))

	api_v1 := router.Group("/api")

	stats := api_v1.Group("/stats")
	{
		stats.GET(
			"/dashboard",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.GetAdminStats(db),
		)
	}

	results := api_v1.Group("/results")
	{
		results.GET(
			"/edit/:uname/:cname/:eid/:rsname",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.GetResultsQuery(db),
		)
		results.GET(
			"/compete/:cid/:eid",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.GetResultsByIdAndEvent(db),
		)
		results.POST(
			"/save",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.PostResults(db, envMap),
		)
		results.POST(
			"/save-validation",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.PostResultsValidation(db),
		)
		results.GET(
			"/save-validation",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.GetResultsValidation(db),
		)
		results.GET("/rankings", controllers.GetRankings(db))
		results.GET("/records", controllers.GetRecords(db))
		results.GET("/regions/grouped", controllers.GetRegionsGrouped(db))
		results.GET("/profile/:id", controllers.GetProfileResults(db))
		results.POST(
			"/averageinfo",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.GetAverageInfo(db),
		)
		results.POST(
			"/averageinfo/records",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.GetAverageInfoRecords(db),
		)
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
		competitions.GET("/wca", controllers.GetUpcomingWCACompetitions(db))
		competitions.GET("/wca/regions/grouped", controllers.GetWCARegionGroups(db))
		competitions.GET(
			"/wca/subscriptions/positions",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.GetWCACompAnnouncementsPositionSubscriptions(db),
		)
		competitions.GET(
			"/wca/subscriptions",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.GetWCACompAnnouncementSubscriptions(db),
		)
		competitions.POST(
			"/wca/subscribe",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.UpdateWCAAnnouncementSubscriptions(db),
		)
		competitions.POST(
			"/wca/subscribe/position/upsert",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.UpdateWCAAnnouncementsPositionSubscriptions(db),
		)
		competitions.DELETE(
			"/wca/subscribe/position/delete",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.DeleteWCAAnnouncementsPositionSubscriptions(db),
		)
		competitions.POST(
			"/",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.PostCompetition(db, envMap),
		)
		competitions.PUT(
			"/",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.PutCompetition(db, envMap),
		)
		competitions.GET("/results/:cid/:eid", controllers.GetResultsFromCompetition(db))
	}

	users := api_v1.Group("/users")
	{
		users.GET(
			"/manage-roles",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.GetManageRolesUsers(db),
		)
		users.PUT(
			"/manage-roles",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.PutManageRolesUsers(db),
		)
		users.POST("/login", controllers.PostLogIn(db, envMap))
		users.GET("/search", controllers.GetSearchUsers(db))
		users.GET("/map", controllers.GetUserMapData(db))
		users.GET(
			"/auth/admin",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			func(c *gin.Context) { c.IndentedJSON(http.StatusAccepted, "authorized") },
		)
	}

	tags := api_v1.Group("/tags")
	{
		tags.GET("/", controllers.GetTags(db))
	}

	announcements := api_v1.Group("/announcements")
	{
		announcements.GET("/id/:id", controllers.GetAnnouncementById(db, envMap))
		announcements.GET(
			"/read/:id",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.ReadAnnouncement(db),
		)
		announcements.POST(
			"/react/:id",
			middlewares.AuthMiddleWare(db, envMap),
			controllers.ReactToAnnouncement(db),
		)
		announcements.DELETE(
			"/delete/:id",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.DeleteAnnouncement(db),
		)
		announcements.GET("/", controllers.GetAnnouncements(db, envMap))
		announcements.POST(
			"/",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.PostAnnouncement(db, envMap),
		)
		announcements.PUT(
			"/",
			middlewares.AuthMiddleWare(db, envMap),
			middlewares.AdminMiddleWare(),
			controllers.PutAnnouncement(db, envMap),
		)
		announcements.GET("/noOfNew", controllers.GetNoOfNewAnnouncements(db, envMap))
	}

	if err := router.Run("0.0.0.0:8000"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
