package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Sql struct {
	Db       *sqlx.DB
	Host     string
	Port     string
	UserName string
	Password string
	DbName   string
	Logger   *zap.Logger
}

func (s *Sql) Connect() {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		s.Host, s.Port, s.UserName, s.Password, s.DbName)

	s.Db = sqlx.MustConnect(os.Getenv("DB_DRIVER"), dataSource)
	if err := s.Db.Ping(); err != nil {
		s.Logger.Error("Kết nối không thành công tới postgres ", zap.Error(err))
		return
	}
	s.Logger.Info("Kết nối thành công tới postgres")
}

func (s *Sql) Close() {
	s.Db.Close()
}
