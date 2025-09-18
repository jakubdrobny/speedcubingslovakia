package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

type Country struct {
	CountryId   string `json:"id"`
	Name        string `json:"name"`
	Iso2        string `json:"iso2"`
	ContinentId string `json:"-"`
}

func GetCountries(ctx context.Context, db interfaces.DB) ([]Country, error) {
	rows, err := db.Query(
		ctx,
		`SELECT c.country_id, c.name, c.iso2, c.continent_id FROM countries c;`,
	)
	if err != nil {
		return []Country{}, fmt.Errorf("%w: when querying countries", err)
	}
	defer rows.Close()

	countries := make([]Country, 0)
	for rows.Next() {
		var country Country
		err = rows.Scan(&country.CountryId, &country.Name, &country.Iso2, &country.ContinentId)
		if err != nil {
			return []Country{}, fmt.Errorf("%w: when scanning country", err)
		}

		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		return []Country{}, fmt.Errorf("%w: when iterating over rows", err)
	}

	return countries, nil
}

func (c *Country) Get(ctx context.Context, db interfaces.DB, name string) error {
	err := db.QueryRow(
		ctx,
		`SELECT c.country_id, c.name, c.iso2, c.continent_id FROM countries c WHERE c.name = $1;`,
		name,
	).Scan(&c.CountryId, &c.Name, &c.Iso2, &c.ContinentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("country with name=%s not found", name)
		}
		return fmt.Errorf("%w: when querying country with name=%s", err, name)
	}

	return nil
}

func (c Country) Insert(ctx context.Context, db interfaces.DB) error {
	_, err := db.Exec(ctx, `INSERT INTO countries (country_id, name, continent_id, iso2) VALUES ($1, $2, $3, $4)`, c.CountryId, c.Name, c.ContinentId, c.Iso2)
	if err != nil {
		return fmt.Errorf("%w: when executing insert continent statement for country=%+v", err, c)
	}

	return nil
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
