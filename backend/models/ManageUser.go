package models

type ManageUser struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	WcaId       string `json:"wca_id"`
	IsAdmin     bool   `json:"is_admin"`
	Country     string `json:"country_name"`
	CountryIso2 string `json:"country_iso2"`
}
