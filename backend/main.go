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
	"nira/memory"
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

	db, err := memory.NewDatabase(config.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	memManager, err := memory.NewManager(db)
	if err != nil {
		log.Fatalf("Failed to initialize memory manager: %v", err)
	}

	ollamaClient := NewOllamaClient(config.OllamaEndpoint, config.DefaultModel)

 toolRegistry := tools.NewRegistry()
 fileReadTool := tools.NewFileReadTool(config.AllowedPaths)
 toolRegistry.Register(fileReadTool)

 fileWriteTool := tools.NewFileWriteTool(config.AllowedPaths)
 toolRegistry.Register(fileWriteTool)

 // RAG foundation tools
 listDirTool := tools.NewListDirectoryTool(config.AllowedPaths)
 toolRegistry.Register(listDirTool)
 searchByNameTool := tools.NewSearchFilesByNameTool(config.AllowedPaths)
 toolRegistry.Register(searchByNameTool)
 fileMetaTool := tools.NewFileMetadataTool(config.AllowedPaths)
 toolRegistry.Register(fileMetaTool)

	// Register WebSearchTool
	tools.RegisterWebSearchTool(toolRegistry.Tools)

	server := NewServer(config.WebSocketPort, ollamaClient, toolRegistry, logger, memManager)

	log.Println("Starting NIRA backend...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
