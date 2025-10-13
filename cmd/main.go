package main

import (
	"os"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/app"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/handlers"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadEnvironmentVars()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))
	redisAddress := os.Getenv("REDIS_ADDRESS")
	rdb := app.InitializeRedisConnection(redisAddress)
	redisContext := utils.InitializeRedisContext()
	router.POST("/submit", handlers.InitSubmitHandler(rdb, redisContext))
	app.InitializeServer(router, rdb, redisContext)
}
