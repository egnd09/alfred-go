package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egnd09/alfred-go/internal/api"
	"github.com/egnd09/alfred-go/internal/config"
	"github.com/egnd09/alfred-go/internal/db"
	"github.com/egnd09/alfred-go/internal/k8s"
	"github.com/egnd09/alfred-go/internal/util"
	"github.com/egnd09/alfred-go/internal/ws"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := util.NewLogger()
	defer logger.Sync()

	// Connect to MongoDB
	mongoClient, err := db.ConnectMongo(cfg.DBMongoURI)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", "error", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Connect to Redis
	redisClient := db.NewRedisClient(cfg.RedisURL)
	defer redisClient.Close()

	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient("")
	if err != nil {
		logger.Warn("Kubernetes client not available", "error", err)
		// Continue without k8s client for local development
	}

	// Initialize WebSocket hub
	hub := ws.NewHub(mongoClient, redisClient, k8sClient, logger)
	go hub.Run()

	// Setup router
	router := api.SetupRouter(cfg, mongoClient, redisClient, hub, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServicePort),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.ServicePort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited properly")
}
