package models

type AdminStatsCollection struct {
	Max       int     `json:"max"`
	Total     int     `json:"total"`
	Median    float64 `json:"median"`
	Average   float64 `json:"average"`
	ChartData struct {
		ColumnNames []string   `json:"columnNames"`
		Data        [][]string `json:"data"`
	} `json:"chartData"`
}
