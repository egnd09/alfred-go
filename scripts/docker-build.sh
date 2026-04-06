#!/bin/bash
set -e

# Build Docker image
echo "Building Docker image..."

docker build -t alfred-go:latest .

echo "Docker image built: alfred-go:latest"
