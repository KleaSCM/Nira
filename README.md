# NIRA - Local AI Assistant

NIRA is a locally-run, tool-empowered AI assistant using Ollama, with a Flutter/Dart GUI, Go backend, and SQLite memory. The system is modular, extensible, and capable of orchestrating tools safely and intelligently.

## Architecture

- **Frontend**: Flutter/Dart GUI
- **Backend Core**: Go service managing models, tools, memory, permissions, processes
- **Model Runtime**: Ollama (initially HammerAI/mythomax-l2; model slot interchangeable)
- **Memory Layer**: SQLite for long-term, compressed, vector-supported memory
- **Tool Framework**: Safe, typed, bidirectional tools Nira can call

## Features

- Responsive GUI for interacting with Nira
- Dynamic tool calling (file I/O, system commands, APIs, etc.)
- Persistent context and long-term memory
- Simple model swapping in Ollama
- Local RAG (Retrieval Augmented Generation) capabilities
- Roleplay engine with character and story management

## Project Structure

```
Nira/
├── backend/          # Go backend service
├── frontend/         # Flutter/Dart GUI
├── Docs/             # Project documentation
└── Tests/            # Test files
```

## Development

See `Scope.txt` for detailed project requirements and `StyleGuide.txt` for coding standards.

## Author

Author: KleaSCM  
Email: KleaSCM@gmail.com

