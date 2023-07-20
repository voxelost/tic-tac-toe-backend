package gameserver

import (
	"main/message"
	"main/utils"
)

func (g *TicTacToeGame) TrySetChar(m *message.Message) bool {
	if g.CurrentPlayer.GetId() != m.Origin.GetId() {
		return false
	}

	// TODO: refactor
	coords, ok := m.Payload.([]interface{})
	if !ok {
		return false
	}

	coordYf, ok := coords[0].(float64)
	if !ok {
		return false
	}

	coordXf, ok := coords[1].(float64)
	if !ok {
		return false
	}

	coordY := int(coordYf)
	coordX := int(coordXf)

	ok = g.PutChar(coordY, coordX, g.CurrentPlayer.Payload["game_char"])
	if !ok {
		return false
	}

	g.unlockMove()
	return false
}

// Run a Client Message through a censoring middleware and broadcast it to all players
func (g *TicTacToeGame) BroadcastClientMessage(message *message.Message) bool {
	preprocessedPayload, ok := utils.PreprocessChatPayload(message.Payload.(string))
	if !ok {
		return false
	}

	message.Payload = preprocessedPayload
	return true
}

func (g *TicTacToeGame) DumbForward(message *message.Message) bool {
	return true
}
