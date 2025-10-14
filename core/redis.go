package core

import (
	"log"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/models"
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

	err = rdb.XGroupCreateMkStream(
		utils.InitializeRedisContext(), 
		utils.GetRedisStreamName(), 
		utils.GetRedisGroupName(), "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Error creating Redis Consumer Group: %v", err)
		return nil
	}
	log.Printf("Consumer group '%s' is ready.\n", utils.GetRedisGroupName())
	return rdb
}

func SubscribeToApprovalEvents(rdb *redis.Client, command chan models.ClientCommand) {
	redisChannelName := utils.GetRedisChannelName()
	pubsub := rdb.Subscribe(utils.InitializeRedisContext(), redisChannelName)
	defer pubsub.Close()

	ch := pubsub.Channel()
	log.Println("Subscribed to Redis channel:", redisChannelName)
	for msg := range ch {
		jobID := msg.Payload
		log.Printf("Received completion event for job: %s", jobID)
		// Send a command to the manager to forward the message to the client.
		command <- models.ClientCommand {
			Action: models.SendMessage,
			JobID: jobID, 
			Message: "approved",
		}
	}
}