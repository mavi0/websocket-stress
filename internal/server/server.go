package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for testing purposes
	},
}

// Client represents a single websocket connection
type Client struct {
	conn *websocket.Conn
}

// Manager manages all connected clients
type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewManager creates a new client manager
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Start begins the client manager's main loop
func (manager *Manager) Start() {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client] = true
			log.Printf("Client connected. Total clients: %d", len(manager.clients))
		case client := <-manager.unregister:
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				log.Printf("Client disconnected. Total clients: %d", len(manager.clients))
				// No need to broadcast disconnect message here, handled client-side
			}
		case message := <-manager.broadcast:
			for client := range manager.clients {
				err := client.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					client.conn.Close()
					manager.unregister <- client
				}
			}
		}
	}
}

// SendUpdates sends periodic updates to all clients
func (manager *Manager) SendUpdates() {
	ticker := time.NewTicker(1 * time.Second)

	// Prepare a large payload (e.g., 1MB of 'a' characters)
	payloadSize := 1 * 1024 * 1024 // 1 MB
	largePayload := make([]byte, payloadSize)
	for i := range largePayload {
		largePayload[i] = 'a'
	}

	for range ticker.C {
		timestamp := time.Now().Format(time.RFC3339)
		clientCount := len(manager.clients)
		// Combine timestamp, client count, and the large payload
		header := fmt.Sprintf("%s | Connected clients: %d | Payload size: %d bytes\n", timestamp, clientCount, payloadSize)
		message := append([]byte(header), largePayload...)

		manager.broadcast <- message
	}
}

// ServeWs handles websocket connections
func ServeWs(manager *Manager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{conn: conn}
	manager.register <- client

	// Listen for close events
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				manager.unregister <- client
				conn.Close()
				break
			}
		}
	}()
}

// ServeHome handles the home page
func ServeHome(w http.ResponseWriter, r *http.Request, rootDir string) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	templatePath := filepath.Join(rootDir, "web/templates/index.html")
	http.ServeFile(w, r, templatePath)
}
