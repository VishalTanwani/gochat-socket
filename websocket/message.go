package websocket

//Message struct
type Message struct {
	Body     string `json:"body,omitempty"`
	Image     string `json:"image,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Type     string `json:"type,omitempty"`
	Room     string `json:"room,omitempty"`
	RoomID   string `json:"room_id,omitempty"`
	Token    string `json:"token,omitempty"`
}
