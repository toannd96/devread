package repository

import (
	"context"

	"devread/model"
	"devread/model/req"
)

type UserRepo interface {
	CheckSignIn(context context.Context, signinReq req.ReqSignIn) (model.User, error)
	CheckEmail(context context.Context, emailReq req.ReqSignUp) (model.User, error)
	UpdateUser(context context.Context, user model.User) (model.User, error)
	UpdatePassword(context context.Context, user model.User) (model.User, error)
	UpdateVerify(context context.Context, user model.User) (model.User, error)
	SaveUser(context context.Context, user model.User) (model.User, error)
	SelectUserByID(context context.Context, userID string) (model.User, error)
}
