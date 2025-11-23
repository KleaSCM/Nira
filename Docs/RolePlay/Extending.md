# Extending RP

This guide explains how to extend the RolePlay (RP) subsystem: adding new widgets, fields, and backend integrations while keeping the codebase maintainable.

Principles
- Keep UI components focused and composable
- Centralize data logic in the repository/manager layers
- Prefer incremental, schema-compatible changes to models
- Document new features in Docs/RolePlay and reference them from the root README

Add a new RP feature (example: "Factions")
1) Model
   - Update `frontend\lib\RolePlay\roleplay_models.dart` with a new `Faction` class
   - Add fields like id, name, description, reputation, memberIds
2) Repository
   - Extend `roleplay_repository.dart` with get/save/delete methods for Faction
   - Provide in-memory storage initially; later wire to SQLite via backend
3) UI
   - Create `FactionList.dart` and `FactionEditor.dart`
   - Link from `RolePlayDashboard.dart`
4) Session impacts
   - Update `SessionManager.dart` if sessions can be aligned with factions (e.g., active faction, rep modifiers)
5) Tests/Docs
   - Add usage docs to this folder (e.g., Docs/RolePlay/Factions.md)
   - Update `frontend/README.md` and root README RP section for discoverability

Adding RP-specific tools
- Define a backend tool in `backend/tools` (e.g., `rp_context_inject.go`)
- Register it in `backend/main.go`
- Extend `server.go` system prompt with RP-specific instructions when RP tab/session is active (future enhancement)
- Expose a launcher on the frontend if needed; otherwise let the model invoke it

Persisting RP data
- Short term: use repository stubs to keep UI responsive
- Mid term: add backend RPCs or REST endpoints to save/load via SQLite
- Consider a migration path for existing in-memory data

Style and UX
- Follow `StyleGuide.txt` and prefer existing color/typography tokens in ChatScreen
- Keep interactions consistent with normal chat (e.g., streaming behavior, error toasts)

Performance
- Lazy-load large lists (characters, cards)
- Batch save where possible; debounce editor updates

Security
- Avoid executing external content; sanitize markdown if rendering
- Gate file and network access by reusing the existing tool sandbox
