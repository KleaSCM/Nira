package tools

import (
    "fmt"
    "io/fs"
    "path/filepath"
    "strings"
    "time"
)

// SearchFilesByNameTool searches for files (and optionally directories) by name pattern under a root.
type SearchFilesByNameTool struct {
    AllowedPaths []string
    checker      PathChecker
}

func NewSearchFilesByNameTool(allowedPaths []string) *SearchFilesByNameTool { return &SearchFilesByNameTool{AllowedPaths: allowedPaths} }
func NewSearchFilesByNameToolWithChecker(allowedPaths []string, checker PathChecker) *SearchFilesByNameTool {
    return &SearchFilesByNameTool{AllowedPaths: allowedPaths, checker: checker}
}

func (t *SearchFilesByNameTool) Name() string { return "search_files_by_name" }

func (t *SearchFilesByNameTool) Description() string {
    return "Search for files by name under a root directory. Args: root (string), pattern (string, substring or glob), max_results (int), include_dirs (bool), case_sensitive (bool)."
}

func (t *SearchFilesByNameTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "root": map[string]interface{}{
                    "type":        "string",
                    "description": "Root directory to search within",
                },
                "pattern": map[string]interface{}{
                    "type":        "string",
                    "description": "Substring or glob pattern to match against names",
                },
                "max_results": map[string]interface{}{
                    "type":        "integer",
                    "description": "Maximum number of results (default 200)",
                },
                "include_dirs": map[string]interface{}{
                    "type":        "boolean",
                    "description": "Include directories in results (default false)",
                },
                "case_sensitive": map[string]interface{}{
                    "type":        "boolean",
                    "description": "Case-sensitive match (default false)",
                },
            },
            "required": []string{"root", "pattern"},
        },
    }
}

func (t *SearchFilesByNameTool) Execute(args map[string]interface{}) (interface{}, error) {
    root, ok := args["root"].(string)
    if !ok {
        return nil, fmt.Errorf("root argument is required and must be a string")
    }
    pattern, ok := args["pattern"].(string)
    if !ok {
        return nil, fmt.Errorf("pattern argument is required and must be a string")
    }
    if !t.isPathAllowed(root) {
        return nil, fmt.Errorf("root '%s' is not in allowed directories", root)
    }

    maxResults := 200
    if v, ok := args["max_results"]; ok {
        switch n := v.(type) {
        case float64:
            maxResults = int(n)
        case int:
            maxResults = n
        }
    }
    includeDirs := false
    if v, ok := args["include_dirs"].(bool); ok {
        includeDirs = v
    }
    caseSensitive := false
    if v, ok := args["case_sensitive"].(bool); ok {
        caseSensitive = v
    }

    useGlob := strings.ContainsAny(pattern, "*?")
    pat := pattern
    if !caseSensitive && !useGlob {
        pat = strings.ToLower(pattern)
    }

    results := make([]map[string]interface{}, 0, 32)
    count := 0

    match := func(name string) bool {
        if useGlob {
            ok, _ := filepath.Match(pattern, name)
            return ok
        }
        if caseSensitive {
            return strings.Contains(name, pattern)
        }
        return strings.Contains(strings.ToLower(name), pat)
    }

    filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil // skip unreadable entries
        }
        // Validate we remain inside allowed paths even when following nested entries
        if !t.isPathAllowed(p) {
            if d.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }
        if d.IsDir() && !includeDirs {
            // still traverse
        }
        name := d.Name()
        if match(name) {
            if !d.IsDir() || includeDirs {
                info, _ := d.Info()
                var size int64
                if info != nil && !info.IsDir() {
                    size = info.Size()
                }
                var mod string
                if info != nil {
                    mod = info.ModTime().UTC().Format(time.RFC3339)
                }
                results = append(results, map[string]interface{}{
                    "name":    name,
                    "path":    p,
                    "is_dir":  d.IsDir(),
                    "size":    size,
                    "mod_time": mod,
                })
                count++
                if count >= maxResults {
                    return filepath.SkipDir
                }
            }
        }
        return nil
    })

    return results, nil
}

func (t *SearchFilesByNameTool) isPathAllowed(path string) bool {
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
