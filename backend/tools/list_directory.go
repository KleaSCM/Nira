package tools

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

// ListDirectoryTool lists directory entries with optional recursion and filters.
type ListDirectoryTool struct {
    AllowedPaths []string
    checker      PathChecker
}

func NewListDirectoryTool(allowedPaths []string) *ListDirectoryTool { return &ListDirectoryTool{AllowedPaths: allowedPaths} }
func NewListDirectoryToolWithChecker(allowedPaths []string, checker PathChecker) *ListDirectoryTool {
    return &ListDirectoryTool{AllowedPaths: allowedPaths, checker: checker}
}

func (t *ListDirectoryTool) Name() string { return "list_directory" }

func (t *ListDirectoryTool) Description() string {
    return "Lists files and folders under a directory. Args: path (string), recursive (bool, optional), include_files (bool), include_dirs (bool), max_items (int)."
}

func (t *ListDirectoryTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "path": map[string]interface{}{
                    "type":        "string",
                    "description": "Directory path to list",
                },
                "recursive": map[string]interface{}{
                    "type":        "boolean",
                    "description": "Recursively list contents",
                },
                "include_files": map[string]interface{}{
                    "type":        "boolean",
                    "description": "Include files in results (default true)",
                },
                "include_dirs": map[string]interface{}{
                    "type":        "boolean",
                    "description": "Include directories in results (default true)",
                },
                "max_items": map[string]interface{}{
                    "type":        "integer",
                    "description": "Maximum number of items to return (default 1000)",
                },
            },
            "required": []string{"path"},
        },
    }
}

func (t *ListDirectoryTool) Execute(args map[string]interface{}) (interface{}, error) {
    path, ok := args["path"].(string)
    if !ok {
        return nil, fmt.Errorf("path argument is required and must be a string")
    }
    if !t.isPathAllowed(path) {
        return nil, fmt.Errorf("path '%s' is not in allowed directories", path)
    }

    recursive := false
    if v, ok := args["recursive"].(bool); ok {
        recursive = v
    }
    includeFiles := true
    if v, ok := args["include_files"].(bool); ok {
        includeFiles = v
    }
    includeDirs := true
    if v, ok := args["include_dirs"].(bool); ok {
        includeDirs = v
    }
    maxItems := 1000
    if v, ok := args["max_items"]; ok {
        switch n := v.(type) {
        case float64:
            maxItems = int(n)
        case int:
            maxItems = n
        }
    }

    results := make([]map[string]interface{}, 0, 64)
    count := 0

    push := func(p string, info os.FileInfo) {
        if count >= maxItems {
            return
        }
        isDir := info.IsDir()
        if (isDir && !includeDirs) || (!isDir && !includeFiles) {
            return
        }
        var size int64
        if !isDir {
            size = info.Size()
        }
        mod := info.ModTime().UTC().Format(time.RFC3339)
        results = append(results, map[string]interface{}{
            "name":    info.Name(),
            "path":    p,
            "is_dir":  isDir,
            "size":    size,
            "mod_time": mod,
        })
        count++
    }

    if recursive {
        err := filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if p == path {
                // skip the root directory itself in results
                return nil
            }
            if count >= maxItems {
                return filepath.SkipDir
            }
            push(p, fi)
            return nil
        })
        if err != nil && err != filepath.SkipDir {
            return nil, fmt.Errorf("failed to list directory: %w", err)
        }
    } else {
        entries, err := os.ReadDir(path)
        if err != nil {
            return nil, fmt.Errorf("failed to read directory: %w", err)
        }
        for _, e := range entries {
            if count >= maxItems {
                break
            }
            info, err := e.Info()
            if err != nil {
                continue
            }
            push(filepath.Join(path, e.Name()), info)
        }
    }

    return results, nil
}

func (t *ListDirectoryTool) isPathAllowed(path string) bool {
    if t.checker != nil { return t.checker.IsAllowed(path) }
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
