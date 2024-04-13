package models

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthDetails struct {
	UserId int `json:"userid"`
	ExpiresIn int64 `json:"expiresin"`
}


func ExtractPayload(tokenString *jwt.Token) (AuthDetails, error) {
	claims, ok := tokenString.Claims.(jwt.MapClaims)

	if !ok { return AuthDetails{}, fmt.Errorf("extracting payload from token failed") }

	var authDetails AuthDetails
	uidFloat, ok := claims["userid"].(float64)
	if !ok { return AuthDetails{}, fmt.Errorf("failed to parse userid") }
	authDetails.UserId = int(uidFloat)

	return authDetails, nil
}

func VerifyJWTToken(tokenString string, secretKey string) (AuthDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return []byte(secretKey), nil
	})
   
	if err != nil { return AuthDetails{}, err }
	if !token.Valid { return AuthDetails{}, fmt.Errorf("invalid token") }
   
	authDetails, err := ExtractPayload(token)
	if err != nil { return AuthDetails{}, err }

	return authDetails, nil
 }

func GetAuthDetailsFromHeader(c *gin.Context, secretKey string) (AuthDetails, error) {
	headers := c.Request.Header["Authorization"]
	if len(headers) <= 0 { return AuthDetails{}, fmt.Errorf("auth header missing") }

	header := strings.Split(headers[0], " ")
	if len(header) < 2 || header[0] != "Bearer" { return AuthDetails{}, fmt.Errorf("bad auth header") }

	authDetails, err := VerifyJWTToken(header[1], secretKey)
	if err != nil { return AuthDetails{}, err }

	return authDetails, nil
}