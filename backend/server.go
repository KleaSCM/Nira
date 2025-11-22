/**
 * WebSocket server module.
 *
 * Handles WebSocket connections from the Flutter frontend, manages
 * message routing, and streams responses back to the client.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: server.go
 * Description: WebSocket server implementation.
 */

package main

import (
	"fmt"
	"net/http"
	"nira/tools"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	Port         int
	Ollama       *OllamaClient
	ToolRegistry *tools.Registry
	ToolHandler  *ToolHandler
	Logger       *Logger
	Conversation []ChatMessage
}

func NewServer(port int, ollama *OllamaClient, registry *tools.Registry, logger *Logger) *Server {
	toolHandler := NewToolHandler(registry, logger)
	return &Server{
		Port:         port,
		Ollama:       ollama,
		ToolRegistry: registry,
		ToolHandler:  toolHandler,
		Logger:       logger,
		Conversation: []ChatMessage{},
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/ws", s.HandleWebSocket)
	address := fmt.Sprintf(":%d", s.Port)
	s.Logger.Info("Server starting on %s", address)
	return http.ListenAndServe(address, nil)
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Error("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	s.Logger.LogWebSocketEvent("connection", "established")

	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			s.Logger.Error("WebSocket read error: %v", err)
			break
		}

		if msg.Type == MessageTypeUser {
			s.handleUserMessage(conn, msg.Content)
		}
	}
}

func (s *Server) handleUserMessage(conn *websocket.Conn, content string) {
	userMsg := ChatMessage{
		Role:    "user",
		Content: content,
	}
	s.Conversation = append(s.Conversation, userMsg)

	systemPrompt := s.buildSystemPrompt()
	systemMsg := ChatMessage{
		Role:    "system",
		Content: systemPrompt,
	}
	messages := append([]ChatMessage{systemMsg}, s.Conversation...)

	maxIterations := 5
	for i := 0; i < maxIterations; i++ {
		assistantContent := ""
		startTime := time.Now()

		err := s.Ollama.Chat(messages, func(chunk string) error {
			assistantContent += chunk
			chunkMsg := WSMessage{
				Type:    MessageTypeChunk,
				Content: chunk,
			}
			return conn.WriteJSON(chunkMsg)
		})

		duration := time.Since(startTime)
		s.Logger.LogOllamaResponse(duration, 0)

		if err != nil {
			s.Logger.Error("Ollama chat error: %v", err)
			errorMsg := WSMessage{
				Type:    MessageTypeError,
				Content: fmt.Sprintf("Error: %v", err),
			}
			conn.WriteJSON(errorMsg)
			return
		}

		// Check for tool calls in the response
		toolCall, hasToolCall := s.ToolHandler.DetectToolCall(assistantContent)
		if !hasToolCall {
			// No tool call, final response
			assistantMsg := ChatMessage{
				Role:    "assistant",
				Content: assistantContent,
			}
			s.Conversation = append(s.Conversation, assistantMsg)

			doneMsg := WSMessage{
				Type:    MessageTypeAssistant,
				Content: assistantContent,
			}
			conn.WriteJSON(doneMsg)
			return
		}

		// Execute tool call
		s.Logger.Info("Detected tool call: %s", toolCall.Name)
		toolResult, err := s.ToolHandler.ExecuteTool(toolCall)
		if err != nil {
			s.Logger.Error("Tool execution failed: %v", err)
			errorMsg := WSMessage{
				Type:    MessageTypeError,
				Content: fmt.Sprintf("Tool error: %v", err),
			}
			conn.WriteJSON(errorMsg)
			return
		}

		// Inject tool result back into conversation
		toolResultStr := s.ToolHandler.FormatToolResult(toolCall.Name, toolResult)
		toolMsg := ChatMessage{
			Role:    "user",
			Content: toolResultStr,
		}
		s.Conversation = append(s.Conversation, toolMsg)

		// Send tool result to frontend
		toolResultMsg := WSMessage{
			Type:    MessageTypeSystem,
			Content: fmt.Sprintf("Tool %s executed: %s", toolCall.Name, toolResultStr),
		}
		conn.WriteJSON(toolResultMsg)

		// Continue conversation with tool result
		messages = append(messages, ChatMessage{
			Role:    "assistant",
			Content: assistantContent,
		}, toolMsg)
	}

	s.Logger.Warn("Maximum tool call iterations reached")
}

func (s *Server) buildSystemPrompt() string {
	prompt := "You are NIRA, a helpful AI assistant. Be concise and friendly.\n\n"
	prompt += "Available tools:\n"

	toolsList := s.ToolRegistry.ListTools()
	for _, tool := range toolsList {
		if name, ok := tool["name"].(string); ok {
			if desc, ok := tool["description"].(string); ok {
				prompt += fmt.Sprintf("- %s: %s\n", name, desc)
			}
		}
	}

	prompt += "\nTo use a tool, respond with a JSON object like: {\"name\": \"tool_name\", \"arguments\": {\"arg1\": \"value1\"}}\n"
	prompt += "Or use the format: tool_name(arg1=\"value1\", arg2=\"value2\")\n"

	return prompt
}
