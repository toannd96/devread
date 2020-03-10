package security

import (
	"backend-viblo-trending/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const SECRET_KEY = "d@ct0@n96"

func CreateAccessToken(user model.User) (string, error) {
	accessClaims := &model.JwtCustomClaims{
		UserId: user.UserId,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	t, err := accessToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return t, nil
}

func CreateRefeshToken(user model.User) (string, error) {
	refreshClaims := &model.JwtCustomClaims{
		UserId: user.UserId,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rt, err := refreshToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return rt, nil
}
