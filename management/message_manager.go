package management

import "symflower/livechat_gqlgen/models"

var MessageCount int64 = 0

var messages = []models.ChatMessage{}

func AddMessage(message models.ChatMessage) {
	messages = append(messages, message)
	MessageCount++
}

func GetMessages() ([]models.ChatMessage) {
	return messages
}

