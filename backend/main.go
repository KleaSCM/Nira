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

 // Initialize Allowed Directories store and seed from config
 allowedStore, err := memory.NewAllowedDirsStore(db)
 if err != nil {
     log.Fatalf("Failed to initialize allowed directories store: %v", err)
 }
 if err := allowedStore.EnsureSeed(config.AllowedPaths); err != nil {
     log.Printf("Warning: failed to seed allowed directories: %v", err)
 }
 memManager.AllowedDirs = allowedStore

	ollamaClient := NewOllamaClient(config.OllamaEndpoint, config.DefaultModel)

 toolRegistry := tools.NewRegistry()
 // Use centralized AllowedDirs store for permission checks
 fileReadTool := tools.NewFileReadToolWithChecker(config.AllowedPaths, allowedStore)
 toolRegistry.Register(fileReadTool)

 fileWriteTool := tools.NewFileWriteToolWithChecker(config.AllowedPaths, allowedStore)
 toolRegistry.Register(fileWriteTool)

 // RAG foundation tools
 listDirTool := tools.NewListDirectoryToolWithChecker(config.AllowedPaths, allowedStore)
 toolRegistry.Register(listDirTool)
 searchByNameTool := tools.NewSearchFilesByNameToolWithChecker(config.AllowedPaths, allowedStore)
 toolRegistry.Register(searchByNameTool)
 fileMetaTool := tools.NewFileMetadataToolWithChecker(config.AllowedPaths, allowedStore)
 toolRegistry.Register(fileMetaTool)

 // Allowed directory management tools
 toolRegistry.Register(tools.NewAllowedDirsListTool(allowedStore))
 toolRegistry.Register(tools.NewAllowedDirsAddTool(allowedStore))
 toolRegistry.Register(tools.NewAllowedDirsRemoveTool(allowedStore))

 // Basic RAG indexing and retrieval tools
 ragIndex := memory.NewRagIndex(db)
 toolRegistry.Register(tools.NewRagIndexFolderTool(allowedStore, ragIndex))
 toolRegistry.Register(tools.NewRagSearchTool(ragIndex, allowedStore))

	// Register WebSearchTool
	tools.RegisterWebSearchTool(toolRegistry.Tools)

	server := NewServer(config.WebSocketPort, ollamaClient, toolRegistry, logger, memManager)

	log.Println("Starting NIRA backend...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
