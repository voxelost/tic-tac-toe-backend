package gameserver

import (
	"context"
	"main/client"
	"main/message"
	"main/status"
	"main/utils"
)

// GameBase struct represents a base for every Game object
type GameBase struct {
	utils.ID
	Players      []*client.Player
	Status       status.GameStatus
	EventManager *message.EventManager
}

// return new GameBase
func NewGameBase(ctx context.Context, clients []*client.Client) *GameBase {
	players := []*client.Player{}
	for _, client_ := range clients {
		players = append(players, client.NewPlayer(client_))
	}

	id_ := utils.NewId()
	gb := &GameBase{
		ID:           *id_,
		Players:      players,
		EventManager: message.NewEventManager(message.NewOrigin(message.Game, id_)),
	}

	return gb
}

// Update game status and broadcast the change
func (gb *GameBase) UpdateStatus(status_ status.GameStatus) {
	gb.Status = status_
	gb.BroadcastMessage(message.NewMessage(message.GameStatusUpdate, status_))
}

// Broadcast given message to all of game's clients using message.Messenger
func (gb *GameBase) BroadcastMessage(message_ *message.Message) {
	gbId := gb.GetId()
	message_.SetOrigin(message.NewOrigin(message.Game, &gbId))
	gb.EventManager.Receive(message_)
}

// Run post game cleanups, shutdown local EventManager, unplug all Clients from GameBase's EventManager
func (gb *GameBase) Destroy() {
	for _, client := range gb.Players {
		if client.Valid {
			client.PopCommunicator()
		}
	}
}
