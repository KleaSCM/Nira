package tools

import "fmt"

// RagSearcher defines the search API provided by memory.RagIndex
type RagSearcher interface {
    Search(query string, limit int, pathPrefix string) ([]map[string]interface{}, error)
}

type RagSearchTool struct {
    search RagSearcher
    checker PathChecker
}

func NewRagSearchTool(search RagSearcher, checker PathChecker) *RagSearchTool {
    return &RagSearchTool{search: search, checker: checker}
}

func (t *RagSearchTool) Name() string { return "rag_search" }
func (t *RagSearchTool) Description() string {
    return "Searches the lightweight index for files/snippets. Args: query (string), limit (int, optional), path_prefix (string, optional)."
}
func (t *RagSearchTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]interface{}{"type": "string", "description": "Query text"},
                "limit": map[string]interface{}{"type": "integer", "description": "Max results (default 10)"},
                "path_prefix": map[string]interface{}{"type": "string", "description": "Restrict to paths under this prefix"},
            },
            "required": []string{"query"},
        },
    }
}

func (t *RagSearchTool) Execute(args map[string]interface{}) (interface{}, error) {
    q, ok := args["query"].(string)
    if !ok || q == "" {
        return nil, fmt.Errorf("query is required")
    }
    limit := 10
    if v, ok := args["limit"]; ok {
        switch n := v.(type) { case float64: limit = int(n); case int: limit = n }
    }
    pathPrefix, _ := args["path_prefix"].(string)
    if pathPrefix != "" && t.checker != nil && !t.checker.IsAllowed(pathPrefix) {
        return nil, fmt.Errorf("path_prefix '%s' is not in allowed directories", pathPrefix)
    }
    return t.search.Search(q, limit, pathPrefix)
}
