package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
)

func Reverse[S ~[]E, E any](s S)  {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

func ParseSolveToMilliseconds(s string) int {
	if s == "DNF" { return constants.DNF }
	if s == "DNS" { return constants.DNS }

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

func CompareSolves(t1 *int, s2 string) {
	t2 := ParseSolveToMilliseconds(s2)
	if *t1 > t2 { *t1 = t2 }
}

func GetWorldRecords(eventName string) (int, int, error) {
	c := colly.NewCollector()

	single, average := constants.DNS, constants.DNS
	var err error

	c.OnHTML("div#results-list h2 a", func(e *colly.HTMLElement) {
		hrefSplit := strings.Split(e.Attr("href"), "/")
		if len(hrefSplit) > 3 && hrefSplit[3] == eventName {
			parentH2 := e.DOM.Parent()
			nextTable := parentH2.Next()
			singleTd := nextTable.Find("td.result").First()

			single = ParseSolveToMilliseconds(strings.Trim(singleTd.Text(), " "))

			// TODO: handle 333mbf parsing
			if eventName != "333mbf" {
				averageTd := singleTd.Parent().Next().Find("td.result").First()
				average = ParseSolveToMilliseconds(strings.Trim(averageTd.Text(), " "))
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

func CreateToken(userid int, secretKey string) (string, error) {
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

func LeftPad(s string, cnt int, ch string) string {
	for ; len(s) < cnt; { s = ch + s }
	return s
}

func FormatTime(timeInMiliseconds int) string {
	if timeInMiliseconds == constants.DNF { return "DNF" }
	if timeInMiliseconds == constants.DNS { return "DNS" }

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

func RegenerateImageForScramble(db *pgxpool.Pool, scrambleId int, scramble string, scramblingcode string) (error) {
	scramble = strings.ReplaceAll(scramble, "\n", " ")
	url := strings.ReplaceAll(fmt.Sprintf("http://localhost:2014/api/v0/view/%s/svg?scramble=%s", scramblingcode, scramble), " ", "%20")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { return err }

	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil { return err }

	_, err = db.Exec(context.Background(), `UPDATE scrambles SET svgimg = $1 WHERE scramble_id = $2;`, string(respBody), scrambleId)
	if err != nil { return err }

	return nil
}