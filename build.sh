#!/bin/bash

# Build the Go project
go build -o codexgigantus ./cmd/main.go

# Make the binary executable
chmod +x codexgigantus

# Add the binary to the PATH (Optional: Adjust the path as per your setup)
echo "export PATH=\$PATH:$(pwd)" >> ~/.bashrc

# For zsh, also add it to ~/.zshrc
echo "export PATH=\$PATH:$(pwd)" >> ~/.zshrc

# Source both .bashrc and .zshrc
source ~/.bashrc
source ~/.zshrc
