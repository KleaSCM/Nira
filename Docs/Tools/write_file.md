Tool: write_file

Overview
- Writes text content to a file path. Parent directories are created as needed.
- Enforces a filesystem sandbox using AllowedPaths from backend/config.go.
- Implemented in backend/tools/file_write.go.

Identifier
- name: write_file

Arguments
- path (string, required): Absolute or relative path to the target file.
- content (string, required): The text to write. Existing files are overwritten.

Returns
- Success: string confirmation, e.g., "Successfully wrote to C:\\path\\to\\file.txt".
- Failure: error propagated to WebSocket as a message of type "error".

Security and sandboxing
- Only paths under AllowedPaths are accepted; others are rejected.
- Default AllowedPaths is ["."] (project root). Add absolute directories in backend/config.go to expand the sandbox.
- Absolute path resolution and relative checks mitigate directory traversal.

Behavior notes
- Overwrite semantics: os.WriteFile() replaces the file’s contents. There is no append mode yet.
- Creates parent directories with os.MkdirAll(dir, 0755) when needed.
- Intended for text content. Writing binary via this tool is not supported.

Frontend usage (direct call)
- The UI prompts after file selection: choose Write, enter content, then sends
  { "name": "write_file", "arguments": { "path": "<file>", "content": "<text>" } }.
- The backend returns a success message which streams to the chat.

AI-initiated usage
- The model can reply with a tool call like:
  {"name":"write_file","arguments":{"path":"./notes/todo.txt","content":"Buy milk"}}
- The server executes and injects a contextual message with the result.

Common errors
- path/content argument missing or wrong type.
- path '<p>' is not in allowed directories.
- failed to create directories: <system error>.
- failed to write file: <system error>.

Testing checklist
- Write to a new file under project root: should create directories and succeed.
- Overwrite an existing file: should succeed and replace content.
- Attempt outside AllowedPaths: should error.

Related configuration
- backend/config.go → AllowedPaths

Source
- backend/tools/file_write.go
