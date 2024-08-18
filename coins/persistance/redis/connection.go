package redisRepository

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"my-stocks/coins/config"
)

var db *redis.Client

func GetConnection(cfg config.Redis) *redis.Client {
	if db != nil {
		return db
	}
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	db = conn

	return db
}
