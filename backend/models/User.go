package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	CountryId string `json:"country_id"`
	ContinentId string `json:"continent_id"`
	Sex string `json:"sex"`
	WcaId string `json:"wcaid"`
	IsAdmin bool `json:"isadmin"`
	Url string `json:"url"`
	AvatarUrl string `json:"avatarurl"`
}

func (u *User) Exists(db *pgxpool.Pool) (bool, error) {
	rows, err := db.Query(context.Background(), `SELECT u.user_id, u.isadmin FROM users u WHERE u.wcaid = $1 AND u.name = $2;`, u.WcaId, u.Name)
	if err != nil { return false, err }

	found := false
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.IsAdmin)
		if err != nil { return false, err }
		found = true
	}

	return found, nil
}

func (u *User) Update(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `UPDATE users SET country_id = $1, sex = $2, url = $3, avatarurl = $4, isadmin = $5, timestamp = CURRENT_TIMESTAMP WHERE wcaid = $6 AND name = $7;`, u.CountryId, u.Sex, u.Url, u.AvatarUrl, u.IsAdmin, u.WcaId, u.Name)
	if err != nil { return err }

	return nil
}

func (u *User) CreateAllAnnouncementReadConnection(db *pgxpool.Pool) error {
	rows, err := db.Query(context.Background(), `SELECT a.announcement_id FROM announcements a;`)
	if err != nil { return err }

	for rows.Next() {
		var announcementId int
		err = rows.Scan(&announcementId)
		if err != nil { return nil }

		_, err = db.Exec(context.Background(), `INSERT INTO announcement_read (announcement_id, user_id, read, read_timestamp) VALUES ($1,$2,FALSE,NULL);`, announcementId, u.Id)
		if err != nil { return err }
	}

	return nil
}

func (u *User) Insert(db *pgxpool.Pool) error {
	err := db.QueryRow(context.Background(), `INSERT INTO users (name, country_id, sex, url, avatarurl, wcaid, isadmin) VALUES ($1,$2,$3,$4,$5,$6,false) RETURNING user_id;`, u.Name, u.CountryId, u.Sex, u.Url, u.AvatarUrl, u.WcaId).Scan(&u.Id)
	if err != nil { return err }
	exists, err := u.Exists(db)
	if !exists || err != nil { return fmt.Errorf("%s %t", err, exists) }

	err = u.CreateAllAnnouncementReadConnection(db)
	if err != nil { return err }

	return nil
}


func GetUserById(db *pgxpool.Pool, uid int) (User, error) {
	rows, err := db.Query(context.Background(), `SELECT u.user_id, u.name, u.country_id, u.sex, u.wcaid, u.isadmin, u.url, u.avatarurl FROM users u WHERE u.user_id = $1;`, uid);
	if err != nil { return User{}, err }

	var user User
	found := false
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.CountryId, &user.Sex, &user.WcaId, &user.IsAdmin, &user.Url, &user.AvatarUrl)
		if err != nil { return User{}, err }
		found = true
	}

	if !found { return User{}, err }
	
	return user, nil
}

func GetUserByWCAID(db *pgxpool.Pool, wcaid string) (int, error) {
	rows, err := db.Query(context.Background(), `SELECT u.user_id FROM users u WHERE u.wcaid = $1;`, wcaid);
	if err != nil { return 0, err }

	var uid int
	for rows.Next() {
		err = rows.Scan(&uid)
		if err != nil { return 0, err }
	}

	return uid, nil
}

func GetUserByName(db *pgxpool.Pool, name string) (int, error) {
	rows, err := db.Query(context.Background(), `SELECT u.user_id FROM users u WHERE u.name = $1;`, name);
	if err != nil { return 0, err }

	var uid int
	for rows.Next() {
		err = rows.Scan(&uid)
		if err != nil { return 0, err }
	}

	return uid, nil
}

func GetUserInfoFromWCA(authInfo *AuthorizationInfo, envMap map[string]string) (User, error) {
	bearer := "Bearer " + authInfo.AccessToken
	req, err := http.NewRequest("GET", envMap["WCA_API_ME_URL"], nil)
	if err != nil { return User{}, err }

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK { return User{}, err }
	defer res.Body.Close()
	type Country struct { Id string `json:"id"`}
	type Avatar struct { Url string `json:"url"`}
	type ME struct {
		Name string `json:"name"`
		WcaId string `json:"wca_id"`
		Sex string `json:"gender"`
		Url string `json:"url"`
		Country Country `json:"country"`
		Avatar Avatar `json:"avatar"`
	}
	type WCAApiMe struct { Me ME `json:"me"` }

	var apiMe WCAApiMe
	err = json.NewDecoder(res.Body).Decode(&apiMe)
	if err != nil { return User{}, err }

	user := User{}
	user.Name = apiMe.Me.Name
	user.CountryId = apiMe.Me.Country.Id
	user.Sex = apiMe.Me.Sex
	user.WcaId = apiMe.Me.WcaId
	user.IsAdmin = false
	user.Url = apiMe.Me.Url
	user.AvatarUrl = apiMe.Me.Avatar.Url

	return user, nil
}

func (u *User) LoadContinent(db *pgxpool.Pool) (error) {
	rows, err := db.Query(context.Background(), `SELECT continents.continent_id FROM continents JOIN countries ON countries.continent_id = continents.continent_id WHERE countries.country_id = $1;`, u.CountryId);
	if err != nil { return err }

	for rows.Next() {
		err = rows.Scan(&u.ContinentId)
		if err != nil { return err }
	}

	return nil
}