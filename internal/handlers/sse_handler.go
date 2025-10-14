package handlers

import (
	"io"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/models"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func InitSSEHandler(command chan models.ClientCommand) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("job_id")
		messageChan := make(chan string)

		// Register the new client with the manager.
		command <- models.ClientCommand {
			Action: utils.AddClient,
			JobID: jobID,
			Channel: messageChan,
		}
		defer func() {
			command <- models.ClientCommand {
				Action: utils.RemoveClient, JobID:  jobID,
			}
		}()

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		
		c.Stream(func (w io.Writer) bool {
			select {
			case msg, ok := <- messageChan:
				if !ok {
					return false // Channel was closed by the manager.
				}
				c.SSEvent("message", msg)
				return true // Keep connection open.
			case <- c.Request.Context().Done():
				return false // Connection closed by client.
			}
		})
	}
}