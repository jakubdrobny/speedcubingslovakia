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
	"github.com/joho/godotenv"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/cube"
)

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func ParseMultiToMilliseconds(s string) int {
	r, _ := regexp.Compile("[0-9]{1,2}[/][0-9]{1,2}[ ][0-9]{0,2}[:]{0,1}[0-5][0-9]:[0-5][0-9]")
	if !r.MatchString(s) {
		return constants.DNS
	}

	if s == "0/0 00:00:00" {
		return constants.DNS
	}

	cubesPart := strings.Split(strings.Split(s, " ")[0], "/")
	timePart := strings.Split(strings.Split(s, " ")[1], ":")

	solved, _ := strconv.Atoi(cubesPart[0])
	attempted, _ := strconv.Atoi(cubesPart[1])
	points := solved - (attempted - solved)

	if solved < 2 || points < 0 {
		return constants.DNF
	}

	res := (points + 1) * 7200 * 1000

	var hours, minutes, seconds int
	if len(timePart) == 3 {
		hours, _ = strconv.Atoi(timePart[0])
		minutes, _ = strconv.Atoi(timePart[1])
		seconds, _ = strconv.Atoi(timePart[2])
	} else {
		hours = 0
		minutes, _ = strconv.Atoi(timePart[0])
		seconds, _ = strconv.Atoi(timePart[1])
	}

	res -= hours * 3600 * 1000
	res -= minutes * 60 * 1000
	res -= seconds * 1000

	return -res
}

func ParseSolveToMilliseconds(s string, isfmc bool, scramble string) int {
	if s == "DNF" {
		return constants.DNF
	}
	if s == "DNS" {
		return constants.DNS
	}

	if isfmc {
		return cube.ParseFMCSolutionToMilliseconds(scramble, s)
	}

	if idx := strings.Index(s, "/"); idx != -1 {
		return ParseMultiToMilliseconds(s)
	}

	if !strings.Contains(s, ".") {
		s += ".00"
	}

	split := strings.Split(s, ".")

	wholePart := strings.Split(split[0], ":")
	decimalPart, err := strconv.Atoi(split[1])

	if err != nil {
		return constants.DNS
	}

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

	if err != nil {
		return constants.DNS
	}

	return res
}

func CompareSolves(t1 *int, s2 string, isfmc bool, scramble string) {
	t2 := ParseSolveToMilliseconds(s2, isfmc, scramble)
	if *t1 > t2 {
		*t1 = t2
	}
}

func GetWorldRecords(eventName string) (int, int, error) {
	single := map[string]string{
		"333":                    "3.13",
		"222":                    "0.43",
		"444":                    "15.71",
		"555":                    "31.60",
		"666":                    "58.03",
		"777":                    "1:34.15",
		"333bf":                  "12.00",
		"333fm":                  "16",
		"333oh":                  "5.66",
		"clock":                  "1.97",
		"minx":                   "23.18",
		"pyram":                  "0.73",
		"skewb":                  "0.75",
		"sq1":                    "3.41",
		"444bf":                  "51.96",
		"555bf":                  "2:04.41",
		"333mbf":                 "62/65 57:47",
		"unofficial-222bf":       "1.00",
		"unofficial-666bf":       "4:00.00",
		"unofficial-777bf":       "8:00.00",
		"333ft":                  "13.00",
		"unofficial-333mts":      "10.00",
		"unofficial-234relay":    "20.00",
		"unofficial-2345relay":   "55.00",
		"unofficial-23456relay":  "2:00.00",
		"unofficial-234567relay": "3:40.00",
		"unofficial-kilominx":    "8.00",
		"unofficial-miniguild":   "2:30.00",
		"unofficial-redi":        "2.70",
		"unofficial-mpyram":      "10.00",
		"unofficial-15puzzle":    "5.00",
		"unofficial-mirror":      "8.00",
		"unofficial-fto":         "11.00",
	}
	average := map[string]string{

		"333":                 "4.09",
		"222":                 "0.78",
		"444":                 "19.38",
		"555":                 "34.76",
		"666":                 "1:05.66",
		"777":                 "1:39.68",
		"333bf":               "14.05",
		"333fm":               "20.00",
		"333oh":               "8.09",
		"clock":               "2.39",
		"minx":                "26.84",
		"pyram":               "1.27",
		"skewb":               "1.52",
		"sq1":                 "4.81",
		"444bf":               "1:06.46",
		"555bf":               "2:27.63",
		"unofficial-222bf":    "3.00",
		"333ft":               "17.00",
		"unofficial-333mts":   "12.00",
		"unofficial-kilominx": "10.00",
		"unofficial-redi":     "4.00",
		"unofficial-mpyram":   "16.00",
		"unofficial-15puzzle": "8.00",
		"unofficial-mirror":   "11.50",
		"unofficial-fto":      "15.00",
	}

	var retSingle, retAverage int
	recSingle, ok := single[eventName]
	if !ok {
		retSingle = constants.VERY_SLOW
	} else {
		retSingle = ParseSolveToMilliseconds(recSingle, false, "")
	}

	recAverage, ok := average[eventName]
	if !ok {
		retAverage = constants.VERY_SLOW
	} else {
		retAverage = ParseSolveToMilliseconds(recAverage, false, "")
	}

	return retSingle, retAverage, nil
}

func DOES_NOT_WORK_CURRENTLY_FOR_SOME_REASON_GetWorldRecords(eventName string) (int, int, error) {
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

	c.Visit("https://www.worldcubeassociation.org/results/records")
	if err != nil {
		return constants.VERY_SLOW, constants.VERY_SLOW, err
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
			"exp":    time.Now().Add(time.Hour * time.Duration(expiresInSeconds)).Unix(),
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
	for len(s) < cnt {
		s = ch + s
	}
	return s
}

func FormatMultiTime(timeInMiliseconds int) string {
	res := ""

	return res
}

func FormatMultiTimes(solves []string) []string {
	for idx := 0; idx < len(solves); idx++ {
		timeInMilliseconds := ParseSolveToMilliseconds(solves[idx], false, "")
		if timeInMilliseconds == constants.DNS {
			solves[idx] = "DNS"
		}
		if timeInMilliseconds == constants.DNF {
			solves[idx] = "DNF"
		}
	}

	return solves
}

func FormatTime(timeInMiliseconds int, isfmc bool) string {
	if timeInMiliseconds == constants.DNF {
		return "DNF"
	}
	if timeInMiliseconds == constants.DNS {
		return "DNS"
	}

	if isfmc {
		return fmt.Sprintf("%.2f", float64(timeInMiliseconds)/1000)
	}

	if timeInMiliseconds < 0 {
		return FormatMultiTime(timeInMiliseconds)
	}

	if timeInMiliseconds%10 >= 5 {
		timeInMiliseconds += 10 - (timeInMiliseconds % 10)
	}

	res := make([]string, 0)

	pw := 1000 * 60 * 60 * 24
	for _, mul := range []int{24, 60, 60, 1000, 1} {
		toPush := fmt.Sprint(timeInMiliseconds / pw)
		if mul == 1 {
			toPush = LeftPad(toPush, 3, "0")
		}
		res = append(res, toPush)
		timeInMiliseconds %= pw
		pw /= mul
	}

	res[len(res)-1] = res[len(res)-1][:len(res[len(res)-1])-1]
	sliceIdx := 0
	for ; sliceIdx < len(res)-2 && res[sliceIdx][0] == '0'; sliceIdx++ {
	}
	res = res[sliceIdx:]

	resString := ""
	resIdx := 0
	for ; resIdx < len(res)-1; resIdx++ {
		if resIdx > 0 {
			resString += LeftPad(res[resIdx], 2, "0")
		} else {
			resString += res[resIdx]
		}

		if resIdx == len(res)-2 {
			resString += "."
		} else {
			resString += ":"
		}
	}
	resString += LeftPad(res[resIdx], 2, "0")

	return resString
}

func RandSeq(n int) string {
	b := make([]byte, n+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : n+2]
}

func SaveScrambleImg(img_id string, svg_content string) error {
	envMap, err := godotenv.Read(
		fmt.Sprintf(".env.%s", os.Getenv("SPEEDCUBINGSLOVAKIA_BACKEND_ENV")),
	)
	if err != nil {
		return err
	}

	folder_path := envMap["SCRAMBLE_IMAGES_PATH"]
	f, err := os.Create(fmt.Sprintf("%s/%s", folder_path, img_id))
	if err != nil {
		return err
	}

	n, err := f.WriteString(svg_content)
	if err != nil {
		return err
	}

	log.Printf("Wrote %d bytes.", n)
	return nil
}

func RegenerateImageForScramble(
	db *pgxpool.Pool,
	scrambleId int,
	scramble string,
	scramblingcode string,
) (string, error) {
	url := fmt.Sprintf(
		"http://localhost:2014/api/v0/view/%s/svg?scramble=%s",
		scramblingcode,
		url.QueryEscape(scramble),
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	imgId := RandSeq(64) + ".svg" // with extension
	err = SaveScrambleImg(imgId, string(respBody))
	if err != nil {
		return "", err
	}

	_, err = db.Exec(
		context.Background(),
		`UPDATE scrambles SET img = $1 WHERE scramble_id = $2;`,
		imgId,
		scrambleId,
	)
	if err != nil {
		return "", err
	}

	if scramblingcode == "sq1" {
		_, err = db.Exec(
			context.Background(),
			`UPDATE scrambles SET scramble = $1 WHERE scramble_id = $2;`,
			imgId,
			scrambleId,
		)
		if err != nil {
			return "", err
		}
	}

	return imgId, nil
}

func GetContinents(db *pgxpool.Pool) ([]string, error) {
	rows, err := db.Query(context.Background(), `SELECT c.name FROM continents c;`)
	if err != nil {
		return []string{}, err
	}

	continents := make([]string, 0)
	for rows.Next() {
		var continentName string
		err = rows.Scan(&continentName)
		if err != nil {
			return []string{}, err
		}

		continents = append(continents, continentName)
	}

	return continents, nil
}

func GetCountries(db *pgxpool.Pool) ([]string, error) {
	rows, err := db.Query(context.Background(), `SELECT c.name FROM countries c;`)
	if err != nil {
		return []string{}, err
	}

	countries := make([]string, 0)
	for rows.Next() {
		var countryName string
		err = rows.Scan(&countryName)
		if err != nil {
			return []string{}, err
		}

		countries = append(countries, countryName)
	}

	return countries, nil
}

func NextMonday() time.Time {
	res := time.Now()

	offset := (7 - int(res.Weekday()-time.Monday)) % 7
	res = res.AddDate(0, 0, offset)

	res = time.Date(res.Year(), res.Month(), res.Day(), 0, 0, 0, 0, res.Location())

	return res
}

func IsFMC(eventiconcode string) bool {
	return eventiconcode == "333fm"
}

func GetSolve(solve string, isfmc bool, scramble string) string {
	if !isfmc {
		return solve
	}
	return FormatTime(ParseSolveToMilliseconds(solve, isfmc, scramble), isfmc)
}

func GetScramblesByResultEntryId(db *pgxpool.Pool, eid int, cid string) ([]string, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT scramble FROM scrambles WHERE event_id = $1 AND competition_id = $2 ORDER BY "order";`,
		eid,
		cid,
	)
	if err != nil {
		return []string{}, err
	}

	scrambles := make([]string, 5)
	idx := 0
	for rows.Next() {
		var scramble string
		err = rows.Scan(&scramble)
		if err != nil {
			return []string{}, err
		}

		scrambles[idx] = scramble
		idx++

		if idx == 5 {
			rows.Close()
			break
		}
	}

	return scrambles, nil
}

// endsWith: 1. to  1st, 2. to 2nd, 3. to 3rd, 4-9. to 4-9th, ...
// except: 11-19. to 11-19th
func PlaceFromDotToEnglish(from string) string {
	from = strings.Join(strings.Split(from, "."), "")
	fromInt, _ := strconv.ParseInt(from, 10, 64)

	if fromInt >= 10 && fromInt < 20 {
		return from + "th"
	}

	rem := fromInt % 10
	if rem == 1 {
		return from + "st"
	}
	if rem == 2 {
		return from + "nd"
	}
	if rem == 3 {
		return from + "rd"
	}

	return from + "th"
}
