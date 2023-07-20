package main

type GameProcess struct {
	PreGameHook         func()
	MainGameProcessHook func()
	PostGameHook        func()
}

// Return a new GameProcess object
func NewGameProcess(preGameHook, mainGameProcess, postGameHook func()) *GameProcess {
	return &GameProcess{
		PreGameHook:         preGameHook,
		MainGameProcessHook: mainGameProcess,
		PostGameHook:        postGameHook,
	}
}
