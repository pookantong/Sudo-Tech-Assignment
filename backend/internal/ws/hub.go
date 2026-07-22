package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"cinema-booking-backend/internal/model"
)

type SeatUpdate struct {
	ShowtimeID primitive.ObjectID `json:"showtime_id"`
	SeatID     primitive.ObjectID `json:"seat_id"`
	SeatLabel  string             `json:"seat_label"`
	Status     model.SeatStatus   `json:"status"`
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = true
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[conn]; !exists {
		return
	}

	delete(h.clients, conn)

	if err := conn.Close(); err != nil {
		log.Printf("ws close error: %v", err)
	}
}

func (h *Hub) Broadcast(update SeatUpdate) {
	body, err := json.Marshal(update)
	if err != nil {
		log.Printf("ws broadcast marshal error: %v", err)
		return
	}

	h.mu.RLock()

	clients := make([]*websocket.Conn, 0, len(h.clients))

	for conn := range h.clients {
		clients = append(clients, conn)
	}

	h.mu.RUnlock()

	for _, conn := range clients {
		if err := conn.WriteMessage(
			websocket.TextMessage,
			body,
		); err != nil {
			log.Printf("ws write error: %v", err)

			h.Unregister(conn)
		}
	}
}