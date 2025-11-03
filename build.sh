#!/bin/bash

# Build script for CodexGigantus
# This script builds the Go project and optionally installs it to a system path

set -e

echo "Building CodexGigantus..."

# Build the project
go build -o CodexGigantus

# Make the binary executable
chmod +x CodexGigantus

echo "✓ Build successful: CodexGigantus binary created"
echo ""
echo "To install shell completions, run:"
echo "  ./CodexGigantus install-completion"
echo ""
echo "To add the binary to your PATH, you can:"
echo "  1. Copy it to /usr/local/bin (requires sudo):"
echo "     sudo cp CodexGigantus /usr/local/bin/"
echo "  2. Or add this directory to your PATH:"
echo "     export PATH=\$PATH:$(pwd)"
