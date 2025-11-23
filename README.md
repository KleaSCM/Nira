# NIRA - Local AI Assistant

NIRA is a locally-run, tool-empowered AI assistant that uses Ollama for LLM inference, a Flutter/Dart GUI for the frontend, and a Go backend with a SQLite memory layer. It is designed to be modular, extensible, and safe, with a typed tool framework that the model can leverage to perform actions such as reading/writing files and searching the web.

## Architecture Overview

- Frontend: Flutter/Dart GUI (desktop/web) communicating over WebSocket
- Backend Core: Go service managing models, tools, memory, and conversation state
- Model Runtime: Ollama (default model: HammerAI/mythomax-l2; easily swappable)
- Memory Layer: SQLite for long-term memory and conversation history
- Tool Framework: Safe, typed, bidirectional tools that the assistant can call

System diagram

```
┌──────────────────────────────┐            ws://localhost:8080/ws            ┌──────────────────────────────────┐
│        Flutter Frontend      │  ◀━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━▶ │            Go Backend            │
│  - Chat UI (ChatScreen)      │                                              │  - WebSocket server (server.go)  │
│  - RP UI (RolePlay/*)        │                                              │  - Tool registry (tools/tool.go) │
│  - WebSocketService          │                                              │  - Tool handler (tool_handler.go)│
└──────────────┬───────────────┘                                              │  - Memory (memory/*, SQLite)     │
               │                                                               │  - Ollama client (ollama.go)     │
               │                                                               └───────────────┬─────────────────┘
               │                                                                 HTTP (REST)   │
               │                                                                                ▼
               │                                                                        ┌─────────────┐
               │                                                                        │   Ollama    │
               │                                                                        └─────────────┘
               │
               │  Direct tool calls (JSON): { name, arguments }
               ▼
      [File Picker / Web Search]

Notes:
- The frontend can invoke tools directly (e.g., read_file, write_file, web_search) by sending a JSON tool call.
- The LLM may also request tools; tool_handler detects and executes them, injecting results back into conversation.
- Memory persists user/assistant messages and tool results in SQLite for continuity across sessions.
```

## Key Features

- Modern, responsive chat UI
- Dynamic tool calling (file read/write, web search; more coming)
- Persistent conversation history via SQLite
- Pluggable LLM via Ollama
- Early RAG foundations and RolePlay (RP) mode groundwork

## Getting Started

Prerequisites
- Go 1.21+
- Flutter SDK (for web/desktop)
- Ollama running locally (default: http://localhost:11434)

1) Backend
- Open a terminal in backend
- go run .
  - The server listens on ws://localhost:8080/ws

2) Frontend
- Open a terminal in frontend
- flutter pub get
- flutter run -d chrome

If you use a different Ollama endpoint, update the configuration (see below).

## Configuration

Configuration is currently defined in backend/config.go. Defaults:
- OllamaEndpoint: http://localhost:11434
- DefaultModel: HammerAI/mythomax-l2
- DatabasePath: ./nira.db
- WebSocketPort: 8080
- AllowedPaths: ["."] (sandbox for file tools; restricts to project directory by default)

To permit file tools to access other directories, add their absolute paths to AllowedPaths in config.go. Keep security in mind and prefer the minimum necessary scope.

## Using NIRA

- Chat: Type a message in the input and hit Enter or click Send.
- Web Search: Click the Web Search tool button, enter a query; results stream into chat.
- Files: Click the File tool button to choose a file, then select Read or Write.
  - Read: Streams file content into chat.
  - Write: Prompts for text and writes it to the selected path (overwrites existing content).

### RolePlay (RP) Mode

The frontend includes an experimental RolePlay (RP) workspace with dedicated screens and models under frontend/lib/RolePlay/.

What you can do today:
- Browse the RP dashboard and switch to the RP tab in the app header.
- Create and edit Characters and Story Cards (UI scaffolding available; logic will expand in later phases).
- Use a separate RP chat surface (RPChatScreen) that will later integrate with RP-specific memory and rules.

Key RP components (frontend/lib/RolePlay/):
- RolePlayDashboard.dart: Entry point for RP features (tab 2 in the UI)
- RPChatScreen.dart: Chat surface dedicated to RP sessions
- CharacterEditor.dart / CharacterList.dart: Create and browse characters
- StoryCardEditor.dart / StoryCardList.dart: Create and browse story beats/lore
- SessionManager.dart: Session lifecycle utilities (in-progress)
- roleplay_models.dart: Data classes (Character, StoryCard, Session, etc.)
- roleplay_repository.dart: Placeholder repository (will back with SQLite in future)

Read the full RP documentation for details, flows, and extension points:
- Docs/RolePlay/Overview.md
- Docs/RolePlay/Components.md
- Docs/RolePlay/DataModels.md
- Docs/RolePlay/Flows.md
- Docs/RolePlay/Extending.md

## Tools Documentation

Each tool has its own dedicated documentation under Docs/Tools/:
- Docs/Tools/read_file.md
- Docs/Tools/write_file.md
- Docs/Tools/web_search.md

Quick summary
- read_file: Reads text from a file within AllowedPaths.
- write_file: Writes text to a file (creates parent directories; overwrites).
- web_search: Performs a web search and returns a list of results.

Refer to the per-tool docs above for arguments, return formats, examples, and security notes.

## Memory Layer

NIRA saves conversation history and basic memory constructs in SQLite. A deeper Phase 2 memory design is captured here:
- Docs/Phase2_Memory_Design.md

## Project Structure

A more complete view of the repository to help you navigate quickly.

```
Nira/
├── backend/                                  # Go backend service
│   ├── main.go                               # Entrypoint (config, registry, server)
│   ├── server.go                             # WebSocket server, streaming, tool & chat loop
│   ├── tool_handler.go                       # Detects/executes AI-initiated tool calls
│   ├── config.go                             # Runtime configuration (Ollama, DB, AllowedPaths)
│   ├── logger.go                             # Structured logging helpers
│   ├── ollama.go                             # Minimal Ollama client wrapper
│   ├── memory/                               # Conversation + memory persistence (SQLite)
│   │   ├── database.go                       # DB connection and init
│   │   ├── manager.go                        # Manager orchestrating memory operations
│   │   ├── conversation.go                   # Conversation message storage
│   │   ├── memory.go                         # Memory interfaces/types
│   ├── tools/                                # Tool framework + implementations
│   │   ├── tool.go                           # Tool interface and registry
│   │   ├── file_read.go                      # read_file tool (sandboxed by AllowedPaths)
│   │   ├── file_write.go                     # write_file tool (sandboxed by AllowedPaths)
│   │   └── web_search.go                     # web_search tool
│   └── tests/                                # Backend tests
│       ├── database_test.go
│       ├── conversation_store_test.go
│       ├── integration_test.go
│       ├── manager_test.go
│       ├── memory_store_test.go
│       └── memory_test.go
│
├── frontend/                                 # Flutter/Dart GUI
│   ├── lib/
│   │   ├── main.dart                         # App bootstrap
│   │   ├── App.dart                          # App root widget
│   │   ├── ChatScreen.dart                   # Main chat UI with tool buttons
│   │   ├── WebSocketService.dart             # Typed WebSocket client
│   │   └── RolePlay/                         # RP workspace (experimental)
│   │       ├── RolePlayDashboard.dart        # RP entry point (tab 2)
│   │       ├── RPChatScreen.dart             # RP chat surface
│   │       ├── CharacterList.dart / CharacterEditor.dart
│   │       ├── StoryCardList.dart / StoryCardEditor.dart
│   │       ├── SessionManager.dart           # Session lifecycle utilities
│   │       ├── roleplay_models.dart          # Character, StoryCard, Session, etc.
│   │       └── roleplay_repository.dart      # Repository abstraction (stubs)
│   ├── README.md                             # Frontend-specific guide
│   └── web/windows/...                       # Flutter platform scaffolding
│
├── Docs/                                     # Documentation hub
│   ├── Tools/                                # Per-tool docs
│   │   ├── read_file.md
│   │   ├── write_file.md
│   │   └── web_search.md
│   ├── RolePlay/                             # RP documentation suite
│   │   ├── Overview.md
│   │   ├── Components.md
│   │   ├── DataModels.md
│   │   ├── Flows.md
│   │   └── Extending.md
│   └── Phase2_Memory_Design.md               # Future memory design notes
│
├── Tests/                                    # External test client & harness
│   ├── test_client.go
│   ├── test_client.exe
│   ├── go.mod
│   └── go.sum
│
├── progress.txt                              # Roadmap/progress checklist
├── Scope.txt                                 # Scope notes
├── StyleGuide.txt                            # Style guidelines
└── README.md                                 # This file
```

## Roadmap and Status

See progress.txt for a living checklist of planned phases (tools, memory, UI, RP, RAG).

## Author

Author: KleaSCM
Email: KleaSCM@gmail.com

