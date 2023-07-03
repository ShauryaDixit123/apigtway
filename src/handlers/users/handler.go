package users

import (
	"apigtway/configs"

	"github.com/go-redis/redis"
)

type Handler struct {
	rcl *redis.Client
}

func InitUserRoutes(init configs.InitRoutes) {
	r := init.Eng
	tkn := r.Group("/tokens")
	hnd := Handler{
		rcl: init.Rcl,
	}
	{
		tkn.GET("/gen", hnd.GenAuth)
	}
}
