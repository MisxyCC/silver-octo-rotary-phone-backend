package utils

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func InitializeWorkerContext() (context.Context, context.CancelFunc) {
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

func GetRedisChannelName() string {
	return "approval_events"
}

func GetRedisGroupName() string {
	return "approval_workers"
}

func LoadEnvironmentVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
