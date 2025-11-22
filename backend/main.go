/**
 * NIRA Backend - Main entry point.
 *
 * Orchestrates the AI assistant service, managing Ollama integration,
 * tool execution, memory persistence, and WebSocket communication with
 * the Flutter frontend.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: main.go
 * Description: Initializes and runs the NIRA backend service.
 */

package main

import (
	"log"
	"nira/tools"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) > 1 && os.Args[1] == "version" {
		log.Println("NIRA Backend v0.1.0")
		os.Exit(0)
	}

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := NewLogger(LogLevelInfo)

	ollamaClient := NewOllamaClient(config.OllamaEndpoint, config.DefaultModel)

	toolRegistry := tools.NewRegistry()
	fileReadTool := tools.NewFileReadTool(config.AllowedPaths)
	toolRegistry.Register(fileReadTool)

	server := NewServer(config.WebSocketPort, ollamaClient, toolRegistry, logger)

	log.Println("Starting NIRA backend...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
