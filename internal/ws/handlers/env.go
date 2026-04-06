package handlers

import (
	"context"
	"encoding/json"

	"github.com/egnd09/alfred-go/internal/models"
	"github.com/egnd09/alfred-go/internal/ws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// EnvHandler handles environment operations
type EnvHandler struct {
	mongoClient *mongo.Client
	logger      *zap.Logger
}

// NewEnvHandler creates a new environment handler
func NewEnvHandler(mongoClient *mongo.Client, logger *zap.Logger) *EnvHandler {
	return &EnvHandler{
		mongoClient: mongoClient,
		logger:      logger,
	}
}

// HandleMessage routes environment messages to appropriate handlers
func (h *EnvHandler) HandleMessage(ctx context.Context, client *ws.Client, msg ws.Message) {
	switch msg.Type {
	case "new_env":
		h.handleNewEnv(ctx, client, msg)
	case "delete_env":
		h.handleDeleteEnv(ctx, client, msg)
	case "default_env":
		h.handleDefaultEnv(ctx, client, msg)
	case "env_list":
		h.handleEnvList(ctx, client, msg)
	default:
		client.SendError(msg.Type, "unknown env message type")
	}
}

func (h *EnvHandler) handleNewEnv(ctx context.Context, client *ws.Client, msg ws.Message) {
	var req models.Environment
	data, _ := json.Marshal(msg.Data)
	if err := json.Unmarshal(data, &req); err != nil {
		client.SendError(msg.Type, "invalid environment data")
		return
	}

	h.logger.Info("creating new environment",
		zap.String("name", req.Name),
		zap.String("user", client.Username),
	)

	// Insert into MongoDB
	collection := h.mongoClient.Database("alfred").Collection("envs")
	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		h.logger.Error("failed to create environment", zap.Error(err))
		client.SendError(msg.Type, "failed to create environment")
		return
	}

	client.SendMessage("env_created", map[string]interface{}{
		"name":     req.Name,
		"services": req.Services,
		"stable":   req.Stable,
	})

	// Broadcast to all clients
	client.Hub.BroadcastToRoom("", ws.Message{
		Type: "env_added",
		Data: req,
	})
}

func (h *EnvHandler) handleDeleteEnv(ctx context.Context, client *ws.Client, msg ws.Message) {
	envName, ok := msg.Data.(string)
	if !ok {
		data, _ := json.Marshal(msg.Data)
		var req struct {
			Name string `json:"name"`
		}
		json.Unmarshal(data, &req)
		envName = req.Name
	}

	h.logger.Info("deleting environment",
		zap.String("name", envName),
		zap.String("user", client.Username),
	)

	// Delete from MongoDB
	collection := h.mongoClient.Database("alfred").Collection("envs")
	filter := bson.M{"name": envName}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		h.logger.Error("failed to delete environment", zap.Error(err))
		client.SendError(msg.Type, "failed to delete environment")
		return
	}

	client.SendMessage("env_deleted", map[string]string{
		"name": envName,
	})

	// Broadcast deletion
	client.Hub.BroadcastToRoom("", ws.Message{
		Type: "env_removed",
		Data: map[string]string{"name": envName},
	})
}

func (h *EnvHandler) handleDefaultEnv(ctx context.Context, client *ws.Client, msg ws.Message) {
	// Get or set default environment
	envName, ok := msg.Data.(string)
	if ok && envName != "" {
		// Set as default for user
		// TODO: Update user's default environment in database
		client.SendMessage("default_env_set", map[string]string{
			"env": envName,
		})
	} else {
		// Get default environment
		// TODO: Query user's default environment from database
		client.SendMessage("default_env", map[string]string{
			"env": "production",
		})
	}
}

func (h *EnvHandler) handleEnvList(ctx context.Context, client *ws.Client, msg ws.Message) {
	h.logger.Debug("listing environments")

	// Query MongoDB for all environments
	collection := h.mongoClient.Database("alfred").Collection("envs")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		h.logger.Error("failed to list environments", zap.Error(err))
		client.SendError(msg.Type, "failed to list environments")
		return
	}
	defer cursor.Close(ctx)

	var envs []models.Environment
	if err := cursor.All(ctx, &envs); err != nil {
		h.logger.Error("failed to decode environments", zap.Error(err))
		client.SendError(msg.Type, "failed to decode environments")
		return
	}

	client.SendMessage("env_list", map[string]interface{}{
		"environments": envs,
		"count":        len(envs),
	})
}
