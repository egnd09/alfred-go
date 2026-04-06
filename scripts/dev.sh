#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting Alfred Development Environment${NC}"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

# Load environment variables
if [ -f .env ]; then
    echo -e "${GREEN}Loading environment from .env${NC}"
    export $(cat .env | grep -v '^#' | xargs)
else
    echo -e "${YELLOW}No .env file found, using defaults${NC}"
    export SERVICE_PORT=${SERVICE_PORT:-5500}
    export DB_MONGO_URI=${DB_MONGO_URI:-mongodb://localhost:27017/alfred}
    export REDIS_URL=${REDIS_URL:-redis://localhost:6379}
    export LOG_LEVEL=${LOG_LEVEL:-debug}
fi

# Start infrastructure
echo -e "${GREEN}Starting MongoDB and Redis...${NC}"
docker compose up -d mongo redis

# Wait for services to be ready
echo -e "${YELLOW}Waiting for services...${NC}"
sleep 3

# Check if go is installed
if command -v go &> /dev/null; then
    echo -e "${GREEN}Go found, running locally...${NC}"
    
    # Download dependencies
    echo -e "${GREEN}Downloading Go dependencies...${NC}"
    go mod tidy
    go mod download
    
    # Run the server
    echo -e "${GREEN}Starting Alfred server on port ${SERVICE_PORT}...${NC}"
    go run ./cmd/server
else
    echo -e "${YELLOW}Go not found, using Docker...${NC}"
    
    # Build and run with Docker
    echo -e "${GREEN}Building Docker image...${NC}"
    docker compose build server
    
    echo -e "${GREEN}Starting Alfred server...${NC}"
    docker compose up server
fi