package model

type Post struct {
	Name       string `json:"name" db:"name,omitempty"`
	Link       string `json:"link" db:"link,omitempty"`
	Tags       string `json:"tags" db:"tags,omitempty"`
	Bookmarked bool   `json:"-,omitempty"`
}
