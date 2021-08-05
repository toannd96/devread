package db

import (
	"devread/log"

	"github.com/go-redis/redis"

	"go.uber.org/zap"
)

type RedisDB struct {
	Client *redis.Client
	Host   string
	Port   string
}

func (rd *RedisDB) NewRedisDB() {
	log := log.WriteLog()
	rd.Client = redis.NewClient(&redis.Options{
		Addr: rd.Host + ":" + rd.Port,
	})
	_, err := rd.Client.Ping().Result()
	if err != nil {
		log.Error("Kết nối không thành công tới redis ", zap.Error(err))
	}

	log.Info("Kết nối thành công tới redis")
}
