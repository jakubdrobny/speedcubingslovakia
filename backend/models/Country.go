package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Country struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Iso2 string `json:"iso2"`
}

func GetCountries(db *pgxpool.Pool) ([]Country, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT c.country_id, c.name, c.iso2 FROM countries c;`,
	)
	if err != nil {
		return []Country{}, err
	}

	countries := make([]Country, 0)
	for rows.Next() {
		var country Country
		err = rows.Scan(&country.Id, &country.Name, &country.Iso2)
		if err != nil {
			return []Country{}, err
		}

		countries = append(countries, country)
	}

	return countries, nil
}
