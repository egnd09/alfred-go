package api

import (
	"net/http"
	"time"

	"github.com/egnd09/alfred-go/internal/config"
	"github.com/egnd09/alfred-go/internal/db"
	"github.com/egnd09/alfred-go/internal/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// SetupRouter initializes all routes with dependencies
func SetupRouter(cfg *config.Config, mongoClient *mongo.Client, redisClient *db.RedisClient, hub *ws.Hub, logger *zap.SugaredLogger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/", HealthCheck)
	router.GET("/health", HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		api.POST("/login", Login(cfg, logger.Desugar()))
	}

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		serveWS(hub, c.Writer, c.Request, cfg, logger)
	})

	return router
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configure appropriately for production
	},
}

// serveWS handles WebSocket connections
func serveWS(hub *ws.Hub, w http.ResponseWriter, r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection", "error", err)
		return
	}

	client := ws.NewClient(hub, conn, logger)
	client.ReadPump()
}
