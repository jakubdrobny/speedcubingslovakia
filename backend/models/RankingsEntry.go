package models

import (
	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

type RankingsEntry struct {
	Place           string   `json:"place"`
	Username        string   `json:"username"`
	WcaId           string   `json:"wca_id"`
	CountryISO2     string   `json:"country_iso2"`
	CountryName     string   `json:"country_name"`
	Result          string   `json:"result"`
	CompetitionId   string   `json:"competitionId"`
	CompetitionName string   `json:"competitionName"`
	Times           []string `json:"times"`
}

func (rankingsEntry *RankingsEntry) Load(single, persons bool, resultsEntry ResultEntry, isfmc, ismbld bool, scrambles []string, rankings *[]RankingsEntry, regionType, regionPrecise string) (string, string, error) {
	if single {
		if persons {
			rankingsEntry.Result = resultsEntry.SingleFormatted(isfmc, scrambles)
			if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") >= constants.VERY_SLOW {
				return "", "", nil
			}
			rankingsEntry.Times = make([]string, 0)

			*rankings = append(*rankings, *rankingsEntry)
		} else {
			rankingsEntry.Times = make([]string, 0)

			noOfSolves, _ := utils.GetNoOfSolves(resultsEntry.Format)

			for idx, solve := range []string{resultsEntry.Solve1, resultsEntry.Solve2, resultsEntry.Solve3, resultsEntry.Solve4, resultsEntry.Solve5} {
				if idx >= noOfSolves {
					break
				}

				result := utils.ParseSolveToMilliseconds(solve, isfmc, scrambles[idx])
				if ismbld {
					rankingsEntry.Result = solve
				} else {
					rankingsEntry.Result = utils.FormatTime(result, isfmc)
				}
				if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") < constants.VERY_SLOW {
					*rankings = append(*rankings, *rankingsEntry)
				}
			}
		}
	} else if !ismbld && resultsEntry.Format != "bo1" {
		resultFormatted, err := resultsEntry.AverageFormatted(isfmc, scrambles)
		if err != nil {
			return "ERR AverageFormatted in rankingsEntry.Load (" + regionType + "+" + regionPrecise + "): " + err.Error(), "Failed to calculate average in rankings entry.", err
		}
		rankingsEntry.Result = resultFormatted
		if utils.ParseSolveToMilliseconds(rankingsEntry.Result, false, "") >= constants.VERY_SLOW {
			return "", "", nil
		}
		rankingsEntry.Times, _ = resultsEntry.GetFormattedTimes(isfmc, scrambles)

		*rankings = append(*rankings, *rankingsEntry)
	}

	return "", "", nil
}
