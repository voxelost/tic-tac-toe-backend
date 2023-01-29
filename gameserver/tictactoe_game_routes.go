package gameserver

import (
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
