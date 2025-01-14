package models

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
)

type Country struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Iso2        string `json:"iso2"`
	ContinentId string `json:"-"`
}

func GetCountries(db *pgxpool.Pool) ([]Country, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT c.country_id, c.name, c.iso2, c.continent_id FROM countries c;`,
	)
	if err != nil {
		log.Println("ERR db.Query in GetCountries: " + err.Error())
		return []Country{}, err
	}

	countries := make([]Country, 0)
	for rows.Next() {
		var country Country
		err = rows.Scan(&country.Id, &country.Name, &country.Iso2, &country.ContinentId)
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

func CountriesArrayToMap(countriesArr []Country) map[string][]Country {
	countriesMap := make(map[string][]Country)
	for _, countryGroupIso2 := range constants.COUNTRY_GROUPS_ISO2 {
		countriesMap[countryGroupIso2] = make([]Country, 0)
	}

	for _, country := range countriesArr {
		// add to its country
		countriesMap[country.Iso2] = []Country{country}

		// add to its continent
		countriesMap[constants.CONTINENT_ID_TO_COUNTRY_GROUP_ISO2[country.ContinentId]] =
			append(
				countriesMap[constants.CONTINENT_ID_TO_COUNTRY_GROUP_ISO2[country.ContinentId]],
				country,
			)

		// add to americas if should
		if country.ContinentId == "_North America" || country.ContinentId == "_South America" {
			countriesMap["XN"] = append(countriesMap["XN"], country)
			countriesMap["XS"] = append(countriesMap["XS"], country)
		}

		// add to world
		countriesMap["XW"] = append(countriesMap["XW"], country)
	}

	return countriesMap
}
