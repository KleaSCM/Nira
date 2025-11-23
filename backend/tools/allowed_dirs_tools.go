package tools

import (
    "fmt"
    "os"
)

// AllowedDirsProvider provides access to allowed directories list management.
type AllowedDirsProvider interface {
    List() []string
    Add(path string) error
    Remove(path string) error
}

// allowed_dirs_list
type AllowedDirsListTool struct{ store AllowedDirsProvider }

func NewAllowedDirsListTool(store AllowedDirsProvider) *AllowedDirsListTool { return &AllowedDirsListTool{store: store} }
func (t *AllowedDirsListTool) Name() string        { return "allowed_dirs_list" }
func (t *AllowedDirsListTool) Description() string { return "Lists directories NIRA is allowed to access. Args: none." }
func (t *AllowedDirsListTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type":       "object",
            "properties": map[string]interface{}{},
        },
    }
}
func (t *AllowedDirsListTool) Execute(args map[string]interface{}) (interface{}, error) {
    return map[string]interface{}{"allowed": t.store.List()}, nil
}

// allowed_dirs_add
type AllowedDirsAddTool struct{ store AllowedDirsProvider }

func NewAllowedDirsAddTool(store AllowedDirsProvider) *AllowedDirsAddTool { return &AllowedDirsAddTool{store: store} }
func (t *AllowedDirsAddTool) Name() string        { return "allowed_dirs_add" }
func (t *AllowedDirsAddTool) Description() string { return "Adds a directory to the allowed list. Args: path (string)." }
func (t *AllowedDirsAddTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "path": map[string]interface{}{"type": "string", "description": "Directory path to allow"},
            },
            "required": []string{"path"},
        },
    }
}
func (t *AllowedDirsAddTool) Execute(args map[string]interface{}) (interface{}, error) {
    p, ok := args["path"].(string)
    if !ok || p == "" {
        return nil, fmt.Errorf("path argument is required and must be a string")
    }
    // Simple existence check
    info, err := os.Stat(p)
    if err != nil || !info.IsDir() {
        return nil, fmt.Errorf("path must be an existing directory")
    }
    if err := t.store.Add(p); err != nil {
        return nil, err
    }
    return map[string]interface{}{"allowed": t.store.List()}, nil
}

// allowed_dirs_remove
type AllowedDirsRemoveTool struct{ store AllowedDirsProvider }

func NewAllowedDirsRemoveTool(store AllowedDirsProvider) *AllowedDirsRemoveTool { return &AllowedDirsRemoveTool{store: store} }
func (t *AllowedDirsRemoveTool) Name() string        { return "allowed_dirs_remove" }
func (t *AllowedDirsRemoveTool) Description() string { return "Removes a directory from the allowed list. Args: path (string)." }
func (t *AllowedDirsRemoveTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "path": map[string]interface{}{"type": "string", "description": "Directory path to remove"},
            },
            "required": []string{"path"},
        },
    }
}
func (t *AllowedDirsRemoveTool) Execute(args map[string]interface{}) (interface{}, error) {
    p, ok := args["path"].(string)
    if !ok || p == "" {
        return nil, fmt.Errorf("path argument is required and must be a string")
    }
    if err := t.store.Remove(p); err != nil {
        return nil, err
    }
    return map[string]interface{}{"allowed": t.store.List()}, nil
}
