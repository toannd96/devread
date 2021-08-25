package db

import (
	"github.com/go-redis/redis"

	"go.uber.org/zap"
)

type RedisDB struct {
	Client *redis.Client
	Host   string
	Port   string
	Logger *zap.Logger
}

func (rd *RedisDB) NewRedisDB() {
	rd.Client = redis.NewClient(&redis.Options{
		Addr: rd.Host + ":" + rd.Port,
	})
	_, err := rd.Client.Ping().Result()
	if err != nil {
		rd.Logger.Error("Kết nối không thành công tới redis ", zap.Error(err))
	}

	rd.Logger.Info("Kết nối thành công tới redis")
}
