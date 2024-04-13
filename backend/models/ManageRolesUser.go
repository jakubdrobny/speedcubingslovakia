package models

type ManageRolesUser struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Isadmin bool `json:"isadmin"`
}