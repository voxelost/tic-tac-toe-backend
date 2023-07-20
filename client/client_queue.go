package client

import (
	"sync"
)

// ClientQueue is the queue of players waiting for a game
type ClientQueue struct {
	sync.Mutex
	Queue chan *Client
}

// Return new ClientQueue
func NewClientQueue() *ClientQueue {
	return &ClientQueue{
		Queue: make(chan *Client, 4096),
	}
}

// Register a Client as waiting for Game
func (cq *ClientQueue) RegisterClient(c *Client) {
	cq.Lock()
	defer cq.Unlock()
	cq.Queue <- c
}

// Remove all clients that are no longer valid.
// WARNING: this method does not lock ClientQueue and needs to be run within a context that does
func (cq *ClientQueue) GarbageCollect() {
	clients := []*Client{}

	for len(cq.Queue) > 0 {
		client := <-cq.Queue
		if client.Valid {
			clients = append(clients, client)
		}
	}

	for _, client := range clients {
		cq.Queue <- client
	}
}

// Lock Game Queue, try to gather N clients and return them. return ok=false if queue doesn't have enough clients
func (cq *ClientQueue) GetNClients(n int) (clients []*Client, ok bool) {
	cq.Lock()
	defer cq.Unlock()
	cq.GarbageCollect()
	if len(cq.Queue) >= n {
		clients = []*Client{}
		for i := 0; i < n; i++ {
			newClient := <-cq.Queue
			clients = append(clients, newClient)
		}

		return clients, true
	}

	return clients, false
}
