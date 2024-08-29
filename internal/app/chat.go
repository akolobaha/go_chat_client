package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go_chat_client/internal/connection"
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

func createChat(newChatUserId string) {
	// Create a new request with the Authorization header

	// Connect to the WebSocket server

	// Create a new chat chatReq
	chatReq := NewChatReq{
		Type: "new_chat",
	}
	chatReq.Data.Users = []string{newChatUserId}

	// Marshal the chatReq to JSON
	jsonMessage, err := json.Marshal(chatReq)
	if err != nil {
		fmt.Println("Error marshalling chatReq to JSON:", err)
		return
	}

	// Send the JSON chatReq
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
