package tablespostgresql

import "github.com/google/uuid"

func NewTestContinent() Continent {
	return Continent{
		Name:       uuid.NewString(),
		RecordName: uuid.NewString(),
	}
}

func NewCountry() Country {
	return Country{
		Name: uuid.NewString(),
		Iso2: uuid.NewString(),
	}
}

func NewUser() User {
	return User{
		Name: uuid.NewString(),
	}
}
