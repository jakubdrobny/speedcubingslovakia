package models

import (
	"context"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

func TestInsertContinent(ctx context.Context, db interfaces.DB) (Continent, error) {
	c := NewTestContinent()
	err := c.Insert(ctx, db)
	if err != nil {
		return Continent{}, err
	}

	return c, nil
}

func TestInsertCountry(ctx context.Context, db interfaces.DB) (Country, Continent, error) {
	continent, err := TestInsertContinent(ctx, db)
	if err != nil {
		return Country{}, Continent{}, err
	}

	country := NewTestCountry(continent.ContinentId)
	err = country.Insert(ctx, db)
	if err != nil {
		return Country{}, Continent{}, err
	}

	return country, continent, nil
}
