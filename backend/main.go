package main

import (
	"fmt"
	"net/http"
	"time"

	"math/rand"

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
        Displayname: "4x4x4",
        Format: "ao5",
        Iconcode: "444",
        Puzzlecode: "4x4x4",
    },
	{
        Id: 4,
        Displayname: "5x5x5",
        Format: "ao5",
        Iconcode: "555",
        Puzzlecode: "5x5x5",
    },
    {
        Id: 5,
        Displayname: "6x6x6",
        Format: "mo3",
        Iconcode: "666",
        Puzzlecode: "6x6x6",
    },
	{
		Id: 6,
		Displayname: "7x7x7",
        Format: "mo3",
        Iconcode: "777",
        Puzzlecode: "7x7x7",
	},
    {
        Id: 7,
        Displayname: "3BLD",
        Format: "bo3",
        Iconcode: "333bf",
        Puzzlecode: "3x3x3",
    },
    {
        Id: 8,
        Displayname: "FMC",
        Format: "mo3",
        Iconcode: "333fm",
    	Puzzlecode: "3x3x3",
    },
	{
        Id: 9,
        Displayname: "OH",
        Format: "ao5",
        Iconcode: "333oh",
    	Puzzlecode: "3x3x3",
    },
	{
		Id: 10,
		Displayname: "Clock",
        Format: "ao5",
        Iconcode: "clock",
        Puzzlecode: "clock",
	},
    {
        Id: 11,
        Displayname: "Mega",
        Format: "ao5",
        Iconcode: "minx",
        Puzzlecode: "megaminx",
    },
    {
        Id: 12,
    	Displayname: "Pyra",
        Format: "ao5",
        Iconcode: "pyram",
    	Puzzlecode: "pyraminx",
    },
	{
        Id: 13,
    	Displayname: "Skewb",
        Format: "ao5",
        Iconcode: "skewb",
    	Puzzlecode: "skewb",
    },
	{
        Id: 14,
    	Displayname: "Sq-1",
        Format: "ao5",
        Iconcode: "sq1",
    	Puzzlecode: "square1",
    },
	{
        Id: 15,
        Displayname: "4BLD",
        Format: "bo3",
        Iconcode: "444bf",
        Puzzlecode: "4x4x4",
    },
	{
        Id: 16,
        Displayname: "5BLD",
        Format: "bo3",
        Iconcode: "555bf",
        Puzzlecode: "5x5x5",
    },
}

type ResultEntry struct {
	Id int `json:"id"`
	Userid int `json:"userid"`
	Username string `json:"username"`
	Competitionid string `json:"competitionid"`
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

var results = []ResultEntry {
    {
        Id: 1,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: "WeeklyCompetition4",
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
    {
        Id: 2,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: "WeeklyCompetition4",
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
    {
        Id: 3,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: "WeeklyCompetition4",
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
    {
        Id: 4,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: "WeeklyCompetition4",
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
    {
        Id: 5,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: "WeeklyCompetition4",
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
    {
        Id: 6,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: "WeeklyCompetition4",
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
    {
        Id: 7,
        Userid: 1,
        Username: "Janko Hrasko",
        Competitionid: "WeeklyCompetition4",
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

type ScrambleSet struct {
	Event CompetitionEvent `json:"event"`
	Scrambles []string `json:"scrambles"`
}

var scrambles = []ScrambleSet {
    {
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "3x3x3"} )],
		Scrambles: []string{
        	"R2 U B2 D' R2 U L2 B' D U' L' F2 U' L F' D'",
        	"F2 B U2 F2 D' B D2 L R2 U' F2 D F2 U' L2 U R2 U2 B'",
        	"L D2 R2 B2 U' R2 D B2 U2 L2 R2 F2 R F' D' R' B U2 B",
        	"B' L2 F' L2 R2 D2 R2 F' L2 B' L2 B R' U B' F2 D' R U' B' R",
        	"U L2 B2 D' L2 F2 L2 U2 L2 U R2 U2 L' D' R2 B' D2 B2 D2",
		},
	},
    {
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "2x2x2"} )],
		Scrambles: []string{
			"R' U2 R F2 U' R' U2 R' F2",
			"F R2 U2 R F2 U F' R' F2",
			"U R2 U F' R2 U2 F' U2 R'",
			"R U2 R2 U' F R2 U F2 U2",
			"F2 U2 R' F R F2 U' F R2 F2",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "4x4x4"} )],
		Scrambles: []string{
			"B2 L2 U2 F' R2 D2 F B R2 B' L2 B' R F D2 B2 U F L' F2 Fw2 Rw2 D L' Fw2 Uw2 L2 D2 L' B2 L Fw' U B R Uw2 B' D2 Rw' B2 Fw Uw2 L B2 Fw2",
			"L2 F' B2 R' D2 L F2 B2 L2 U2 L U' L U' F U2 R2 D Uw2 Fw2 R2 Uw2 B L F' Uw2 B Rw2 B2 L2 Uw Rw2 U' Uw' L U2 Fw Rw D' Fw R' F2 L'",
			"B D B2 D2 L2 U' B2 U L2 F2 U' R2 F' L F U R U' R2 D' F2 Uw2 B U Rw2 B' D Fw2 D Rw2 U B' Uw2 F' L Fw2 U Rw U' L' Fw R Rw' Fw2 Uw' F Rw'", 
			"U B' L B2 R D2 R2 F2 B2 R' B2 R' U2 F U' L' D' B' D2 R2 L2 Uw2 L' F' Fw2 R2 B Fw2 Uw2 Rw2 F' R' D2 Uw' R L B2 L Fw' D Rw2 Fw' D' Fw'",
			"B2 R' U R2 F2 U2 R2 F R2 F' R2 F2 U2 B L' U2 D' L2 D' L B2 Rw2 Uw2 F2 L B2 L2 Rw2 D2 Uw2 B' R B' Fw2 Uw' L2 D F2 Rw F2 D2 Rw' Uw F Rw2 B'",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "5x5x5"} )],
		Scrambles: []string{
			"F2 Rw U' F' D Dw' Uw2 Lw' Fw Bw2 L2 Lw2 Bw D' B D2 B' Dw' Rw' D2 B' Bw F' Dw' Uw Bw' D2 F Uw Bw2 F' Dw2 R B R F Uw L' R2 Fw R' F Lw' F R' Dw D R2 Rw' Uw2 B' R D2 Lw' R' D Lw2 Uw2 R' Lw'",
			"R Dw' Fw' Dw2 B2 Dw2 L' Uw2 B' Fw2 Bw' R2 Uw2 Rw Fw2 R' F' Dw Lw' Bw' Fw' Rw' F2 Bw B' D2 R2 Uw' Rw2 Fw F Dw2 Fw Rw' R2 Dw2 D2 Fw Bw2 U2 Rw' U2 D L2 Uw2 U Lw L2 B2 D' B2 F' R2 Dw' F2 D2 Dw B2 Fw L",
			"Fw2 Uw' Lw2 Fw' Rw' Lw' Fw2 B' Rw Bw2 R Rw' U2 B D U2 Bw Fw2 L Fw2 B Lw2 L' U2 Uw R2 Lw2 Dw' D2 U' F' Lw' Fw B2 D' Lw2 L2 Dw R' Uw2 B Bw2 U D Dw L U' Uw2 R2 F Bw' U2 L Dw' Lw2 Fw Dw' Rw Dw' Bw2",
			"Uw Lw2 F' Lw2 R2 U2 Fw Lw2 Uw2 D2 B Bw2 R2 Dw' Rw' Dw' B2 Bw2 Rw2 Lw' Uw2 Rw2 U Uw' F2 Lw R' B' Lw' D' U' L' R' F2 Lw' Uw' Bw2 Rw' L' Uw2 B' F2 L F2 Fw Uw2 Bw L2 F2 B2 L R Dw' R2 F' Uw B' Bw Rw' Fw'",
			"Rw' Dw' Rw' Bw2 Uw D R Rw2 B' R Uw' L2 Dw R' B Rw D' U' Bw Fw2 D L2 Uw2 B2 Bw Fw2 U' Fw2 B' Bw Dw2 D F Dw' Uw' Lw R2 B Lw R2 B' D L Lw' R' D2 Dw U2 L2 Rw' Uw Lw Dw Bw' Fw R Fw Lw' Bw' Uw'",

		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "6x6x6"} )],
		Scrambles: []string{
			"R' U2 Uw2 3Rw2 Fw' 3Fw' D' Fw 3Fw2 R' Uw2 Lw Dw R D' Bw2 R2 U2 Rw2 3Rw U F' L 3Fw2 R' F2 3Rw2 D Dw' Lw' B R' Fw Bw2 3Uw2 Fw' U2 3Fw' Fw' D L2 F2 Uw 3Fw2 3Uw' Bw Uw2 R2 Rw' 3Fw2 R Lw B Dw2 U2 Bw 3Rw R2 3Fw Fw R' 3Uw' Fw Uw 3Rw2 L2 Lw' U2 Lw U2 Bw' F 3Fw Dw R2 Rw2 L' 3Rw 3Fw Fw2",
        	"L' R2 Bw F 3Uw D' 3Fw Lw2 Rw' Bw' R Bw2 D2 Bw' D2 F' D2 L2 Rw2 Lw' 3Rw' F Bw' D2 3Uw Bw2 Lw' U 3Uw Rw Bw2 Lw' F' B Bw2 U 3Fw' F2 R2 Bw' Fw 3Rw2 Uw Fw R F2 Lw U2 Bw2 Uw' B Uw' Lw 3Uw2 Dw F Uw' F2 L2 3Fw Dw' Bw2 Rw2 Lw' Dw' F' Lw' B' Rw' D' Dw2 Fw Lw 3Fw Dw' D2 F' D2 3Fw2 Fw'",
        	"R' Fw' D Fw Uw' U2 F' L2 Rw' 3Rw' Lw' Dw U2 Lw' 3Rw' R' B2 3Uw2 Uw2 F' 3Uw Rw2 F2 R2 Lw 3Uw2 Uw' 3Fw2 Fw2 D' 3Uw Fw 3Rw' Fw Dw2 3Rw2 Lw L F Lw' B2 Uw' 3Fw2 Dw D Lw' F2 R Bw' Rw' Fw' 3Rw' Fw 3Rw' F2 Fw 3Fw2 D2 F L B2 Lw' L2 D2 3Fw' 3Uw Uw Rw' Uw F' Rw' L' U Fw2 U Uw 3Uw F R 3Rw'",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "7x7x7"} )],
		Scrambles: []string{
			"3Lw' 3Uw2 3Lw2 Lw' Dw' Uw2 Fw D' 3Uw 3Rw2 Lw2 R2 3Fw' Uw2 B' Dw' Fw' Rw2 R Lw Bw' 3Uw2 3Rw Bw2 3Lw2 3Fw Fw D' B' Uw2 3Bw D' Rw' L D' Uw 3Uw' Fw2 L Dw Fw' 3Fw R U' D2 3Uw Bw' Lw' Rw' 3Bw2 Uw' Dw 3Uw 3Lw2 3Rw2 Bw' Dw2 3Uw2 B2 3Fw Dw2 Uw' 3Bw 3Lw 3Rw2 F Lw' 3Rw2 Dw' Fw' 3Bw2 Rw R D2 B' U Lw2 3Lw 3Fw 3Bw' U2 F D Uw2 3Bw Fw B 3Rw 3Dw' 3Bw 3Dw 3Lw' D' Fw2 B2 3Uw2 Bw' B 3Fw L2",
			"3Fw' Uw2 3Uw Fw2 F2 3Dw2 Lw 3Rw 3Uw' Uw' B2 R 3Fw2 3Bw' Lw' U2 3Fw2 3Rw2 Uw Bw U Lw2 3Fw Dw 3Lw' B 3Lw2 Lw 3Bw' 3Fw2 Bw D2 3Lw' U R D Bw2 U' Bw 3Fw' 3Bw 3Rw Rw2 3Bw D Lw' D' 3Dw2 R 3Uw2 3Rw2 3Bw U2 Lw2 3Fw B2 Rw R B2 3Dw' L' Uw' 3Bw2 3Dw2 Fw2 Bw' 3Bw2 F 3Fw' R2 Uw2 Fw2 3Lw 3Bw 3Lw2 U R' 3Rw' 3Lw2 Rw2 D2 U2 Dw' Bw' 3Uw2 B Bw 3Dw' 3Rw 3Bw2 3Fw2 D2 R' D' Uw2 Fw2 Uw Rw 3Lw' B",
			"Dw D' Bw2 3Lw 3Dw' Dw 3Rw Rw' Dw' 3Rw' 3Fw' 3Rw' Dw2 Fw Rw2 3Fw' Rw Lw R2 3Lw' Dw' Uw 3Lw L Dw 3Rw Rw 3Uw2 Uw2 D' 3Rw2 3Dw2 Dw' 3Bw2 U2 Uw D' Rw' 3Bw2 3Fw 3Uw' U' Fw2 3Lw2 3Dw' 3Uw Bw2 3Uw 3Bw' Rw' 3Lw Dw Rw2 D 3Uw Bw2 Uw2 3Fw' U 3Uw2 Dw2 3Bw Rw2 3Fw' 3Bw R2 Bw' Lw' 3Uw 3Lw 3Bw2 Fw' 3Lw' Bw' Lw2 L 3Lw2 Rw U2 D2 Rw2 3Bw2 D2 3Uw' 3Rw2 Fw' D Uw 3Uw U 3Lw 3Dw 3Rw 3Lw Fw' Rw2 3Fw2 3Uw2 U B2",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "3BLD"} )],
		Scrambles: []string{
			"F2 L2 D' F2 D2 L2 U B2 U F2 U B' U L' B R2 D2 F R' D2 Rw' Uw'",
			"L2 B D2 B2 U2 L2 B' R2 F' D2 F L2 D' B' F' U' L R' B U' F Rw Uw",
			"R' D L2 B2 U2 R2 B2 L' U2 L2 B2 D2 R' F' L U2 R2 B F D' Rw",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "FMC"} )],
		Scrambles: []string{
			"R' U' F D2 U2 L2 B U2 B2 L2 D2 U2 F' R B2 F R U F U' B2 F2 R' U' F",
			"R' U' F D2 F2 D2 B2 L2 F2 L2 D F2 B' R F' L F U2 F' D R' U' F",
			"R' U' F L U2 R' D2 L B2 R F2 L B2 U2 F2 B' L F' R D F U L2 B2 R' U' F",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "OH"} )],
		Scrambles: []string{
			"L2 D U2 R2 F' D2 L2 F R2 D2 L2 D2 L' B2 U R F' D' L2",
			"D2 L' B2 D2 L2 D2 U2 F2 U2 L2 R' B D B F R U2 R D'",
			"L2 D2 F D2 F U2 F' R2 F L2 F U F' U R' B' D2 B L F'",
			"D2 U F U2 B D2 F2 R2 F R2 D2 B L' U' R' F' D' L2 D2 U'",
			"B L U2 B2 D' R2 D R2 U2 F2 D L2 U2 F' D' F2 R B2 R U",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "Clock"} )],
		Scrambles: []string{
			"UR5+ DR4- DL3+ UL1+ U3+ R4+ D5+ L0+ ALL5- y2 U5- R5- D3+ L2- ALL4+",
			"UR2- DR3- DL2- UL4- U4+ R2- D4- L0+ ALL3+ y2 U3+ R5- D5+ L1+ ALL5-",
			"UR5+ DR1+ DL3- UL2- U1+ R0+ D6+ L4+ ALL2+ y2 U3- R1- D5- L6+ ALL5+",
			"UR5+ DR3- DL5+ UL0+ U1+ R1- D4- L6+ ALL4+ y2 U5+ R6+ D0+ L4- ALL1-",
			"UR5+ DR1+ DL5+ UL1+ U2- R5- D4- L3- ALL3- y2 U0+ R2+ D5+ L3- ALL1-",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "Mega"} )],
		Scrambles: []string{
			"R-- D-- R-- D-- R-- D-- R++ D++ R-- D-- U'\n  R++ D-- R-- D-- R-- D++ R++ D-- R++ D-- U'\n  R++ D++ R++ D++ R++ D++ R++ D-- R++ D++ U\n  R++ D++ R-- D-- R++ D-- R++ D++ R++ D++ U \n  R++ D++ R++ D-- R-- D++ R-- D++ R-- D++ U \n  R++ D-- R-- D++ R-- D++ R-- D-- R++ D-- U'\n  R++ D++ R-- D++ R++ D++ R-- D-- R-- D-- U'\n",
			"R-- D++ R-- D-- R++ D-- R++ D++ R++ D-- U'\n  R++ D++ R++ D-- R-- D-- R-- D++ R-- D++ U \n  R-- D-- R-- D-- R-- D-- R++ D++ R++ D-- U'\n  R++ D-- R++ D++ R-- D++ R-- D-- R++ D-- U'\n  R++ D-- R-- D-- R++ D-- R-- D-- R++ D-- U'\n  R++ D-- R++ D-- R-- D++ R-- D-- R++ D++ U \n  R++ D-- R++ D-- R++ D-- R++ D-- R-- D-- U'\n",
			"R-- D++ R-- D++ R-- D-- R-- D-- R-- D++ U \n  R++ D-- R++ D-- R-- D-- R-- D++ R++ D++ U \n  R-- D-- R-- D-- R++ D++ R-- D-- R-- D++ U \n  R-- D++ R-- D-- R++ D-- R++ D-- R-- D-- U'\n  R-- D++ R++ D-- R++ D++ R-- D++ R++ D++ U \n  R-- D++ R++ D++ R++ D-- R++ D++ R++ D-- U'\n  R++ D++ R++ D++ R++ D++ R++ D++ R-- D-- U'\n",
			"R-- D++ R++ D-- R++ D++ R-- D++ R-- D-- U'\n  R-- D++ R++ D-- R-- D++ R-- D-- R-- D-- U'\n  R-- D-- R-- D-- R-- D-- R++ D-- R-- D-- U'\n  R-- D++ R++ D++ R++ D++ R-- D-- R-- D++ U \n  R++ D++ R-- D-- R-- D-- R++ D-- R++ D++ U \n  R-- D++ R++ D++ R-- D++ R++ D++ R++ D-- U'\n  R++ D++ R++ D++ R-- D-- R-- D-- R++ D-- U'\n",
			"R++ D-- R-- D++ R-- D-- R-- D-- R++ D-- U'\n  R-- D-- R-- D++ R-- D-- R-- D-- R++ D-- U'\n  R-- D++ R++ D-- R++ D-- R++ D++ R-- D-- U'\n  R-- D-- R-- D++ R-- D++ R++ D++ R-- D++ U \n  R++ D++ R-- D++ R++ D++ R-- D-- R++ D++ U \n  R-- D-- R++ D++ R-- D++ R++ D-- R-- D-- U'\n  R++ D-- R++ D++ R-- D-- R++ D++ R++ D-- U'\n",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "Pyra"} )],
		Scrambles: []string{
			"R U B U' B' L' U R' l' r' u'",
			"B' U' R B R' L B' U r u'",
			"L' U L R' L B' R B' l b u",
			"B' R U B L U' B' L l' r' b' u",
			"B L U' L B' R' U' B' l r b' u",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "Skewb"} )],
		Scrambles: []string{
			"U R L B L U' R' L' R'",
			"U L' R U' R' B R' B U' R'",
			"B' L R L U' B R' L B'",
			"R B' R' B R' B L' R L'",
			"L' U' B' U B' R L' U",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "Sq-1"} )],
		Scrambles: []string{
			"(0,5)/ (-2,-5)/ (0,3)/ (-3,0)/ (-1,-4)/ (3,-3)/ (4,0)/ (3,0)/ (3,0)/ (5,0)/ (2,0)/ (2,-5)/ (2,0)/",
			"(0,-4)/ (-3,3)/ (1,-2)/ (3,0)/ (6,0)/ (-3,-4)/ (3,0)/ (-4,0)/ (1,0)/ (2,0)/ (-1,0)/ (-4,0)/ (0,-2)/ (5,0)/",
			"(-5,0)/ (0,3)/ (-1,2)/ (4,-2)/ (-1,-4)/ (6,0)/ (-2,0)/ (-3,-3)/ (-5,0)/ (0,-2)/ (0,-2)/ (5,0)/ (2,0)/ (-3,0)",
			"(3,-1)/ (3,0)/ (-5,-2)/ (3,0)/ (-4,-1)/ (3,0)/ (-5,0)/ (0,-3)/ (0,-2)/ (3,0)/ (0,-4)",
			"(1,3)/ (5,-4)/ (0,-3)/ (4,-2)/ (5,-4)/ (6,-2)`/` (-3,0)/ (-2,-3)/ (-2,0)/ (0,-1)/ (2,0)/ (-1,0)",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "4BLD"} )],
		Scrambles: []string{
			"L' U' B R2 D2 B D2 B L2 D2 F' R2 D2 F2 R' B2 L2 D2 B U' F Uw2 Fw2 D Fw2 R U R2 Fw2 L F2 B2 L' D2 Fw Uw2 L F Fw D F2 B' Uw R Rw' Uw L' z y'",
			"F' B2 L2 B' D2 B2 U2 R2 F' U2 F2 D L2 D F' D' R2 L U' Uw2 B' Rw2 U R2 Fw2 R2 D' F' R2 D' Rw2 B' Rw' F2 Rw F' R' L2 B' Fw Rw' B2 U Rw'",
			"R' F' R2 F' B2 U2 R F' D' L2 B2 D L2 D2 R2 L2 U' L2 D2 R' Uw2 Rw2 U Fw2 L' U' F2 Uw2 L' F2 Fw2 L' Rw2 Fw R2 U F2 R' Fw Uw Fw' Uw2 Rw' F Fw2 Rw2  x' y2",
		},
    },
	{
		Event: events[slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == "5BLD"} )],
		Scrambles: []string{
			"D2 Rw Uw2 Dw' U' F' Uw F' B Rw2 Dw2 Fw' R2 F2 Fw' Bw2 Rw2 R2 Bw R U2 Dw Fw R U2 Lw L2 Rw' Uw2 F' L' F2 Uw2 Fw U Fw R2 Uw' U' F2 Dw B Rw Fw' D2 Bw Fw' U' B L Bw F U2 Bw Dw2 D Bw Uw2 R B2 3Rw 3Uw2",
			"D2 L2 Fw' Bw' B2 U2 Rw2 L2 D' U2 Fw' D L' Fw' L2 Bw2 L F2 Fw R Lw' U' Rw2 B2 Uw2 Rw' F D2 L2 R' Uw2 Rw' Dw2 F' Uw2 Lw2 Uw Lw Uw2 L Uw D2 F2 Lw B2 L Rw' Dw2 D2 Rw2 Dw2 D Rw2 Dw' U2 Lw' Fw' D Rw' U' 3Fw' 3Uw",
			"Fw2 B Rw Fw F Bw2 Lw2 Uw Fw F' U Rw Bw B2 U2 Fw' Bw' L' Lw' D U2 F2 Lw2 Fw2 L Dw2 Lw Dw2 B2 U' Dw Uw' R D2 U2 B' Dw' Lw' U2 Fw2 Rw Lw' Fw' B L D' F' R' U F2 R' D Rw2 R Bw' R Dw2 Bw2 B2 Lw2 3Rw' 3Uw2",
		},
    },
}

type CompetitionData struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Startdate time.Time `json:"startdate"`
	Enddate time.Time `json:"enddate"`
	Events []CompetitionEvent `json:"events"`
	Scrambles []ScrambleSet `json:"scrambles"`
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

var competitions = allCompetitionData()

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
	router.GET("/api/results/edit/:uname/:cname/:ename", getResultsQuery)
	router.GET("/api/results/:id/:event", getResultsByIdAndEvent)
	router.POST("api/results/save", postResults);
	router.GET("/api/events", getEvents)
	router.GET("/api/competitions/:filter", getFilteredCompetitions)
	router.GET("/api/competition/:id", getCompetitionById)
	router.POST("/api/competition", postCompetition)
	router.PUT("/api/competition", putCompetition)
	router.GET("/api/users/manage-roles", getManageRolesUsers)
	router.PUT("/api/users/manage-roles", putManageRolesUsers)

	router.Run("localhost:8080")
}

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "pong")
}

func getResultsQuery(c *gin.Context) {
	competitorName := c.Param("uname")
	competitionName := c.Param("cname")
	eventName := c.Param("ename")

	r := make([]ResultEntry, 0, len(results))

	for _, result := range results {
		if (result.Username == competitorName || competitorName == "_") && (result.Competitionname == competitionName || competitionName == "_") && result.Eventname == eventName {
			r = append(r, result)
		}
	}

	c.IndentedJSON(http.StatusOK, r)
}

func getResultsByIdAndEvent(c *gin.Context) {
	eventDisplayname := c.Param("event")
	cid := c.Param("id")

	resultsIdx := slices.IndexFunc(results, func (e ResultEntry) bool { return e.Eventname == eventDisplayname })
	var r ResultEntry

	if resultsIdx == -1 {
		eventIdx := slices.IndexFunc(events, func (e CompetitionEvent) bool { return e.Displayname == eventDisplayname })
		event := events[eventIdx]

		r = ResultEntry{
			Id: rand.Int(),
			Userid: 1,
			Username: "Janko Hrasko",
			Competitionid: cid,
			Competitionname: cid,
			Eventid: event.Id,
			Eventname: event.Displayname,
			Iconcode: event.Iconcode,
			Format: event.Format,
			Solve1: "",
			Solve2: "",
			Solve3: "",
			Solve4: "",
			Solve5: "",
			Comment: "",
			Status: approvedResultsStatus,
		}

		results = append(results, r)
	} else {
		r = results[resultsIdx]
	}

	c.IndentedJSON(http.StatusOK, r)
}

func postResults(c *gin.Context) {
	var resultEntry ResultEntry

	if err := c.BindJSON(&resultEntry); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "");
		return;
	}
	
	resultsIdx := slices.IndexFunc(results, func (r ResultEntry) bool { return r.Id == resultEntry.Id })
	
	if resultsIdx == -1 {
		resultEntry.Id = rand.Int()
		results = append(results, resultEntry)
	} else {
		results[resultsIdx] = resultEntry
	}

	c.IndentedJSON(http.StatusCreated, resultEntry)
}

func getEvents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, events);
}

func getFilteredCompetitions(c *gin.Context) {
	filter := c.Param("filter")
	
	result := make([]CompetitionData, 0);

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
	
	idx := slices.IndexFunc(competitions, func (c CompetitionData) bool { return c.Id == id })
	if idx == -1 {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	
	result := competitions[idx]
	c.IndentedJSON(http.StatusOK, result)
}

func postCompetition(c *gin.Context) {
	var competition CompetitionData

	if err := c.BindJSON(&competition); err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("competition", competition)

	competitions = append(competitions, competition)

	c.IndentedJSON(http.StatusCreated, competition)
}

func putCompetition(c *gin.Context) {
	var competition CompetitionData

	if err := c.BindJSON(&competition); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	idx := slices.IndexFunc(competitions, func (c CompetitionData) bool { return c.Id == competition.Id })
	if idx == -1 {
		c.IndentedJSON(http.StatusInternalServerError, "Invalid competetion.")
		return
	}

	competitions[idx] = competition

	c.IndentedJSON(http.StatusCreated, competition)
}

func getManageRolesUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, manageRolesUsers)
}

func putManageRolesUsers(c *gin.Context) {
	var newManageRolesUsers []ManageRolesUser

	if err := c.BindJSON(&newManageRolesUsers); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error());
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