#!/bin/bash
set -e

echo "Running tests..."

# Run all tests with verbose output
go test -v ./...

echo "All tests passed!"
