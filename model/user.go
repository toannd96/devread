package model

import "time"

type User struct {
	UserID       string    `json:"-" db:"user_id, omitempty"`
	FullName     string    `json:"full_name,omitempty" db:"full_name, omitempty"`
	Email        string    `json:"email,omitempty" db:"email, omitempty"`
	Password     string    `json:"-" db:"password, omitempty"`
	Token        string    `json:"token,omitempty"`
	CreateAt     time.Time `json:"-" db:"create_at, omitempty"`
	UpdateAt     time.Time `json:"-" db:"update_at, omitempty"`
	Verify       bool      `json:"-" db:"verify, omitempty"`
}
