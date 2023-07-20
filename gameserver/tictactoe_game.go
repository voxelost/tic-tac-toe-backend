package gameserver

import (
	"context"
	"fmt"
	"main/client"
	"main/message"
	"main/status"
	"strings"
)

var (
	CharNone = " "
	CharO    = "O"
	CharX    = "X"
)

// TicTacToeBoard represents a 3x3 board of chars
type TicTacToeBoard struct {
	Chars [3][3]string
}

func NewTicTacToeBoard() *TicTacToeBoard {
	return &TicTacToeBoard{
		[3][3]string{
			{CharNone, CharNone, CharNone},
			{CharNone, CharNone, CharNone},
			{CharNone, CharNone, CharNone},
		}}
}

// The TicTacToeGame struct represents the game that's being played
type TicTacToeGame struct {
	GameBase
	CurrentPlayer *client.Player
	Round         int
	BoardStates   []TicTacToeBoard

	moveDone chan bool // a channel for blocking the main game loop until a player makes a correct move
}

// Return a new TicTacToeGame object
func NewTicTacToeGame(ctx context.Context, clients []*client.Client) *TicTacToeGame {
	if len(clients) != 2 {
		return nil
	}

	g := &TicTacToeGame{
		BoardStates: []TicTacToeBoard{*NewTicTacToeBoard()},
		moveDone:    make(chan bool),
	}

	g.GameBase = *NewGameBase(ctx, clients)

	for _, client := range g.Players {
		// setup CancelGame hooks
		client.CancelGame = g.Cancel

		// setup EventManager callbacks
		g.EventManager.SubscribeMessenger(client.Messenger)
	}

	r := g.EventManager.Router

	// debug
	r.Route(message.Game, message.Debug, g.TestRoute)

	// game control
	r.Route(message.Client, message.GameAction, g.TrySetChar)

	// chat
	r.Route(message.Client, message.Chat, g.BroadcastClientMessage)

	g.UpdateStatus(status.WaitingForExecutor)
	return g
}

// Pre tictactoe game hook
func (g *TicTacToeGame) PreGameHook() {
	g.UpdateStatus(status.Starting)
	g.Players[0].Payload["game_char"] = CharO
	g.Players[1].Payload["game_char"] = CharX

	clientNames := []string{}
	for _, client := range g.Players {
		clientNames = append(clientNames, fmt.Sprintf("[%s] - [%s]", client.GetId(), client.Payload))
	}

	g.UpdateStatus(status.Started)
	g.BroadcastMessagef("the players are:\n%s", strings.Join(clientNames, "\n"))
}

// Main tictactoe game loop
func (g *TicTacToeGame) MainGameProcessHook() {
	for {
		g.NextRound()
		g.BroadcastGameState()

		if g.CheckWin(g.CurrentPlayer.Payload["game_char"]) {
			break
		}

		g.LockUntilMoveDone()

		if g.CheckFinished() {
			break
		}
	}

	g.BroadcastMessagef("PLAYER %s [%s] WON", g.CurrentPlayer.GetId(), g.CurrentPlayer.Payload["game_char"])
	g.UpdateStatus(status.Finished)
}

// Post tictactoe game hook
func (g *TicTacToeGame) PostGameHook() {
	g.BroadcastMessagef("game has ended")
	g.UpdateStatus(status.ShuttingDown)
	g.BroadcastMessagef("game instance is shutting down...")

	g.Destroy()
}

// Setup next round
func (g *TicTacToeGame) NextRound() {
	g.Round++
	g.CurrentPlayer = g.Players[g.Round%len(g.Players)]
}

// Get last set board state
func (g *TicTacToeGame) GetLastState() *TicTacToeBoard {
	return &g.BoardStates[len(g.BoardStates)-1]
}

// Cancel the current game
func (g *TicTacToeGame) Cancel() {
	for len(g.moveDone) > 0 {
		<-g.moveDone
	} // flush chan
	g.moveDone <- false // unlock main process
	g.UpdateStatus(status.Cancelled)
	g.Destroy()
}

// Block until a value is present in the moveDone channel
func (g *TicTacToeGame) LockUntilMoveDone() {
	<-g.moveDone
}

// Run GameBase.Destroy and cancel own context
func (g *TicTacToeGame) Destroy() {
	g.GameBase.Destroy()

	for _, player := range g.Players {
		player.CancelGame = nil
	}
}

// Put a char on the board and push the board into g.BoardStates
func (g *TicTacToeGame) PutChar(y, x int, c string) (ok bool) {
	lastState := g.GetLastState()
	if lastState.Chars[y][x] != CharNone {
		return false
	}

	newState := lastState
	newState.Chars[y][x] = c
	g.BoardStates = append(g.BoardStates, *newState)
	return true
}

func (g *TicTacToeGame) BroadcastGameState() {
	message := message.NewMessage(message.Debug, map[string]interface{}{
		"round":          g.Round,
		"current_player": g.CurrentPlayer.GetId(),
		"board":          g.GetLastState(),
	})

	g.BroadcastMessage(message)
}

// Check if the given char won
func (g *TicTacToeGame) CheckWin(c string) bool {
	lastState := g.GetLastState().Chars
	flagDiagonal1 := true
	flagDiagonal2 := true

	for y := 0; y < 3; y++ {
		flagHorizontal := true
		flagVertical := true
		for x := 0; x < 3; x++ {
			if lastState[y][x] != c {
				flagHorizontal = false
			}
			if lastState[x][y] != c {
				flagVertical = false
			}
		}

		if flagHorizontal || flagVertical {
			return true
		}

		if lastState[y][y] != c {
			flagDiagonal1 = false
		}

		if lastState[y][2-y] != c {
			flagDiagonal2 = false
		}
	}

	return flagDiagonal1 || flagDiagonal2
}
