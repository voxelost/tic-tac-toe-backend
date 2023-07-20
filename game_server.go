package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type GameServer struct {
	GlobalClientQueue ClientQueue
	GlobalGameQueue   GameQueue
	ExecutorPool      WorkerPool
}

// Return new GameServer
func NewGameServer(ctx context.Context) *GameServer {
	gameQueue := NewGameQueue()
	return &GameServer{
		GlobalClientQueue: *NewClientQueue(),
		GlobalGameQueue:   *gameQueue,
		ExecutorPool:      *NewWorkerPool(ctx, gameQueue.Queue, 10), // TODO: defulat pool size to 512
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

	client := NewClient(context.Background(), conn)
	gs.GlobalClientQueue.RegisterClient(client)
	gs.GlobalGameQueue.TryRegisterGame(context.Background(), &gs.GlobalClientQueue)
}
