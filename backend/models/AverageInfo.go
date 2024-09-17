package models

import "github.com/jakubdrobny/speedcubingslovakia/backend/constants"

type AverageInfo struct {
	Single              string   `json:"single"`
	Average             string   `json:"average"`
	Times               []string `json:"times"`
	Bpa                 string   `json:"bpa"`
	Wpa                 string   `json:"wpa"`
	ShowPossibleAverage bool     `json:"showPossibleAverage"`
	FinishedCompeting   bool     `json:"finishedCompeting"`
	Place               string   `json:"place"`
	SingleRecord        string   `json:"singleRecord"`
	SingleRecordColor   string   `json:"singleRecordColor"`
	AverageRecord       string   `json:"averageRecord"`
	AverageRecordColor  string   `json:"averageRecordColor"`
}

func (averageInfo *AverageInfo) LoadRecords(personalBestEntry ProfileTypePersonalBests) {
	if averageInfo.Single == personalBestEntry.Single.Value {
		if personalBestEntry.Single.WR == "1" {
			averageInfo.SingleRecord = "WR"
			averageInfo.SingleRecordColor = constants.WR_COLOR
		} else if personalBestEntry.Single.CR == "1" {
			averageInfo.SingleRecord = "CR"
			averageInfo.SingleRecordColor = constants.CR_COLOR
		} else if personalBestEntry.Single.NR == "1" {
			averageInfo.SingleRecord = "NR"
			averageInfo.SingleRecordColor = constants.NR_COLOR
		} else {
			averageInfo.SingleRecord = "PB"
			averageInfo.SingleRecordColor = constants.PR_COLOR
		}
	}

	if averageInfo.Average == personalBestEntry.Average.Value {
		if personalBestEntry.Average.WR == "1" {
			averageInfo.AverageRecord = "WR"
			averageInfo.AverageRecordColor = constants.WR_COLOR
		} else if personalBestEntry.Average.CR == "1" {
			averageInfo.AverageRecord = "CR"
			averageInfo.AverageRecordColor = constants.CR_COLOR
		} else if personalBestEntry.Average.NR == "1" {
			averageInfo.AverageRecord = "NR"
			averageInfo.AverageRecordColor = constants.NR_COLOR
		} else {
			averageInfo.AverageRecord = "PB"
			averageInfo.AverageRecordColor = constants.PR_COLOR
		}
	}
}
