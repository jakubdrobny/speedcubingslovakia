package models

type WCACompAnnouncementSubscriptions struct {
	Id          int    `json:"-"`
	CountryId   string `json:"countryId"`
	CountryName string `json:"countryName"`
	State       string `json:"state"`
	Subscribed  bool   `json:"subscribed"`
}

type UpdateWCAAnnouncementSubscriptionsRequestBody struct {
	CountryName string `json:"countryName"`
	State       string `json:"state"`
	Subscribed  bool   `json:"subscribed"`
}
