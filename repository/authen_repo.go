package repository

type AuthenRepo interface {
	CreateTokenMail(token string, userID string) error
	CreateTokenVerify(token string, email string) error
	FetchToken(accessUUID string) (string, error)
	FetchTokenMail(token string) (string, error)
	DeleteTokenMail(token string) error
	InsertTokenMail(newKey string) error
}
