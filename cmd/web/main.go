// Package main provides the web GUI entry point for CodexGigantus
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/baditaflorin/codexgigantus/internal/gui"
	"github.com/baditaflorin/codexgigantus/pkg/env"
)

func main() {
	// Load environment configuration
	envConfig, err := env.Load()
	if err != nil {
		log.Fatalf("Failed to load environment configuration: %v", err)
	}

	// Create GUI server
	server, err := gui.NewServer()
	if err != nil {
		log.Fatalf("Failed to create GUI server: %v", err)
	}

	// Print startup information
	fmt.Println("CodexGigantus Web GUI")
	fmt.Println("====================")
	fmt.Printf("Starting web server on %s:%d\n", envConfig.WebHost, envConfig.WebPort)
	fmt.Println("Open your browser and navigate to the address above")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the server")

	// Start server
	if err := server.Start(envConfig.WebHost, envConfig.WebPort); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
