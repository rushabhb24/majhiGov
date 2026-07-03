package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader upgrades HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Origin validation is handled by CorsMiddleware before WS upgrade
		return true
	},
}

// NotificationMessage is the payload sent over WebSocket to clients
type NotificationMessage struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// wsClient represents a single connected WebSocket session
type wsClient struct {
	userID int
	conn   *websocket.Conn
	send   chan []byte
}

// Hub manages all active WebSocket connections by userID
type Hub struct {
	clients    map[int][]*wsClient // userID → slice (multi-tab support)
	mu         sync.RWMutex
	register   chan *wsClient
	unregister chan *wsClient
}

// NewHub creates a new WebSocket hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int][]*wsClient),
		register:   make(chan *wsClient, 64),
		unregister: make(chan *wsClient, 64),
	}
}

// Run processes hub events. Must be run in a dedicated goroutine: go hub.Run()
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c.userID] = append(h.clients[c.userID], c)
			h.mu.Unlock()
			log.Printf("WS: User %d connected (total connections for user: %d)", c.userID, len(h.clients[c.userID]))

		case c := <-h.unregister:
			h.mu.Lock()
			conns := h.clients[c.userID]
			for i, conn := range conns {
				if conn == c {
					h.clients[c.userID] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			if len(h.clients[c.userID]) == 0 {
				delete(h.clients, c.userID)
			}
			h.mu.Unlock()
			log.Printf("WS: User %d disconnected", c.userID)
		}
	}
}

// SendToUser delivers a notification to all active connections of a specific user
func (h *Hub) SendToUser(userID int, notification NotificationMessage) {
	h.mu.RLock()
	conns := make([]*wsClient, len(h.clients[userID]))
	copy(conns, h.clients[userID])
	h.mu.RUnlock()

	if len(conns) == 0 {
		return
	}

	data, err := json.Marshal(notification)
	if err != nil {
		return
	}

	for _, c := range conns {
		select {
		case c.send <- data:
		default:
			// Client buffer full — skip silently
		}
	}
}

// Broadcast sends a notification to ALL connected users
func (h *Hub) Broadcast(notification NotificationMessage) {
	data, err := json.Marshal(notification)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conns := range h.clients {
		for _, c := range conns {
			select {
			case c.send <- data:
			default:
			}
		}
	}
}

// ServeWS upgrades the HTTP connection to WebSocket and registers it with the hub
func (h *Hub) ServeWS(userID int, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade error for user %d: %v", userID, err)
		return
	}

	c := &wsClient{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
	}

	h.register <- c

	// Write pump — sends buffered messages to the WebSocket connection
	go func() {
		defer func() {
			h.unregister <- c
			conn.Close()
		}()
		for msg := range c.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()

	// Read pump — keeps connection alive and detects client disconnection
	go func() {
		defer func() {
			h.unregister <- c
			conn.Close()
		}()
		conn.SetReadLimit(512)
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			return nil
		})
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()

	// Ping pump — sends periodic pings to keep the connection alive through proxies
	go func() {
		ticker := time.NewTicker(45 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()
}
