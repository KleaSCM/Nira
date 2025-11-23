Tool: file_metadata

Overview
- Returns basic metadata about a given file or directory.
- Enforced by AllowedPaths sandbox (backend/config.go).
- Implemented in backend/tools/file_metadata.go.

Identifier
- name: file_metadata

Arguments
- path (string, required): Path to a file or directory within AllowedPaths.

Returns
- Success: JSON object with
  - name: string
  - path: string (as provided)
  - abs_path: string (absolute path)
  - is_dir: boolean
  - size: integer bytes (0 for directories)
  - mod_time: string (RFC3339 UTC)
- Failure: error propagated to WebSocket as a message of type "error".

Security and sandboxing
- The path must be under one of AllowedPaths (default: ["."]).

Usage examples
- Frontend direct call:
  { "name": "file_metadata", "arguments": { "path": "./backend/server.go" } }
- Model-initiated call:
  {"name":"file_metadata","arguments":{"path":"./Docs/Tools"}}

Common errors
- path argument missing/invalid
- failed to stat path: <system error>
- path not in allowed directories

Testing checklist
- File path → returns correct size and mod_time
- Directory path → is_dir=true, size=0
- Outside AllowedPaths → error

Source
- backend/tools/file_metadata.go
