# RP Data Models

All RP models live in `frontend\lib\RolePlay\roleplay_models.dart`. Names and fields may evolve; the goal here is to document intent and common extensions.

Character
- id: String (uuid)
- name: String
- summary: String (short bio / pitch)
- traits: List<String> (adjectives, quirks)
- background: String (long form)
- goals: List<String>
- tags: List<String> (search/routing hints)
- notes: String
- Future: relationships, inventory, custom stats/attributes

StoryCard
- id: String (uuid)
- title: String
- kind: String (lore | scene | item | location | event | custom)
- content: String (markdown-friendly text)
- tags: List<String>
- links: List<String> (ids of related characters or cards)
- Future: embedding metadata, retrieval hints, time/place fields, images

Session
- id: String (uuid)
- title: String
- characterIds: List<String> (active party)
- createdAt: DateTime
- updatedAt: DateTime
- summary: String (rolling recap)
- state: Map<String, dynamic> (extensible store for scene flags, inventory, etc.)

Message (RP context)
- id: String (uuid)
- sessionId: String
- role: String (user | assistant | system)
- content: String
- timestamp: DateTime
- Future: speaker character id, mood, scene id

Repository Abstraction (roleplay_repository.dart)
- getCharacters(), saveCharacter(), deleteCharacter()
- getStoryCards(), saveStoryCard(), deleteStoryCard()
- getSessions(), saveSession(), deleteSession()
- getMessages(sessionId), appendMessage(sessionId, message)

Persistence Strategy (planned)
- Short term: in-memory or local file-backed (simple JSON) for rapid iteration
- Mid term: leverage backend SQLite with dedicated RP tables
- Long term: indexing + embeddings for RP-aware retrieval
