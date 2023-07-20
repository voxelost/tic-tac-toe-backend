package main

import (
	"context"
	"log"
	"main/gameserver"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	if false {
		f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

	GameServer := gameserver.NewGameServer(context.Background())

	r := http.NewServeMux()

	// Only log requests to our admin dashboard to stdout
	r.HandleFunc("/", GameServer.ListenAndServe)

	// Wrap our server with our gzip handler to gzip compress all responses.
	http.ListenAndServe(":8000", handlers.CompressHandler(r))
}
