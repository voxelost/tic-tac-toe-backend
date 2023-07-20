package client

import (
	"main/utils"
)

// ClientQueue is the queue of players waiting for a game
type ClientQueue struct {
	*utils.ModifiableQueue[*Client]
}

// Return new ClientQueue
func NewClientQueue(queueSize int) *ClientQueue {
	return &ClientQueue{
		ModifiableQueue: utils.NewModifiableQueue[*Client](queueSize),
	}
}

// Register a Client as waiting for Game
func (cq *ClientQueue) Register(c *Client) {
	cq.Push(c)
}

// Unregister a Client from the ClientQueue
func (cq *ClientQueue) Unregister(c *Client) {
	cq.Delete(c.GetId())
}

// Remove all clients that are no longer valid.
func (cq *ClientQueue) GarbageCollect() {
	clientIds := cq.GetIds()
	idsToDelete := []utils.ID{}
	for _, id := range clientIds {
		c := cq.Get(id)
		if !c.Valid {
			idsToDelete = append(idsToDelete, id)
		}
	}

	cq.RemoveMultiple(idsToDelete)
}

// Lock Game Queue, try to gather N clients and return them. return ok=false if queue doesn't have enough clients
func (cq *ClientQueue) GetNClients(n int) (clients []*Client, ok bool) {
	cq.GarbageCollect()

	cq.Lock()
	defer cq.Unlock()

	if cq.Len() < n {
		return clients, false
	}

	for i := 0; i < n; i++ {
		client, ok := cq.PopNonBlocking()
		if !ok {
			for _, client_ := range clients {
				cq.Push(client_)
			}

			return []*Client{}, false
		}

		clients = append(clients, client)
	}

	return clients, true
}
