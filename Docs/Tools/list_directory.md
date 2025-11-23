Tool: list_directory

Overview
- Lists files and folders under a specified directory.
- Optional recursion and filters for files/dirs.
- Enforced by AllowedPaths sandbox (backend/config.go).
- Implemented in backend/tools/list_directory.go.

Identifier
- name: list_directory

Arguments
- path (string, required): Directory to list.
- recursive (boolean, optional, default=false): Recurse into subdirectories.
- include_files (boolean, optional, default=true): Include files in results.
- include_dirs (boolean, optional, default=true): Include directories in results.
- max_items (integer, optional, default=1000): Limit number of returned entries.

Returns
- Success: array of entries, each with
  - name: string
  - path: string
  - is_dir: boolean
  - size: integer (0 for directories)
  - mod_time: string (RFC3339 UTC)
- Failure: error propagated to WebSocket as a message of type "error".

Security and sandboxing
- The path must be under one of AllowedPaths (default: ["."]).
- Absolute resolution + relative checks mitigate traversal.

Usage examples
- Frontend direct call:
  { "name": "list_directory", "arguments": { "path": "./Docs", "recursive": false } }
- Model-initiated call:
  {"name":"list_directory","arguments":{"path":"./frontend/lib","include_files":true,"include_dirs":false}}

Common errors
- path argument missing/invalid
- failed to read directory: <system error>
- path not in allowed directories

Testing checklist
- List project root (non-recursive) → returns files/dirs
- Recursive list with max_items small → truncates appropriately
- Path outside AllowedPaths → error

Source
- backend/tools/list_directory.go
