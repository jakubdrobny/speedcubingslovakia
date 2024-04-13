package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/golang-jwt/jwt"
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
)

func Reverse[S ~[]E, E any](s S)  {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

func TryParseSolveToMilliseconds(s string) int {
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
	t2 := TryParseSolveToMilliseconds(s2)
	if float64(*t1 - t2) > 1e-9 {
		*t1 = t2;
	}
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

			single = TryParseSolveToMilliseconds(strings.Trim(singleTd.Text(), " "))

			// TODO: handle 333mbf parsing
			if eventName != "333mbf" {
				averageTd := singleTd.Parent().Next().Find("td.result").First()
				average = TryParseSolveToMilliseconds(strings.Trim(averageTd.Text(), " "))
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

func FormatTime(timeInMiliseconds int) string {
	if timeInMiliseconds == constants.DNF { return "DNF" }
	if timeInMiliseconds == constants.DNS { return "DNS" }

	wholePart := ""
	decimalPart := ""

	pw := 1000 * 60 * 60 * 24
	for _, mul := range []int{1, 24, 60, 60, 1} {
		times := timeInMiliseconds / pw
		if times > 0 || wholePart != "" {
			strTimes := fmt.Sprintf("%02d", times)
			if len(wholePart) > 0 { wholePart += ":" }
			wholePart += strTimes
		}

		fmt.Println(timeInMiliseconds, times, wholePart)

		timeInMiliseconds %= pw
		pw /= mul
	}

	if wholePart == "" { wholePart += "0" }

	decimalPart = fmt.Sprintf("%02d", timeInMiliseconds / 10)

	fmt.Println("hahahahaha", wholePart, decimalPart)

	return wholePart + "." + decimalPart
}