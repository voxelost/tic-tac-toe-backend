package gameserver

import (
	"context"
	"main/client"
	"main/message"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// GameServer struct represents a game server responsible for registering new clients, handling global EventManager
// events and controlling Clients
type GameServer struct {
	Clients      *client.ClientCache
	ClientQueue  *client.ClientQueue
	GameQueue    *GameQueue
	EventManager *message.EventManager
	ExecutorPool *WorkerPool
}

// Return new GameServer
func NewGameServer(ctx context.Context) *GameServer {
	gameQueue := NewGameQueue(1024)
	gs := &GameServer{
		Clients:      client.NewClientCache(),
		ClientQueue:  client.NewClientQueue(1024),
		GameQueue:    gameQueue,
		EventManager: message.NewEventManager(message.NewOrigin(message.Server, nil)),
		ExecutorPool: NewWorkerPool(ctx, &gameQueue.ModifiableQueue, 512),
	}

	gs.InitNotificationRoutine(ctx, 500*time.Millisecond)
	r := gs.EventManager.Router

	// client control messages
	r.Route(message.Client, message.RegisterForClientQueue, gs.RegisterForClientQueue)
	r.Route(message.Client, message.UnregisterFromClientQueue, gs.UnregisterFromClientQueue)

	// chat messages
	r.Route(message.Client, message.Chat, gs.BroadcastClientMessage)

	// game server meta messages
	r.Route(message.Server, message.GameServerMeta, gs.DumbForward)
	return gs
}

// handle new users and put them in queue,
func (gs *GameServer) ListenAndServe(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := client.NewClient(context.Background(), conn, gs.ForgetClient)
	gs.EventManager.SubscribeMessenger(client.Messenger)
	gs.Clients.Register(client)
}

func (gs *GameServer) InitNotificationRoutine(ctx context.Context, notificationFrequency time.Duration) {
	go func() {
		for {
			time.Sleep(notificationFrequency)
			select {
			case <-ctx.Done():
				return
			default:
				gs.SendNotifications()
			}
		}
	}()
}

func (gs *GameServer) ForgetClient(c *client.Client) {
	gs.Clients.Unregister(c)
	gs.ClientQueue.Unregister(c)
	gs.EventManager.UnsubscribeMessenger(c.Messenger)
}
