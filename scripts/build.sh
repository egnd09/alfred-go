#!/bin/bash
set -e

# Build the Go binary
echo "Building Alfred server..."

# Build for current platform
go build -o bin/server ./cmd/server

echo "Build complete: bin/server"
