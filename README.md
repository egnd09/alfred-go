# Alfred-Go

Kubernetes Development Environment Orchestrator - Go Implementation

## Overview

Alfred is a development environment orchestrator that creates and manages environments for development teams in a Kubernetes cluster. It provides a real-time web interface for provisioning, building, and deploying microservices.

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│   Server    │────▶│  Kubernetes │
│  (React)    │◀────│   (Go+WS)   │◀────│   (EKS)     │
└─────────────┘     └─────────────┘     └─────────────┘
                          │
           ┌──────────────┼──────────────┐
           ▼              ▼              ▼
       ┌────────┐   ┌─────────┐   ┌─────────┐
       │MongoDB │   │ Redis   │   │  Vault  │
       └────────┘   └─────────┘   └─────────┘
```

## Stack

- **Server**: Go 1.22+ (Gin + Gorilla WebSocket)
- **Client**: React 16 + Redux (unchanged)
- **Database**: MongoDB
- **Cache**: Redis
- **Container**: Docker + Kubernetes (EKS)
- **CD**: ArgoCD
- **Secrets**: HashiCorp Vault
- **CI**: GitHub Actions
- **Cloud**: AWS (ECR, EKS)

## Project Structure

```
alfred-go/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── api/
│   │   ├── router.go           # Gin router setup
│   │   ├── health.go           # Health check endpoint
│   │   └── auth.go             # GitHub OAuth
│   ├── ws/
│   │   ├── hub.go              # WebSocket hub
│   │   ├── client.go           # WS client with read/write pumps
│   │   ├── messages.go         # Message types
│   │   └── handlers/
│   │       ├── ci.go           # Build/deploy handlers
│   │       └── env.go          # Environment handlers
│   ├── k8s/
│   │   ├── client.go           # K8s client wrapper
│   │   └── pods.go             # Pod operations
│   ├── db/
│   │   ├── mongo.go            # MongoDB connection
│   │   └── redis.go            # Redis connection
│   ├── config/
│   │   └── config.go           # Viper config loader
│   ├── models/
│   │   ├── environment.go      # Environment model
│   │   ├── job.go              # Job model
│   │   └── user.go             # User model
│   └── util/
│       └── logger.go           # Zap structured logger
├── scripts/
│   ├── dev.sh                  # Local development
│   ├── build.sh                # Build binary
│   ├── test.sh                 # Run tests
│   └── docker-build.sh         # Build Docker image
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- kubectl (for K8s operations)
- MongoDB (local or remote)
- Redis (local or remote)

### Local Development

1. **Clone and setup environment**
```bash
git clone https://github.com/egnd09/alfred-go.git
cd alfred-go
cp .env.example .env
# Edit .env with your configuration
```

2. **Start infrastructure**
```bash
./scripts/dev.sh
```

This will:
- Start MongoDB and Redis via Docker
- Download Go dependencies
- Run the server on port 5500

3. **Or use Docker Compose directly**
```bash
docker compose up --build
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVICE_PORT` | Server port | `5500` |
| `DB_MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017/alfred` |
| `REDIS_URL` | Redis URL | `redis://localhost:6379` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `GITHUB_CLIENTID` | GitHub OAuth app ID | - |
| `GITHUB_CLIENTSECRET` | GitHub OAuth secret | - |
| `AWS_ACCESS_KEY_ID` | AWS access key | - |
| `AWS_SECRET_ACCESS_KEY` | AWS secret | - |
| `AWS_DEFAULT_REGION` | AWS region | `us-west-2` |
| `EKS_CLUSTER_NAME` | EKS cluster name | - |
| `KUBECONFIG` | Kubeconfig path (optional) | - |

## API Endpoints

### REST

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Health check |
| GET | `/health` | Health check |
| POST | `/api/login` | GitHub OAuth login |

### WebSocket

Connect to `ws://localhost:5500/ws` with JWT token in query params.

**Events:**

| Event | Description |
|-------|-------------|
| `new_env` | Create environment |
| `delete_env` | Remove environment |
| `default_env` | Select active environment |
| `env_list` | List environments |
| `container_list` | List containers |
| `container_status` | Get container status |
| `kill_pod` | Kill a pod |
| `get_docker_logs` | Stream pod logs |
| `new_build` | Start build |
| `cancel_build` | Cancel build |
| `get_tags` | Get repo tags/branches |
| `get_last_builds` | Get build history |

## Development

### Build

```bash
./scripts/build.sh
# or
go build -o bin/server ./cmd/server
```

### Test

```bash
./scripts/test.sh
# or
go test -v ./...
```

### Docker Build

```bash
./scripts/docker-build.sh
# or
docker build -t alfred-go:latest .
```

## Frontend Integration

The React frontend from the original Alfred project is compatible with this Go backend.

### Setup Frontend

1. Copy the client directory from the original project:
```bash
cp -r ../alfred/client ./client
```

2. Update WebSocket URL in frontend to connect to Go server:
```javascript
// In your socket configuration
const socket = io('http://localhost:5500', {
  transports: ['websocket'],
  query: { token: yourJwtToken }
});
```

3. Start the frontend:
```bash
cd client
npm install
npm start
```

## Docker Compose

```yaml
version: '3.8'
services:
  server:
    build: .
    ports:
      - "5500:5500"
    environment:
      - DB_MONGO_URI=mongodb://mongo:27017/alfred
      - REDIS_URL=redis://redis:6379
    depends_on:
      - mongo
      - redis

  mongo:
    image: mongo:7
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  mongo_data:
```

## Migration Notes

This is a complete rewrite from Node.js to Go. Key changes:

1. **HTTP Framework**: Koa → Gin
2. **WebSocket**: Socket.IO → Gorilla WebSocket
3. **Config**: dotenv → Viper
4. **Logging**: Winston → Zap
5. **Database**: Mongoose → mongo-driver
6. **Architecture**: Event-driven with clean separation

## License

MIT
