package gameserver

import (
	"context"
	"fmt"
	"log"
	"main/client"
	"main/message"
	"main/status"

	"github.com/google/uuid"
)

// GameBase struct represents a base for every Game object
type GameBase struct {
	*message.Messenger
	Id           string
	Clients      []*client.Client
	Status       status.GameStatus
	EventManager message.EventManager
}

// return new GameBase
func NewGameBase(ctx context.Context, clients []*client.Client) *GameBase {
	gb := &GameBase{
		Id:           uuid.New().String(),
		Clients:      clients,
		EventManager: *message.NewEventManager(),
	}

	gb.Messenger = message.NewMessenger(gb.ReceiveMessage)
	gb.EventManager.SubscribeMessenger(gb.Messenger)
	return gb
}

// Update game status and broadcast the change
func (gb *GameBase) UpdateStatus(status_ status.GameStatus) {
	if gb.CheckFinished() {
		return
	}

	gb.Status = status_
	gb.BroadcastMessagef("game status updated: %s", gb.Status)
}

// Receive a message from EventManager
func (gb *GameBase) ReceiveMessage(message *message.Message) {
	log.Printf("GameBase %s received a message: %s\n", gb.Id, message.Payload)
}

// Broadcast given message to all of game's clients using message.Messenger
func (gb *GameBase) BroadcastMessage(message_ *message.Message) {
	message_.SetOrigin(message.NewOrigin(message.Game, &gb.Id))
	gb.Messenger.SendToEventManager(message_)
}

// Broadcast given formatted message to all of game's clients using message.Messenger
func (gb *GameBase) BroadcastMessagef(format string, a ...any) {
	gb.BroadcastMessage(message.NewMessage(message.Debug, map[string]interface{}{"data": fmt.Sprintf(format, a...)}))
}

// Return true if the game is in finished state, no matter the cause
func (gb *GameBase) CheckFinished() bool {
	return gb.Status == status.Cancelled || gb.Status == status.Finished
}

// Run post game cleanups, shutdown local EventManager, replug Clients to global EventManager
func (gb *GameBase) Destroy() {
	gb.EventManager.UnsubscribeMessenger(gb.Messenger)
}
