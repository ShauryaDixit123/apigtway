package app

import (
	"apigtway/configs"
	"fmt"

	"github.com/go-redis/redis"
)

func NewRedisClient(config *configs.RedisConfig) *redis.Client {
	fmt.Println(config.Host + ":" + config.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		DB:       config.DB,
		Password: config.Password,
	})
	return client
}
