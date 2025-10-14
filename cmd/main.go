package main

import (
	"os"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/core"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/handlers"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/models"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// commands is the central channel to communicate with the manager.
var commands = make(chan models.ClientCommand)

func main() {
	utils.LoadEnvironmentVars()
	redisAddress := os.Getenv("REDIS_ADDRESS")
	rdb := core.InitializeRedisConnection(redisAddress)
	redisContext := utils.InitializeRedisContext()
	go core.ClientChannelManager(commands)
	go core.SubscribeToApprovalEvents(rdb, commands)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))
	router.POST("/submit", handlers.InitSubmitHandler(rdb, redisContext))
	router.GET("/events/:job_id", handlers.InitSSEHandler(commands))
	core.InitializeServer(router, rdb, redisContext)
}
