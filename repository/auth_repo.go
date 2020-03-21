package repository

import (
	"backend-viblo-trending/model"
)

type AuthRepo interface {
	CreateAuth(userID string, tokenDetails *model.TokenDetails) error
	FetchAuth(accessUUID string) (string, error)
	DeleteAccessToken(accessUUID string) error
	DeleteRefreshToken(refresUUID string) error
}
