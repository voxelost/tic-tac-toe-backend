package main

import (
	"context"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	GameServer := NewGameServer(context.Background())

	r := http.NewServeMux()

	// Only log requests to our admin dashboard to stdout
	r.HandleFunc("/", GameServer.ListenAndServe)

	// Wrap our server with our gzip handler to gzip compress all responses.
	http.ListenAndServe(":8000", handlers.CompressHandler(r))
}
