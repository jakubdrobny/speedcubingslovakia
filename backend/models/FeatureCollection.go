package models

import (
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MapDataUser struct {
	Username string `json:"username"`
	WcaId    string `json:"wcaid"`
	Score    string `json:"score"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates"`
}

type Property struct {
	Name        string        `json:"name"`
	CountryIso2 string        `json:"countryIso2"`
	Users       []MapDataUser `json:"users"`
}

type Feature struct {
	Type       string   `json:"type"`
	Properties Property `json:"properties"`
	Geometry   Geometry `json:"geometry"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// returns: map of OveralResults by country, logMsg, retMsg, error
func GetUsersByCountryWithKinchScore(
	db *pgxpool.Pool,
) (map[string][]MapDataUser, string, string, error) {
	usersByCountry := make(map[string][]MapDataUser)

	overallResults, err := GetOverallResults(db, "", "World", "World")
	if err != nil {
		return map[string][]MapDataUser{}, "ERR GetOverallResults in GetUsersByCountryWithKinchScore: " + err.Error(), "Failed to get user scores.", err
	}

	for _, result := range overallResults {
		var mapDataUser MapDataUser
		mapDataUser.Username = result.Username
		mapDataUser.WcaId = result.WcaId
		mapDataUser.Score = result.Score

		usersByCountry[result.CountryIso2] = append(usersByCountry[result.CountryIso2], mapDataUser)
	}

	for k := range usersByCountry {
		sort.Slice(usersByCountry[k], func(i, j int) bool {
			val1, val2 := usersByCountry[k][i].Score, usersByCountry[k][j].Score
			val1, val2 = fmt.Sprintf("%06s", val1), fmt.Sprintf("%06s", val2)
			if val1 == val2 {
				return usersByCountry[k][i].Username < usersByCountry[k][j].Username
			}
			return val1 > val2
		})
	}

	return usersByCountry, "", "", nil
}
