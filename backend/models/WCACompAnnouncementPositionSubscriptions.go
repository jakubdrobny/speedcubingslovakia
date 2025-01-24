package models

type WCACompAnnouncementsPositionSubscriptions struct {
	Id               string  `json:"id"`
	UserId           string  `json:"-"`
	LatitudeDegrees  float64 `json:"lat"`
	LongitudeDegrees float64 `json:"long"`
	Radius           float64 `json:"radius"`
	New              bool    `json:"new"`
	Open             bool    `json:"boolean"`
}
