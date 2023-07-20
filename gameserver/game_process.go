package gameserver

import "main/utils"

type GameProcess interface {
	GetId() utils.ID
	PreGameHook()
	MainGameProcessHook()
	PostGameHook()
}
