package repo_impl

import (
	"backend-viblo-trending/db"
	"backend-viblo-trending/model"
	"backend-viblo-trending/repository"
	"time"
)

type AuthRepoImpl struct {
	client *db.RedisDB
}

func NewAuthRepo(client *db.RedisDB) repository.AuthRepo {
	return &AuthRepoImpl{
		client: client,
	}
}

func (au *AuthRepoImpl) CreateAuth(userID string, tokenDetails *model.TokenDetails) error {
	at := time.Unix(tokenDetails.AtExpires, 0) // converting Unix to UTC(to Time object)
	rt := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	errAccess := au.client.Client.Set(tokenDetails.AccessUUID, userID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := au.client.Client.Set(tokenDetails.RefreshUUID, userID, rt.Sub(now)).Err()

	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (au *AuthRepoImpl) FetchAuth(accessUUID string) (string, error) {
	userID, err := au.client.Client.Get(accessUUID).Result()
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (au *AuthRepoImpl) DeleteAccessToken(accessUUID string) error {
	deleteAt, err := au.client.Client.Del(accessUUID).Result()
	if err != nil || deleteAt != 1 {
		return err
	}
	return nil
}

func (au *AuthRepoImpl) DeleteRefreshToken(refresUUID string) error {
	deleteRt, err := au.client.Client.Del(refresUUID).Result()
	if err != nil || deleteRt != 1 {
		return err
	}
	return nil
}
