package core

import (
	"log"
	"net"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/models"
)

func SSEConnectionManager(connManagerCommands chan models.ConnManagerCommand) {
	connections := make(map[string]net.Conn)
	log.Println("SSE Connection Manager started.")

	for cmd := range connManagerCommands {
		switch cmd.Action {
		case models.ConnAdd:
			connections[cmd.JobID] = cmd.Conn
		case models.ConnRemove:
			delete(connections, cmd.JobID)
		case models.ConnCloseAll:
			log.Printf("Closing %d active SSE connections...", len(connections))
			for jobID, conn := range connections {
				conn.Close()
				delete(connections, jobID)
			}
			// Signal that the close all operation is complete.
			if cmd.Wg != nil {
				cmd.Wg.Done()
			}
		}
	}

}
