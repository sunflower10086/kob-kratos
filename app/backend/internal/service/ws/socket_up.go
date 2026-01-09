package ws

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Solve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(hub *Hub, userId int64, w http.ResponseWriter, r *http.Request, logger *log.Helper) {
	if userId == 0 {
		// Try to get from header or other means if needed
		// For now simple query param or generate for guest?
		// userId = "guest"
	}

	ServeWs(hub, userId, w, r, logger)
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, userId int64, w http.ResponseWriter, r *http.Request, logger *log.Helper) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256), UserId: userId}
	client.Hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
