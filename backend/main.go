package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/alexsergivan/transliterator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

type CompetitionEvent struct {
	Id int `json:"id"`
	Displayname string `json:"displayname"`
	Format string `json:"format"`
	Iconcode string `json:"iconcode"`
	Puzzlecode string `json:"puzzlecode"`
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

type ScrambleSet struct {
	Event CompetitionEvent `json:"event"`
	Scrambles []string `json:"scrambles"`
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

type ManageRolesUser struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Isadmin bool `json:"isadmin"`
}

type AuthorizationInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	WcaId string `json:"wcaid"`
	AvatarUrl string `json:"avatarUrl"`
	IsAdmin bool `json:"isadmin"`
}

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	CountryId string `json:"country_id"`
	Sex string `json:"sex"`
	WcaId string `json:"wcaid"`
	IsAdmin bool `json:"isadmin"`
	Url string `json:"url"`
	AvatarUrl string `json:"avatarurl"`
}


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
	
	results := api_v1.Group("/results", authMiddleWare(db, envMap))
	{
		results.GET("/edit/:uname/:cid/:eid", adminMiddleWare(), getResultsQuery(db))
		results.GET("/compete/:cid/:eid", getResultsByIdAndEvent(db))
		results.POST("/save", postResults(db))
		results.POST("/save-validation", adminMiddleWare(), postResultsValidation(db))
	}

	events := api_v1.Group("/events")
	{
		events.GET("/", getEvents(db))
	}

	competitions := api_v1.Group("/competitions")
	{
		competitions.GET("/filter/:filter", getFilteredCompetitions(db))
		competitions.GET("/id/:id", getCompetitionById(db))
		competitions.POST("/", authMiddleWare(db, envMap), adminMiddleWare(), postCompetition(db))
		competitions.PUT("/", authMiddleWare(db, envMap), adminMiddleWare(), putCompetition(db))
	}

	users := api_v1.Group("/users")
	{
		users.GET("/manage-roles", authMiddleWare(db, envMap), getManageRolesUsers(db))
		users.PUT("/manage-roles", authMiddleWare(db, envMap), adminMiddleWare(), putManageRolesUsers(db))
	}

	router.GET("/api/auth/admin", authMiddleWare(db, envMap), adminMiddleWare(), func(c *gin.Context) { c.IndentedJSON(http.StatusAccepted, "authorized")});

	router.POST("/api/login", postLogIn(db, envMap))

	router.Run("localhost:8080")
}

func adminMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		isadmin := c.MustGet("isadmin").(bool)
		if !isadmin {
			c.IndentedJSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

type AuthDetails struct {
	UserId int `json:"userid"`
	ExpiresIn int64 `json:"expiresin"`
}

func createToken(userid int, secretKey string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
			"userid": userid, 
			"exp": time.Now().Add(time.Hour * 12).Unix(),
        })

    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
    	return "", err
    }

 	return tokenString, nil
}

func extractPayload(tokenString *jwt.Token) (AuthDetails, error) {
	claims, ok := tokenString.Claims.(jwt.MapClaims)

	if !ok { return AuthDetails{}, fmt.Errorf("extracting payload from token failed") }

	var authDetails AuthDetails
	uidFloat, ok := claims["userid"].(float64)
	if !ok { return AuthDetails{}, fmt.Errorf("failed to parse userid") }
	authDetails.UserId = int(uidFloat)

	return authDetails, nil
}

func verifyJWTToken(tokenString string, secretKey string) (AuthDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return []byte(secretKey), nil
	})
   
	if err != nil { return AuthDetails{}, err }
	if !token.Valid { return AuthDetails{}, fmt.Errorf("invalid token") }
   
	authDetails, err := extractPayload(token)
	if err != nil { return AuthDetails{}, err }

	return authDetails, nil
 }

func getAuthDetailsFromHeader(c *gin.Context, secretKey string) (AuthDetails, error) {
	headers := c.Request.Header["Authorization"]
	if len(headers) <= 0 { return AuthDetails{}, fmt.Errorf("auth header missing") }

	header := strings.Split(headers[0], " ")
	if len(header) < 2 || header[0] != "Bearer" { return AuthDetails{}, fmt.Errorf("bad auth header") }

	authDetails, err := verifyJWTToken(header[1], secretKey)
	if err != nil { return AuthDetails{}, err }

	return authDetails, nil
}

func authMiddleWare(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authDetails, err := getAuthDetailsFromHeader(c, envMap["JWT_SECRET_KEY"])
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, err)
			c.Abort()
			return
		}
		
		user, err := getUserById(db, authDetails.UserId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.Set("uid", user.Id)
		c.Set("isadmin", user.IsAdmin)

		c.Next()
	}
}

func getAuthInfo(code string, envMap map[string]string) (AuthorizationInfo, error) {
	res, err := http.PostForm(envMap["WCA_TOKEN_URL"], url.Values{
		"grant_type": {"authorization_code"},
		"client_id": {envMap["WCA_CLIENT_ID"]},
		"client_secret": {envMap["WCA_CLIENT_SECRET"]},
		"code": {code},
		"redirect_uri": {envMap["WCA_REDIRECT_URI"]},
	})
	if err != nil || res.StatusCode != http.StatusOK { return AuthorizationInfo{}, err }
	defer res.Body.Close()

	var authInfo AuthorizationInfo
	err = json.NewDecoder(res.Body).Decode(&authInfo)
	if err != nil { return AuthorizationInfo{}, err }

	return authInfo, nil
}

func (u *User) Exists(db *pgxpool.Pool) (bool, error) {
	rows, err := db.Query(context.Background(), `SELECT u.user_id, u.isadmin FROM users u WHERE u.wcaid = $1;`, u.WcaId)
	if err != nil { return false, err }

	found := false
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.IsAdmin)
		if err != nil { return false, err }
		found = true
	}

	return found, nil
}

func (u *User) Update(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `UPDATE users SET name = $1, country_id = $2, sex = $3, url = $4, avatarurl = $5, isadmin = $6, timestamp = CURRENT_TIMESTAMP WHERE wcaid = $7;`, u.Name, u.CountryId, u.Sex, u.Url, u.AvatarUrl, u.IsAdmin, u.WcaId)
	if err != nil { return err }

	return nil
}

func (u *User) Insert(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `INSERT INTO users (name, country_id, sex, url, avatarurl, wcaid, isadmin) VALUES ($1,$2,$3,$4,$5,$6,false);`, u.Name, u.CountryId, u.Sex, u.Url, u.AvatarUrl, u.WcaId)
	if err != nil { return err }

	return nil
}

func getUserInfoFromWCA(authInfo *AuthorizationInfo, envMap map[string]string) (User, error) {
	bearer := "Bearer " + authInfo.AccessToken
	req, err := http.NewRequest("GET", envMap["WCA_API_ME_URL"], nil)
	if err != nil { return User{}, err }

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK { return User{}, err }
	defer res.Body.Close()
	type Country struct { Id string `json:"id"`}
	type Avatar struct { Url string `json:"url"`}
	type ME struct {
		Name string `json:"name"`
		WcaId string `json:"wca_id"`
		Sex string `json:"gender"`
		Url string `json:"url"`
		Country Country `json:"country"`
		Avatar Avatar `json:"avatar"`
	}
	type WCAApiMe struct { Me ME `json:"me"` }

	var apiMe WCAApiMe
	err = json.NewDecoder(res.Body).Decode(&apiMe)
	if err != nil { return User{}, err }

	user := User{}
	user.Name = apiMe.Me.Name
	user.CountryId = apiMe.Me.Country.Id
	user.Sex = apiMe.Me.Sex
	user.WcaId = apiMe.Me.WcaId
	user.IsAdmin = false
	user.Url = apiMe.Me.Url
	user.AvatarUrl = apiMe.Me.Avatar.Url

	return user, nil
}

func postLogIn(db *pgxpool.Pool, envMap map[string]string) gin.HandlerFunc {
	return func (c *gin.Context) {
		reqBodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		code := string(reqBodyBytes)
		authInfo, err := getAuthInfo(code, envMap)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		user, err := getUserInfoFromWCA(&authInfo, envMap)
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
		authInfo.AccessToken, err = createToken(user.Id, envMap["JWT_SECRET_KEY"])
		authInfo.IsAdmin = user.IsAdmin
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		
		c.IndentedJSON(http.StatusOK, authInfo)
	}
}

func getResultsQuery(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		userName := c.Param("uname")
		competitionId := c.Param("cid")
		eventId, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		
		var resultEntries []ResultEntry

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

				resultEntry, err := getResultEntryById(db, resultEntryId)
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

				resultEntry, err := getResultEntryById(db, resultEntryId)
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

				resultEntry, err := getResultEntryById(db, resultEntryId)
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

				resultEntry, err := getResultEntryById(db, resultEntryId)
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

func getAvailableEvents(db *pgxpool.Pool) ([]CompetitionEvent, error) {
	rows, err := db.Query(context.Background(), "SELECT e.event_id, e.displayname, e.format, e.iconcode, e.puzzlecode FROM events e ORDER BY e.event_id;");
	if err != nil { return []CompetitionEvent{}, err }

	var events []CompetitionEvent
	for rows.Next() {
		var event CompetitionEvent
		err = rows.Scan(&event.Id, &event.Displayname, &event.Format, &event.Iconcode, &event.Puzzlecode)
		if err != nil { return []CompetitionEvent{}, err }
		events = append(events, event)
	}

	return events, nil
}

func getEventById(db *pgxpool.Pool, eventId int) (CompetitionEvent, error) {
	rows, err := db.Query(context.Background(), "SELECT e.event_id, e.displayname, e.format, e.iconcode, e.puzzlecode FROM events e WHERE e.event_id = $1;", eventId);
	if err != nil { return CompetitionEvent{}, err }

	var event CompetitionEvent
	found := false
	for rows.Next() {
		err = rows.Scan(&event.Id, &event.Displayname, &event.Format, &event.Iconcode, &event.Puzzlecode)
		if err != nil { return CompetitionEvent{}, err }
		found = true
	}

	if !found { return CompetitionEvent{}, fmt.Errorf("event not found by id") }

	return event, nil
}

func getCompetitionByIdObject(db *pgxpool.Pool, id string) (CompetitionData, error) {
	rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c WHERE c.competition_id = $1;`, id)
	if err != nil { return CompetitionData{}, err }

	var competition CompetitionData
	found := false
	for rows.Next() {
		err = rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
		if err != nil { return CompetitionData{}, err }
		found = true
	}

	if !found { return CompetitionData{}, err }

	return competition, nil
}

func getResultEntry(db *pgxpool.Pool, competitorId int, competitionId string, eventId int) (ResultEntry, error) {
	rows, err := db.Query(context.Background(), `SELECT re.result_id, re.competition_id, re.user_id, re.event_id, re.solve1, re.solve2, re.solve3, re.solve4, re.solve5, re.comment, re.status_id FROM results re WHERE re.user_id = $1 AND re.competition_id = $2 AND re.event_id = $3;`, competitorId, competitionId, eventId)
	if err != nil { return ResultEntry{}, err }

	var resultEntry ResultEntry
	found := false
	for rows.Next() {
		err = rows.Scan(&resultEntry.Id, &resultEntry.Competitionid, &resultEntry.Userid, &resultEntry.Eventid, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Comment, &resultEntry.Status.Id)
		if err != nil { return ResultEntry{}, err }
		found = true
	}

	if !found { return ResultEntry{}, fmt.Errorf("not found") }

	return resultEntry, nil
}

func getUserById(db *pgxpool.Pool, uid int) (User, error) {
	rows, err := db.Query(context.Background(), `SELECT u.user_id, u.name, u.country_id, u.sex, u.wcaid, u.isadmin, u.url, u.avatarurl FROM users u WHERE u.user_id = $1;`, uid);
	if err != nil { return User{}, err }

	var user User
	found := false
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.CountryId, &user.Sex, &user.WcaId, &user.IsAdmin, &user.Url, &user.AvatarUrl)
		if err != nil { return User{}, err }
		found = true
	}

	if !found { return User{}, err }
	
	return user, nil
}

func getResultsByIdAndEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		eventId, err := strconv.Atoi(c.Param("eid"))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		competitionId := c.Param("cid")
		userId := c.MustGet("uid").(int)

		user, err := getUserById(db, userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		event, err := getEventById(db, eventId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		competition, err := getCompetitionByIdObject(db, competitionId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return	
		}

		resultEntry, err := getResultEntry(db, userId, competitionId, eventId)

		if err != nil {
			if err.Error() != "not found" {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			} else {
				approvedResultsStatus, err := getResultsStatus(db, 3)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return	
				}

				resultEntry = ResultEntry{
					Id: rand.Int(),
					Userid: userId,
					Username: user.Name,
					Competitionid: competitionId,
					Competitionname: competition.Name,
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

				err = resultEntry.Insert(db)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
				}
			}
		} else {
			currentStatus, err := getResultsStatus(db, resultEntry.Status.Id)
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

func tryParseSolve(s string) (float64, error) {
	if !strings.Contains(s, ".") {
		s += ".00"
	}

	split := strings.Split(s, ".")
	
	wholePart := strings.Split(split[0], ":")
	decimalPart, err := strconv.ParseFloat(split[1], 64)

	if err != nil {
		return 0, fmt.Errorf("invalid time or DNF/DNS");
	}

	res := decimalPart * 10 // to milliseconds

	var add float64
	slices.Reverse(wholePart)
	if len(wholePart) > 0 { 
		add, err = strconv.ParseFloat(wholePart[0], 64)
		res += add * 1000
	}

	if len(wholePart) > 1 { 
		add, err = strconv.ParseFloat(wholePart[1], 64)
		res += 60 * add * 1000
	}

	if len(wholePart) > 2 { 
		add, err = strconv.ParseFloat(wholePart[2], 64)
		res += 60 * 60 * add * 1000
	}

	if len(wholePart) > 3 { 
		add, err = strconv.ParseFloat(wholePart[3], 64)
		res += 24 * 60 * 60 * add * 1000
	}

	if err != nil {
		return 0, fmt.Errorf("invalid time in formatted")
	}

	return res, nil
}

func compareSolves(t1 *float64, s2 string) {
	t2, err := tryParseSolve(s2)
	if err == nil && *t1 - t2 > 1e-9 {
		*t1 = t2;
	}
}

func (r *ResultEntry) single() float64 {
	res := math.MaxFloat64

	compareSolves(&res, r.Solve1)
	compareSolves(&res, r.Solve2)
	compareSolves(&res, r.Solve3)
	compareSolves(&res, r.Solve4)
	compareSolves(&res, r.Solve5)

	return res
}

func (r *ResultEntry) getSolvesFromResultEntry() []float64 {
	values := make([]float64, 0)

	t1, err1 := tryParseSolve(r.Solve1)
	if err1 != nil {
		values = append(values, math.MaxFloat64)
	} else {
		values = append(values, t1)
	}

	t2, err2 := tryParseSolve(r.Solve2)
	if err2 != nil {
		values = append(values, math.MaxFloat64)
	} else {
		values = append(values, t2)
	}
	
	t3, err3 := tryParseSolve(r.Solve3)
	if err3 != nil {
		values = append(values, math.MaxFloat64)
	} else {
		values = append(values, t3)
	}

	t4, err4 := tryParseSolve(r.Solve4)
	if err4 != nil {
		values = append(values, math.MaxFloat64)
	} else {
		values = append(values, t4)
	}

	t5, err5 := tryParseSolve(r.Solve5)
	if err5 != nil {
		values = append(values, math.MaxFloat64)
	} else {
		values = append(values, t5)
	}

	return values;
}

func (r *ResultEntry) average(noOfSolves int) float64 {
	solves := r.getSolvesFromResultEntry()
	slices.Sort(solves)

	sum := 0.
	cntBad := 0

	for idx, solve := range solves {
		if idx >= noOfSolves {
			break
		}

		if solve == math.MaxFloat64 {
			cntBad++
			if (noOfSolves == 5 && cntBad > 1) || (noOfSolves == 3 && cntBad > 0) {
				return math.MaxFloat64
			}
		}

		if noOfSolves == 3 || (noOfSolves == 5 && idx > 0 && idx < 4) {
			sum += solve
		}
	}

	return float64(sum) / float64(3)
}

func getWorldRecords(eventName string) (float64, float64, error) {
	c := colly.NewCollector()

	single, average := math.MaxFloat64, math.MaxFloat64
	var err error

	c.OnHTML("div#results-list h2 a", func(e *colly.HTMLElement) {
		hrefSplit := strings.Split(e.Attr("href"), "/")
		if len(hrefSplit) > 3 && hrefSplit[3] == eventName {
			parentH2 := e.DOM.Parent()
			nextTable := parentH2.Next()
			singleTd := nextTable.Find("td.result").First()

			single, err = tryParseSolve(strings.Trim(singleTd.Text(), " "))

			if eventName != "333mbf" {
				averageTd := singleTd.Parent().Next().Find("td.result").First()
				average, err = tryParseSolve(strings.Trim(averageTd.Text(), " "))
			}
		}
	})

	c.OnError(func(_ *colly.Response, er error) {
		err = er
	})

	err = c.Visit("https://www.worldcubeassociation.org/results/records")
	if err != nil {
		return 0, 0, err
	}

	return single, average, nil
}

func getNoOfSolves(format string) (int, error) {
	match := regexp.MustCompile(`\d+$`).FindString(format)
	res, err := strconv.Atoi(match)

	if err != nil {
		return 0, fmt.Errorf("did not find a number at the end of format")
	}

	return res, nil
}

func (r *ResultEntry) isSuspicous() bool {
	noOfSolves, err := getNoOfSolves(r.Format)
	if err != nil {
		return false
	}

	curSingle, curAverage := r.single(), r.average(noOfSolves)

	recSingle, recAverage, err := getWorldRecords(r.Iconcode)
	if err != nil { return false }

	return recSingle - curSingle > 1e-9 || recAverage - curAverage > 1e-9;
}

func getResultsStatus(db *pgxpool.Pool, statusId int) (ResultsStatus, error) {
	rows, err := db.Query(context.Background(), `SELECT rs.results_status_id, rs.approvalfinished, rs.approved, rs.visible, rs.displayname FROM results_status rs WHERE rs.results_status_id = $1;`, statusId);
	if err != nil {
		return ResultsStatus{}, err
	}

	var resultsStatus ResultsStatus
	found := false
	for rows.Next() {
		err = rows.Scan(&resultsStatus.Id, &resultsStatus.ApprovalFinished, &resultsStatus.Approved, &resultsStatus.Visible, &resultsStatus.Displayname)
		if err != nil {
			return ResultsStatus{}, err
		}
		found = true
	}

	if !found {
		return ResultsStatus{}, err
	}

	return resultsStatus, nil
}

func (r *ResultEntry) Insert(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `INSERT INTO results (competition_id, user_id, event_id, solve1, solve2, solve3, solve4, solve5, comment, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`, r.Competitionid, 1, r.Eventid, r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5, r.Comment, r.Status.Id)
	if err != nil {
		return err
	}

	return nil;
}

func (r *ResultEntry) Validate(db *pgxpool.Pool) (error) {

	var err error
	if (r.isSuspicous()) {
		r.Status, err = getResultsStatus(db, 1) // waitingForApproval
		if err != nil { return err }
	} else {
		r.Status, err = getResultsStatus(db, 3) // approved
		if err != nil { return err }
	}

	return nil
}

func (r *ResultEntry) Update(db *pgxpool.Pool, valid ...bool) error {
	if r.Solve1 != "" || r.Solve2 != "" || r.Solve3 != "" || r.Solve4 != "" || r.Solve5 != "" {
		r.Solve1 = "DNS"
		r.Solve2 = "DNS"
		r.Solve3 = "DNS"
		r.Solve4 = "DNS"
		r.Solve5 = "DNS"
	}

	if len(valid) == 0 || (len(valid) > 0 && !valid[0]) {
		err := r.Validate(db)
		if err != nil { return err }
	}
	
	_, err := db.Exec(context.Background(), `UPDATE results SET solve1 = $1, solve2 = $2, solve3 = $3, solve4 = $4, solve5 = $5, comment = $6, status_id = $7, timestamp = CURRENT_TIMESTAMP;`, r.Solve1, r.Solve2, r.Solve3, r.Solve4, r.Solve5, r.Comment, r.Status.Id)
	if err != nil { return err }

	return nil;
}

func postResults(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		var resultEntry ResultEntry
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

func getResultEntryById(db *pgxpool.Pool, resultId int) (ResultEntry, error) {
	rows, err := db.Query(context.Background(), `SELECT re.result_id, re.competition_id, re.user_id, re.event_id, re.solve1, re.solve2, re.solve3, re.solve4, re.solve5, re.comment, re.status_id, c.name, e.displayname, rs.approvalfinished, rs.approved, rs.visible, rs.displayname, u.name, e.format, e.iconcode FROM results re JOIN competitions c ON c.competition_id = re.competition_id JOIN events e ON e.event_id = re.event_id JOIN results_status rs ON results_status_id = re.status_id JOIN users u ON u.user_id = re.user_id WHERE re.result_id = $1;`, resultId)
	if err != nil { return ResultEntry{}, err }

	var resultEntry ResultEntry
	found := false
	for rows.Next() {
		err = rows.Scan(&resultEntry.Id, &resultEntry.Competitionid, &resultEntry.Userid, &resultEntry.Eventid, &resultEntry.Solve1, &resultEntry.Solve2, &resultEntry.Solve3, &resultEntry.Solve4, &resultEntry.Solve5, &resultEntry.Comment, &resultEntry.Status.Id, &resultEntry.Competitionname, &resultEntry.Eventname, &resultEntry.Status.ApprovalFinished, &resultEntry.Status.Approved, &resultEntry.Status.Visible, &resultEntry.Status.Displayname, &resultEntry.Username, &resultEntry.Format, &resultEntry.Iconcode)
		if err != nil { return ResultEntry{}, err }
		found = true
	}

	if !found { return ResultEntry{}, err }

	return resultEntry, nil
}

func postResultsValidation(db *pgxpool.Pool) gin.HandlerFunc {
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

		resultEntry, err := getResultEntryById(db, body.ResultId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
			return;
		}

		statusId := 3
		if !body.Verdict { statusId = 2 }
		resultStatus, err := getResultsStatus(db, statusId)
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

func getEvents(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		events, err := getAvailableEvents(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			c.IndentedJSON(http.StatusOK, events);
		}
	}
}

func getAllCompetitions(db *pgxpool.Pool) ([]CompetitionData, error) {
	rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c;`)
	if err != nil { return []CompetitionData{}, err }

	competitions := make([]CompetitionData, 0)

	for rows.Next() {
		var competition CompetitionData
		err = rows.Scan(&competition.Id, &competition.Name, &competition.Startdate, &competition.Enddate)
		if err != nil { return []CompetitionData{}, err }
		competition.getEvents(db)
		competitions = append(competitions, competition)
	}

	return competitions, nil
}

func getFilteredCompetitions(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		filter := c.Param("filter")
		
		result := make([]CompetitionData, 0);
		competitions, err := getAllCompetitions(db)
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

func (c *CompetitionData) getEvents(db *pgxpool.Pool) (error) {
	events := make([]CompetitionEvent, 0)

	rows, err := db.Query(context.Background(), `SELECT e.event_id, e.displayname, e.format, e.iconcode, e.puzzlecode FROM competition_events ce JOIN events e ON ce.event_id = e.event_id WHERE ce.competition_id = $1 ORDER BY e.event_id`, c.Id)
	if err != nil { return err }

	for rows.Next() {
		var event CompetitionEvent
		err := rows.Scan(&event.Id, &event.Displayname, &event.Format, &event.Iconcode, &event.Puzzlecode)
		if err != nil { return err }
		events = append(events, event)
	}

	c.Events = events

	return nil
}

func (s *ScrambleSet) addScramble(scramble string) {
	s.Scrambles = append(s.Scrambles, scramble)
}

func (c *CompetitionData) getScrambles(db *pgxpool.Pool) (error) {
	scrambleSets := make([]ScrambleSet, 0)

	for _, event := range c.Events {
		rows, err := db.Query(context.Background(), `SELECT s.scramble, e.event_id, e.displayname, e.format, e.iconcode, e.puzzlecode FROM scrambles s LEFT JOIN events e ON s.event_id = e.event_id WHERE s.competition_id = $1 AND s.event_id = $2 ORDER BY e.event_id, s."order";`, c.Id, event.Id)
		if err != nil { return err }

		var scrambleSet ScrambleSet
		for rows.Next() {
			var scramble string
			err := rows.Scan(&scramble, &scrambleSet.Event.Id, &scrambleSet.Event.Displayname, &scrambleSet.Event.Format, &scrambleSet.Event.Iconcode, &scrambleSet.Event.Puzzlecode)
			if err != nil { return err }
			scrambleSet.addScramble(scramble)
		}

		scrambleSets = append(scrambleSets, scrambleSet)
	}

	c.Scrambles = scrambleSets

	return nil
}

func getCompetitionById(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")	
		
		rows, err := db.Query(context.Background(), `SELECT c.competition_id, c.name, c.startdate, c.enddate FROM competitions c WHERE c.competition_id = $1;`, id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}

		var competition CompetitionData
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

		err = competition.getEvents(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}

		err = competition.getScrambles(db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)	
			return;
		}
		
		c.IndentedJSON(http.StatusOK, competition)
	}
}

func (competition *CompetitionData) recomputeCompetitionId() {
	trans := transliterator.NewTransliterator(nil)
	competition.Id = trans.Transliterate(strings.Join(strings.Split(competition.Name, " "), ""), "")
}

func postCompetition(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var competition CompetitionData

		if err := c.BindJSON(&competition); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "what bro")
			return
		}

		competition.recomputeCompetitionId()

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

type CompetitionEvents struct {
	Id int
	Competition_id string
	Event_id int
}

func (c *CompetitionData) removeAllEvents(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `DELETE FROM competition_events WHERE competition_id = $1;`, c.Id)
	return err
}

func (c *CompetitionData) addEvents(db *pgxpool.Pool) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	for _, event := range c.Events {
		_, err = tx.Exec(context.Background(), `INSERT INTO competition_events (competition_id, event_id) VALUES ($1, $2);`, c.Id, event.Id)
		if err != nil {
			tx.Rollback(context.Background())
			return err 
		}
	}

	return tx.Commit(context.Background());
}

func updateCompetitionEvents(competition *CompetitionData, db *pgxpool.Pool) error {
	if err := competition.removeAllEvents(db); err != nil { return err }
	if err := competition.addEvents(db); err != nil { return err }
	return nil;
}

func putCompetition(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var competition CompetitionData

		if err := c.BindJSON(&competition); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		_, err := db.Exec(context.Background(), `UPDATE competitions SET name = $1, startdate = $2, enddate = $3, timestamp = CURRENT_TIMESTAMP WHERE competition_id = $4;`, competition.Name, competition.Startdate, competition.Enddate, competition.Id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		err = updateCompetitionEvents(&competition, db)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusCreated, competition)
	}
}

func getManageRolesUsers(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		manageRolesUsers := make([]ManageRolesUser, 0)

		rows, err := db.Query(context.Background(), `SELECT u.user_id, u.name, u.isadmin FROM users u;`)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		uid := c.MustGet("uid").(int)

		for rows.Next() {
			var manageRolesUser ManageRolesUser
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

func putManageRolesUsers(db *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		var manageRolesUsers []ManageRolesUser

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