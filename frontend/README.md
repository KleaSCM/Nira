# NIRA Frontend (Flutter)

This is the Flutter/Dart user interface for NIRA. It provides:
- A modern chat UI for normal assistant interactions
- A dedicated RolePlay (RP) workspace with its own chat surface and editors
- Tool launchers (Web Search, File Read/Write) that send direct tool calls to the backend

## Quick Start

Prerequisites
- Flutter SDK installed and on PATH
- Backend running locally on ws://localhost:8080/ws (see project root README)

Install and run
1. flutter pub get
2. flutter run -d chrome

If you use a different backend port/host, update the WebSocket URL in your app initialization (ChatScreen uses ws://localhost:8080/ws by default).

## Key Files
- lib/ChatScreen.dart: Main chat UI (tab 1), WebSocket connection, tool buttons
- lib/WebSocketService.dart: Typed WebSocket messaging and helper methods
- lib/RolePlay/RolePlayDashboard.dart: RP entry point (tab 2)
- lib/RolePlay/RPChatScreen.dart: RP chat surface
- lib/RolePlay/CharacterList.dart & CharacterEditor.dart: Character management
- lib/RolePlay/StoryCardList.dart & StoryCardEditor.dart: Story card management
- lib/RolePlay/SessionManager.dart: Session lifecycle helpers (in-progress)
- lib/RolePlay/roleplay_models.dart: Data models
- lib/RolePlay/roleplay_repository.dart: Repository stub

## Using Tools from the UI
- Web Search: Click the Web Search button → enter query → results stream into chat
- File Read/Write: Click File → pick a file → choose Read or Write
  - Read: Sends { name: "read_file", arguments: { path } }
  - Write: Prompts for text then sends { name: "write_file", arguments: { path, content } }

See Docs/Tools/*.md for details on each tool.

## RolePlay (RP) Mode
- Switch to the RP tab in the app header
- Explore the Dashboard, create Characters and Story Cards
- Open RP Chat to converse in an RP context (dedicated surface)
- For deeper docs, start here:
  - Docs/RolePlay/Overview.md
  - Docs/RolePlay/Components.md
  - Docs/RolePlay/DataModels.md
  - Docs/RolePlay/Flows.md
  - Docs/RolePlay/Extending.md

## Development Tips
- The UI uses Google Fonts and flutter_animate for polish; ensure these packages are installed
- When iterating on WebSocket behavior, check backend/server.go logs for message flow
- Keep UI logic in widgets light; prefer moving stateful data to repositories/managers

## Troubleshooting
- No connection? Ensure backend is running and reachable at ws://localhost:8080/ws
- File tool errors? Verify backend AllowedPaths in backend/config.go include your target directory
- Web search issues? Check network and backend tool registration
