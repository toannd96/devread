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

func (au *AuthRepoImpl) CreateAuthMail(token string, userID string) error {

	// 5 minute
	errAccess := au.client.Client.Set(token, userID, 300000000000).Err()
	if errAccess != nil {
		return errAccess
	}

	return nil
}

func (au *AuthRepoImpl) CreateAuthVerify(token string, email string) error {

	// 1 day
	errAccess := au.client.Client.Set(token, email, 86400000000000).Err()
	if errAccess != nil {
		return errAccess
	}

	return nil
}

func (au *AuthRepoImpl) InsertAuthMail(newKey string) error {
	count, errCount := au.client.Client.DbSize().Result()
	if errCount != nil {
		return errCount
	}

	allKey, err := au.client.Client.Keys("*").Result()
	if err != nil {
		return err
	}

	if count >= 2 {
		errToken := au.client.Client.Rename(allKey[:][0], newKey).Err()
		if errToken != nil {
			return errToken
		}

		count, errCount := au.client.Client.DbSize().Result()
		if errCount != nil {
			return errCount
		}

		if count >= 2 {
			errToken := au.client.Client.Rename(allKey[:][1], newKey).Err()
			if errToken != nil {
				return errToken
			}
		}
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

func (au *AuthRepoImpl) FetchAuthMail(token string) (string, error) {
	result, err := au.client.Client.Get(token).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (au *AuthRepoImpl) DeleteAccessToken(accessUUID string) error {
	deleteAt, err := au.client.Client.Del(accessUUID).Result()
	if err != nil || deleteAt != 1 {
		return err
	}
	return nil
}

func (au *AuthRepoImpl) DeleteTokenMail(token string) error {
	deleteAt, err := au.client.Client.Del(token).Result()
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

// func (au *AuthRepoImpl) InsertToken(token1 string, token2 string) error {
// 	totalToken, err := au.client.Client.DbSize().Result()
// 	if err != nil {
// 		return err
// 	}
// 	if totalToken == 1 {
// 		return
// 	}

// }
