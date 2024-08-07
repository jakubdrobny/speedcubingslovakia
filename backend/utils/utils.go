package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/cube"
	"github.com/joho/godotenv"
)

func Reverse[S ~[]E, E any](s S)  {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

func ParseMultiToMilliseconds(s string) int {
	r, _ := regexp.Compile("[0-9]{1,2}[/][0-9]{1,2}[ ][0-9]{0,2}[:]{0,1}[0-5][0-9]:[0-5][0-9]")
	if !r.MatchString(s) {
		return constants.DNS
	}

	if s == "0/0 00:00:00" { return constants.DNS }

	cubesPart := strings.Split(strings.Split(s, " ")[0], "/")
	timePart := strings.Split(strings.Split(s, " ")[1], ":")

	solved, _ := strconv.Atoi(cubesPart[0])
	attempted, _ := strconv.Atoi(cubesPart[1])
	points := solved - (attempted - solved)

	if solved < 2 || points < 0 {
		return constants.DNF
	}

	res := points * 7200 * 1000
	
	var hours, minutes, seconds int
	if len(timePart) == 3 {
		hours, _ = strconv.Atoi(timePart[0])
		minutes, _ = strconv.Atoi(timePart[1])
		seconds, _ = strconv.Atoi(timePart[2])
	} else {
		hours  = 0
		minutes, _ = strconv.Atoi(timePart[0])
		seconds, _ = strconv.Atoi(timePart[1])
	}

	res -= hours * 3600 * 1000
	res -= minutes * 60 * 1000
	res -= seconds * 1000

	return -res
}

func ParseSolveToMilliseconds(s string, isfmc bool, scramble string) int {
	if s == "DNF" { return constants.DNF }
	if s == "DNS" { return constants.DNS }

	if isfmc { return cube.ParseFMCSolutionToMilliseconds(scramble, s) }

	if idx := strings.Index(s, "/"); idx != -1 {
		return ParseMultiToMilliseconds(s)
	}

	if !strings.Contains(s, ".") { s += ".00" }

	split := strings.Split(s, ".")
	
	wholePart := strings.Split(split[0], ":")
	decimalPart, err := strconv.Atoi(split[1])

	if err != nil { return constants.DNS }

	res := decimalPart * 10 // to milliseconds

	var add int
	Reverse(wholePart)

	if len(wholePart) > 0 { 
		add, err = strconv.Atoi(wholePart[0])
		res += add * 1000
	}

	if len(wholePart) > 1 { 
		add, err = strconv.Atoi(wholePart[1])
		res += 60 * add * 1000
	}

	if len(wholePart) > 2 { 
		add, err = strconv.Atoi(wholePart[2])
		res += 60 * 60 * add * 1000
	}

	if len(wholePart) > 3 { 
		add, err = strconv.Atoi(wholePart[3])
		res += 24 * 60 * 60 * add * 1000
	}

	if err != nil { return constants.DNS }

	return res
}

func CompareSolves(t1 *int, s2 string, isfmc bool, scramble string) {
	t2 := ParseSolveToMilliseconds(s2, isfmc, scramble)
	if *t1 > t2 { *t1 = t2 }
}

func GetWorldRecords(eventName string) (int, int, error) {
	c := colly.NewCollector()

	single, average := constants.VERY_SLOW, constants.VERY_SLOW
	var err error

	c.OnHTML("div#results-list h2 a", func(e *colly.HTMLElement) {
		hrefSplit := strings.Split(e.Attr("href"), "/")
		if len(hrefSplit) > 3 && hrefSplit[3] == eventName {
			parentH2 := e.DOM.Parent()
			nextTable := parentH2.Next()
			singleTd := nextTable.Find("td.result").First()

			single = ParseSolveToMilliseconds(strings.Trim(singleTd.Text(), " "), false, "")

			if eventName != "333mbf" {
				averageTd := singleTd.Parent().Next().Find("td.result").First()
				average = ParseSolveToMilliseconds(strings.Trim(averageTd.Text(), " "), false, "")
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

func GetNoOfSolves(format string) (int, error) {
	match := regexp.MustCompile(`\d+$`).FindString(format)
	res, err := strconv.Atoi(match)

	if err != nil {
		return 0, fmt.Errorf("did not find a number at the end of format")
	}

	return res, nil
}

func CreateToken(userid int, secretKey string, expiresInSeconds int) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
			"userid": userid, 
			"exp": time.Now().Add(time.Hour * time.Duration(expiresInSeconds)).Unix(),
        })

    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
    	return "", err
    }

 	return tokenString, nil
}

// returns true if the format is ok, eg. not >= 60 seconds, not >= 60 minutes, ...
// otherwise returns false
func CheckFormat(solve string) bool {
	timeInMiliseconds := ParseSolveToMilliseconds(solve, false, "")
	formattedTime := FormatTime(timeInMiliseconds, false)
	return solve == formattedTime
}

func LeftPad(s string, cnt int, ch string) string {
	for ; len(s) < cnt; { s = ch + s }
	return s
}

func FormatMultiTime(timeInMiliseconds int) string {
	res := ""

	return res
}

func FormatMultiTimes(solves []string) []string {
	for idx := 0; idx < len(solves); idx++ {
		timeInMilliseconds := ParseSolveToMilliseconds(solves[idx], false, "")
		if timeInMilliseconds == constants.DNS { solves[idx] = "DNS" }
		if timeInMilliseconds == constants.DNF { solves[idx] = "DNF" }
	}

	return solves
}

func FormatTime(timeInMiliseconds int, isfmc bool) string {
	if timeInMiliseconds == constants.DNF { return "DNF" }
	if timeInMiliseconds == constants.DNS { return "DNS" }

	if isfmc { return fmt.Sprint(int(timeInMiliseconds / 1000)) + ".00" }

	if timeInMiliseconds < 0 {
		return FormatMultiTime(timeInMiliseconds)
	}

	if timeInMiliseconds % 10 >= 5 {
		timeInMiliseconds += 10 - (timeInMiliseconds % 10)
	}

	res := make([]string, 0)

	pw := 1000 * 60 * 60 * 24
	for _, mul := range []int{24, 60, 60, 1000, 1} {
		toPush := fmt.Sprint(timeInMiliseconds / pw)
		if mul == 1 { toPush = LeftPad(toPush, 3, "0") }	
		res = append(res, toPush)
		timeInMiliseconds %= pw
		pw /= mul
	}

	res[len(res) - 1] = res[len(res) - 1][:len(res[len(res) - 1]) - 1]
	sliceIdx := 0
	for ; sliceIdx < len(res) - 2 && res[sliceIdx][0] == '0'; sliceIdx++ {}
	res = res[sliceIdx:]

	resString := ""
	resIdx := 0
	for ; resIdx < len(res) - 1; resIdx++ {
		if resIdx > 0 {
			resString += LeftPad(res[resIdx], 2, "0")
		} else {
			resString += res[resIdx] 
		}

		if resIdx == len(res) - 2 {
			resString += "."
		} else {
			resString += ":" 
		}
	}
	resString += LeftPad(res[resIdx], 2, "0")

	return resString
}

func RandSeq(n int) string {
    b := make([]byte, n + 2)
    rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : n + 2]
}

func SaveScrambleImg(img_id string, svg_content string) error {
	envMap, err := godotenv.Read(fmt.Sprintf(".env.%s", os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV")))
	if err != nil { return err }

	folder_path := envMap["SCRAMBLE_IMAGES_PATH"]
	f, err := os.Create(fmt.Sprintf("%s/%s", folder_path, img_id))
	if err != nil { return err }

	n, err := f.WriteString(svg_content)
	if err != nil { return err }

	log.Printf("Wrote %d bytes.", n)
	return nil
}

func RegenerateImageForScramble(db *pgxpool.Pool, scrambleId int, scramble string, scramblingcode string) (string, error) {
	url := fmt.Sprintf("http://localhost:2014/api/v0/view/%s/svg?scramble=%s", scramblingcode, url.QueryEscape(scramble))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { return "", err }

	resp, err := http.DefaultClient.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil { return "", err }

	imgId := RandSeq(64) + ".svg"  // with extension
	err = SaveScrambleImg(imgId, string(respBody))
	if err != nil { return "", err }

	_, err = db.Exec(context.Background(), `UPDATE scrambles SET img = $1 WHERE scramble_id = $2;`, imgId, scrambleId)
	if err != nil { return "", err }

	if scramblingcode == "sq1" {
		_, err = db.Exec(context.Background(), `UPDATE scrambles SET scramble = $1 WHERE scramble_id = $2;`, imgId, scrambleId)
		if err != nil { return "", err }
	}

	return imgId, nil
}

func GetContinents(db *pgxpool.Pool) ([]string, error) {
	rows, err := db.Query(context.Background(), `SELECT c.name FROM continents c;`)
	if err != nil { return []string{}, err }

	continents := make([]string, 0)
	for rows.Next() {
		var continentName string
		err = rows.Scan(&continentName)
		if err != nil { return []string{}, err }

		continents = append(continents, continentName)
	}
	
	return continents, nil
}

func GetCountries(db *pgxpool.Pool) ([]string, error) {
	rows, err := db.Query(context.Background(), `SELECT c.name FROM countries c;`)
	if err != nil { return []string{}, err }

	countries := make([]string, 0)
	for rows.Next() {
		var countryName string
		err = rows.Scan(&countryName)
		if err != nil { return []string{}, err }

		countries = append(countries, countryName)
	}
	
	return countries, nil
}

func NextMonday() (time.Time) {
	res := time.Now()

	offset := (7 - int(res.Weekday() - time.Monday)) % 7
	res = res.AddDate(0, 0, offset)

	res = time.Date(res.Year(), res.Month(), res.Day(), 0, 0, 0, 0, res.Location())

	return res
}

func IsFMC(eventiconcode string) bool {
	return eventiconcode == "333fm"
}

func GetSolve(solve string, isfmc bool, scramble string) (string) {
	if !isfmc { return solve }
	return FormatTime(ParseSolveToMilliseconds(solve, isfmc, scramble), isfmc)
}

func GetScramblesByResultEntryId(db *pgxpool.Pool, eid int, cid string) ([]string, error) {
	rows, err := db.Query(context.Background(), `SELECT scramble FROM scrambles WHERE event_id = $1 AND competition_id = $2 ORDER BY "order";`, eid, cid)
	if err != nil { return []string{}, err }

	scrambles := make([]string, 5)
	idx := 0
	for rows.Next() {
		var scramble string
		err = rows.Scan(&scramble)
		if err != nil { return []string{}, err }

		scrambles[idx] = scramble
		idx++

		if idx == 5 {
			rows.Close()
			break
		}
	}

	return scrambles, nil
}