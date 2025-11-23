package tools

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

// FileMetadataTool returns basic metadata for a given path.
type FileMetadataTool struct {
    AllowedPaths []string
}

func NewFileMetadataTool(allowedPaths []string) *FileMetadataTool {
    return &FileMetadataTool{AllowedPaths: allowedPaths}
}

func (t *FileMetadataTool) Name() string { return "file_metadata" }

func (t *FileMetadataTool) Description() string {
    return "Returns basic metadata for a file or directory. Args: path (string)."
}

func (t *FileMetadataTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "path": map[string]interface{}{
                    "type":        "string",
                    "description": "Path to file or directory",
                },
            },
            "required": []string{"path"},
        },
    }
}

func (t *FileMetadataTool) Execute(args map[string]interface{}) (interface{}, error) {
    path, ok := args["path"].(string)
    if !ok {
        return nil, fmt.Errorf("path argument is required and must be a string")
    }
    if !t.isPathAllowed(path) {
        return nil, fmt.Errorf("path '%s' is not in allowed directories", path)
    }
    info, err := os.Stat(path)
    if err != nil {
        return nil, fmt.Errorf("failed to stat path: %w", err)
    }
    abs, _ := filepath.Abs(path)
    mod := info.ModTime().UTC().Format(time.RFC3339)
    res := map[string]interface{}{
        "name":     info.Name(),
        "path":     path,
        "abs_path": abs,
        "is_dir":   info.IsDir(),
        "size":     func() int64 { if info.IsDir() { return 0 }; return info.Size() }(),
        "mod_time": mod,
    }
    return res, nil
}

func (t *FileMetadataTool) isPathAllowed(path string) bool {
    if len(t.AllowedPaths) == 0 { return false }
    absPath, err := filepath.Abs(path)
    if err != nil { return false }
    for _, allowed := range t.AllowedPaths {
        absAllowed, err := filepath.Abs(allowed)
        if err != nil { continue }
        rel, err := filepath.Rel(absAllowed, absPath)
        if err != nil { continue }
        if rel != ".." && !filepath.IsAbs(rel) {
            return true
        }
    }
    return false
}
