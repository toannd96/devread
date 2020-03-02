package repository

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/model/requests"
	"context"
)

type UserRepo interface {
	CheckLogin(context context.Context, loginReq requests.RequestSignIn) (model.User, error)
	SaveUser(context context.Context, user model.User) (model.User, error)
	SelectUserById(context context.Context, userId string) (model.User, error)
	UpdateUser(context context.Context, user model.User) (model.User, error)
}
