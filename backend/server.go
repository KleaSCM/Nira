/**
 * WebSocket server module.
 *
 * Handles WebSocket connections from the Flutter frontend, manages
 * message routing, and streams responses back to the client.
 * Supports both AI-initiated and user-initiated tool calls.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: server.go
 * Description: WebSocket server implementation.
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nira/memory"
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
	Memory       *memory.Manager
	Conversation []ChatMessage
}

// DirectToolCall represents a tool call directly from the frontend
type DirectToolCall struct {
    ID        string                 `json:"id,omitempty"`
    Name      string                 `json:"name"`
    Arguments map[string]interface{} `json:"arguments"`
}

func NewServer(port int, ollama *OllamaClient, registry *tools.Registry, logger *Logger, mem *memory.Manager) *Server {
	toolHandler := NewToolHandler(registry, logger)
	return &Server{
		Port:         port,
		Ollama:       ollama,
		ToolRegistry: registry,
		ToolHandler:  toolHandler,
		Logger:       logger,
		Memory:       mem,
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

	// Load recent conversation history
	recentMessages, err := s.Memory.LoadRecentMessages(50)
	if err != nil {
		s.Logger.Warn("Failed to load recent messages: %v", err)
	} else {
		s.Conversation = []ChatMessage{}
		for _, msg := range recentMessages {
			s.Conversation = append(s.Conversation, ChatMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
		s.Logger.Info("Loaded %d messages from history", len(recentMessages))
	}

	for {
		_, rawMsg, err := conn.ReadMessage()
		if err != nil {
			s.Logger.Error("WebSocket read error: %v", err)
			break
		}

		// CRITICAL DEBUG: Log exactly what we received
		s.Logger.Info("📨 Received raw message: %s", string(rawMsg))

		// Try to parse as direct tool call first
		var directToolCall DirectToolCall
		if err := json.Unmarshal(rawMsg, &directToolCall); err == nil && directToolCall.Name != "" {
			s.Logger.Info("✅ Parsed as direct tool call: %s", directToolCall.Name)
			s.handleDirectToolCall(conn, &directToolCall)
			continue
		}

		// Otherwise, parse as regular WSMessage
		var msg WSMessage
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			s.Logger.Error("❌ Failed to parse as WSMessage: %v", err)
			continue
		}

		s.Logger.Info("✅ Parsed WSMessage - Type: '%s', Content: '%s'", msg.Type, msg.Content)
		s.Logger.Info("🔍 Comparing msg.Type ('%s') with MessageTypeUser ('%s')", msg.Type, MessageTypeUser)

		if msg.Type == MessageTypeUser {
			s.Logger.Info("✅ Message type matches! Calling handleUserMessage")
			s.handleUserMessage(conn, msg.Content)
		} else {
			s.Logger.Warn("⚠️ Message type '%s' does not match expected type '%s'", msg.Type, MessageTypeUser)
		}
	}
}

func (s *Server) handleDirectToolCall(conn *websocket.Conn, toolCall *DirectToolCall) {
    s.Logger.Info("Executing direct tool call: %s with args: %v", toolCall.Name, toolCall.Arguments)

	// Get the tool from registry
	tool, exists := s.ToolRegistry.Tools[toolCall.Name]
	if !exists {
		s.Logger.Error("Tool not found: %s", toolCall.Name)
		errorMsg := WSMessage{
			Type:    MessageTypeError,
			Content: fmt.Sprintf("Tool '%s' not found", toolCall.Name),
		}
		conn.WriteJSON(errorMsg)
		return
	}

 // Execute the tool
 result, err := tool.Execute(toolCall.Arguments)
 if err != nil {
     s.Logger.Error("Tool execution failed: %v", err)
     // If this was a silent call with an ID, return a single error message tagged with the ID
     if toolCall.Arguments != nil {
         if silent, ok := toolCall.Arguments["_silent"].(bool); ok && silent {
             errorMsg := WSMessage{
                 Type:    MessageTypeError,
                 Content: fmt.Sprintf("Tool execution failed: %v", err),
                 ID:      toolCall.ID,
             }
             conn.WriteJSON(errorMsg)
             return
         }
     }
     errorMsg := WSMessage{
         Type:    MessageTypeError,
         Content: fmt.Sprintf("Tool execution failed: %v", err),
     }
     conn.WriteJSON(errorMsg)
     return
 }

	s.Logger.Info("Tool executed successfully, result type: %T", result)

 // Format the result
    var resultText string
    switch v := result.(type) {
 case []tools.WebSearchResult:
     resultText = s.formatWebSearchResults(v)
 case string:
     resultText = v
 case map[string]interface{}:
     // Prefer showing primary content field if present
     if content, ok := v["content"].(string); ok {
         resultText = content
     } else {
         // Fallback to JSON
         if jsonBytes, err := json.MarshalIndent(v, "", "  "); err == nil {
             resultText = string(jsonBytes)
         } else {
             resultText = fmt.Sprintf("%v", v)
         }
     }
 default:
        // Try to JSON serialize any other type
        if jsonBytes, err := json.MarshalIndent(result, "", "  "); err == nil {
            resultText = string(jsonBytes)
        } else {
            resultText = fmt.Sprintf("%v", result)
        }
    }

    // If this was a silent call with an ID, send a single system message containing the JSON (or text) and return.
    if toolCall.Arguments != nil {
        if silent, ok := toolCall.Arguments["_silent"].(bool); ok && silent {
            // Prefer to send JSON for structured results
            var payload string
            switch res := result.(type) {
            case string:
                // Wrap plain strings as JSON string literal
                b, _ := json.Marshal(res)
                payload = string(b)
            default:
                if jb, err := json.Marshal(res); err == nil {
                    payload = string(jb)
                } else {
                    // fallback to resultText
                    b, _ := json.Marshal(resultText)
                    payload = string(b)
                }
            }
            reply := WSMessage{
                Type:    MessageTypeSystem,
                Content: payload,
                ID:      toolCall.ID,
            }
            _ = conn.WriteJSON(reply)
            return
        }
    }

    // Add tool result to conversation context (generic wording)
    header := fmt.Sprintf("[Tool %s result]", toolCall.Name)
    if q, ok := toolCall.Arguments["query"]; ok {
        header = fmt.Sprintf("[Tool %s for '%v']", toolCall.Name, q)
    } else if p, ok := toolCall.Arguments["path"]; ok {
        header = fmt.Sprintf("[Tool %s: %v]", toolCall.Name, p)
    }
 toolResultMsg := fmt.Sprintf("%s\n%s", header, resultText)

	s.Conversation = append(s.Conversation, ChatMessage{
		Role:    "user",
		Content: toolResultMsg,
	})

	// Save to memory
	if err := s.Memory.SaveMessage("user", toolResultMsg, "tool_result"); err != nil {
		s.Logger.Warn("Failed to save tool result: %v", err)
	}

	// Stream the result back to the frontend as chunks (for smooth UX)
	s.streamText(conn, resultText)

	// Send completion signal
	doneMsg := WSMessage{
		Type:    MessageTypeAssistant,
		Content: "",
	}
	conn.WriteJSON(doneMsg)
}

func (s *Server) streamText(conn *websocket.Conn, text string) {
	// Stream text in small chunks for better UX
	chunkSize := 50
	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunk := text[i:end]

		chunkMsg := WSMessage{
			Type:    MessageTypeChunk,
			Content: chunk,
		}
		conn.WriteJSON(chunkMsg)
		time.Sleep(10 * time.Millisecond) // Small delay for smooth streaming
	}
}

func (s *Server) formatWebSearchResults(results []tools.WebSearchResult) string {
	if len(results) == 0 {
		return "No search results found."
	}

	output := fmt.Sprintf("🔍 Found %d search results:\n\n", len(results))
	for i, result := range results {
		output += fmt.Sprintf("%d. **%s**\n", i+1, result.Title)
		if result.Snippet != "" && result.Snippet != result.Title {
			output += fmt.Sprintf("   %s\n", result.Snippet)
		}
		output += fmt.Sprintf("   🔗 %s\n\n", result.URL)
	}
	return output
}

func (s *Server) handleUserMessage(conn *websocket.Conn, content string) {
	s.Logger.Info("🎯 handleUserMessage called with content: '%s'", content)

	userMsg := ChatMessage{
		Role:    "user",
		Content: content,
	}
	s.Conversation = append(s.Conversation, userMsg)

	if err := s.Memory.SaveMessage("user", content, ""); err != nil {
		s.Logger.Warn("Failed to save user message: %v", err)
	}

	systemPrompt := s.buildSystemPrompt()
	s.Logger.Info("📝 System prompt length: %d chars", len(systemPrompt))

	systemMsg := ChatMessage{
		Role:    "system",
		Content: systemPrompt,
	}
	messages := append([]ChatMessage{systemMsg}, s.Conversation...)
	s.Logger.Info("📨 Total messages to send to Ollama: %d", len(messages))

	maxIterations := 5
	for i := 0; i < maxIterations; i++ {
		s.Logger.Info("🔄 Iteration %d/%d", i+1, maxIterations)
		assistantContent := ""
		startTime := time.Now()
		chunkCount := 0

		s.Logger.Info("🚀 Calling Ollama.Chat()...")
		err := s.Ollama.Chat(messages, func(chunk string) error {
			chunkCount++
			assistantContent += chunk
			if chunkCount <= 3 {
				s.Logger.Info("📦 Chunk %d: '%s'", chunkCount, chunk)
			}
			chunkMsg := WSMessage{
				Type:    MessageTypeChunk,
				Content: chunk,
			}
			return conn.WriteJSON(chunkMsg)
		})

		duration := time.Since(startTime)
		s.Logger.Info("⏱️ Ollama response completed in %v, received %d chunks", duration, chunkCount)
		s.Logger.LogOllamaResponse(duration, 0)

		if err != nil {
			s.Logger.Error("❌ Ollama chat error: %v", err)
			errorMsg := WSMessage{
				Type:    MessageTypeError,
				Content: fmt.Sprintf("Ollama Error: %v", err),
			}
			conn.WriteJSON(errorMsg)
			return
		}

		s.Logger.Info("✅ Assistant content length: %d chars", len(assistantContent))

		// Check for tool calls in the response (AI-initiated)
		toolCall, hasToolCall := s.ToolHandler.DetectToolCall(assistantContent)
		if !hasToolCall {
			// No tool call, final response
			assistantMsg := ChatMessage{
				Role:    "assistant",
				Content: assistantContent,
			}
			s.Conversation = append(s.Conversation, assistantMsg)

			if err := s.Memory.SaveMessage("assistant", assistantContent, ""); err != nil {
				s.Logger.Warn("Failed to save assistant message: %v", err)
			}

			doneMsg := WSMessage{
				Type:    MessageTypeAssistant,
				Content: assistantContent,
			}
			conn.WriteJSON(doneMsg)
			return
		}

		// Execute tool call (AI-initiated)
		s.Logger.Info("Detected AI tool call: %s", toolCall.Name)
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


		// 1. Add what the assistant just said (the tool call request)
		assistantMsg := ChatMessage{
			Role:    "assistant",
			Content: assistantContent,
		}
		s.Conversation = append(s.Conversation, assistantMsg)


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

		messages = append(messages, assistantMsg)
		messages = append(messages, toolMsg)

		// Loop continues now with updated 'messages'...
	}

	s.Logger.Warn("Maximum tool call iterations reached")
}

func (s *Server) buildSystemPrompt() string {
    prompt := "You are NIRA, a helpful local AI assistant. Be concise and friendly. You can call tools to work with the user's local files.\n\n"
    prompt += "Available tools (name: description):\n"

    // Iterate and format tool schemas
    for _, tool := range s.ToolRegistry.ListTools() {
        prompt += fmt.Sprintf("- %s: %s\n", tool["name"], tool["description"])
    }

    prompt += "\nGeneral rules for tool use:\n"
    prompt += "- Always emit tool calls as a single JSON object: {\"name\":\"tool_name\",\"arguments\":{...}} with no extra text.\n"
    prompt += "- After a tool result is injected back into context, read it and continue the task. If the task requires multiple steps, call additional tools.\n"
    prompt += "- Paths are relative to the project root unless the user provides an absolute path. Prefer ./<folder> style.\n"
    prompt += "- If a file or folder is unclear or not found, ask a brief clarifying question before proceeding.\n"
    prompt += "- For edits, read the file first, compute the full new content, then write it using write_file (overwrite semantics).\n"

    prompt += "\nAllowed directory system:\n"
    prompt += "- You may only access files within the user's allowed directories.\n"
    prompt += "- If access is denied or a path is outside allowed roots, ask the user to allow the directory, or call allowed_dirs_add with their confirmation.\n"
    prompt += "- You can inspect current permissions using allowed_dirs_list.\n"

    prompt += "\nHow to handle common requests:\n"
    prompt += "1) ‘Tell me what files are in <dir>’ → Call list_directory with {path:\"./<dir>\", recursive:false}.\n"
    prompt += "2) ‘Summarize <file> in <dir>’ → If you don't know the exact path:\n   a) Call search_files_by_name with {root:\"./<dir>\", pattern:\"<file>\"}.\n   b) Pick the best match, then call read_file with {path}.\n   c) Write a concise summary as assistant text (no further tool call).\n"
    prompt += "3) ‘Make <change> to <file> in <dir>’ →\n   a) search_files_by_name to find the file,\n   b) read_file to load content,\n   c) modify the text deterministically,\n   d) write_file with the full updated content.\n"

    prompt += "\nIndexing and retrieval (basic local RAG):\n"
    prompt += "- To index a folder of text files for faster search, call rag_index_folder with {root, patterns:[\"*.md\",\"*.txt\"], max_size_mb, max_files}.\n"
    prompt += "- To retrieve relevant files/snippets, call rag_search with {query:\"...\", limit, path_prefix}.\n"
    prompt += "- Always ensure the root/path_prefix is within allowed directories; if not, request permission first.\n"

    prompt += "\nRolePlay (RP) data management (backend-owned):\n"
    prompt += "- Characters: use rp_character_list, rp_character_get, rp_character_save, rp_character_delete.\n"
    prompt += "- Story cards: use rp_storycard_list, rp_storycard_get, rp_storycard_save, rp_storycard_delete.\n"
    prompt += "- When saving, provide full fields; the backend persists them in SQLite.\n"
    prompt += "- IDs are strings. If you omit id on save, a new one will be generated.\n"

    prompt += "\nFew-shot examples (copy the JSON exactly when calling tools):\n"
    prompt += "User: tell me what files are in Docs directory\n"
    prompt += "Assistant: {\"name\":\"list_directory\",\"arguments\":{\"path\":\"./Docs\",\"recursive\":false}}\n\n"

    prompt += "User: summarise README in root directory\n"
    prompt += "Assistant: {\"name\":\"search_files_by_name\",\"arguments\":{\"root\":\".\",\"pattern\":\"README.md\"}}\n\n"

    prompt += "User: make the title more exciting in README.md in root\n"
    prompt += "Assistant: {\"name\":\"search_files_by_name\",\"arguments\":{\"root\":\".\",\"pattern\":\"README.md\"}}\n\n"

    prompt += "User: allow access to D:\\\\Notes and index it\n"
    prompt += "Assistant: {\"name\":\"allowed_dirs_add\",\"arguments\":{\"path\":\"D:\\\\\\Notes\"}}\n\n"

    prompt += "User: find notes about vector search in my notes folder\n"
    prompt += "Assistant: {\"name\":\"rag_search\",\"arguments\":{\"query\":\"vector search\",\"path_prefix\":\"D:\\\\\\Notes\"}}\n\n"

    // RP few-shots
    prompt += "User: create a character named Aria with traits brave and curious\n"
    prompt += "Assistant: {\"name\":\"rp_character_save\",\"arguments\":{\"name\":\"Aria\",\"traits\":[\"brave\",\"curious\"]}}\n\n"

    prompt += "User: list my story cards about the academy\n"
    prompt += "Assistant: {\"name\":\"rp_storycard_list\",\"arguments\":{\"query\":\"academy\",\"limit\":50}}\n\n"

    return prompt
}
