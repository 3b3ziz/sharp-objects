#!/bin/bash

# Sharp Objects - Build and Deploy Script
# Builds the binary and installs it to ~/.local/bin

set -e

echo "Building sharp-objects..."
go build -o sharp-objects

echo "Installing to ~/.local/bin..."
mkdir -p ~/.local/bin
cp sharp-objects ~/.local/bin/
chmod +x ~/.local/bin/sharp-objects

echo "✓ Deployed successfully!"
echo "Run 'sharp-objects' from anywhere to start."
