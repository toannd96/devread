package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type Sql struct {
	Db       *sqlx.DB
	Host     string
	Port     string
	UserName string
	Password string
	DbName   string
}

func (s *Sql) Connect() {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.UserName, s.Password, s.DbName)

	s.Db = sqlx.MustConnect(os.Getenv("DB_DRIVER"), dataSource)
	if err := s.Db.Ping(); err != nil {
		return
	}
	log.Println("Kết nối thành công tới postgres")
}

func (s *Sql) Close() {
	s.Db.Close()
}
