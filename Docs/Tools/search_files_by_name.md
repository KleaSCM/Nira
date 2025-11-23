Tool: search_files_by_name

Overview
- Searches for files (and optionally directories) by name under a root directory.
- Supports simple substring matching or glob patterns (*, ?).
- Enforced by AllowedPaths sandbox (backend/config.go).
- Implemented in backend/tools/search_files_by_name.go.

Identifier
- name: search_files_by_name

Arguments
- root (string, required): Root directory to search within.
- pattern (string, required): Substring or glob pattern to match (e.g., "notes", "*.md").
- max_results (integer, optional, default=200): Maximum number of matches to return.
- include_dirs (boolean, optional, default=false): Include matching directories in results.
- case_sensitive (boolean, optional, default=false): If true, matching is case-sensitive for substring mode.

Returns
- Success: array of entries, each with
  - name: string
  - path: string
  - is_dir: boolean
  - size: integer (0 for directories)
  - mod_time: string (RFC3339 UTC)
- Failure: error propagated to WebSocket as a message of type "error".

Security and sandboxing
- The root and traversed paths must remain under AllowedPaths. Paths outside are skipped.

Usage examples
- Frontend direct call:
  { "name": "search_files_by_name", "arguments": { "root": ".", "pattern": "*.md", "max_results": 100 } }
- Model-initiated call:
  {"name":"search_files_by_name","arguments":{"root":"./Docs","pattern":"RAG","case_sensitive":false}}

Common errors
- root/pattern argument missing/invalid
- root not in allowed directories

Testing checklist
- Search with glob (*.md) under project root → returns markdown files
- Search with substring ("server") including dirs=false → ignores directories named server
- Set include_dirs=true → directories that match appear
- Set a small max_results → results truncated

Source
- backend/tools/search_files_by_name.go
