package models

import (
	"context"
	"fmt"

	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

type ManageUser struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	WcaId       string `json:"wca_id"`
	IsAdmin     bool   `json:"is_admin"`
	CountryName string `json:"country_name"`
	CountryIso2 string `json:"country_iso2"`
}

func (u *ManageUser) UpdateRole(ctx context.Context, db interfaces.DB) error {
	_, err := db.Exec(
		ctx,
		`UPDATE users SET isadmin = $1 WHERE user_id = $2;`,
		u.IsAdmin, u.Id,
	)
	if err != nil {
		return fmt.Errorf("%w: when executing updating user role, user_id=%d, isadmin=%t", err, u.Id, u.IsAdmin)
	}

	return nil
}

func ViewManageUsers(ctx context.Context, db interfaces.DB) ([]ManageUser, error) {
	rows, err := db.Query(ctx, `
		SELECT u.user_id, u.name, u.wcaid, u.isadmin, c.name, c.iso2 
		FROM users u
		JOIN countries c ON c.country_id = u.country_id
		ORDER BY timestamp`,
	)
	if err != nil {
		return []ManageUser{}, fmt.Errorf("%w: when querying users", err)
	}
	defer rows.Close()

	manageUsers := make([]ManageUser, 0)
	for rows.Next() {
		manageUser := ManageUser{}
		err := rows.Scan(&manageUser.Id, &manageUser.Name, &manageUser.WcaId, &manageUser.IsAdmin, &manageUser.CountryName, &manageUser.CountryIso2)
		if err != nil {
			return []ManageUser{}, fmt.Errorf("%w: when scanning user", err)
		}

		manageUsers = append(manageUsers, manageUser)
	}

	if err := rows.Err(); err != nil {
		return []ManageUser{}, fmt.Errorf("%w: when iterating over rows", err)
	}

	return manageUsers, nil
}
