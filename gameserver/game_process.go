package gameserver

import (
	"context"
	"main/utils"
)

type GameProcess interface {
	GetId() utils.ID
	PreGameHook(context.CancelFunc)
	MainGameProcessHook(context.Context)
	PostGameHook()
}
