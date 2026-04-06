package ws

import (
	"encoding/json"
	"sync"

	"github.com/egnd09/alfred-go/internal/k8s"
	"github.com/egnd09/alfred-go/internal/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]map[*Client]bool // environment name -> clients
	broadcast  chan *BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex

	mongoClient *mongo.Client
	redisClient interface{}
	k8sClient   *k8s.Client
	logger      *zap.SugaredLogger
}

// BroadcastMessage represents a message to be broadcast
type BroadcastMessage struct {
	room    string
	message []byte
}

// NewHub creates a new WebSocket hub
func NewHub(mongoClient *mongo.Client, redisClient interface{}, k8sClient *k8s.Client, logger *zap.SugaredLogger) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		rooms:       make(map[string]map[*Client]bool),
		broadcast:   make(chan *BroadcastMessage, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		mongoClient: mongoClient,
		redisClient: redisClient,
		k8sClient:   k8sClient,
		logger:      logger,
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Info("Client connected", "total", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.removeFromRooms(client)
				close(client.send)
			}
			h.mu.Unlock()
			h.logger.Info("Client disconnected", "total", len(h.clients))

		case msg := <-h.broadcast:
			h.mu.RLock()
			room, ok := h.rooms[msg.room]
			if ok {
				for client := range room {
					select {
					case client.send <- msg.message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(room string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*Client]bool)
	}
	h.rooms[room][client] = true
}

// LeaveRoom removes a client from a room
func (h *Hub) LeaveRoom(room string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if roomClients, ok := h.rooms[room]; ok {
		delete(roomClients, client)
		if len(roomClients) == 0 {
			delete(h.rooms, room)
		}
	}
}

// removeFromRooms removes client from all rooms
func (h *Hub) removeFromRooms(client *Client) {
	for room, clients := range h.rooms {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.rooms, room)
		}
	}
}

// Broadcast sends a message to all clients in a room
func (h *Hub) Broadcast(room string, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Failed to marshal broadcast message", "error", err)
		return
	}
	h.broadcast <- &BroadcastMessage{room: room, message: data}
}

// ValidateToken validates JWT token and returns claims
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key-change-in-production"), nil
	})
}

// HandleMessage routes messages to appropriate handlers
func (h *Hub) HandleMessage(client *Client, msg Message) {
	switch msg.Type {
	case "new_env":
		h.handleNewEnv(client, msg)
	case "delete_env":
		h.handleDeleteEnv(client, msg)
	case "default_env":
		h.handleDefaultEnv(client, msg)
	case "env_list":
		h.handleEnvList(client, msg)
	case "new_build":
		h.handleNewBuild(client, msg)
	case "cancel_build":
		h.handleCancelBuild(client, msg)
	case "get_tags":
		h.handleGetTags(client, msg)
	case "get_last_builds":
		h.handleGetLastBuilds(client, msg)
	case "container_list":
		h.handleContainerList(client, msg)
	case "container_status":
		h.handleContainerStatus(client, msg)
	case "kill_pod":
		h.handleKillPod(client, msg)
	case "get_docker_logs":
		h.handleGetDockerLogs(client, msg)
	case "join_room":
		h.handleJoinRoom(client, msg)
	default:
		h.logger.Warn("Unknown message type", "type", msg.Type)
	}
}

// handleJoinRoom handles room join requests
func (h *Hub) handleJoinRoom(client *Client, msg Message) {
	room, ok := msg.Data["room"].(string)
	if !ok {
		return
	}
	h.JoinRoom(room, client)
	client.room = room
}
