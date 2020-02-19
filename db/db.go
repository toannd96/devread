package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

type Sql struct {
	Db       *sqlx.DB
	Host     string
	Port     int
	UserName string
	Password string
	DbName   string
}

func (s *Sql) Connect() {
	dataSource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.UserName, s.Password, s.DbName)

	s.Db = sqlx.MustConnect("postgres", dataSource)
	if err := s.Db.Ping(); err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println("Connect db ok")
}

func (s *Sql) Close() {
	s.Db.Close()
}
