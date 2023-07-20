package gameserver

type GameProcess interface {
	PreGameHook()
	MainGameProcessHook()
	PostGameHook()
}
