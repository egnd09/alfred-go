package ws

import "encoding/json"

// MessageType represents WebSocket message types
type MessageType string

const (
	// Environment messages
	TypeNewEnv      MessageType = "new_env"
	TypeDeleteEnv   MessageType = "delete_env"
	TypeDefaultEnv  MessageType = "default_env"
	TypeEnvList     MessageType = "env_list"
	TypeEnvCreated  MessageType = "env_created"
	TypeEnvDeleted  MessageType = "env_deleted"
	TypeEnvError    MessageType = "env_error"

	// CI/Build messages
	TypeNewBuild      MessageType = "new_build"
	TypeCancelBuild   MessageType = "cancel_build"
	TypeGetTags       MessageType = "get_tags"
	TypeGetLastBuilds MessageType = "get_last_builds"
	TypeBuildStarted  MessageType = "build_started"
	TypeBuildComplete MessageType = "build_complete"
	TypeBuildError    MessageType = "build_error"

	// Container messages
	TypeContainerList    MessageType = "container_list"
	TypeContainerStatus  MessageType = "container_status"
	TypeKillPod          MessageType = "kill_pod"
	TypeGetDockerLogs    MessageType = "get_docker_logs"
	TypeContainerEvent   MessageType = "container_event"
	TypePodKilled        MessageType = "pod_killed"
	TypeLogs             MessageType = "logs"

	// Auth messages
	TypeAuthSuccess MessageType = "auth_success"
	TypeAuthError   MessageType = "auth_error"

	// Error
	TypeError MessageType = "error"
)

// Message represents a WebSocket message
type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// NewEnvData represents data for creating a new environment
type NewEnvData struct {
	Name     string   `json:"name"`
	Services []string `json:"services,omitempty"`
	Stable   bool     `json:"stable,omitempty"`
}

// DeleteEnvData represents data for deleting an environment
type DeleteEnvData struct {
	Name string `json:"name"`
}

// EnvListData represents environment list response
type EnvListData struct {
	Environments []EnvInfo `json:"environments"`
}

// EnvInfo represents environment information
type EnvInfo struct {
	Name     string   `json:"name"`
	Services []string `json:"services,omitempty"`
	Stable   bool     `json:"stable"`
}

// NewBuildData represents data for starting a new build
type NewBuildData struct {
	Repo    string `json:"repo"`
	Branch  string `json:"branch,omitempty"`
	Tag     string `json:"tag,omitempty"`
	EnvName string `json:"envName"`
}

// CancelBuildData represents data for canceling a build
type CancelBuildData struct {
	BuildID string `json:"buildId"`
}

// GetTagsData represents data for fetching tags
type GetTagsData struct {
	Repo string `json:"repo"`
}

// ContainerListData represents data for listing containers
type ContainerListData struct {
	EnvName string `json:"envName"`
}

// KillPodData represents data for killing a pod
type KillPodData struct {
	EnvName string `json:"envName"`
	PodName string `json:"podName"`
}

// GetDockerLogsData represents data for fetching logs
type GetDockerLogsData struct {
	EnvName string `json:"envName"`
	PodName string `json:"podName"`
	Follow  bool   `json:"follow,omitempty"`
	Tail    int64  `json:"tail,omitempty"`
}

// ErrorData represents error message data
type ErrorData struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// ToJSON converts message to JSON bytes
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// ParseMessage parses JSON bytes into a Message
func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// NewMessage creates a new message
func NewMessage(t MessageType, data interface{}) *Message {
	return &Message{
		Type: t,
		Data: data,
	}
}

// ErrorMessage creates an error message
func ErrorMessage(msg string, code int) *Message {
	return &Message{
		Type: TypeError,
		Data: ErrorData{
			Message: msg,
			Code:    code,
		},
	}
}
