package model

type Post struct {
	Name       string `json:"name" db:"name,omitempty"`
	Link       string `json:"link" db:"link,omitempty"`
	Tag        string `json:"tag" db:"tag,omitempty"`
	Bookmarked bool   `json:"bookmarked"`
}
