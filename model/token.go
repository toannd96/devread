package model

import "github.com/dgrijalva/jwt-go"

type TokenDetails struct {
	UserID       string
	jwt.StandardClaims

}
