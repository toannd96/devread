package model

type Role int

const (
	MEMBER Role = iota
	ADMIN
	ADMIN1
	ADMIN2
)

func (r Role) String() string {
	return []string{"MEMBER", "ADMIN"}[r]
}
