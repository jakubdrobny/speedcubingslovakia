package models

import "github.com/google/uuid"

func NewTestContinent() Continent {
	return Continent{ContinentId: uuid.NewString(), Name: uuid.NewString(), RecordName: uuid.NewString()}
}

func NewTestCountry(continentId string) Country {
	return Country{CountryId: uuid.NewString(), Name: uuid.NewString(), Iso2: uuid.NewString(), ContinentId: continentId}
}
