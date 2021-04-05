package config

import (
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	rdb *redis.Client
}

func (rc RedisConfig) Init() {
	rc.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (rc RedisConfig) GetConnection() *redis.Client {
	return rc.rdb
}
