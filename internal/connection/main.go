package connection

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go_chat_client/config"
	"net/http"
	"time"
)

var WebSocket *websocket.Conn

type RawRequest struct {
	Type string     `json:"type"`
	Data ReqMessage `json:"data"`
}

type ReqMessage struct {
	MsgID  string    `json:"msg_id"`
	Body   string    `json:"body"`
	TDate  time.Time `json:"t_date"`
	FromID string    `json:"from_id"`
}

func Establish() {
	var err error
	header := http.Header{}
	cfg := config.Get()

	header.Add("Authorization", cfg.Token)

	WebSocket, _, err = websocket.DefaultDialer.Dial(cfg.ServerUrl, header)

	//Send ping messages periodically
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if err := WebSocket.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				return
			}
		}
	}()

	if err != nil {
		fmt.Println("Error connecting to WebSocket:", err)
		return
	}

}

func Close() {
	err := WebSocket.Close()
	if err != nil {
		return
	}
}
