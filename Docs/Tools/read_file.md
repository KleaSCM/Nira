Tool: read_file

Overview
- Reads the contents of a text file and returns it as UTF-8 text.
- Enforces a filesystem sandbox using AllowedPaths from backend/config.go.
- Implemented in backend/tools/file_read.go.

Identifier
- name: read_file

Arguments
- path (string, required): Absolute or relative path to the target file.

Returns
- Success: JSON object with keys
  - content: string, full text content read from the file
  - path: string, the path that was read
- Failure: error propagated to WebSocket as a message of type "error".

Security and sandboxing
- Path is permitted only if it resides under one of AllowedPaths.
- Default AllowedPaths is ["."] (project root). Adjust in backend/config.go to permit additional directories.
- Absolute resolution plus relative checks prevent directory traversal.

Behavior notes
- Reads the entire file into memory (io.ReadAll). For very large files, consider future chunked IO.
- Intended for text; binary content may appear garbled in the chat UI.

Frontend usage (direct call)
- The frontend sends a JSON object: { name: "read_file", arguments: { path: "<file path>" } }.
- The backend prefers the "content" field when streaming back to the UI, so the chat shows the file text directly.

AI-initiated usage
- The model can reply with {"name":"read_file","arguments":{"path":"./Docs/Phase2_Memory_Design.md"}}.
- The server executes and injects a contextual message with the tool result.

Common errors
- path argument is required and must be a string.
- path not in allowed directories.
- failed to open file: <system error>.
- failed to read file: <system error>.

Testing checklist
- Read a file within project root: should stream content.
- Attempt to read outside AllowedPaths: should error.
- Read a non-existent file: should error.

Related configuration
- backend/config.go → AllowedPaths

Source
- backend/tools/file_read.go
