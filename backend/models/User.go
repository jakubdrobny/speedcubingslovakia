package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakubdrobny/speedcubingslovakia/backend/email"
	"github.com/jakubdrobny/speedcubingslovakia/backend/interfaces"
)

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CountryId   string `json:"country_id"`
	ContinentId string `json:"continent_id"`
	Sex         string `json:"sex"`
	WcaId       string `json:"wcaid"`
	IsAdmin     bool   `json:"isadmin"`
	Url         string `json:"url"`
	AvatarUrl   string `json:"avatarurl"`
	Email       string `json:"-"`
}

func (u *User) Exists(db *pgxpool.Pool) (bool, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT u.user_id, u.isadmin FROM users u WHERE u.wcaid = $1 AND u.name = $2;`,
		u.WcaId,
		u.Name,
	)
	if err != nil {
		return false, err
	}

	found := false
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.IsAdmin)
		if err != nil {
			return false, err
		}
		found = true
	}

	return found, nil
}

func (u *User) Update(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE users SET country_id = $1, sex = $2, url = $3, avatarurl = $4, isadmin = $5, timestamp = CURRENT_TIMESTAMP, email = $6 WHERE wcaid = $7 AND name = $8;`,
		u.CountryId,
		u.Sex,
		u.Url,
		u.AvatarUrl,
		u.IsAdmin,
		u.Email,
		u.WcaId,
		u.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) CreateAllAnnouncementReadConnection(db *pgxpool.Pool) error {
	rows, err := db.Query(context.Background(), `SELECT a.announcement_id FROM announcements a;`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var announcementId int
		err = rows.Scan(&announcementId)
		if err != nil {
			return nil
		}

		_, err = db.Exec(
			context.Background(),
			`INSERT INTO announcement_read (announcement_id, user_id, read, read_timestamp) VALUES ($1,$2,FALSE,NULL);`,
			announcementId,
			u.Id,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) Insert(db *pgxpool.Pool) error {
	err := db.QueryRow(context.Background(), `INSERT INTO users (name, country_id, sex, url, avatarurl, wcaid, isadmin, email) VALUES ($1,$2,$3,$4,$5,$6,false,$7) RETURNING user_id;`, u.Name, u.CountryId, u.Sex, u.Url, u.AvatarUrl, u.WcaId, u.Email).
		Scan(&u.Id)
	if err != nil {
		return err
	}
	exists, err := u.Exists(db)
	if !exists || err != nil {
		return fmt.Errorf("%s %t", err, exists)
	}

	err = u.CreateAllAnnouncementReadConnection(db)
	if err != nil {
		return err
	}

	return nil
}

func GetUserById(db interfaces.DB, uid int) (User, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT u.user_id, u.name, u.country_id, u.sex, u.wcaid, u.isadmin, u.url, u.avatarurl, u.email FROM users u WHERE u.user_id = $1;`,
		uid,
	)
	if err != nil {
		return User{}, err
	}

	var user User
	found := false
	for rows.Next() {
		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.CountryId,
			&user.Sex,
			&user.WcaId,
			&user.IsAdmin,
			&user.Url,
			&user.AvatarUrl,
			&user.Email,
		)
		if err != nil {
			return User{}, err
		}
		found = true
	}

	if !found {
		return User{}, err
	}

	return user, nil
}

func GetUserByWCAID(db *pgxpool.Pool, wcaid string) (int, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT u.user_id FROM users u WHERE u.wcaid = $1;`,
		wcaid,
	)
	if err != nil {
		return 0, err
	}

	var uid int
	for rows.Next() {
		err = rows.Scan(&uid)
		if err != nil {
			return 0, err
		}
	}

	return uid, nil
}

func GetUserByName(db *pgxpool.Pool, name string) (int, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT u.user_id FROM users u WHERE u.name = $1;`,
		name,
	)
	if err != nil {
		return 0, err
	}

	var uid int
	for rows.Next() {
		err = rows.Scan(&uid)
		if err != nil {
			return 0, err
		}
	}

	return uid, nil
}

func GetEmailByWCAID(db *pgxpool.Pool, wcaid string) (string, error) {
	var email string
	err := db.QueryRow(context.Background(), `SELECT u.email FROM users u WHERE u.wcaid = $1;`, wcaid).
		Scan(&email)
	return email, err
}

func GetUserInfoFromWCA(authInfo *AuthorizationInfo, envMap map[string]string) (User, error) {
	bearer := "Bearer " + authInfo.AccessToken
	req, err := http.NewRequest("GET", envMap["WCA_API_ME_URL"], nil)
	if err != nil {
		return User{}, err
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return User{}, err
	}
	defer res.Body.Close()
	type Country struct {
		Id string `json:"id"`
	}
	type Avatar struct {
		Url string `json:"url"`
	}
	type ME struct {
		Name    string  `json:"name"`
		WcaId   string  `json:"wca_id"`
		Sex     string  `json:"gender"`
		Url     string  `json:"url"`
		Country Country `json:"country"`
		Avatar  Avatar  `json:"avatar"`
		Email   string  `json:"email"`
	}
	type WCAApiMe struct {
		Me ME `json:"me"`
	}

	var apiMe WCAApiMe
	err = json.NewDecoder(res.Body).Decode(&apiMe)
	if err != nil {
		return User{}, err
	}

	user := User{}
	user.Name = apiMe.Me.Name
	user.CountryId = apiMe.Me.Country.Id
	user.Sex = apiMe.Me.Sex
	user.WcaId = apiMe.Me.WcaId
	user.IsAdmin = false
	user.Url = apiMe.Me.Url
	user.AvatarUrl = apiMe.Me.Avatar.Url
	user.Email = apiMe.Me.Email

	return user, nil
}

func (u *User) LoadContinent(db *pgxpool.Pool) error {
	rows, err := db.Query(
		context.Background(),
		`SELECT continents.continent_id FROM continents JOIN countries ON countries.continent_id = continents.continent_id WHERE countries.country_id = $1;`,
		u.CountryId,
	)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&u.ContinentId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u User) SendNewUserMailAsync(ctx context.Context, db interfaces.DB, envMap map[string]string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in SendSuspicousMailAsync goroutine", r)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Request was cancelled. Aborting suspicous mail send.")
		return nil
	default:
	}

	var order int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&order)
	if err != nil {
		return fmt.Errorf("%w: when querying all users from db", err)
	}

	var mailSubject string
	if envMap["NODE_ENV"] == "development" {
		mailSubject += "DEVELOPMENT: "
	}
	mailSubject += "New user registered!!!"

	profileLink := u.WcaId
	if u.WcaId == "" {
		profileLink = u.Name
	}
	content :=
		"<b>Username + WCA ID:</b> <a href=\"" + envMap["WEBSITE_HOME"] + "/profile/" + profileLink + "\">" + u.Name + "</a> (" + profileLink + ")<br>" +
			"<b>Email:</b> " + u.Email + "<br>" +
			"<b>Sex:</b> " + u.Sex + "<br>" +
			"<b>Country:</b> " + u.CountryId + "<br>" +
			"<b>User no. " + strconv.Itoa(order) + "</b>"

	err = email.SendMail(
		envMap["MAIL_USERNAME"],
		envMap["MAIL_USERNAME"],
		mailSubject,
		content,
		envMap,
	)
	if err != nil {
		return fmt.Errorf("%w: when sending email about new user", err)
	}

	log.Println("Successfully sent mail about new user.")

	return nil
}

func FindFuzzyDuplicateUser(ctx context.Context, db interfaces.DB, userID int) (User, bool, error) {
	var duplicate User
	err := db.QueryRow(ctx, `
		SELECT u2.user_id, u2.name, u2.country_id, u2.sex, u2.wcaid, u2.isadmin, u2.url, u2.avatarurl, u2.email
		FROM users u1
		JOIN users u2 ON u1.user_id < u2.user_id
		WHERE u1.user_id = $1
		  AND (unaccent(u1.name) ILIKE '%' || unaccent(u2.name) || '%' OR unaccent(u2.name) ILIKE '%' || unaccent(u1.name) || '%')
		  AND u1.name <> '' AND u2.name <> ''
		ORDER BY u2.user_id
		LIMIT 1;
	`, userID).Scan(
		&duplicate.Id,
		&duplicate.Name,
		&duplicate.CountryId,
		&duplicate.Sex,
		&duplicate.WcaId,
		&duplicate.IsAdmin,
		&duplicate.Url,
		&duplicate.AvatarUrl,
		&duplicate.Email,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, false, nil
		}
		return User{}, false, fmt.Errorf("%w: when querying error for duplicate user", err)
	}

	return duplicate, true, nil
}

func MergeUsers(ctx context.Context, db interfaces.DB, oldUserID, newUserID int) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: when starting db transaction", err)
	}
	defer tx.Rollback(ctx)

	var areNamesSimilar bool
	validationQuery := `
		SELECT
			(unaccent(u1.name) ILIKE '%' || unaccent(u2.name) || '%' OR
			 unaccent(u2.name) ILIKE '%' || unaccent(u1.name) || '%')
		FROM users u1, users u2
		WHERE u1.user_id = $1 AND u2.user_id = $2;`

	err = tx.QueryRow(ctx, validationQuery, oldUserID, newUserID).Scan(&areNamesSimilar)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("validation failed: one or both user IDs (%d, %d) do not exist", oldUserID, newUserID)
		}
		return fmt.Errorf("%w: when executing validation query", err)
	}

	if !areNamesSimilar {
		return fmt.Errorf("validation failed: users %d and %d do not have similar names", oldUserID, newUserID)
	}

	type foreignKeyInfo struct {
		ChildTable  string
		ChildColumn string
	}

	var fks []foreignKeyInfo
	fkQuery := `
		SELECT
			conrelid::regclass AS child_table,
			a.attname AS child_column
		FROM
			pg_constraint AS c
		JOIN pg_attribute AS a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
		WHERE
			c.contype = 'f'
			AND c.confrelid = 'users'::regclass
	`
	rows, err := tx.Query(ctx, fkQuery)
	if err != nil {
		return fmt.Errorf("%w: when querying foreign key constraints", err)
	}

	for rows.Next() {
		var fk foreignKeyInfo
		if err := rows.Scan(&fk.ChildTable, &fk.ChildColumn); err != nil {
			rows.Close()
			return fmt.Errorf("%w: when scanning fk record", err)
		}
		fks = append(fks, fk)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return fmt.Errorf("%w: when iteraing through rows", rows.Err())
	}

	for _, fk := range fks {
		updateQuery := fmt.Sprintf(`UPDATE %s SET %s = $1 WHERE %s = $2`, fk.ChildTable, fk.ChildColumn, fk.ChildColumn)
		if _, err := tx.Exec(ctx, updateQuery, newUserID, oldUserID); err != nil {
			return fmt.Errorf("%w: when executing update child table %s query", err, fk.ChildTable)
		}
	}

	if _, err := tx.Exec(ctx, `DELETE FROM users WHERE user_id = $1`, oldUserID); err != nil {
		return fmt.Errorf("%w: when executing delete old user", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: when commiting transaction", err)
	}

	return nil
}

func SearchUsers(ctx context.Context, db interfaces.DB, query string) ([]ManageUser, error) {
	sql := `
		SELECT
			u.user_id,
			u.name,
			u.wcaid,
			c.name,
			c.iso2,
			u.isadmin
		FROM
			users u
		LEFT JOIN
			countries c ON u.country_id = c.country_id
		WHERE
			unaccent(u.name) ILIKE unaccent($1)
		ORDER BY
			u.name
		LIMIT 20;`

	searchPattern := "%" + query + "%"

	rows, err := db.Query(ctx, sql, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("%w: when executing user search query", err)
	}
	defer rows.Close()

	users := make([]ManageUser, 0)
	for rows.Next() {
		var user ManageUser
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.WcaId,
			&user.CountryName,
			&user.CountryIso2,
			&user.IsAdmin,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: when scanning user row", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []ManageUser{}, fmt.Errorf("%w: when iterating over rows", err)
	}

	return users, nil
}
