package gameserver

import (
	"context"
	"main/client"
	"main/message"
	"main/status"
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

	unlockMove context.CancelFunc
	cancelGame context.CancelFunc
	// moveDone sync.Mutex // a lock for blocking the main game loop until a player makes a correct move
}

// Return a new TicTacToeGame object
func NewTicTacToeGame(ctx context.Context, clients []*client.Client) *TicTacToeGame {
	if len(clients) != 2 {
		return nil
	}

	g := &TicTacToeGame{
		GameBase:    *NewGameBase(ctx, clients),
		BoardStates: []TicTacToeBoard{*NewTicTacToeBoard()},
	}

	for _, client := range g.Players {
		// setup CancelGame hooks
		client.CancelGame = g.Cancel

		// setup EventManager callbacks
		g.EventManager.SubscribeMessenger(client.Messenger)
	}

	r := g.EventManager.Router

	// debug
	r.Route(message.Game, message.Debug, g.DumbForward)

	// game control
	r.Route(message.Client, message.GameAction, g.TrySetChar)
	r.Route(message.Game, message.GameState, g.DumbForward)
	r.Route(message.Game, message.GameMeta, g.DumbForward)
	r.Route(message.Game, message.GameStatusUpdate, g.DumbForward)

	// chat
	r.Route(message.Client, message.Chat, g.BroadcastClientMessage)

	g.UpdateStatus(status.WaitingForExecutor)
	return g
}

// Pre tictactoe game hook
func (g *TicTacToeGame) PreGameHook(cancelGame context.CancelFunc) {
	g.cancelGame = cancelGame
	g.UpdateStatus(status.Starting)
	g.Players[0].Payload["game_char"] = CharO
	g.Players[1].Payload["game_char"] = CharX

	g.UpdateStatus(status.Started)
	g.BroadcastMessage(message.NewMessage(message.GameMeta, map[string]string{
		string(g.Players[0].GetId()): CharO,
		string(g.Players[1].GetId()): CharX,
	}))
}

// Main tictactoe game loop
func (g *TicTacToeGame) MainGameProcessHook(_ctx context.Context) {
	for {
		select {
		case <-_ctx.Done():
			return
		default:
			g.NextRound()
			g.BroadcastGameState()

			ctx, cancelCtx := context.WithCancel(_ctx)
			g.unlockMove = cancelCtx
			g.LockUntilMoveDone(ctx)

			if g.CheckWin(g.CurrentPlayer.Payload["game_char"]) {
				g.BroadcastGameState()
				return
			}
		}
	}
}

// Post tictactoe game hook
func (g *TicTacToeGame) PostGameHook() {
	victoryStatus := status.VictoryO
	if g.CurrentPlayer.Payload["game_char"] == CharX && g.CurrentPlayer.Valid {
		victoryStatus = status.VictoryX
	}

	g.UpdateStatus(victoryStatus)
	g.UpdateStatus(status.Finished)
	g.UpdateStatus(status.ShuttingDown)
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
	g.Status = status.Cancelled
	g.cancelGame()
	g.Destroy()
}

// Block until a value is present in the moveDone channel
func (g *TicTacToeGame) LockUntilMoveDone(ctx context.Context) {
	<-ctx.Done()
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
	message := message.NewMessage(message.GameState, map[string]interface{}{
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
