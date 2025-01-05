package models

type WCACompAnnouncementSubscriptions struct {
	Id          int    `json:"-"`
	CountryId   string `json:"countryId"`
	CountryName string `json:"countryName"`
	Subscribed  bool   `json:"subscribed"`
}

type UpdateWCAAnnouncementSubscriptionsRequestBody struct {
	CountryName string `json:"countryName"`
	Subscribed  bool   `json:"subscribed"`
}
