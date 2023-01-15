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

	c.Messenger = message.NewMessenger(c.ReceiveEventManagerMessage)
	c.Connection = connection.NewConnection(ctx, conn, c.ReceiveConnectionMessage, c.Invalidate)
	return c
}

// Receive a message from EventManager
func (c *Client) ReceiveEventManagerMessage(message_ *message.Message) {
	log.Printf("Client %s received a message from event manager: %s\n", c.Id, message_.Payload)
	c.Connection.SendMessage(message_)
}

// Receive a connection from remote client
func (c *Client) ReceiveConnectionMessage(message_ *message.Message) {
	log.Printf("Client %s received a message from remote connection: %s\n", c.Id, message_.Payload)
	message_.SetOrigin(message.NewOrigin(message.Client, &c.Id))

	log.Printf("Forwarding to EM\n")
	c.Messenger.SendToEventManager(message_)
}

// Invalidate a client
func (c *Client) Invalidate() {
	c.Destroy()
	c.Valid = false
}

// Destroy a client. If they are in a game, cancel that game
// TODO: refactor
func (c *Client) Destroy() {
	if c.CancelGame != nil {
		c.CancelGame()
		c.CancelGame = nil
	}
}
