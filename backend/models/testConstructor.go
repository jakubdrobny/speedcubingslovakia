package models

import (
	"math/rand/v2"

	"github.com/google/uuid"
)

func NewTestContinent() Continent {
	return Continent{Id: uuid.NewString(), Name: uuid.NewString(), RecordName: uuid.NewString()}
}

func NewTestCountry(continentId string) Country {
	return Country{Id: uuid.NewString(), Name: uuid.NewString(), Iso2: uuid.NewString(), ContinentId: continentId}
}

func NewTestUser(countryId string) User {
	return User{Name: uuid.NewString(), WcaId: uuid.NewString(), Url: uuid.NewString(), AvatarUrl: uuid.NewString(), CountryId: countryId, Sex: "m", IsAdmin: false, Email: uuid.NewString()}
}

func NewTestWCACompAnnouncementsPositionSubscription(userId int) WCACompAnnouncementsPositionSubscription {
	return WCACompAnnouncementsPositionSubscription{UserId: userId, LatitudeDegrees: rand.Float64()*360 - 180, LongitudeDegrees: rand.Float64()*360 - 180, Radius: rand.IntN(20000)}
}

func NewTestWCACompAnnouncementsSubscription(userId int, countryId string) WCACompAnnouncementsSubscription {
	return WCACompAnnouncementsSubscription{UserId: userId, CountryId: countryId, State: uuid.NewString()}
}
