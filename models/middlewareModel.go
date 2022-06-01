package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	User_id   int    `json:"user_id"`
	User_name string `json:"user_name"`
	jwt.StandardClaims
}

type Signings struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	AccessTokenExpiration string `json:"access_token_expiration"`
}
