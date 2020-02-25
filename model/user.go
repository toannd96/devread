package model

import (
	"time"
)

type User struct {
	UserId   string
	FullName string
	Email    string
	Password string
	Role     string
	CreateAt time.Time
	UpdateAt time.Time
	Token    string
}
