Tool: web_search

Overview
- Performs a web search and returns a list of results with title, URL, and optional snippet.
- Implemented in backend/tools/web_search.go and registered from backend/main.go via tools.RegisterWebSearchTool().

Identifier
- name: web_search

Arguments
- query (string, required): The search query text.

Returns
- Success: array of objects (WebSearchResult)
  - title: string
  - url: string
  - snippet: string (may be empty)
- Failure: error propagated to WebSocket as a message of type "error".

Frontend usage (direct call)
- The UI shows a dialog, then sends { "name": "web_search", "arguments": { "query": "<text>" } }.
- The backend formats results into a readable list before streaming to the UI.

AI-initiated usage
- The model can reply with {"name":"web_search","arguments":{"query":"best go websocket packages"}}.
- The server executes and injects the formatted results back into the conversation.

Formatting and streaming
- The server formats results with numbering, title, snippet (when available), and a link on separate lines.
- Example (abbreviated):
  1. Example Title
     Example snippet
     🔗 https://example.com

Common errors
- query argument missing or not a string.
- Upstream search failure or networking issues.

Testing checklist
- Successful query: should return a numbered list.
- Empty query or invalid type: should error.
- Simulate network failure: verify error handling.

Source
- backend/tools/web_search.go
- backend/server.go (formatting in formatWebSearchResults)
