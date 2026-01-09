package ws

import (
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

// Hub maintains the set of active Clients and broadcasts messages to the
// Clients.
type Hub struct {
	// Registered Clients.
	clients map[int64]*Client
	rwMu    sync.RWMutex

	// Inbound messages from the Clients.
	broadcast chan []byte

	logger *log.Helper
}

func NewHub(logger log.Logger) *Hub {
	return &Hub{
		broadcast: make(chan []byte),
		clients:   make(map[int64]*Client, 10000),
		logger:    log.NewHelper(log.With(logger, "module", "service/ws/hub")),
	}
}

// Run hub event loop
func (h *Hub) Run() {
	for {
		select {
		case message, ok := <-h.broadcast:
			if !ok {
				return
			}
			h.logger.Infof("broadcast %s", string(message))
			h.rwMu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					// We need to upgrade lock to remove client, but we are iterating.
					// For simplicity in this loop, we might just skip or handle removal lazily.
					// Or more safely: collect dead clients and remove them after loop.
					// However, typical pattern with mutex is to just rely on client read/write failure to trigger Unregister.
					// But here we are sending. If send buffer is full, we assume dead?
					// Let's just launch a goroutine to unregister to avoid deadlock in RLock.
					go h.Unregister(client)
				}
			}
			h.rwMu.RUnlock()
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.rwMu.Lock()
	defer h.rwMu.Unlock()
	h.clients[client.UserId] = client
	h.logger.Infof("%v connected", client.UserId)
}

func (h *Hub) Unregister(client *Client) {
	h.rwMu.Lock()
	defer h.rwMu.Unlock()
	if _, ok := h.clients[client.UserId]; ok {
		delete(h.clients, client.UserId)
		close(client.Send)
		h.logger.Infof("%v disconnected", client.UserId)
	}
}

func (h *Hub) BroadcastTo(userId int64, message []byte) error {
	h.rwMu.RLock()
	defer h.rwMu.RUnlock()
	if client, ok := h.clients[userId]; ok {
		select {
		case client.Send <- message:
			return nil
		default:
			return fmt.Errorf("client %d send buffer full", userId)
		}
	}
	return fmt.Errorf("client %d not found", userId)
}

func (h *Hub) IsOnline(userId int64) bool {
	h.rwMu.RLock()
	defer h.rwMu.RUnlock()
	_, ok := h.clients[userId]
	return ok
}

func (h *Hub) GetOnlineUsers() []int64 {
	h.rwMu.RLock()
	defer h.rwMu.RUnlock()
	users := make([]int64, 0, len(h.clients))
	for userId := range h.clients {
		users = append(users, userId)
	}
	return users
}
