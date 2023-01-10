package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientID string

// Client struct represents a server client
type Client struct {
	Id         ClientID
	Conn       *websocket.Conn
	CancelGame context.CancelFunc
	ConnIn     chan []byte
	ConnOut    chan []byte
	Valid      bool // if the client is still connected // TODO: RETHINK AS ClientStatus
}

// Return a new Client object
func NewClient(ctx context.Context, connection *websocket.Conn) *Client {
	c := &Client{
		Id:      ClientID(uuid.New().String()),
		Conn:    connection,
		ConnIn:  make(chan []byte, 1024),
		ConnOut: make(chan []byte, 1024),
		Valid:   true,
	}

	c.StartCommunicationRoutines(ctx)
	c.SendMessage("hi, your id is: %s", c.Id)

	// reminderContext, _ := context.WithCancel(ctx_)
	// go c.StartPeriodicReminder(reminderContext)
	return c
}

// TODO: RETHINK
func (c *Client) StartPeriodicReminder(ctx context.Context) {
	sleepTime := 5
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(sleepTime) * time.Second)
			c.SendMessage("you are waiting in queue")

			sleepTime += 2
		}
	}
}

// Push a message to the user's ConnOut channel.
func (c *Client) SendMessage(format string, a ...any) {
	c.ConnOut <- []byte(fmt.Sprintf(format, a...))
}

// Start communication routines
func (c *Client) StartCommunicationRoutines(ctx context.Context) {
	go c.ReadMessageRoutine(ctx)
	go c.WriteMessageRoutine(ctx)
}

// Run a routine that reads messages from websocket and pushes them into ConnIn channel
func (c *Client) ReadMessageRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, p, err := c.Conn.ReadMessage()
			if err != nil {
				switch err.(type) {
				case *websocket.CloseError:
					fmt.Println("connection closed by remote client")
					c.Invalidate()
					return
				default:
					log.Println(err)
					return
				}
			}

			switch messageType {
			case websocket.TextMessage:
				fmt.Println("i just received", string(p))
				c.ConnIn <- p
				c.ConnOut <- append([]byte("echoed: "), p...) // echo message back; TODO: REMOVE

			case websocket.CloseMessage:
				fmt.Println("connection was closed")
				c.Invalidate()
				return
			default:
				fmt.Println("dunno what happened:", messageType, p)
			}
		}
	}
}

// Run a routine that writes messages from ConnOut to websocket
func (c *Client) WriteMessageRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.SendMessage("game is shutting down")
			return
		case sendData := <-c.ConnOut:
			fmt.Println("i am sending: ", string(sendData))
			if err := c.Conn.WriteMessage(websocket.TextMessage, sendData); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Client) Invalidate() {
	c.Valid = false
	if c.CancelGame != nil {
		c.CancelGame()
	}
}
