package message

type MessageType string

var (
	// GameControl = "game_control_message"
	GameAction = MessageType("game_action")
	GameState  = MessageType("game_state")
	Debug      = MessageType("debug")

	// Chat = "chat"
	// TODO?: live server diagnostics feed
)
