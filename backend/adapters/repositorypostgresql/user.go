package repository

import (
	"context"
	"fmt"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/views"
)

type User struct {
	DB models.DbExecutor
}

func NewUser(db models.DbExecutor) *User {
	return &User{
		DB: db,
	}
}

func (r *User) UpdateRole(db models.DbExecutor, user *models.User) error {
	_, err := db.ExecContext(
		context.Background(),
		`UPDATE users u SET u.isadmin = $1 WHERE u.user_id = $2;`,
		user.IsAdmin, user.Id,
	)
	if err != nil {
		return fmt.Errorf("%w: when executing updating user role, user_id=%d, isadmin=%t", err, user.Id, user.IsAdmin)
	}

	return nil
}

func (r *User) ViewManageUsers(ctx context.Context, db models.DbExecutor) ([]views.ManageUser, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT u.user_id, u.name, u.wcaid, u.isadmin, c.name, c.iso2 
		FROM users u
		JOIN countries c ON c.country_id = u.country_id
		ORDER BY timestamp`,
	)
	if err != nil {
		return []views.ManageUser{}, fmt.Errorf("%w: when querying users", err)
	}
	defer rows.Close()

	manageUsers := make([]views.ManageUser, 0)
	for rows.Next() {
		manageUser := views.ManageUser{}
		err := rows.Scan(&manageUser.Id, &manageUser.Name, &manageUser.WcaId, &manageUser.IsAdmin, &manageUser.Country, &manageUser.CountryIso2)
		if err != nil {
			return []views.ManageUser{}, fmt.Errorf("%w: when scanning user", err)
		}

		manageUsers = append(manageUsers, manageUser)
	}

	if err := rows.Err(); err != nil {
		return []views.ManageUser{}, fmt.Errorf("%w: when iterating over rows", err)
	}

	return manageUsers, nil
}
