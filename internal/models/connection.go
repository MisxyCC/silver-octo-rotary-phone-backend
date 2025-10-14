package models

import (
	"net"
	"sync"
)

type ConnManagerCommand struct {
	Action ConnManagerAction
	JobID  string
	Conn   net.Conn
	Wg     *sync.WaitGroup // Used to signal when CloseAll is complete
}

const (
	AddClient CommandAction = iota
	RemoveClient
	SendMessage
)

type ConnManagerAction int

const (
	ConnAdd ConnManagerAction = iota
	ConnRemove
	ConnCloseAll
)