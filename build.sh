#!/bin/bash

# Build the Go project
go build -o llm-codeinject ./cmd/main.go

# Make the binary executable
chmod +x llm-codeinject

# Add the binary to the PATH (Optional: Adjust the path as per your setup)
echo "export PATH=\$PATH:$(pwd)" >> ~/.bashrc
source ~/.bashrc
