package models

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type AuthorizationInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	WcaId string `json:"wcaid"`
	AvatarUrl string `json:"avatarUrl"`
	IsAdmin bool `json:"isadmin"`
}


func GetAuthInfo(code string, envMap map[string]string) (AuthorizationInfo, error) {
	res, err := http.PostForm(envMap["WCA_TOKEN_URL"], url.Values{
		"grant_type": {"authorization_code"},
		"client_id": {envMap["WCA_CLIENT_ID"]},
		"client_secret": {envMap["WCA_CLIENT_SECRET"]},
		"code": {code},
		"redirect_uri": {envMap["WCA_REDIRECT_URI"]},
	})
	if err != nil || res.StatusCode != http.StatusOK { return AuthorizationInfo{}, err }
	defer res.Body.Close()

	var authInfo AuthorizationInfo
	err = json.NewDecoder(res.Body).Decode(&authInfo)
	if err != nil { return AuthorizationInfo{}, err }

	return authInfo, nil
}