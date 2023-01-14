package gameserver

import (
	"context"
	"log"
	"main/client"
	"main/message"
	"net/http"

	"github.com/gorilla/websocket"
)

// GameServer struct represents a game server responsible for registering new clients, handling global EventManager
// events and controlling Clients
type GameServer struct {
	ClientQueue  *client.ClientQueue
	GameQueue    *GameQueue
	EventManager *message.EventManager
	ExecutorPool *WorkerPool
}

// Return new GameServer
func NewGameServer(ctx context.Context) *GameServer {
	gameQueue := NewGameQueue()
	return &GameServer{
		ClientQueue:  client.NewClientQueue(),
		GameQueue:    gameQueue,
		EventManager: message.NewEventManager(),
		ExecutorPool: NewWorkerPool(ctx, gameQueue.Queue, 10), // TODO: default pool size to 512
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// handle new users and put them in queue,
func (gs *GameServer) ListenAndServe(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := client.NewClient(context.Background(), conn)
	gs.EventManager.SubscribeMessenger(client.Messenger)
	gs.ClientQueue.RegisterClient(client)
	// todo: don't force game registration, trigger this when users explicitly request a game
	gs.GameQueue.TryRegisterGame(context.Background(), gs.ClientQueue)
}
