package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func InitSSEHandler(clientCommand chan models.ClientCommand, connManagerCommands chan models.ConnManagerCommand) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("job_id")
		if jobID == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Job ID is required"})
			return
		}
		// 1. Hijack the underlying network connection for manual control.
		hijacker, ok := c.Writer.(http.Hijacker)
		if !ok {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server does not support hijacking"})
			return
		}

		// Hijack() returns the raw network connection and a buffered reader/writer
		conn, _, err := hijacker.Hijack()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to hijack connection"})
			return
		}

		// 2. Command the connection manager to ADD this new connection.
		connManagerCommands <- models.ConnManagerCommand{Action: models.ConnAdd, JobID: jobID, Conn: conn}

		// 3. Defer actions to clean up when the handler exits (client disconnects or server shuts down).
		defer conn.Close()
		defer func() {
			// Command the connection manager to REMOVE this connection.
			connManagerCommands <- models.ConnManagerCommand{Action: models.ConnRemove, JobID: jobID}
		}()
		// 4. Create a channel for receiving messages for THIS specific client.
		messageChan := make(chan string)
		// Command the message delivery manager to ADD this client's message channel.
		clientCommand <- models.ClientCommand{Action: models.AddClient, JobID: jobID, Channel: messageChan}
		// Defer the command to REMOVE the message channel.
		defer func() {
			clientCommand <- models.ClientCommand{Action: models.RemoveClient, JobID: jobID}
		}()
		// 5. Manually write the HTTP headers required for an SSE connection.
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n")
		fmt.Fprintf(conn, "Content-Type: text/event-stream\r\n")
		fmt.Fprintf(conn, "Cache-Control: no-cache\r\n")
		fmt.Fprintf(conn, "Connection: keep-alive\r\n")
		fmt.Fprintf(conn, "Access-Control-Allow-Origin: *\r\n")
		fmt.Fprintf(conn, "\r\n")
		// 6. Loop indefinitely, waiting for messages or a client disconnect signal.
		for {
			select {
			case msg, ok := <-messageChan:
				if !ok {
					log.Println("The message channel was closed by the manager")
					return
				}
				// Write the SSE formatted message directly to the connection.
				fmt.Fprintf(conn, "event: message\ndata: %s\n\n", msg)

			case <-c.Request.Context().Done():
				// The client has closed the connection from their end.
				log.Printf("Client for job %s disconnected.", jobID)
				return
			}
		}
	}
}
