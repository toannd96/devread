package repository

import "tech_posts_trending/model"

type AuthRepo interface {
	CreateAuth(userID string, tokenDetails *model.TokenDetails) error
	CreateAuthMail(token string, userID string) error
	CreateAuthVerify(token string, email string) error
	FetchAuth(accessUUID string) (string, error)
	FetchAuthMail(token string) (string, error)
	DeleteAccessToken(accessUUID string) error
	DeleteRefreshToken(refresUUID string) error
	DeleteTokenMail(token string) error
	InsertAuthMail(newKey string) error
}
