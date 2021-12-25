package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Rooms is room mostructuredel
type Rooms struct {
	Name        string   `json:"name,omitempty" bson:"name,omitempty"`
	Type        string   `json:"type,omitempty" bson:"type,omitempty"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	GroupIcon   string   `json:"group_icon,omitempty" bson:"group_icon,omitempty"`
	Users       []string `json:"users,omitempty" bson:"users,omitempty"`
	Token       string   `json:"token,omitempty" bson:"token,omitempty"`
	CreatedBy   string   `json:"create_by,omitempty" bson:"create_by,omitempty"`
	CreatedAt   int64    `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdatedAt   int64    `json:"update_at,omitempty" bson:"update_at,omitempty"`
}

//User is user model
type User struct {
	Email        string  `json:"email,omitempty" bson:"email,omitempty"`
	Name         string  `json:"name,omitempty" bson:"name,omitempty"`
	ProfileImage string  `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	Status       string  `json:"status,omitempty" bson:"status,omitempty"`
	About        string  `json:"about,omitempty" bson:"about,omitempty"`
	Token        string  `json:"token,omitempty" bson:"token,omitempty"`
	LastLogin    []int64 `json:"last_login,omitempty" bson:"last_login,omitempty"`
	CreatedAt    int64   `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdatedAt    int64   `json:"update_at,omitempty" bson:"update_at,omitempty"`
}

//UserAlreadyInRoom it will give bollean value whether user is in the room or not
func UserAlreadyInRoom(token, roomID string) (bool, error) {
	temp := struct {
		Token  string `json:"token,omitempty"`
		RoomID string `json:"_id,omitempty"`
	}{Token: token, RoomID: roomID}
	jsonReq, err := json.Marshal(temp)
	if err != nil {
		return false, err
	}
	resp, err := http.Post("http://localhost:4000/room/details", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		fmt.Println("error at fetching room details", err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var room Rooms
	json.Unmarshal(bodyBytes, &room)
	temp = struct {
		Token  string `json:"token,omitempty"`
		RoomID string `json:"_id,omitempty"`
	}{Token: token, RoomID: roomID}
	jsonReq, err = json.Marshal(temp)

	resp, err = http.Post("http://localhost:4000/user/profile", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		fmt.Println("error at fetching user details", err)
	}

	defer resp.Body.Close()
	bodyBytes, _ = ioutil.ReadAll(resp.Body)
	var user User
	json.Unmarshal(bodyBytes, &user)
	for _, v := range room.Users {
		if v == user.Email {
			return true, nil
		}
	}
	return false, nil
}

//SendDataToDB will send data to mongo db
func SendDataToDB(msg Message) (string, error) {
	jsonReq, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("http://localhost:4000/message/send", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// Convert response body to string
	bodyString := string(bodyBytes)
	return bodyString, nil
}
