# RP Components

This document describes the primary Flutter widgets that implement the RolePlay (RP) UI under `frontend\lib\RolePlay`.

Top-level
- RolePlayDashboard.dart
  - Purpose: Landing surface for RP. From here users access Characters, Story Cards, RP Chat, and settings.
  - Notes: Designed to be the second tab in the main app (`Chat` | `RP`).

Chat
- RPChatScreen.dart
  - Purpose: A chat surface dedicated to RP sessions. It will integrate session context (active character(s), scene, lore) in future phases.
  - Key responsibilities:
    - Render RP conversation messages
    - Provide input for user messages in RP context
    - Later: show session state (active scene, mood, party)

Characters
- CharacterList.dart
  - Purpose: Lists existing characters with basic filters and actions (edit, select for session).
  - Interacts with: `roleplay_repository.dart` for data source.
- CharacterEditor.dart
  - Purpose: Create/edit a character’s sheet (name, traits, background, goals, etc.).
  - Future fields: inventory, relationships, secrets, custom attributes.

Story Cards
- StoryCardList.dart
  - Purpose: Browse/search story cards (lore entries, scenes, items, locations, events).
  - Interacts with: `roleplay_repository.dart` for data source.
- StoryCardEditor.dart
  - Purpose: Create/edit story cards; attach tags/links to characters or sessions.
  - Future fields: embedding metadata, retrieval hints, scene triggers.

Session and Data Layers
- SessionManager.dart
  - Purpose: Helper utilities for starting/stopping RP sessions, switching active character, and tracking session metadata.
  - Future: Summaries, checkpoints, scene transitions.
- roleplay_models.dart
  - Purpose: Data models for RP domain: Character, StoryCard, Session, etc.
  - Keeps fields concise now; extensible for future RAG/RP requirements.
- roleplay_repository.dart
  - Purpose: Repository abstraction for persistence. Current state may be in-memory or simplified; will back onto SQLite (via backend) in later phases.

Integration with the rest of the app
- The main `ChatScreen` offers a TabBar; the `RP` tab displays `RolePlayDashboard`.
- WebSocket connectivity is shared with the normal chat; RP will later add RP-specific prompts and memory.
