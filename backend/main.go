package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ResultEntry struct {
	Id int `json:"id"`
	Userid int `json:"userid"`
	Username string `json:"username"`
	Competitionid int `json:"competitionid"`
	Competitionname string `json:"competitionname"`
	Eventid int `json:"eventid"`
	Eventname string `json:"eventname"`
	Iconcode string `json:"iconcode"`
	Format string `json:"format"`
	Solve1 string `json:"solve1"`
	Solve2 string `json:"solve2"`
	Solve3 string `json:"solve3"`
	Solve4 string `json:"solve4"`
	Solve5 string `json:"solve5"`
	Comment string `json:"comment"`
	Status ResultsStatus `json:"status"`
}

type ResultsStatus struct {
	Id int `json:"id"`
	ApprovalFinished bool `json:"approvalFinished"`
	Approved bool `json:"approved"`
	Visible bool `json:"visible"`
	Displayname string `json:"displayname"`
}

var waitingForApprovalResultsStatus = ResultsStatus {
	Id: 1,
    ApprovalFinished: false,
    Visible: false,
    Displayname: "Waiting for approval",
}

var deniedResultsStatus = ResultsStatus {
	Id: 2,
    ApprovalFinished: true,
	Approved: false,
    Visible: false,
    Displayname: "Denied",
}

var approvedResultsStatus = ResultsStatus {
	Id: 3,
    ApprovalFinished: true,
	Approved: true,
    Visible: true,
    Displayname: "Approved",
}

var results = map[string]ResultEntry {
    "3x3x3": {
        Id: 1,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: 4,
        Competitionname: "Weekly Competition 4",
        Eventid: 1,
        Eventname: "3x3x3",
        Iconcode: "333",
        Format: "ao5",
        Solve1: "12.55",
        Solve2: "10.14",
        Solve3: "8.81",
        Solve4: "DNF",
        Solve5: "14.43",
        Comment: "",
        Status: waitingForApprovalResultsStatus,
    },
    "2x2x2": {
        Id: 2,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: 4,
        Competitionname: "Weekly Competition 4",
        Eventid: 2,
        Eventname: "2x2x2",
        Iconcode: "222",
        Format: "ao5",
        Solve1: "2.55",
        Solve2: "1.14",
        Solve3: "8.81",
        Solve4: "2.00",
        Solve5: "1.43",
        Comment: "",
        Status: approvedResultsStatus,
    },
    "6x6x6": {
        Id: 3,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: 4,
        Competitionname: "Weekly Competition 4",
        Eventid: 3,
        Eventname: "6x6x6",
        Iconcode: "666",
        Format: "mo3",
        Solve1: "2:00.55",
        Solve2: "1:59.14",
        Solve3: "1:58.80",
        Solve4: "",
        Solve5: "",
        Comment: "",
        Status: approvedResultsStatus,
    },
    "Mega": {
        Id: 4,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: 4,
        Competitionname: "Weekly Competition 4",
        Eventid: 5,
        Eventname: "Megaminx",
        Iconcode: "mega",
        Format: "ao5",
        Solve1: "42.55",
        Solve2: "41.14",
        Solve3: "48.81",
        Solve4: "42.00",
        Solve5: "41.43",
        Comment: "",
        Status: deniedResultsStatus,
    },
    "Pyra": {
        Id: 5,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: 4,
        Competitionname: "Weekly Competition 4",
        Eventid: 5,
        Eventname: "Pyraminx",
        Iconcode: "pyra",
        Format: "ao5",
        Solve1: "2.13",
        Solve2: "1.01",
        Solve3: "2.99",
        Solve4: "2.00",
        Solve5: "2.69",
        Comment: "",
        Status: waitingForApprovalResultsStatus,
    },
    "3BLD": {
        Id: 6,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: 4,
        Competitionname: "Weekly Competition 4",
        Eventid: 6,
        Eventname: "3BLD",
        Iconcode: "333bld",
        Format: "bo3",
        Solve1: "DNF",
        Solve2: "1:00.05",
        Solve3: "DNS",
        Solve4: "",
        Solve5: "",
        Comment: "",
        Status: deniedResultsStatus,
    },
    "FMC": {
        Id: 7,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: 4,
        Competitionname: "Weekly Competition 4",
        Eventid: 7,
        Eventname: "FMC",
        Iconcode: "fmc",
        Format: "mo3",
        Solve1: "R U R' U'",
        Solve2: "abc",
        Solve3: "",
        Solve4: "",
        Solve5: "",
        Comment: "",
        Status: approvedResultsStatus,
    },
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type"},
        ExposeHeaders: []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

	router.GET("/api/ping", ping)
	router.GET("api/results", getResults)

	router.Run("localhost:8080")
}

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "pong")
}

func getResults(c *gin.Context) {
	r := make([]ResultEntry, 0, len(results))

	for _, value := range results {
		r = append(r, value)
	}

	c.IndentedJSON(http.StatusOK, r)
}