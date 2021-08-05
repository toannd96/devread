package repository

type AuthenRepo interface {
	CreateTokenMail(token string, userID string) error
	FetchTokenMail(token string) (string, error)
	DeleteTokenMail(token string) error
	InsertTokenMail(newKey string) error
}
