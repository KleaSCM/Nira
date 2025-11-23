# RolePlay (RP) Overview

This document introduces the RolePlay (RP) subsystem implemented in the Flutter frontend under `frontend\lib\RolePlay`. RP mode provides an immersive, character-driven interaction space that will increasingly diverge from the normal chat experience as memory, scene state, and lore tools are added.

Goals
- Dedicated UI and data model for character sheets, story cards, and sessions.
- Clean separation between Normal Chat and RP Chat contexts.
- Future: RP-specific memory isolation and retrieval, scene engine, world state.

High-level Architecture
- UI: Flutter widgets in `frontend\lib\RolePlay\*` (Dashboard, Editors, Lists, Chat).
- Data Models: `roleplay_models.dart` for Character, StoryCard, Session, etc.
- Repository: `roleplay_repository.dart` abstracts persistence (in-progress; currently in-memory or stubs).
- Session Management: `SessionManager.dart` coordinates the active RP session.
- Backend: normal WebSocket backend; later phases may add RP-specific endpoints and memory tables.

Current Status (Phase foundation)
- RP tab in the app header routes to `RolePlayDashboard`.
- Character and StoryCard editors/lists are scaffolded and functional at the UI level.
- An `RPChatScreen` exists as a dedicated chat surface for RP.
- Persistence is minimal; advanced features will roll out with the Memory/RP phases.

Roadmap (abridged)
1) RP Mode toggle and context separation
2) Character sheets and story cards persisted in SQLite
3) Scene engine with session summaries and world state
4) RP-aware retrieval (inject relevant character/lore into prompts)
5) Visual builders and viewers (full character builder, story card editor, scene view, recap browser)

See also
- Components: `Docs/RolePlay/Components.md`
- Data Models: `Docs/RolePlay/DataModels.md`
- Flows: `Docs/RolePlay/Flows.md`
- Extending: `Docs/RolePlay/Extending.md`
