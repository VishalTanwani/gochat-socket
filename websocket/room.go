package websocket

import (
	"fmt"
)

//Room struct
type Room struct {
	ID         string
	Name       string
	Register   chan *Client
	UnRegister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

//NewRoom is for creatng new pool
func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (room *Room) writeMessageForRoom(msg Message) {
	for c := range room.Clients {
		c.Write(msg)
	}
}

//StartRoom is for sending different message to all client which are in room
func (room *Room) StartRoom() {
	for {
		select {
		case client := <-room.Register:
			room.Clients[client] = true
		case client := <-room.UnRegister:
			if ok := room.Clients[client]; ok {
				delete(room.Clients, client)
			}
		case msg := <-room.Broadcast:
			_, err := SendDataToDB(msg)
			if err != nil {
				fmt.Println("at startroom ", err)
			}
			room.writeMessageForRoom(msg)
		}
	}
}
