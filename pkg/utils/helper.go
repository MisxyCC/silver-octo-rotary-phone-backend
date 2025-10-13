package utils

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)


func InitializeWorkerContext() (context.Context, context.CancelFunc){
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	return workerCtx, cancelWorker
}

func InitializeServerContext(timeDuration time.Duration) (context.Context, context.CancelFunc) {
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), timeDuration)
	return ctxShutdown, cancelShutdown
}

func InitializeRedisContext() context.Context {
	return context.Background()
}

func GetRedisStreamName() string {
	return "approval_stream"
}

func StartWorker(workerCtx context.Context, workerID string, wg *sync.WaitGroup, rdb redis.Client, redisContext context.Context) {
	defer wg.Done()
	log.Printf("Worker: [%s] started...\n", workerID)
	for {
		select {
		case <- workerCtx.Done():
			log.Printf("Worker: [%s] Shutting down...\n", workerID)
			return
		default:
			streams, err := rdb.XRead(redisContext, &redis.XReadArgs{
				Streams: []string {GetRedisStreamName(), "$"},
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
				}
			}
		}
	}
}

func LoadEnvironmentVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}