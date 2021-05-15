package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type RedisConfig struct {
	rdb *redis.Client
}

func (rc *RedisConfig) Init() {
	//Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rc.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (rc *RedisConfig) GetConnection() *redis.Client {
	return rc.rdb
}
