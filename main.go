package main

import (
	"fmt"
	"os"
	"github.com/VishalTanwani/gochat/socket/websocket"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "home page")
}

func wsServer(server *websocket.Server, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrader(w, r)
	if err != nil {
		fmt.Println("at main.go wsServer 1", err)
		return
	}

	client := websocket.NewClient("", conn, server)
	go client.Read()
	server.Register <- client
}

func setRoutes() {
	server := websocket.NewServer()
	go server.Start()
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsServer(server, w, r)
	})
}

func main() {
	fmt.Println("hello chat")
	setRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Printf("server is running at %s port\n",port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
