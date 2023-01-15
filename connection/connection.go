package connection

import (
	"context"
	"log"
	"main/message"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn                   *websocket.Conn
	MessageReceiveCallback func(*message.Message)
	SendMessage            func(*message.Message)
	CancelCallback         func()
}

func NewConnection(ctx context.Context, conn *websocket.Conn, messageReceive func(*message.Message), cancelCallback func()) *Connection {
	c := &Connection{
		Conn:                   conn,
		MessageReceiveCallback: messageReceive,
		CancelCallback:         cancelCallback,
	}

	c.SendMessage = c.WebsocketWriteRoutine

	go c.WebsocketReadRoutine(ctx)

	return c
}

// Run a routine that reads messages from websocket and pushes them into ConnIn channel
func (c *Connection) WebsocketReadRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, p, err := c.Conn.ReadMessage()
			if err != nil {
				switch err.(type) {
				case *websocket.CloseError:
					c.CancelCallback()
					return
				default:
					c.CancelCallback()
					return
				}
			}

			switch messageType {
			case websocket.TextMessage:
				c.MessageReceiveCallback(message.NewMessageFromBytes(p))
			case websocket.CloseMessage:
				c.CancelCallback()
				return
			default:
				continue
			}
		}
	}
}

// Write to websocket
func (c *Connection) WebsocketWriteRoutine(message *message.Message) {
	messageBytes, _ := message.Marshal()
	if err := c.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		log.Println(err)
		return
	}
}
