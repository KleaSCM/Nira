package tools

import (
    "fmt"
    "nira/memory"
    "time"
)

// ---- Character tools ----

type RPCharacterListTool struct{ store *memory.RPStore }
func NewRPCharacterListTool(store *memory.RPStore) *RPCharacterListTool { return &RPCharacterListTool{store: store} }
func (t *RPCharacterListTool) Name() string        { return "rp_character_list" }
func (t *RPCharacterListTool) Description() string { return "Lists RP characters. Args: query (string, optional), limit (int), offset (int)." }
func (t *RPCharacterListTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]interface{}{"type": "string"},
                "limit": map[string]interface{}{"type": "integer"},
                "offset": map[string]interface{}{"type": "integer"},
            },
        },
    }
}
func (t *RPCharacterListTool) Execute(args map[string]interface{}) (interface{}, error) {
    q, _ := args["query"].(string)
    limit := intFrom(args["limit"], 100)
    offset := intFrom(args["offset"], 0)
    return t.store.ListCharacters(q, limit, offset)
}

type RPCharacterGetTool struct{ store *memory.RPStore }
func NewRPCharacterGetTool(store *memory.RPStore) *RPCharacterGetTool { return &RPCharacterGetTool{store: store} }
func (t *RPCharacterGetTool) Name() string        { return "rp_character_get" }
func (t *RPCharacterGetTool) Description() string { return "Gets a character by id. Args: id (string)." }
func (t *RPCharacterGetTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "id": map[string]interface{}{"type": "string"},
            },
            "required": []string{"id"},
        },
    }
}
func (t *RPCharacterGetTool) Execute(args map[string]interface{}) (interface{}, error) {
    id, ok := args["id"].(string)
    if !ok || id == "" { return nil, fmt.Errorf("id is required") }
    return t.store.GetCharacter(id)
}

type RPCharacterSaveTool struct{ store *memory.RPStore }
func NewRPCharacterSaveTool(store *memory.RPStore) *RPCharacterSaveTool { return &RPCharacterSaveTool{store: store} }
func (t *RPCharacterSaveTool) Name() string        { return "rp_character_save" }
func (t *RPCharacterSaveTool) Description() string { return "Creates or updates a character. Args: id (string, optional), name (string), summary (string), traits ([string]), background (string), goals ([string]), tags ([string]), notes (string)." }
func (t *RPCharacterSaveTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "id": map[string]interface{}{"type": "string"},
                "name": map[string]interface{}{"type": "string"},
                "summary": map[string]interface{}{"type": "string"},
                "traits": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
                "background": map[string]interface{}{"type": "string"},
                "goals": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
                "tags": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
                "notes": map[string]interface{}{"type": "string"},
            },
            "required": []string{"name"},
        },
    }
}
func (t *RPCharacterSaveTool) Execute(args map[string]interface{}) (interface{}, error) {
    name, ok := args["name"].(string)
    if !ok || name == "" { return nil, fmt.Errorf("name is required") }
    c := &memory.RPCharacter{
        ID:        stringFrom(args["id"], genID()),
        Name:      name,
        Summary:   stringFrom(args["summary"], ""),
        Background: stringFrom(args["background"], ""),
        Notes:     stringFrom(args["notes"], ""),
    }
    c.Traits = stringSlice(args["traits"]) 
    c.Goals = stringSlice(args["goals"]) 
    c.Tags = stringSlice(args["tags"]) 
    if err := t.store.SaveCharacter(c); err != nil { return nil, err }
    return c, nil
}

type RPCharacterDeleteTool struct{ store *memory.RPStore }
func NewRPCharacterDeleteTool(store *memory.RPStore) *RPCharacterDeleteTool { return &RPCharacterDeleteTool{store: store} }
func (t *RPCharacterDeleteTool) Name() string        { return "rp_character_delete" }
func (t *RPCharacterDeleteTool) Description() string { return "Deletes a character. Args: id (string)." }
func (t *RPCharacterDeleteTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "id": map[string]interface{}{"type": "string"},
            },
            "required": []string{"id"},
        },
    }
}
func (t *RPCharacterDeleteTool) Execute(args map[string]interface{}) (interface{}, error) {
    id, ok := args["id"].(string)
    if !ok || id == "" { return nil, fmt.Errorf("id is required") }
    if err := t.store.DeleteCharacter(id); err != nil { return nil, err }
    return map[string]interface{}{"deleted": id}, nil
}

// ---- Story card tools ----

type RPStoryCardListTool struct{ store *memory.RPStore }
func NewRPStoryCardListTool(store *memory.RPStore) *RPStoryCardListTool { return &RPStoryCardListTool{store: store} }
func (t *RPStoryCardListTool) Name() string        { return "rp_storycard_list" }
func (t *RPStoryCardListTool) Description() string { return "Lists story cards. Args: query (string, optional), kind (string, optional), limit (int), offset (int)." }
func (t *RPStoryCardListTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]interface{}{"type": "string"},
                "kind": map[string]interface{}{"type": "string"},
                "limit": map[string]interface{}{"type": "integer"},
                "offset": map[string]interface{}{"type": "integer"},
            },
        },
    }
}
func (t *RPStoryCardListTool) Execute(args map[string]interface{}) (interface{}, error) {
    q, _ := args["query"].(string)
    kind, _ := args["kind"].(string)
    limit := intFrom(args["limit"], 100)
    offset := intFrom(args["offset"], 0)
    return t.store.ListStoryCards(q, kind, limit, offset)
}

type RPStoryCardGetTool struct{ store *memory.RPStore }
func NewRPStoryCardGetTool(store *memory.RPStore) *RPStoryCardGetTool { return &RPStoryCardGetTool{store: store} }
func (t *RPStoryCardGetTool) Name() string        { return "rp_storycard_get" }
func (t *RPStoryCardGetTool) Description() string { return "Gets a story card by id. Args: id (string)." }
func (t *RPStoryCardGetTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "id": map[string]interface{}{"type": "string"},
            },
            "required": []string{"id"},
        },
    }
}
func (t *RPStoryCardGetTool) Execute(args map[string]interface{}) (interface{}, error) {
    id, ok := args["id"].(string)
    if !ok || id == "" { return nil, fmt.Errorf("id is required") }
    return t.store.GetStoryCard(id)
}

type RPStoryCardSaveTool struct{ store *memory.RPStore }
func NewRPStoryCardSaveTool(store *memory.RPStore) *RPStoryCardSaveTool { return &RPStoryCardSaveTool{store: store} }
func (t *RPStoryCardSaveTool) Name() string        { return "rp_storycard_save" }
func (t *RPStoryCardSaveTool) Description() string { return "Creates or updates a story card. Args: id (string, optional), title (string), kind (string), content (string), tags ([string]), links ([string])." }
func (t *RPStoryCardSaveTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "id": map[string]interface{}{"type": "string"},
                "title": map[string]interface{}{"type": "string"},
                "kind": map[string]interface{}{"type": "string"},
                "content": map[string]interface{}{"type": "string"},
                "tags": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
                "links": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
            },
            "required": []string{"title", "kind"},
        },
    }
}
func (t *RPStoryCardSaveTool) Execute(args map[string]interface{}) (interface{}, error) {
    title, ok := args["title"].(string)
    if !ok || title == "" { return nil, fmt.Errorf("title is required") }
    kind, ok := args["kind"].(string)
    if !ok || kind == "" { return nil, fmt.Errorf("kind is required") }
    sc := &memory.RPStoryCard{
        ID:      stringFrom(args["id"], genID()),
        Title:   title,
        Kind:    kind,
        Content: stringFrom(args["content"], ""),
    }
    sc.Tags = stringSlice(args["tags"]) 
    sc.Links = stringSlice(args["links"]) 
    if err := t.store.SaveStoryCard(sc); err != nil { return nil, err }
    return sc, nil
}

type RPStoryCardDeleteTool struct{ store *memory.RPStore }
func NewRPStoryCardDeleteTool(store *memory.RPStore) *RPStoryCardDeleteTool { return &RPStoryCardDeleteTool{store: store} }
func (t *RPStoryCardDeleteTool) Name() string        { return "rp_storycard_delete" }
func (t *RPStoryCardDeleteTool) Description() string { return "Deletes a story card. Args: id (string)." }
func (t *RPStoryCardDeleteTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": t.Name(),
        "description": t.Description(),
        "parameters": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "id": map[string]interface{}{"type": "string"},
            },
            "required": []string{"id"},
        },
    }
}
func (t *RPStoryCardDeleteTool) Execute(args map[string]interface{}) (interface{}, error) {
    id, ok := args["id"].(string)
    if !ok || id == "" { return nil, fmt.Errorf("id is required") }
    if err := t.store.DeleteStoryCard(id); err != nil { return nil, err }
    return map[string]interface{}{"deleted": id}, nil
}

// ---- helpers ----

func intFrom(v interface{}, def int) int {
    switch n := v.(type) {
    case float64:
        return int(n)
    case int:
        return n
    default:
        return def
    }
}

func stringFrom(v interface{}, def string) string {
    if s, ok := v.(string); ok { return s }
    return def
}

func stringSlice(v interface{}) []string {
    out := []string{}
    if arr, ok := v.([]interface{}); ok {
        for _, it := range arr {
            if s, ok := it.(string); ok { out = append(out, s) }
        }
    }
    return out
}

func genID() string {
    // Simple unique id (timestamp-based). Replace with UUID if desired.
    return fmt.Sprintf("rp_%d", time.Now().UTC().UnixNano())
}
