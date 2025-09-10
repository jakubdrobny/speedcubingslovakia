package views

type Continent struct {
	Name       string `json:"continent"`
	RecordName string `json:"record_name"`
}

type Country struct {
	Name      string    `json:"name"`
	Iso2      string    `json:"iso2"`
	Continent Continent `json:"continent"`
}

type ManageUser struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	WcaId   string  `json:"wca_id"`
	IsAdmin bool    `json:"is_admin"`
	Country Country `json:"country"`
}
