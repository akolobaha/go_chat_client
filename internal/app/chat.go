package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go_chat_client/internal/connection"
	"time"
)

type NewChatReq struct {
	Type string `json:"type"`
	Data struct {
		Users []string `json:"users"`
	} `json:"data"`
}

type NewChatResp struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Request struct {
	Type string  `json:"type"`
	Data Message `json:"data"`
}

type Message struct {
	Type    string `json:"type"`
	Message string `json:"msg"`
	ChatId  string `json:"ch_id"`
}

var inputBuffer string

func createChat(newChatUserId string) {
	chatReq := NewChatReq{
		Type: "new_chat",
	}
	chatReq.Data.Users = []string{newChatUserId}

	jsonMessage, err := json.Marshal(chatReq)
	if err != nil {
		fmt.Println("Error marshalling chatReq to JSON:", err)
		return
	}

	err = connection.WebSocket.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		fmt.Println("Error sending chatReq:", err)
		return
	}

	var resp NewChatResp

	_, message, err := connection.WebSocket.ReadMessage()

	err = json.Unmarshal(message, &resp)

	if err != nil {
		fmt.Println("Error reading chatReq:", err)
		return
	}
	fmt.Printf("%s\n", resp.Data)
}

func newMessage(text string, chatId string) {
	messageReq := Request{
		Type: "new_msg",
		Data: Message{
			Type:    "add",
			Message: text,
			ChatId:  chatId,
		},
	}

	var err error

	jsonMessage, err := json.Marshal(messageReq)

	if err != nil {
		fmt.Println("Error marshalling chatReq to JSON:", err)
		return
	}

	err = connection.WebSocket.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		fmt.Println("Error sending chatReq:", err)
		return
	}

}

func GetFormattedTime() string {
	currentTime := time.Now()
	return currentTime.Format("02.01.2006 15:04")
}
