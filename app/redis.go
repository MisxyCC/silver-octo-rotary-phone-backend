package app

import (
	"log"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/pkg/utils"
	"github.com/go-redis/redis/v8"
)


func InitializeRedisConnection(redisAddress string) *redis.Client{
	redisContext := utils.InitializeRedisContext()
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	_, err := rdb.Ping(redisContext).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis instance successfully.")
	return rdb
}