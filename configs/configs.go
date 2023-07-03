package configs

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Config struct {
	Host  string
	Port  string
	Redis RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}
type InitRoutes struct {
	Eng *gin.Engine
	Rcl *redis.Client
}
