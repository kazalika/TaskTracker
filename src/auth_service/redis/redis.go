package redis

import (
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRedisClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // Если требуется
		DB:       0,  // Номер БД в Redis, по умолчанию 0
	})
}

func GetRedisClient() *redis.Client {
	return rdb
}
