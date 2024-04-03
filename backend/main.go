package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

type CompetitionEvent struct {
	Id int `json:"id"`
	Displayname string `json:"displayname"`
	Format string `json:"format"`
	Iconcode string `json:"iconcode"`
	Puzzlecode string `json:"puzzlecode"`
}

var events = []CompetitionEvent {
    {
        Id: 1,
        Displayname: "3x3x3",
        Format: "ao5",
        Iconcode: "333",
        Puzzlecode: "3x3x3",
    },
    {
        Id: 2,
        Displayname: "2x2x2",
        Format: "ao5",
        Iconcode: "222",
        Puzzlecode: "2x2x2",
    },
    {
        Id: 3,
        Displayname: "6x6x6",
        Format: "mo3",
        Iconcode: "666",
        Puzzlecode: "6x6x6",
    },
    {
        Id: 4,
        Displayname: "Mega",
        Format: "ao5",
        Iconcode: "mega",
        Puzzlecode: "megaminx",
    },
    {
        Id: 5,
    	Displayname: "Pyra",
        Format: "ao5",
        Iconcode: "pyra",
    	Puzzlecode: "pyraminx",
    },
    {
        Id: 6,
        Displayname: "3BLD",
        Format: "bo3",
        Iconcode: "333bld",
        Puzzlecode: "3x3x3",
    },
    {
        Id: 7,
        Displayname: "FMC",
        Format: "mo3",
        Iconcode: "fmc",
    	Puzzlecode: "3x3x3",
    },
}

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

var scrambles = [][]string{
    {
        "R2 U B2 D' R2 U L2 B' D U' L' F2 U' L F' D'",
        "F2 B U2 F2 D' B D2 L R2 U' F2 D F2 U' L2 U R2 U2 B'",
        "L D2 R2 B2 U' R2 D B2 U2 L2 R2 F2 R F' D' R' B U2 B",
        "B' L2 F' L2 R2 D2 R2 F' L2 B' L2 B R' U B' F2 D' R U' B' R",
        "U L2 B2 D' L2 F2 L2 U2 L2 U R2 U2 L' D' R2 B' D2 B2 D2",
	},
    {
        "R' U2 R F2 U' R' U2 R' F2",
        "F R2 U2 R F2 U F' R' F2",
        "U R2 U F' R2 U2 F' U2 R'",
        "R U2 R2 U' F R2 U F2 U2",
        "F2 U2 R' F R F2 U' F R2 F2",
    },
    {
        "R' U2 Uw2 3Rw2 Fw' 3Fw' D' Fw 3Fw2 R' Uw2 Lw Dw R D' Bw2 R2 U2 Rw2 3Rw U F' L 3Fw2 R' F2 3Rw2 D Dw' Lw' B R' Fw Bw2 3Uw2 Fw' U2 3Fw' Fw' D L2 F2 Uw 3Fw2 3Uw' Bw Uw2 R2 Rw' 3Fw2 R Lw B Dw2 U2 Bw 3Rw R2 3Fw Fw R' 3Uw' Fw Uw 3Rw2 L2 Lw' U2 Lw U2 Bw' F 3Fw Dw R2 Rw2 L' 3Rw 3Fw Fw2",
        "L' R2 Bw F 3Uw D' 3Fw Lw2 Rw' Bw' R Bw2 D2 Bw' D2 F' D2 L2 Rw2 Lw' 3Rw' F Bw' D2 3Uw Bw2 Lw' U 3Uw Rw Bw2 Lw' F' B Bw2 U 3Fw' F2 R2 Bw' Fw 3Rw2 Uw Fw R F2 Lw U2 Bw2 Uw' B Uw' Lw 3Uw2 Dw F Uw' F2 L2 3Fw Dw' Bw2 Rw2 Lw' Dw' F' Lw' B' Rw' D' Dw2 Fw Lw 3Fw Dw' D2 F' D2 3Fw2 Fw'",
        "R' Fw' D Fw Uw' U2 F' L2 Rw' 3Rw' Lw' Dw U2 Lw' 3Rw' R' B2 3Uw2 Uw2 F' 3Uw Rw2 F2 R2 Lw 3Uw2 Uw' 3Fw2 Fw2 D' 3Uw Fw 3Rw' Fw Dw2 3Rw2 Lw L F Lw' B2 Uw' 3Fw2 Dw D Lw' F2 R Bw' Rw' Fw' 3Rw' Fw 3Rw' F2 Fw 3Fw2 D2 F L B2 Lw' L2 D2 3Fw' 3Uw Uw Rw' Uw F' Rw' L' U Fw2 U Uw 3Uw F R 3Rw'",
    },
    {
        "R-- D-- R-- D-- R-- D-- R++ D++ R-- D-- U'\n  R++ D-- R-- D-- R-- D++ R++ D-- R++ D-- U'\n  R++ D++ R++ D++ R++ D++ R++ D-- R++ D++ U\n  R++ D++ R-- D-- R++ D-- R++ D++ R++ D++ U \n  R++ D++ R++ D-- R-- D++ R-- D++ R-- D++ U \n  R++ D-- R-- D++ R-- D++ R-- D-- R++ D-- U'\n  R++ D++ R-- D++ R++ D++ R-- D-- R-- D-- U'\n",
        "R-- D++ R-- D-- R++ D-- R++ D++ R++ D-- U'\n  R++ D++ R++ D-- R-- D-- R-- D++ R-- D++ U \n  R-- D-- R-- D-- R-- D-- R++ D++ R++ D-- U'\n  R++ D-- R++ D++ R-- D++ R-- D-- R++ D-- U'\n  R++ D-- R-- D-- R++ D-- R-- D-- R++ D-- U'\n  R++ D-- R++ D-- R-- D++ R-- D-- R++ D++ U \n  R++ D-- R++ D-- R++ D-- R++ D-- R-- D-- U'\n",
        "R-- D++ R-- D++ R-- D-- R-- D-- R-- D++ U \n  R++ D-- R++ D-- R-- D-- R-- D++ R++ D++ U \n  R-- D-- R-- D-- R++ D++ R-- D-- R-- D++ U \n  R-- D++ R-- D-- R++ D-- R++ D-- R-- D-- U'\n  R-- D++ R++ D-- R++ D++ R-- D++ R++ D++ U \n  R-- D++ R++ D++ R++ D-- R++ D++ R++ D-- U'\n  R++ D++ R++ D++ R++ D++ R++ D++ R-- D-- U'\n",
        "R-- D++ R++ D-- R++ D++ R-- D++ R-- D-- U'\n  R-- D++ R++ D-- R-- D++ R-- D-- R-- D-- U'\n  R-- D-- R-- D-- R-- D-- R++ D-- R-- D-- U'\n  R-- D++ R++ D++ R++ D++ R-- D-- R-- D++ U \n  R++ D++ R-- D-- R-- D-- R++ D-- R++ D++ U \n  R-- D++ R++ D++ R-- D++ R++ D++ R++ D-- U'\n  R++ D++ R++ D++ R-- D-- R-- D-- R++ D-- U'\n",
        "R++ D-- R-- D++ R-- D-- R-- D-- R++ D-- U'\n  R-- D-- R-- D++ R-- D-- R-- D-- R++ D-- U'\n  R-- D++ R++ D-- R++ D-- R++ D++ R-- D-- U'\n  R-- D-- R-- D++ R-- D++ R++ D++ R-- D++ U \n  R++ D++ R-- D++ R++ D++ R-- D-- R++ D++ U \n  R-- D-- R++ D++ R-- D++ R++ D-- R-- D-- U'\n  R++ D-- R++ D++ R-- D-- R++ D++ R++ D-- U'\n",
    },
    {
        "R U B U' B' L' U R' l' r' u'",
        "B' U' R B R' L B' U r u'",
        "L' U L R' L B' R B' l b u",
        "B' R U B L U' B' L l' r' b' u",
        "B L U' L B' R' U' B' l r b' u",
    },
    {
        "F2 L2 D' F2 D2 L2 U B2 U F2 U B' U L' B R2 D2 F R' D2 Rw' Uw'",
        "L2 B D2 B2 U2 L2 B' R2 F' D2 F L2 D' B' F' U' L R' B U' F Rw Uw",
        "R' D L2 B2 U2 R2 B2 L' U2 L2 B2 D2 R' F' L U2 R2 B F D' Rw",
    },
    {
        "R' U' F D2 U2 L2 B U2 B2 L2 D2 U2 F' R B2 F R U F U' B2 F2 R' U' F",
        "R' U' F D2 F2 D2 B2 L2 F2 L2 D F2 B' R F' L F U2 F' D R' U' F",
        "R' U' F L U2 R' D2 L B2 R F2 L B2 U2 F2 B' L F' R D F U L2 B2 R' U' F",
	},
}

type CompetitionData struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Startdate time.Time `json:"startdate"`
	Enddate time.Time `json:"enddate"`
	Events []CompetitionEvent `json:"events"`
	Scrambles [][]string `json:"scrambles"`
	Results ResultEntry `json:"results"`
}

func allCompetitionData() []CompetitionData {
	res := make([]CompetitionData, 0)
	startdate := time.Now()
	startdate = startdate.AddDate(0, 0, -23)
	enddate := time.Time(startdate)
	enddate = enddate.AddDate(0, 0, 7)

	for i := range [10]int{} {
		if len(res) > 0 {
			startdate = time.Time(res[len(res) - 1].Enddate)
			enddate = time.Time(startdate)
			enddate = enddate.AddDate(0, 0, 7)
		}

		res = append(res, CompetitionData{
			Id: fmt.Sprintf("WeeklyCompetition%d", i + 1),
			Name: fmt.Sprintf("Weekly Competition %d", i + 1),
			Startdate: startdate,
			Enddate: enddate,
			Events: events,
			Scrambles: scrambles,
		})
	}

	return res
}

type ManageRolesUser struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Isadmin bool `json:"isadmin"`
}

var manageRolesUsers = []ManageRolesUser {
	{
		Id: 1,
		Name: "Janko Hrasko",
		Isadmin: true,
	},
	{
		Id: 2,
		Name: "Ferko Mrkvicka",
		Isadmin: false,
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
	router.GET("/api/results", getResults)
	router.GET("/api/results/:id/:event", getResultsByIdAndEvent)
	router.GET("/api/events", getEvents)
	router.GET("/api/competitions/:filter", getFilteredCompetitions)
	router.GET("/api/competition/:id", getCompetitionById)
	router.GET("/api/users/manage-roles", getManageRolesUsers)
	router.PUT("/api/users/manage-roles", putManageRolesUsers)

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

func getResultsByIdAndEvent(c *gin.Context) {
	event := c.Param("event")

	c.IndentedJSON(http.StatusOK, results[event])
}

func getEvents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, events);
}

func getFilteredCompetitions(c *gin.Context) {
	filter := c.Param("filter")
	
	result := make([]CompetitionData, 0);
	competitions := allCompetitionData();

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

func getCompetitionById(c *gin.Context) {
	id := c.Param("id")
	competitions := allCompetitionData()
	
	idx := slices.IndexFunc(competitions, func (c CompetitionData) bool { return c.Id == id })
	result := CompetitionData{}
	if idx != -1 {
		result = competitions[idx]
	}

	c.IndentedJSON(http.StatusOK, result)
}

func getManageRolesUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, manageRolesUsers)
}

func putManageRolesUsers(c *gin.Context) {
	var newManageRolesUsers []ManageRolesUser

	if err := c.BindJSON(&newManageRolesUsers); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "");
		return;
	}

	for _, cur_mru := range newManageRolesUsers {
		idx := slices.IndexFunc(manageRolesUsers, func (mru ManageRolesUser) bool { return mru.Id == cur_mru.Id })
		if idx == -1 {
			c.IndentedJSON(http.StatusInternalServerError, fmt.Sprintf("User %v is not present!", cur_mru.Name));
			return;
		}

		manageRolesUsers[idx] = cur_mru
	}

	c.IndentedJSON(http.StatusCreated, manageRolesUsers)
}