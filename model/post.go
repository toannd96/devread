package model

type Post struct {
	Name          string    `json:"name" db:"name,omitempty"`
	Link          string    `json:"link" db:"link,omitempty"`
	Bookmarked    bool      `json:"-,omitempty"`
}
