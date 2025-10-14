package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)


func InitializeServer(router *gin.Engine, rdb *redis.Client, redisContext context.Context) {
	workerCtx, cancelWorkers := utils.InitializeWorkerContext()
	numWorkers := 3
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		workerID := fmt.Sprintf("worker-%d", i)
		go startWorker(workerCtx, workerID, &wg, *rdb, redisContext)
	}

	srv := &http.Server {
		Addr: ":8080",
		Handler: router,
	}

	// รัน Server ใน Goroutine เพื่อป้องกันการ Block
	go func () {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	handleShutdownGracefully(cancelWorkers, &wg, srv, rdb)
}

func startWorker(workerCtx context.Context, workerID string, wg *sync.WaitGroup, rdb redis.Client, redisContext context.Context) {
	defer wg.Done()
	log.Printf("Worker: [%s] started...\n", workerID)
	for {
		select {
		case <- workerCtx.Done():
			log.Printf("Worker: [%s] Shutting down...\n", workerID)
			return
		default:
			streams, err := rdb.XRead(redisContext, &redis.XReadArgs{
				Streams: []string {utils.GetRedisStreamName(), "$"},
				Count: 1,
				Block: 2  * time.Second,
			}).Result()

			if err != nil {
				if err != redis.Nil {
					log.Printf("Worker: [%s] Error reading from stream: %v\n", workerID, err)
				}
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					jobID := message.Values["id"].(string)
					user := message.Values["user"].(string)
					log.Printf("Worker [%s] : Processing job %s for user %s...\n", workerID, jobID, user)
					time.Sleep(5 * time.Second)
					log.Printf("Worker [%s] : Job %s has been approved.\n", workerID, jobID)

					// Publish completion event to Redis Pub/Sub for real-time notification.
					rdb.Publish(utils.InitializeRedisContext(), utils.GetRedisChannelName(), jobID)
				}
			}
		}
	}
}

func handleShutdownGracefully(cancelWorkers context.CancelFunc, wg *sync.WaitGroup, srv *http.Server, rdb *redis.Client) {
	
	// 1. สร้าง Channel เพื่อรอรับ OS Signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit

	log.Println("Shutdown signal received, initiating graceful shutdown...")
	cancelWorkers()
	wg.Wait()
	log.Printf("All workers have been shut down.")
	ctxShutdown, cancelShutdown := utils.InitializeServerContext(5 * time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server Shutdown Failed: ", err)
	}
	log.Println("Server has been shut down successfully.")

	if err := rdb.Close(); err != nil {
		log.Fatal("Failed to close Redis connection: ", err)
	}
	log.Println("Redis connection has been closed successfully. Exitting...")
}