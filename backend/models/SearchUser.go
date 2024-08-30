package models

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type SearchUser struct {
	Id int `json:"-"`
	Username string `json:"username"`
	WCAID string `json:"wcaid"`
}

func GetUsersFromDB(tx pgx.Tx, query string) ([]SearchUser, int, string, string) {
	searchUsers := make([]SearchUser, 0)

	if query == "_" {
		rows, err := tx.Query(context.Background(), `SELECT u.user_id, u.name, (CASE WHEN u.wcaid LIKE '' THEN u.name ELSE u.wcaid END) AS wcaid FROM users u ORDER BY u.name;`)
		if err != nil {
			return []SearchUser{}, http.StatusInternalServerError, "ERR tx.Query all in GetUsersFromDB (" + query + "): " + err.Error(), "Failed querying all users."
		}

		for rows.Next() {
			var searchUser SearchUser
			err = rows.Scan(&searchUser.Id, &searchUser.Username, &searchUser.WCAID)
			if err != nil {
				return []SearchUser{}, http.StatusInternalServerError, "ERR scanning all in GetUsersFromDB with query (" + query + "): " + err.Error(), "Failed querying all users."
			}

			searchUsers = append(searchUsers, searchUser)
		}
	} else {
		rows, err := tx.Query(context.Background(), `SELECT u.user_id, u.name, (CASE WHEN u.wcaid LIKE '' THEN u.name ELSE u.wcaid END) AS wcaid FROM users u WHERE u.wcaid LIKE $1 ORDER BY u.name;`, query)
		if err != nil {
			return []SearchUser{}, http.StatusInternalServerError, "ERR tx.Query wcaid in GetUsersFromDB (" + query + "): " + err.Error(), "Failed querying users by WCAID."
		}

		for rows.Next() {
			var searchUser SearchUser
			err = rows.Scan(&searchUser.Id, &searchUser.Username, &searchUser.WCAID)
			if err != nil {
				return []SearchUser{}, http.StatusInternalServerError, "ERR scanning wcaid in GetUsersFromDB (" + query + "): " + err.Error(), "Failed querying users by WCAID."
			}

			searchUsers = append(searchUsers, searchUser)
		}

		if len(searchUsers) == 0 {
			rows, err := tx.Query(context.Background(), `SELECT u.user_id, u.name, (CASE WHEN u.wcaid LIKE '' THEN u.name ELSE u.wcaid END) AS wcaid FROM users u WHERE LOWER(u.name) LIKE LOWER('%' || $1 || '%') ORDER BY u.name;`, query)
			if err != nil {
				return []SearchUser{}, http.StatusInternalServerError, "ERR tx.Query name in GetUsersFromDB (" + query + "): " + err.Error(), "Failed querying users by name."
			}

			for rows.Next() {
				var searchUser SearchUser
				err = rows.Scan(&searchUser.Id, &searchUser.Username, &searchUser.WCAID)
				if err != nil {
					return []SearchUser{}, http.StatusInternalServerError, "ERR scanning name in GetUsersFromDB (" + query + "): " + err.Error(), "Failed querying users by name."
				}

				searchUsers = append(searchUsers, searchUser)
			}
		}
	}

	return searchUsers, http.StatusOK, "", ""
}

func (u *SearchUser) CreateAnnouncementReadConnection(tx pgx.Tx, aid int) error {
	_, err := tx.Exec(context.Background(), `INSERT INTO announcement_read (announcement_id, user_id, read, read_timestamp) VALUES ($1,$2,FALSE,NULL);`, aid, u.Id)
	return err
}