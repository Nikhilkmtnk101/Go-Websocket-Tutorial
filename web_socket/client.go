package web_socket

import (
	"github.com/gorilla/websocket"
)

type Clients map[*Client]bool

type Client struct {
	conn           *websocket.Conn
	messageChannel chan Event
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:           conn,
		messageChannel: make(chan Event),
	}
}

func (c *Client) SendMessage(event Event) error {
	err := c.conn.WriteMessage(event.MessageType, event.Payload)
	return err
}
