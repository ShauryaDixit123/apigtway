package main

import (
	"apigtway/app"
	"apigtway/configs"
	"apigtway/configs/local"
	"apigtway/src/handlers"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("app started!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conf := local.GetConfig()
	// init redis
	rcl := app.NewRedisClient(&conf.Redis)
	_, err = rcl.Ping().Result()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	init := configs.InitRoutes{Eng: r, Rcl: rcl}
	handlers.InitilizeRoutes(init)
	fmt.Printf("api started on port!, %s", conf.Port)
	r.Run(":" + conf.Port)
}
