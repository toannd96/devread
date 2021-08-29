package db

import (
	"github.com/go-redis/redis"

	"go.uber.org/zap"
)

type RedisDB struct {
	Client *redis.Client
	Logger *zap.Logger

	// Dùng dưới local
	// Host     string
	// Port     string

	// Dùng trên server heroku
	Url string
}

func (rd *RedisDB) NewRedisDB() {
	// Dùng dưới local
	// rd.Client = redis.NewClient(&redis.Options{
	// 	Addr:     rd.Host + ":" + rd.Port,
	// })

	// Dùng trên server heroku
	opt, _ := redis.ParseURL(rd.Url)
	rd.Client = redis.NewClient(opt)

	_, err := rd.Client.Ping().Result()
	if err != nil {
		rd.Logger.Error("Kết nối không thành công tới redis ", zap.Error(err))
	}

	rd.Logger.Info("Kết nối thành công tới redis")
}
