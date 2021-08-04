package model

import "github.com/golang-jwt/jwt"

type TokenDetails struct {
	UserID string
	jwt.StandardClaims
}
