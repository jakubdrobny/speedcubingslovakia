package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type AuthorizationInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	WcaId       string `json:"wcaid"`
	AvatarUrl   string `json:"avatarUrl"`
	IsAdmin     bool   `json:"isadmin"`
	Username    string `json:"username"`
}

func GetAuthInfo(code string, envMap map[string]string) (AuthorizationInfo, error) {
	res, err := http.PostForm(envMap["WCA_TOKEN_URL"], url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {envMap["WCA_CLIENT_ID"]},
		"client_secret": {envMap["WCA_CLIENT_SECRET"]},
		"code":          {code},
		"redirect_uri":  {envMap["WCA_REDIRECT_URI"]},
	})
	if err != nil {
		slog.Error("failed to send request to WCA token URL", "error", err)
		return AuthorizationInfo{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("failed to read response body", "error", err)
		return AuthorizationInfo{}, err
	}

	slog.Info("Raw response from WCA token endpoint", "status_code", res.StatusCode, "body", string(bodyBytes))

	if res.StatusCode != http.StatusOK {
		return AuthorizationInfo{}, fmt.Errorf("WCA token request failed with status %d: %s", res.StatusCode, string(bodyBytes))
	}

	var authInfo AuthorizationInfo
	err = json.Unmarshal(bodyBytes, &authInfo)
	if err != nil {
		return AuthorizationInfo{}, err
	}

	return authInfo, nil
}
