package connection

import (
	"context"
	"main/message"
	"main/utils"
	"sync"

	"github.com/gorilla/websocket"
)

// Connection holds handles all communication with the remote client
type Connection struct {
	sync.Mutex // for ws writes
	utils.ID
	CancelCallback func()

	conn                   *websocket.Conn
	messageReceiveCallback func(*message.Message)
}

// Return a new Connection
func NewConnection(ctx context.Context, clientId utils.ID, conn *websocket.Conn, messageReceive func(*message.Message), cancelCallback func()) *Connection {
	c := &Connection{
		conn:           conn,
		CancelCallback: cancelCallback,
		ID:             clientId,

		messageReceiveCallback: messageReceive,
	}

	go c.WebsocketReadRoutine(ctx)
	return c
}

// Run a routine that reads messages from websocket and pushes them into ConnIn channel
func (c *Connection) WebsocketReadRoutine(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Cancel()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, p, err := c.conn.ReadMessage()
			if err != nil {
				c.Cancel()
				return
			}

			switch messageType {
			case websocket.TextMessage:
				c.messageReceiveCallback(message.NewMessageFromBytes(p))
			default:
				c.Cancel()
				return
			}
		}
	}
}

// Write a marshalled Message to websocket
func (c *Connection) SendMessage(message *message.Message) {
	defer func() {
		if err := recover(); err != nil {
			c.Cancel()
		}
	}()
	c.Lock()
	defer c.Unlock()

	message.RecipientId = c.GetId()
	messageBytes, _ := message.Marshal()

	if err := c.conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		c.Cancel()
	}
}

func (c *Connection) Cancel() {
	c.Lock()
	defer c.Unlock()

	if c.CancelCallback != nil {
		c.CancelCallback()
		c.CancelCallback = nil
	}
}
