package handlers

import (
	"context"
	"net/http"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/models"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func InitSubmitHandler(rdb *redis.Client, redisContext context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonRequest models.ApprovalRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		job := models.ApprovalJob{
			ID:      uuid.New().String(),
			User:    jsonRequest.User,
			Amount:  jsonRequest.Amount,
			Details: jsonRequest.Details,
			Status:  "pending",
		}

		jobValues, err := job.ToMap()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to prepare job data"})
			return
		}
		_, err = rdb.XAdd(redisContext, &redis.XAddArgs{
			Stream: utils.GetRedisStreamName(),
			Values: jobValues,
		}).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError,
				models.ErrorResponse{Error: "Failed to submit the request"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Request submitted successfully", "job_id": job.ID})
	}
}
