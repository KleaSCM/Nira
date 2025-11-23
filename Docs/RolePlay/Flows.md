# RP UI Flows

This guide outlines common user flows in the RP workspace and how the widgets/repository interact.

1) Create a Character
- Navigate: RP tab → Characters → New
- Fill basic fields (name, summary, traits) and save
- Repository: saveCharacter(character) persists it
- CharacterList refreshes to include the new entry

2) Create a Story Card
- Navigate: RP tab → Story Cards → New
- Choose a kind (lore/scene/item/location/event/custom), add content and tags
- Save and link to relevant characters if needed (via ids in links)

3) Start an RP Session
- Navigate: RP tab → Start Session
- Select character(s) to include in the party
- SessionManager creates a session and sets it active
- RPChatScreen opens bound to the active session

4) RP Chat
- In RPChatScreen, messages are appended to the active session
- User types → message appended with role=user
- Assistant responses stream from backend WebSocket → role=assistant
- Future: RP-specific prompts, memory injection, and scene state displayed in the UI

5) Edit During Session
- Open CharacterEditor from the session to tweak attributes
- Edit StoryCards to adjust lore or scene prompts
- Repository emits updates; active views refresh accordingly

6) End Session
- SessionManager marks session inactive, computes summary, and updates updatedAt
- Future: persist session recap and key memory to SQLite

Backend Interaction
- Uses the same WebSocket as normal chat for now
- Tool calls (web search, file read/write) are available from the main chat; RP-specific tools will arrive later

Notes
- All flows are designed to be offline-first, evolving toward persistent storage
- Keep UI pure; lean on repository and manager for mutations and side effects
