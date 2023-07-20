package gameserver

import (
	"context"
	"fmt"
	"main/chatmod"
	"main/message"
)

func (gs *GameServer) PrintClientDebug(message *message.Message) bool {
	fmt.Printf("client: %s\n", message.Payload)
	return true
}

func (gs *GameServer) PrintServerDebug(message *message.Message) bool {
	fmt.Printf("server: %s\n", message.Payload)
	return true
}

func (gs *GameServer) DumbForward(message *message.Message) bool {
	return true
}

// Run a Client Message through a censoring middleware and broadcast it to all players
func (gs *GameServer) BroadcastClientMessage(message *message.Message) bool {
	defer func() {
		if err := recover(); err != nil {
			return // recover from bad client message
		}
	}()

	censoredPayload, err := chatmod.CensorChatMesage(message.Payload.(string))
	if err != nil {
		return false
	}

	if len(censoredPayload) > 200 {
		censoredPayload = censoredPayload[:200]
	}

	message.Payload = censoredPayload
	return true
}

// Register the client for GameServer.ClientQueue
func (gs *GameServer) RegisterForClientQueue(message *message.Message) bool {
	client := gs.Clients.GetClientForId(message.Origin.GetId())
	gs.ClientQueue.Register(client)
	gs.GameQueue.TryRegisterGame(context.Background(), gs.ClientQueue)
	return false
}

// Unregister the client from GameServer.ClientQueue
func (gs *GameServer) UnregisterFromClientQueue(message *message.Message) bool {
	client := gs.Clients.GetClientForId(message.Origin.GetId())
	gs.ClientQueue.Unregister(client)
	return false
}
