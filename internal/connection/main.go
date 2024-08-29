package connection

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go_chat_client/config"
	"net/http"
)

var WebSocket *websocket.Conn

func Establish() {
	var err error
	header := http.Header{}
	cfg := config.Get()

	header.Add("Authorization", cfg.Token)

	WebSocket, _, err = websocket.DefaultDialer.Dial(cfg.ServerUrl, header)

	if err != nil {
		fmt.Println("Error connecting to WebSocket:", err)
		return
	}

	// TODO: реализовать закрытие соединения при gracefull shutdown
	//defer conn.Close()
}

func Close() {
	err := WebSocket.Close()
	if err != nil {
		return
	}
}
