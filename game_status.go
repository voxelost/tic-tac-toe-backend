package main

type GameStatus string

var (
	WaitingForExecutor = GameStatus("waiting_for_executor")
	Starting           = GameStatus("starting")
	Started            = GameStatus("started")
	InProgress         = GameStatus("in_progress")
	ShuttingDown       = GameStatus("shutting_down")
	Cancelled          = GameStatus("cancelled")
)
