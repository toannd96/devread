package repo_impl

import (
	"backend-viblo-trending/custom_error"
	"backend-viblo-trending/db"
	"backend-viblo-trending/log"
	"backend-viblo-trending/model"
	"backend-viblo-trending/model/requests"
	"backend-viblo-trending/repository"
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type UserRepoImpl struct {
	sql *db.Sql
}

func NewUserRepo(sql *db.Sql) repository.UserRepo {
	return &UserRepoImpl{
		sql: sql,
	}

}

func (u *UserRepoImpl) SaveUser(context context.Context, user model.User) (model.User, error) {
	statement := `
		INSERT INTO users(user_id, email, password, role, full_name, created_at, updated_at)
		VALUES(:user_id, :email, :password, :role, :full_name, :created_at, :updated_at)
	`
	user.CreateAt = time.Now()
	user.UpdateAt = time.Now()

	_, err := u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		log.Error(err.Error())
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return user, custom_error.UserConflict
			}
		}
		return user, custom_error.SignUpFail
	}

	return user, nil

}

func (u *UserRepoImpl) CheckLogin(context context.Context, loginReq requests.RequestSignIn) (model.User, error) {
	var user = model.User{}
	err := u.sql.Db.GetContext(context, &user, "SELECT * FROM users WHERE email=$1", loginReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, custom_error.UserNotFound
		}
		log.Error(err.Error())
		return user, err
	}
	return user, nil
}
