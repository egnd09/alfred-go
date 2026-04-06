package ws

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewHub(t *testing.T) {
	logger := zap.NewNop().Sugar()
	
	hub := NewHub(nil, nil, nil, logger)
	
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}
	if hub.clients == nil {
		t.Error("Hub clients map not initialized")
	}
	if hub.register == nil {
		t.Error("Hub register channel not initialized")
	}
	if hub.unregister == nil {
		t.Error("Hub unregister channel not initialized")
	}
	if hub.broadcast == nil {
		t.Error("Hub broadcast channel not initialized")
	}
}

func TestHubRegister(t *testing.T) {
	logger := zap.NewNop().Sugar()
	hub := NewHub(nil, nil, nil, logger)
	
	// Create mock client
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		env:    "test-env",
	}
	
	// Register client
	go hub.Run()
	hub.register <- client
	
	// Verify client was registered
	<-client.send // Wait for registration confirmation
	
	hub.RLock()
	_, exists := hub.clients[client]
	hub.RUnlock()
	
	if !exists {
		t.Error("Client not registered in hub")
	}
}
