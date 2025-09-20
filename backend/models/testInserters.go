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

	country := NewTestCountry(continent.Id)
	err = country.Insert(ctx, db)
	if err != nil {
		return Country{}, Continent{}, err
	}

	return country, continent, nil
}

func TestInsertUser(ctx context.Context, db interfaces.DB) (User, Country, Continent, error) {
	country, continent, err := TestInsertCountry(ctx, db)
	if err != nil {
		return User{}, Country{}, Continent{}, err
	}

	user := NewTestUser(country.Id)
	err = user.Insert(ctx, db)
	if err != nil {
		return User{}, Country{}, Continent{}, err
	}

	return user, country, continent, nil
}

func TestInsertWCACompAnnouncementsPositionSubscription(ctx context.Context, db interfaces.DB) (WCACompAnnouncementsPositionSubscription, User, Country, Continent, error) {
	u, co, ct, err := TestInsertUser(ctx, db)
	if err != nil {
		return WCACompAnnouncementsPositionSubscription{}, User{}, Country{}, Continent{}, err
	}

	sub := NewTestWCACompAnnouncementsPositionSubscription(u.Id)
	err = sub.Insert(ctx, db)
	if err != nil {
		return WCACompAnnouncementsPositionSubscription{}, User{}, Country{}, Continent{}, err
	}

	return sub, u, co, ct, nil
}

func TestInsertWCACompAnnouncementsSubscription(ctx context.Context, db interfaces.DB) (WCACompAnnouncementsSubscription, User, Country, Continent, error) {
	u, co, ct, err := TestInsertUser(ctx, db)
	if err != nil {
		return WCACompAnnouncementsSubscription{}, User{}, Country{}, Continent{}, err
	}

	sub := NewTestWCACompAnnouncementsSubscription(u.Id, co.Id)
	err = sub.Insert(ctx, db)
	if err != nil {
		return WCACompAnnouncementsSubscription{}, User{}, Country{}, Continent{}, err
	}

	return sub, u, co, ct, nil
}
