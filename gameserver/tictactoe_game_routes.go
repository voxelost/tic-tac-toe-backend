package gameserver

import (
	"encoding/json"
	"fmt"
	"main/chatmod"
	"main/message"
)

func (g *TicTacToeGame) TestRoute(m *message.Message) bool {
	fmt.Println(m)
	return true
}

func (g *TicTacToeGame) TrySetChar(m *message.Message) bool {
	if g.CurrentPlayer.GetId() != m.Origin.GetId() {
		fmt.Printf("unauthorized board edit attempt, current player is %s and not %s\n", g.CurrentPlayer.GetId(), m.Origin.GetId())
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

	g.moveDone <- true
	return false
}

// Run a Client Message through a censoring middleware and broadcast it to all players
func (g *TicTacToeGame) BroadcastClientMessage(message *message.Message) bool {
	payloadBytes, err := json.Marshal(message.Payload)
	if err != nil {
		return false
	}

	censoredPayload, err := chatmod.CensorChatMesage(string(payloadBytes))
	if err != nil {
		return false
	}

	if len(payloadBytes) > 200 {
		censoredPayload = censoredPayload[:200]
	}

	message.Payload = censoredPayload
	return true
}
