package gameserver

import (
	"context"
	"main/message"
	"main/utils"
)

func (gs *GameServer) DumbForward(message *message.Message) bool {
	return true
}

// Run a Client Message through a censoring middleware and broadcast it to all players
func (gs *GameServer) BroadcastClientMessage(message *message.Message) bool {
	preprocessedPayload, ok := utils.PreprocessChatPayload(message.Payload.(string))
	if !ok {
		return false
	}

	message.Payload = preprocessedPayload
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
