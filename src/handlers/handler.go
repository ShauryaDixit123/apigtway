package handlers

import (
	"apigtway/configs"
	"apigtway/src/handlers/users"
)

func InitilizeRoutes(init configs.InitRoutes) {
	users.InitUserRoutes(init)
}
