package models

import (
	"context"
	"log"

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
		log.Println("ERR db.Query in GetCountries: " + err.Error())
		return []Country{}, err
	}

	countries := make([]Country, 0)
	for rows.Next() {
		var country Country
		err = rows.Scan(&country.Id, &country.Name, &country.Iso2)
		if err != nil {
			log.Println("ERR rows.Scan(country) in GetCountries: " + err.Error())
			return []Country{}, err
		}

		countries = append(countries, country)
	}

	return countries, nil
}

func GetCountryByName(db *pgxpool.Pool, name string) (Country, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT c.country_id, c.name, c.iso2 FROM countries c WHERE c.name = $1;`,
		name,
	)
	if err != nil {
		log.Println("ERR db.Query(country) in GetCountryByName: " + err.Error())
		return Country{}, err
	}

	var country Country
	for rows.Next() {
		err = rows.Scan(&country.Id, &country.Name, &country.Iso2)
		if err != nil {
			log.Println("ERR rows.Scan(country) in GetCountryByName: " + err.Error())
			return Country{}, err
		}
	}

	return country, nil
}

func CountriesArrayToMap(countriesArr []Country) map[string]Country {
	countriesMap := make(map[string]Country)

	for _, country := range countriesArr {
		countriesMap[country.Iso2] = country
	}

	return countriesMap
}
