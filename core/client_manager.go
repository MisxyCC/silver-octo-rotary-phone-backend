package core

import (
	"log"

	"github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/models"
	"github.com/MisxyCC/silver-octo-rotary-phone-backend/pkg/utils"
)

func ClientChannelManager(commands <-chan models.ClientCommand) {
	clientChannels := make(map[string]chan string)
	log.Println("Client channel manager started.")
	for cmd := range commands {
		switch cmd.Action {
		case utils.AddClient:
			clientChannels[cmd.JobID] = cmd.Channel
			log.Printf("Manager: Added client for job %s", cmd.JobID)
		case utils.RemoveClient:
			if ch, ok := clientChannels[cmd.JobID]; ok {
				close(ch)
				delete(clientChannels, cmd.JobID)
				log.Printf("Manager: Removed client for job %s", cmd.JobID)
			}
		case utils.SendMessage:
			if ch, ok := clientChannels[cmd.JobID]; ok {
				select {
				case ch <- cmd.Message: // Attempt to send the message
				default:
					log.Printf("Manager: Dropped message for job %s, client channel not ready.", cmd.JobID)
				}
			}
		}
	}
}
