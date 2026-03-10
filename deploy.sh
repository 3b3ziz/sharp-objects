#!/bin/bash

# Sharp Objects - Build and Deploy Script
# Builds the binary and installs it to ~/.local/bin

set -e

# Get version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Building sharp-objects ${VERSION}..."
go build -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o sharp-objects

echo "Installing to ~/.local/bin..."
mkdir -p ~/.local/bin
cp sharp-objects ~/.local/bin/
chmod +x ~/.local/bin/sharp-objects

echo "✓ Deployed ${VERSION} successfully!"
echo "Run 'sharp-objects -v' to verify, or 'sharp-objects' to start."
