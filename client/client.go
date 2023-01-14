package client

import (
	"context"
	"log"
	"main/connection"
	"main/message"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client struct represents a server client
type Client struct {
	*message.Messenger
	Id         string
	CancelGame context.CancelFunc
	Connection *connection.Connection
	Valid      bool // if the client is still connected // TODO: rethink as ClientStatus
	Payload    map[string]interface{}
}

// Return a new Client object
func NewClient(ctx context.Context, conn *websocket.Conn) *Client {
	c := &Client{
		Id:      uuid.New().String(),
		Valid:   true,
		Payload: make(map[string]interface{}),
	}

	c.Messenger = message.NewMessenger(c.ReceiveMessage)
	c.Connection = connection.NewConnection(ctx, conn, c.ReceiveMessage, c.Invalidate)
	return c
}

// Receive a message from EventManager
func (c *Client) ReceiveMessage(message *message.Message) {
	log.Printf("Client %s received a message: %s\n", c.Id, message.Payload)
	c.Connection.SendMessage(message)
}

func (c *Client) Invalidate() {
	c.Destroy()
	c.Valid = false
}

func (c *Client) Destroy() {
	if c.CancelGame != nil {
		c.CancelGame()
		c.CancelGame = nil
	}
}
