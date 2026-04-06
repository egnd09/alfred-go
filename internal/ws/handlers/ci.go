package handlers

import (
	"context"
	"encoding/json"

	"github.com/egnd09/alfred-go/internal/k8s"
	"github.com/egnd09/alfred-go/internal/ws"
	"go.uber.org/zap"
)

// CIHandler handles CI/build operations
type CIHandler struct {
	k8sClient *k8s.Client
	logger    *zap.Logger
}

// NewCIHandler creates a new CI handler
func NewCIHandler(k8sClient *k8s.Client, logger *zap.Logger) *CIHandler {
	return &CIHandler{
		k8sClient: k8sClient,
		logger:    logger,
	}
}



// BuildRequest represents a build request
type BuildRequest struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	Tag    string `json:"tag"`
	Env    string `json:"env"`
}

func (h *CIHandler) handleNewBuild(ctx context.Context, client *ws.Client, msg ws.Message) {
	var req BuildRequest
	data, _ := json.Marshal(msg.Data)
	if err := json.Unmarshal(data, &req); err != nil {
		client.SendError(msg.Type, "invalid build request")
		return
	}

	h.logger.Info("starting new build",
		zap.String("repo", req.Repo),
		zap.String("branch", req.Branch),
		zap.String("user", client.Username),
	)

	// TODO: Implement actual build triggering
	// This would typically:
	// 1. Create a Kubernetes Job or Pod to run the build
	// 2. Update job status in MongoDB
	// 3. Broadcast build started event

	client.SendMessage("build_started", map[string]interface{}{
		"repo":   req.Repo,
		"branch": req.Branch,
		"status": "pending",
	})
}

func (h *CIHandler) handleCancelBuild(ctx context.Context, client *ws.Client, msg ws.Message) {
	buildID, ok := msg.Data.(string)
	if !ok {
		buildID = ""
	}

	h.logger.Info("canceling build",
		zap.String("build_id", buildID),
		zap.String("user", client.Username),
	)

	// TODO: Implement build cancellation
	client.SendMessage("build_canceled", map[string]interface{}{
		"build_id": buildID,
	})
}

func (h *CIHandler) handleGetTags(ctx context.Context, client *ws.Client, msg ws.Message) {
	repo, _ := msg.Data.(string)

	h.logger.Debug("getting tags", zap.String("repo", repo))

	// TODO: Implement tag fetching from Git or container registry
	// This would query GitHub API or container registry
	tags := []string{"v1.0.0", "v1.1.0", "latest"}

	client.SendMessage("tags_list", map[string]interface{}{
		"repo": repo,
		"tags": tags,
	})
}

func (h *CIHandler) handleGetLastBuilds(ctx context.Context, client *ws.Client, msg ws.Message) {
	envName, _ := msg.Data.(string)

	h.logger.Debug("getting last builds", zap.String("env", envName))

	// TODO: Query MongoDB for recent builds
	builds := []map[string]interface{}{
		{
			"id":      "build-001",
			"repo":    "example-app",
			"branch":  "main",
			"status":  "completed",
			"created": "2024-01-15T10:00:00Z",
		},
	}

	client.SendMessage("last_builds", map[string]interface{}{
		"env":    envName,
		"builds": builds,
	})
}

func (h *CIHandler) handleContainerList(ctx context.Context, client *ws.Client, msg ws.Message) {
	envName, _ := msg.Data.(string)

	h.logger.Debug("listing containers", zap.String("env", envName))

	// Get pods from Kubernetes
	pods, err := h.k8sClient.ListPods(ctx, envName)
	if err != nil {
		h.logger.Error("failed to list pods", zap.Error(err))
		client.SendError(msg.Type, "failed to list containers")
		return
	}

	client.SendMessage("container_list", map[string]interface{}{
		"env":   envName,
		"pods":  pods,
		"count": len(pods),
	})
}

func (h *CIHandler) handleContainerStatus(ctx context.Context, client *ws.Client, msg ws.Message) {
	data, _ := json.Marshal(msg.Data)
	var req struct {
		Env     string `json:"env"`
		PodName string `json:"podName"`
	}
	json.Unmarshal(data, &req)

	h.logger.Debug("getting container status",
		zap.String("env", req.Env),
		zap.String("pod", req.PodName),
	)

	// TODO: Get detailed pod status from Kubernetes
	status := map[string]interface{}{
		"name":   req.PodName,
		"status": "Running",
		"ready":  true,
	}

	client.SendMessage("container_status", status)
}

func (h *CIHandler) handleKillPod(ctx context.Context, client *ws.Client, msg ws.Message) {
	data, _ := json.Marshal(msg.Data)
	var req struct {
		Env     string `json:"env"`
		PodName string `json:"podName"`
	}
	json.Unmarshal(data, &req)

	h.logger.Info("killing pod",
		zap.String("env", req.Env),
		zap.String("pod", req.PodName),
		zap.String("user", client.Username),
	)

	if err := h.k8sClient.KillPod(ctx, req.Env, req.PodName); err != nil {
		h.logger.Error("failed to kill pod", zap.Error(err))
		client.SendError(msg.Type, "failed to kill pod")
		return
	}

	client.SendMessage("pod_killed", map[string]interface{}{
		"podName": req.PodName,
		"env":     req.Env,
	})

	// Broadcast to environment room is handled by the hub
}

func (h *CIHandler) handleGetDockerLogs(ctx context.Context, client *ws.Client, msg ws.Message) {
	data, _ := json.Marshal(msg.Data)
	var req struct {
		Env     string `json:"env"`
		PodName string `json:"podName"`
		Follow  bool   `json:"follow"`
	}
	json.Unmarshal(data, &req)

	h.logger.Debug("getting pod logs",
		zap.String("env", req.Env),
		zap.String("pod", req.PodName),
	)

	logs, err := h.k8sClient.GetPodLogs(ctx, req.Env, req.PodName)
	if err != nil {
		h.logger.Error("failed to get pod logs", zap.Error(err))
		client.SendError(msg.Type, "failed to get logs")
		return
	}

	client.SendMessage("docker_logs", map[string]interface{}{
		"podName": req.PodName,
		"logs":    logs,
	})
}
