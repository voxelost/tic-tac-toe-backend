package main

import (
	"strings"
	"time"

	"github.com/google/uuid"
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

// The TicTacToeGame struct represents the game that's being played
type TicTacToeGame struct {
	Id            string
	Clients       []*Client
	CurrentClient *Client
	BoardStates   []TicTacToeBoard
	Status        GameStatus
}

// Return a new TicTacToeGame object
func NewGame(clients []*Client) *TicTacToeGame {
	if len(clients) < 2 {
		return nil
	}

	g := &TicTacToeGame{
		Id:            uuid.New().String(),
		Clients:       clients,
		CurrentClient: clients[0],
		BoardStates:   []TicTacToeBoard{},
	}

	for _, client := range clients {
		// setup CancelGame hooks
		client.CancelGame = g.Cancel()
	}

	g.UpdateStatus(WaitingForExecutor)
	return g
}

func (g *TicTacToeGame) ToProcess() *GameProcess {
	return NewGameProcess(g.PreGameHook, g.MainGameLoop, g.PostGameHook)
}

// Pre tictactoe game hook
func (g *TicTacToeGame) PreGameHook() {
	g.UpdateStatus(Starting)
	clientNames := []string{}
	for _, client_ := range g.Clients {
		clientNames = append(clientNames, string(client_.Id))
	}

	g.BroadcastMessage("hello, your game has started")
	g.BroadcastMessage("the players are: [%s]", strings.Join(clientNames, ", "))
	g.UpdateStatus(Started)
}

// Main tictactoe game loop
func (g *TicTacToeGame) MainGameLoop() {
	loops := 3
	for i := 0; i < loops; i++ {
		g.BroadcastMessage("[%d/%d] the game is ongoing...", i+1, loops)
		time.Sleep(time.Second)
	}
}

// Post tictactoe game hook
func (g *TicTacToeGame) PostGameHook() {
	g.BroadcastMessage("game has ended")
	g.RunPostGameCleanup()
}

// Get last set board state
func (g *TicTacToeGame) GetLastState() TicTacToeBoard {
	return g.BoardStates[len(g.BoardStates)-1]
}

// Update game status and broadcast the change
func (g *TicTacToeGame) UpdateStatus(status GameStatus) {
	g.Status = status
	g.BroadcastMessage("game status updated: %s", g.Status)
}

// Broadcast given formatted message to all of game's clients
func (g *TicTacToeGame) BroadcastMessage(format string, a ...any) {
	for _, client := range g.Clients {
		client.SendMessage(format, a...)
	}
}

// Return a hook func that will cancel the current game
func (g *TicTacToeGame) Cancel() func() {
	return func() {
		g.UpdateStatus(Cancelled)
	}
}

func (g *TicTacToeGame) RunPostGameCleanup() {
	g.UpdateStatus(ShuttingDown)
	g.BroadcastMessage("game instance is shutting down...")
}

// Put a char on the board and push the board into g.BoardStates
func (g *TicTacToeGame) PutChar(y, x int, c string) (ok bool) {
	lastState := g.GetLastState()
	if lastState.Chars[y][x] != CharNone {
		return false
	}

	newState := lastState
	newState.Chars[y][x] = c
	g.BoardStates = append(g.BoardStates, newState)
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

		if lastState[y][y] != c {
			flagDiagonal1 = false
		}

		if lastState[y][2-y] != c {
			flagDiagonal2 = false
		}

		if flagHorizontal || flagVertical {
			return true
		}
	}

	return flagDiagonal1 || flagDiagonal2
}
