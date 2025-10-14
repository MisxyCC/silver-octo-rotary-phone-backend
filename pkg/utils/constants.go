package utils

import "github.com/MisxyCC/silver-octo-rotary-phone-backend/internal/models"

const (
	AddClient models.CommandAction = iota
	RemoveClient
	SendMessage
)