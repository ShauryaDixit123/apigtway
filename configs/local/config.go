package local

import (
	"apigtway/configs"
	"os"
	"strconv"
)

func GetConfig() configs.Config {

	rdb, er := strconv.Atoi(os.Getenv("REDIS_DB"))
	if er != nil {
		rdb = 0
	}
	return configs.Config{
		Host: os.Getenv("TIER"),
		Port: os.Getenv("API_PORT"),
		Redis: configs.RedisConfig{
			Host:     os.Getenv("TIER"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       rdb,
		},
	}
}
