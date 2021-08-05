package repo_impl

import (
	"devread/db"
	"devread/repository"
)

type AuthenRepoImpl struct {
	client *db.RedisDB
}

func NewAuthenRepo(client *db.RedisDB) repository.AuthenRepo {
	return &AuthenRepoImpl{
		client: client,
	}
}

func (au *AuthenRepoImpl) CreateTokenMail(token string, userID string) error {

	// 5 minute = 300000000000
	// 1 day
	errAccess := au.client.Client.Set(token, userID, 86400000000000).Err()
	if errAccess != nil {
		return errAccess
	}

	return nil
}

func (au *AuthenRepoImpl) InsertTokenMail(newKey string) error {
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

func (au *AuthenRepoImpl) FetchTokenMail(token string) (string, error) {
	result, err := au.client.Client.Get(token).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (au *AuthenRepoImpl) DeleteTokenMail(token string) error {
	deleteAt, err := au.client.Client.Del(token).Result()
	if err != nil || deleteAt != 1 {
		return err
	}
	return nil
}
