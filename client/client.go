package client

import (
	"context"
	"main/connection"
	"main/message"
	"main/utils"
	"sync"

	"github.com/gorilla/websocket"
)

// Client struct represents a server client
type Client struct {
	*message.Messenger
	utils.ID
	sync.Mutex

	CancelGame       func()
	RemoveFromServer func()
	Valid            bool // if the client is still connected

	connection *connection.Connection
}

// Return a new Client object
func NewClient(ctx context.Context, conn *websocket.Conn, destroy func(*Client)) *Client {
	c := &Client{
		ID:    *utils.NewId(),
		Valid: true,
	}

	c.Messenger = message.NewMessenger(c.GetId())
	c.Messenger.SetEventManagerReceiveCallback(c.ReceiveEventManagerMessage)

	c.RemoveFromServer = func() { destroy(c) }
	c.connection = connection.NewConnection(ctx, c.GetId(), conn, c.ReceiveConnectionMessage, c.Invalidate)
	return c
}

// Receive a message from EventManager
func (c *Client) ReceiveEventManagerMessage(message_ *message.Message) {
	c.Lock()
	defer c.Unlock()

	if c.Valid {
		c.connection.SendMessage(message_) // forward to remote Client
	}
}

// Receive a connection from remote client
func (c *Client) ReceiveConnectionMessage(message_ *message.Message) {
	if message_ == nil {
		return
	}

	cId := c.GetId()
	message_.SetOrigin(message.NewOrigin(message.Client, &cId))
	c.Messenger.GetActiveCommunicator().SendToEventManager(message_) // forward to EventManager
}

// Invalidate a client. Calls to this function have to be idempotent
func (c *Client) Invalidate() {
	c.Lock()
	defer c.Unlock()

	c.Valid = false
	c.Destroy()
}

// Destroy a client. If they are in a game, cancel that game. Calls to this functions are idempotent
func (c *Client) Destroy() {
	if c.CancelGame != nil {
		c.CancelGame()
		c.CancelGame = nil
	}

	if c.RemoveFromServer != nil {
		c.RemoveFromServer()
		c.RemoveFromServer = nil
	}

	c.Messenger.UnsubFromAllEventManagers()
}
