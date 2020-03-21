package repository

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/model/req"
	"context"
)

type UserRepo interface {
	CheckSignIn(context context.Context, SignInReq req.ReqSignIn) (model.User, error)
	SaveUser(context context.Context, user model.User) (model.User, error)
	SelectUserByID(context context.Context, userID string) (model.User, error)
	UpdateUser(context context.Context, user model.User) (model.User, error)
}
