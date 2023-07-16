package web_socket

import (
	"alekseikormski.com/server-status-monitoring/app/core"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type WebSocketServer struct {
	Upgrader   *websocket.Upgrader
	updateChan chan []*core.Application
}

func NewWebSocket(updateChan chan []*core.Application) *WebSocketServer {
	return &WebSocketServer{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		updateChan: updateChan,
	}
}

func (wss *WebSocketServer) Handle(ws *websocket.Conn) {
	for {
		applications := <-wss.updateChan

		copy := []string{}

		for _, app := range applications {
			appConverted, err := json.Marshal(app)
			if err != nil {
				log.Printf("cannot encode app: %w", err)
				return
			}

			copy = append(copy, string(appConverted))
		}
		if err := ws.WriteJSON(copy); err != nil {
			log.Printf("cannot send json: %v", err)
			return
		}
	}
}

func (wss *WebSocketServer) Reader(ws *websocket.Conn) {
	for {
		_, _, err := ws.ReadMessage() //what is message type?
		if err != nil {
			return
		}
	}
}
