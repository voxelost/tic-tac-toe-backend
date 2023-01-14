package gameserver

import (
	"context"
	"fmt"
	"main/client"
	"main/message"
	"main/status"
	"math/rand"
	"strings"
	"time"
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
	CurrentPlayer *client.Client
	Round         int
	BoardStates   []TicTacToeBoard
}

// Return a new TicTacToeGame object
func NewTicTacToeGame(ctx context.Context, clients []*client.Client) *TicTacToeGame {
	if len(clients) != 2 {
		return nil
	}

	g := &TicTacToeGame{
		BoardStates: []TicTacToeBoard{*NewTicTacToeBoard()},
	}

	g.GameBase = *NewGameBase(ctx, clients)
	g.UpdateStatus(status.WaitingForExecutor)
	return g
}

// Pre tictactoe game hook
func (g *TicTacToeGame) PreGameHook() {
	g.UpdateStatus(status.Starting)

	g.Clients[0].Payload["game_char"] = CharO
	g.Clients[1].Payload["game_char"] = CharX

	clientNames := []string{}
	for _, client := range g.Clients {
		// setup CancelGame hooks
		client.CancelGame = g.Cancel
		// setup EventManager callbacks
		g.EventManager.SubscribeMessenger(client.Messenger)
		clientNames = append(clientNames, fmt.Sprintf("[%s] - [%s]", client.Id, client.Payload))
	}

	g.UpdateStatus(status.Started)
	g.BroadcastMessagef("the players are:\n%s", strings.Join(clientNames, "\n"))
}

// Main tictactoe game loop
func (g *TicTacToeGame) MainGameProcessHook() {
	for {
		if g.CheckFinished() {
			break
		}

		g.NextRound()

		g.PutChar(rand.Intn(3), rand.Intn(3), g.CurrentPlayer.Payload["game_char"].(string))
		message := message.NewMessage(message.Debug, map[string]interface{}{
			"round":          g.Round,
			"current_player": g.CurrentPlayer.Id,
			"board":          g.GetLastState(),
		})

		g.BroadcastMessage(message)

		if g.CheckWin(g.CurrentPlayer.Payload["game_char"].(string)) {
			g.BroadcastMessagef("PLAYER %s [%s] WON", g.CurrentPlayer.Id, g.CurrentPlayer.Payload["game_char"].(string))
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
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
	g.CurrentPlayer = g.Clients[g.Round%len(g.Clients)]
}

// Get last set board state
func (g *TicTacToeGame) GetLastState() *TicTacToeBoard {
	return &g.BoardStates[len(g.BoardStates)-1]
}

// Return a hook func that will cancel the current game
func (g *TicTacToeGame) Cancel() {
	g.UpdateStatus(status.Cancelled)
	g.Destroy()
}

// Run GameBase.Destroy and cancel own context
func (g *TicTacToeGame) Destroy() {
	g.GameBase.Destroy()
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
