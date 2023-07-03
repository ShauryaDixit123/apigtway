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
	ro := r.Group("/tokens")
	hnd := Handler{
		rcl: init.Rcl,
	}
	ro.GET("/gen", hnd.GenAuth)
	ro.POST("/refauth", hnd.RefreshAuth)
	ro.Use(hnd.CheckAuth())
	{
		// secured routes
		ro.GET("/ping", hnd.PingWithAuth)

	}
}
