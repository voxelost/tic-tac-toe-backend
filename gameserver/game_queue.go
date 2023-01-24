package gameserver

import (
	"context"
	"main/client"
	"main/utils"
)

type JobQueue chan GameProcess

// GameQueue is the queue of games waiting for launch
type GameQueue struct {
	utils.ModifiableQueue[GameProcess]
}

// Return new GameQueue
func NewGameQueue(queueSize int) *GameQueue {
	return &GameQueue{
		ModifiableQueue: *utils.NewModifiableQueue[GameProcess](queueSize),
	}
}

// Try to register a new game. If there are at least 2 Clients in the ClientQueue, the game is pushed into GameQueue's
// internal JobQueue and ok=true. Otherwise nothing is done and ok=false
func (gq *GameQueue) TryRegisterGame(ctx context.Context, clientQueue *client.ClientQueue) (ok bool) {
	clients, ok := clientQueue.GetNClients(2)
	if !ok {
		return false
	}

	gq.Push(NewTicTacToeGame(ctx, clients))
	return true
}

// Remove a game from the Queue.
func (gq *GameQueue) UnregisterGame(game *GameBase) {
	gq.Delete(game.GetId())
}
