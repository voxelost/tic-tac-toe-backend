package gameserver

import (
	"context"
	"main/client"
	"testing"
)

func TestCheckHoriz1(t *testing.T) {
	game := NewTicTacToeGame(context.Background(), []*client.Client{{}, {}})
	boardChar := "X"
	newBoard := TicTacToeBoard{
		Chars: [3][3]string{
			{boardChar, boardChar, boardChar},
			{" ", " ", " "},
			{" ", " ", " "},
		},
	}

	game.BoardStates = append(game.BoardStates, newBoard)
	if !game.CheckWin(boardChar) {
		t.Error("should be true")
	}
}

func TestCheckHoriz2(t *testing.T) {
	game := NewTicTacToeGame(context.Background(), []*client.Client{{}, {}})
	boardChar := "X"
	newBoard := TicTacToeBoard{
		Chars: [3][3]string{
			{" ", " ", " "},
			{boardChar, boardChar, boardChar},
			{" ", " ", " "},
		},
	}

	game.BoardStates = append(game.BoardStates, newBoard)
	if !game.CheckWin(boardChar) {
		t.Error("should be true")
	}
}

func TestCheckHoriz3(t *testing.T) {
	game := NewTicTacToeGame(context.Background(), []*client.Client{{}, {}})
	boardChar := "X"
	newBoard := TicTacToeBoard{
		Chars: [3][3]string{
			{" ", " ", " "},
			{" ", " ", " "},
			{boardChar, boardChar, boardChar},
		},
	}

	game.BoardStates = append(game.BoardStates, newBoard)
	if !game.CheckWin(boardChar) {
		t.Error("should be true")
	}
}

func TestCheckVertical1(t *testing.T) {
	game := NewTicTacToeGame(context.Background(), []*client.Client{{}, {}})
	boardChar := "X"
	newBoard := TicTacToeBoard{
		Chars: [3][3]string{
			{boardChar, " ", " "},
			{boardChar, " ", " "},
			{boardChar, " ", " "},
		},
	}

	game.BoardStates = append(game.BoardStates, newBoard)
	if !game.CheckWin(boardChar) {
		t.Error("should be true")
	}
}

func TestCheckVertical2(t *testing.T) {
	game := NewTicTacToeGame(context.Background(), []*client.Client{{}, {}})
	boardChar := "X"
	newBoard := TicTacToeBoard{
		Chars: [3][3]string{
			{" ", boardChar, " "},
			{" ", boardChar, " "},
			{" ", boardChar, " "},
		},
	}

	game.BoardStates = append(game.BoardStates, newBoard)
	if !game.CheckWin(boardChar) {
		t.Error("should be true")
	}
}

func TestCheckVertical3(t *testing.T) {
	game := NewTicTacToeGame(context.Background(), []*client.Client{{}, {}})
	boardChar := "X"
	newBoard := TicTacToeBoard{
		Chars: [3][3]string{
			{" ", " ", boardChar},
			{" ", " ", boardChar},
			{" ", " ", boardChar},
		},
	}

	game.BoardStates = append(game.BoardStates, newBoard)
	if !game.CheckWin(boardChar) {
		t.Error("should be true")
	}
}

func TestCheckDiagonal1(t *testing.T) {
	game := NewTicTacToeGame(context.Background(), []*client.Client{{}, {}})
	boardChar := "X"
	newBoard := TicTacToeBoard{
		Chars: [3][3]string{
			{boardChar, " ", " "},
			{" ", boardChar, " "},
			{" ", " ", boardChar},
		},
	}

	game.BoardStates = append(game.BoardStates, newBoard)
	if !game.CheckWin(boardChar) {
		t.Error("should be true")
	}
}

func TestCheckDiagonal2(t *testing.T) {
	game := NewTicTacToeGame(context.Background(), []*client.Client{{}, {}})
	boardChar := "X"
	newBoard := TicTacToeBoard{
		Chars: [3][3]string{
			{" ", " ", boardChar},
			{" ", boardChar, " "},
			{boardChar, " ", " "},
		},
	}

	game.BoardStates = append(game.BoardStates, newBoard)
	if !game.CheckWin(boardChar) {
		t.Error("should be true")
	}
}
