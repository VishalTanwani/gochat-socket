package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

//Client struct
type Client struct {
	ID     string
	Name   string
	Conn   *websocket.Conn
	Server *Server
	Token  string
	Rooms  map[*Room]bool
}

//NewClient will return new client object
func NewClient(UUID string, conn *websocket.Conn, server *Server) *Client {
	return &Client{
		ID:     UUID,
		Conn:   conn,
		Server: server,
		Rooms:  make(map[*Room]bool),
	}
}

//Read for reading client messages
func (c *Client) Read() {
	defer func() {
		c.Server.UnRegister <- c
		c.Conn.Close()
	}()
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("at client read", err)
			return
		}

		switch msg.Type {
		case "message":
			c.sendMessage(msg)
		case "info":
			c.sendMessage(msg)
		case "joinRoom":
			c.joinRoom(msg)
		case "leaveRoom":
			c.leaveRoom(msg)
		}
	}
}

func (c *Client) sendMessage(msg Message) {
	if room := c.Server.FindRoom(msg.RoomID); room != nil {
		room.Broadcast <- msg
	}
}

func (c *Client) joinRoom(msg Message) {
	room := c.Server.FindRoom(msg.RoomID)
	if room == nil {
		room = c.Server.CreateRoom(msg.RoomID)
	}
	c.Rooms[room] = true
	c.ID = msg.UserID
	c.Name = msg.UserName
	c.Token = msg.Token
	room.Register <- c
	check, err := UserAlreadyInRoom(msg.Token, msg.RoomID)
	if err == nil {
		if !check {
			room.Broadcast <- msg
		}
	}

}

func (c *Client) leaveRoom(msg Message) {
	room := c.Server.FindRoom(msg.RoomID)
	room.UnRegister <- c
	room.Broadcast <- msg
	delete(c.Rooms, room)
}

func (c *Client) Write(msg Message) {
	c.Conn.WriteJSON(msg)
}
