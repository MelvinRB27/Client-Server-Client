package main

import (

	"log"
	"net/http"

	"github.com/MelvinRB27/Client-Server/internal"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))

	log.Println("Websocket listening on: localhost:8080")
	
	ws := internal.NewWebSocketChat()
	http.HandleFunc("/chat", ws.HandlerUserConnection)

	go ws.UserChatManager()
	log.Fatal(http.ListenAndServe(":8080", nil))
}