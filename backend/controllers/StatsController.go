package controllers

import (
	"context"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/utils"
)

func GetAdminStats(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var adminStatsCollection models.AdminStatsCollection

		rows, err := db.Query(
			context.Background(),
			`SELECT r.user_id, r.competition_id, c.enddate FROM results r JOIN competitions c ON c.competition_id = r.competition_id WHERE solve1 != 'DNS' or solve2 != 'DNS' or solve3 != 'DNS' or solve4 != 'DNS' or solve5 != 'DNS';`,
		)
		if err != nil {
			log.Println("ERR db.Query(results) in GetAdminStats: " + err.Error())
			c.IndentedJSON(http.StatusInternalServerError, "Failed to fetch data from db.")
			return
		}

		type Comp struct {
			Id          string
			Enddate     time.Time
			Competitors int
		}

		comps := map[Comp]map[int]bool{}
		total_competitors := map[int]bool{}

		for rows.Next() {
			var uid int
			comp := Comp{}
			err := rows.Scan(&uid, &comp.Id, &comp.Enddate)
			if err != nil {
				log.Println("ERR rows.Scan(uid, cid) in GetAdminStats: " + err.Error())
				c.IndentedJSON(http.StatusInternalServerError, "Faild to parse data from db.")
				return
			}

			if _, ok := comps[comp]; !ok {
				comps[comp] = make(map[int]bool)
			}
			comps[comp][uid] = true
			total_competitors[uid] = true
		}

		adminStatsCollection.Total = len(total_competitors)

		comps_arr := make([]Comp, 0)

		for comp := range comps {
			comp.Competitors = len(comps[comp])
			comps_arr = append(comps_arr, comp)
		}

		sort.Slice(comps_arr, func(i2, j int) bool {
			return comps_arr[i2].Enddate.Before(comps_arr[j].Enddate)
		})

		competitorsPerComp := make([]float64, 0)
		for idx, comp := range comps_arr {
			competitorsPerComp = append(competitorsPerComp, float64(comp.Competitors))
			adminStatsCollection.Max = max(adminStatsCollection.Max, comp.Competitors)
			adminStatsCollection.Average += float64(comp.Competitors)

			dataToAppend := make([]string, 3)
			dataToAppend[0] = comp.Id
			dataToAppend[2] = strconv.Itoa(comp.Competitors)

			last7Comps := make([]float64, 0)
			howMany := min(7, idx+1)
			for i := range howMany {
				last7Comps = append(last7Comps, float64(comps_arr[idx-howMany+1+i].Competitors))
			}
			dataToAppend[1] = strconv.Itoa(int(utils.GetMedian(last7Comps)))

			adminStatsCollection.ChartData.Data = append(
				adminStatsCollection.ChartData.Data,
				dataToAppend,
			)
		}

		sort.Slice(competitorsPerComp, func(i, j int) bool {
			return competitorsPerComp[i] < competitorsPerComp[j]
		})

		adminStatsCollection.Average = math.Round(
			(adminStatsCollection.Average/float64(len(competitorsPerComp)))*100,
		) / 100
		adminStatsCollection.Median = utils.GetMedian(competitorsPerComp)

		adminStatsCollection.ChartData.ColumnNames = []string{
			"Competition",
			"7 day median",
			"competitors",
		}

		c.IndentedJSON(http.StatusOK, adminStatsCollection)
	}
}
