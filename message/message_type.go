package message

type MessageType string

var (
	// GameControl = "game_control_message"
	ClientControl             = MessageType("client_control")
	RegisterForClientQueue    = MessageType("register_for_client_queue")
	UnregisterFromClientQueue = MessageType("unregister_from_client_queue")
	GameAction                = MessageType("game_action")
	GameState                 = MessageType("game_state")
	Chat                      = MessageType("chat")
	Debug                     = MessageType("debug")
	GameServerMeta            = MessageType("game_server_meta")
)
