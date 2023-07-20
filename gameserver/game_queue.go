package gameserver

import (
	"context"
	"main/client"
)

type JobQueue chan GameProcess

// GameQueue is the queue of games waiting for launch
type GameQueue struct {
	Queue JobQueue
}

// Return new GameQueue
func NewGameQueue() *GameQueue {
	return &GameQueue{
		Queue: make(JobQueue, 2048),
	}
}

// Try to register a new game. If there are at least 2 Clients in the ClientQueue, the game is pushed into GameQueue's
// internal JobQueue and ok=true. Otherwise nothing is done and ok=false
func (gq *GameQueue) TryRegisterGame(ctx context.Context, clientQueue *client.ClientQueue) (ok bool) {
	clients, ok := clientQueue.GetNClients(2)
	if !ok {
		return false
	}

	gq.Queue <- NewTicTacToeGame(ctx, clients)

	return true
}
