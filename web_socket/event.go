package web_socket

import "encoding/json"

type Event struct {
	MessageType int             `json:"message_type"`
	Payload     json.RawMessage `json:"payload"`
}

func NewEvent(messageType int, payload json.RawMessage) Event {
	return Event{
		MessageType: messageType,
		Payload:     payload,
	}
}
