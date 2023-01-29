package gameserver

import "main/message"

type ServerStatus struct {
	Executors      int `json:"executors"`
	ActiveGames    int `json:"active_games"`
	PlayersInQueue int `json:"players_in_queue"`
	GamesInQueue   int `json:"games_in_queue"`
	Players        int `json:"players_online"`
}

func (gs *GameServer) NewServerStatus() *ServerStatus {
	return &ServerStatus{
		Executors:      len(gs.ExecutorPool.Workers),
		ActiveGames:    gs.ExecutorPool.GetActiveGames(),
		PlayersInQueue: gs.ClientQueue.Len(),
		GamesInQueue:   gs.GameQueue.ModifiableQueue.Len(),
		Players:        gs.Clients.Count(),
	}
}

func (gs *GameServer) SendNotifications() {
	gs.BroadcastMessage(message.NewMessage(message.GameServerMeta, gs.NewServerStatus()))
}

func (gs *GameServer) BroadcastMessage(message_ *message.Message) {
	message_.Origin = message.NewOrigin(message.Server, nil)
	gs.EventManager.Receive(message_)
}
