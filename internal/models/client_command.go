package models

type CommandAction int

// clientCommand is a message sent to the manager goroutine.
type ClientCommand struct {
	Action  CommandAction
	JobID   string
	Channel chan string
	Message string
}
