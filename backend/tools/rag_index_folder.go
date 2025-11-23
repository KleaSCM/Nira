package tools

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// RagIndexWriter defines minimal upsert API backed by memory.RagIndex
type RagIndexWriter interface {
    Upsert(path, name, modTime string, size int64, content string) error
}

// RagIndexFolderTool indexes text files under an allowed directory into the basic rag_index table.
type RagIndexFolderTool struct {
    checker PathChecker
    index   RagIndexWriter
}

func NewRagIndexFolderTool(checker PathChecker, index RagIndexWriter) *RagIndexFolderTool {
    return &RagIndexFolderTool{checker: checker, index: index}
}

func (t *RagIndexFolderTool) Name() string        { return "rag_index_folder" }
func (t *RagIndexFolderTool) Description() string { return "Indexes text files in a folder. Args: root (string), patterns ([string], optional), max_size_mb (int), max_files (int)." }
func (t *RagIndexFolderTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "root": map[string]interface{}{"type": "string", "description": "Root directory to index"},
                "patterns": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Glob patterns to include (e.g., *.md)"},
                "max_size_mb": map[string]interface{}{"type": "integer", "description": "Max file size in MB (default 2)"},
                "max_files": map[string]interface{}{"type": "integer", "description": "Max files to index (default 500)"},
            },
            "required": []string{"root"},
        },
    }
}

func (t *RagIndexFolderTool) Execute(args map[string]interface{}) (interface{}, error) {
    root, ok := args["root"].(string)
    if !ok || root == "" { return nil, fmt.Errorf("root argument is required and must be a string") }
    if t.checker == nil || !t.checker.IsAllowed(root) {
        return nil, fmt.Errorf("root '%s' is not in allowed directories", root)
    }

    // Extract patterns
    var patterns []string
    if v, ok := args["patterns"].([]interface{}); ok {
        for _, it := range v {
            if s, ok := it.(string); ok && s != "" { patterns = append(patterns, s) }
        }
    }
    if len(patterns) == 0 {
        patterns = []string{"*.md", "*.txt", "*.json", "*.yaml", "*.yml"}
    }
    maxSizeMB := 2
    if v, ok := args["max_size_mb"]; ok {
        switch n := v.(type) { case float64: maxSizeMB = int(n); case int: maxSizeMB = n }
    }
    maxFiles := 500
    if v, ok := args["max_files"]; ok {
        switch n := v.(type) { case float64: maxFiles = int(n); case int: maxFiles = n }
    }
    indexed := 0
    var lastErr error

    err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
        if err != nil { return nil }
        if d.IsDir() { return nil }
        if !t.checker.IsAllowed(p) { return nil }
        name := d.Name()
        // pattern match
        matched := false
        for _, pat := range patterns {
            ok, _ := filepath.Match(pat, name)
            if ok { matched = true; break }
        }
        if !matched { return nil }
        info, ierr := d.Info()
        if ierr != nil { return nil }
        if info.Size() > int64(maxSizeMB)*1024*1024 { return nil }
        // Read file as text (best effort)
        b, rerr := os.Open(p)
        if rerr != nil { return nil }
        defer b.Close()
        data, rr := io.ReadAll(b)
        if rr != nil { return nil }
        content := string(data)
        mod := info.ModTime().UTC().Format(time.RFC3339)
        if err := t.index.Upsert(p, name, mod, info.Size(), content); err != nil { lastErr = err }
        indexed++
        if indexed >= maxFiles { return filepath.SkipDir }
        return nil
    })
    if err != nil && err != filepath.SkipDir { return nil, fmt.Errorf("indexing failed: %w", err) }
    msg := fmt.Sprintf("Indexed %d files under %s (patterns: %s)", indexed, root, strings.Join(patterns, ", "))
    if lastErr != nil { msg += "; last error: " + lastErr.Error() }
    return msg, nil
}
