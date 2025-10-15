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

// clientCommands is the central channel to communicate with the manager.
var clientCommands = make(chan models.ClientCommand)

// connManagerCommands is the channel for managing SSE connections.
var connManagerCommands = make(chan models.ConnManagerCommand)

func main() {
	utils.LoadEnvironmentVars()
	redisAddress := os.Getenv("REDIS_ADDRESS")
	rdb := core.InitializeRedisConnection(redisAddress)
	redisContext := utils.GetRedisContext()
	go core.ClientChannelManager(clientCommands)
	go core.SSEConnectionManager(connManagerCommands) // Start the new connection manager
	go core.SubscribeToApprovalEvents(rdb, clientCommands)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))
	router.POST("/submit", handlers.InitSubmitHandler(rdb, redisContext))
	router.GET("/events/:job_id", handlers.InitSSEHandler(clientCommands, connManagerCommands))
	core.InitializeServer(router, rdb, redisContext, connManagerCommands)
}
