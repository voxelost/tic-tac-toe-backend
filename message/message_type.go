package message

type MessageType string

var (
	ClientControl             = MessageType("client_control")
	RegisterForClientQueue    = MessageType("register_for_client_queue")
	UnregisterFromClientQueue = MessageType("unregister_from_client_queue")
	GameAction                = MessageType("game_action")
	GameState                 = MessageType("game_state")
	GameMeta                  = MessageType("game_meta")
	Chat                      = MessageType("chat")
	Debug                     = MessageType("debug")
	GameServerMeta            = MessageType("game_server_meta")
	GameStatusUpdate          = MessageType("game_status_update")
)
