package repo_impl

import (
	"backend-viblo-trending/custom_error"
	"backend-viblo-trending/db"
	"backend-viblo-trending/model"
	"backend-viblo-trending/model/req"
	"backend-viblo-trending/repository"
	"context"
	"database/sql"
	"github.com/lib/pq"
	"time"
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
		INSERT INTO users(user_id, email, password, full_name, verify, create_at, update_at)
		VALUES(:user_id, :email, :password, :full_name, :verify, :create_at, :update_at)
	`
	user.CreateAt = time.Now()
	user.UpdateAt = time.Now()

	_, err := u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return user, custom_error.UserConflict
			}
		}
		return user, custom_error.SignUpFail
	}

	return user, nil

}

func (u *UserRepoImpl) CheckSignIn(context context.Context, signinReq req.ReqSignIn) (model.User, error) {
	var user = model.User{}
	err := u.sql.Db.GetContext(context, &user, "SELECT * FROM users WHERE email=$1", signinReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, custom_error.UserNotFound
		}
		return user, err
	}
	return user, nil
}

func (u *UserRepoImpl) CheckEmail(context context.Context, emailReq req.ReqEmail) (model.User, error) {
	var user = model.User{}
	err := u.sql.Db.GetContext(context, &user, "SELECT * FROM users WHERE email=$1", emailReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, custom_error.UserNotFound
		}
		return user, err
	}
	return user, nil
}

func (u *UserRepoImpl) CheckEmailSignUp(context context.Context, emailReq req.ReqtSignUp) (model.User, error) {
	var user = model.User{}
	err := u.sql.Db.GetContext(context, &user, "SELECT * FROM users WHERE email=$1", emailReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, custom_error.UserNotFound
		}
		return user, err
	}
	return user, nil
}

func (u *UserRepoImpl) SelectUserByID(context context.Context, userID string) (model.User, error) {
	var user model.User
	err := u.sql.Db.GetContext(context, &user, "SELECT * FROM users WHERE user_id=$1", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, custom_error.UserNotFound
		}
		return user, err
	}
	return user, nil
}

func (u *UserRepoImpl) UpdateUser(context context.Context, user model.User) (model.User, error) {
	statement := `
	UPDATE users
	SET
		full_name = (CASE WHEN LENGTH(:full_name) = 0 THEN full_name ElSE :full_name END),
		password = (CASE WHEN LENGTH(:password) = 0 THEN password ElSE :password END),
		update_at = COALESCE (:update_at, update_at)
	WHERE user_id = :user_id
	`
	user.UpdateAt = time.Now()

	result, err := u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		return user, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return user, custom_error.UserNotUpdated
	}
	if count == 0 {
		return user, custom_error.UserNotUpdated
	}
	return user, nil
}

func (u *UserRepoImpl) UpdatePassword(context context.Context, user model.User) (model.User, error) {
	statement := `
	Update users
	SET
		password = (CASE WHEN LENGTH(:password) = 0 THEN password ElSE :password END),
		update_at = COALESCE (:update_at, update_at)
	WHERE user_id = :user_id
	`
	user.UpdateAt = time.Now()

	result, err := u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		return user, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return user, custom_error.UserNotUpdated
	}
	if count == 0 {
		return user, custom_error.UserNotUpdated
	}
	return user, nil
}

func (u *UserRepoImpl) UpdateVerify(context context.Context, user model.User) (model.User, error) {
	statement := `
	Update users
	SET
		verify = (CASE WHEN LENGTH(:verify) = 0 THEN verify ElSE :verify END),
		update_at = COALESCE (:update_at, update_at)
	WHERE user_id = :user_id
	`
	user.UpdateAt = time.Now()

	result, err := u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		return user, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return user, custom_error.UserNotUpdated
	}
	if count == 0 {
		return user, custom_error.UserNotUpdated
	}
	return user, nil
}
