package views

type WCACompAnnouncementsSubscription struct {
	Id          int
	CountryId   string `json:"countryId"`
	CountryName string `json:"countryName"`
	State       string `json:"state"`
	Subscribed  bool   `json:"subscribed"`
}
