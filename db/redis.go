package db

import (
	"github.com/go-redis/redis"
	"log"
)

type RedisDB struct {
	Client *redis.Client
	Host   string
	Port   string
}

func (rd *RedisDB) NewRedisDB() {
	rd.Client = redis.NewClient(&redis.Options{
		Addr: rd.Host + ":" + rd.Port,
	})
	_, err := rd.Client.Ping().Result()
	if err != nil {
		log.Println(err)
	}
	log.Println("Kết nối thành công tới redis")
}
