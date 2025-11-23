package memory

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

// RPStore provides CRUD helpers for RP entities (characters, story cards).
type RPStore struct { db *Database }

func NewRPStore(db *Database) *RPStore { return &RPStore{db: db} }

// ---- Characters ----
type RPCharacter struct {
    ID        string   `json:"id"`
    Name      string   `json:"name"`
    Summary   string   `json:"summary"`
    Traits    []string `json:"traits"`
    Background string  `json:"background"`
    Goals     []string `json:"goals"`
    Tags      []string `json:"tags"`
    Notes     string   `json:"notes"`
    CreatedAt string   `json:"created_at"`
    UpdatedAt string   `json:"updated_at"`
}

func (s *RPStore) ListCharacters(query string, limit, offset int) ([]RPCharacter, error) {
    if limit <= 0 { limit = 100 }
    if offset < 0 { offset = 0 }
    where := ""
    args := []any{}
    if strings.TrimSpace(query) != "" {
        where = "WHERE name LIKE ? OR summary LIKE ?"
        like := "%" + strings.ReplaceAll(query, "%", "") + "%"
        args = append(args, like, like)
    }
    rows, err := s.db.DB.Query(
        fmt.Sprintf("SELECT id,name,summary,traits_json,background,goals_json,tags_json,notes,created_at,updated_at FROM rp_characters %s ORDER BY updated_at DESC LIMIT ? OFFSET ?", where),
        append(args, limit, offset)...,
    )
    if err != nil { return nil, err }
    defer rows.Close()
    out := []RPCharacter{}
    for rows.Next() {
        var c RPCharacter
        var traitsJS, goalsJS, tagsJS string
        if err := rows.Scan(&c.ID, &c.Name, &c.Summary, &traitsJS, &c.Background, &goalsJS, &tagsJS, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
            continue
        }
        _ = json.Unmarshal([]byte(emptyJSON(traitsJS, "[]")), &c.Traits)
        _ = json.Unmarshal([]byte(emptyJSON(goalsJS, "[]")), &c.Goals)
        _ = json.Unmarshal([]byte(emptyJSON(tagsJS, "[]")), &c.Tags)
        out = append(out, c)
    }
    return out, rows.Err()
}

func (s *RPStore) GetCharacter(id string) (*RPCharacter, error) {
    row := s.db.DB.QueryRow("SELECT id,name,summary,traits_json,background,goals_json,tags_json,notes,created_at,updated_at FROM rp_characters WHERE id=?", id)
    var c RPCharacter
    var traitsJS, goalsJS, tagsJS string
    if err := row.Scan(&c.ID, &c.Name, &c.Summary, &traitsJS, &c.Background, &goalsJS, &tagsJS, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
        if err == sql.ErrNoRows { return nil, nil }
        return nil, err
    }
    _ = json.Unmarshal([]byte(emptyJSON(traitsJS, "[]")), &c.Traits)
    _ = json.Unmarshal([]byte(emptyJSON(goalsJS, "[]")), &c.Goals)
    _ = json.Unmarshal([]byte(emptyJSON(tagsJS, "[]")), &c.Tags)
    return &c, nil
}

func (s *RPStore) SaveCharacter(c *RPCharacter) error {
    now := time.Now().UTC().Format(time.RFC3339)
    if c.CreatedAt == "" { c.CreatedAt = now }
    c.UpdatedAt = now
    traitsJS, _ := json.Marshal(c.Traits)
    goalsJS, _ := json.Marshal(c.Goals)
    tagsJS, _ := json.Marshal(c.Tags)
    _, err := s.db.DB.Exec(`
        INSERT INTO rp_characters(id,name,summary,traits_json,background,goals_json,tags_json,notes,created_at,updated_at)
        VALUES(?,?,?,?,?,?,?,?,?,?)
        ON CONFLICT(id) DO UPDATE SET
            name=excluded.name,
            summary=excluded.summary,
            traits_json=excluded.traits_json,
            background=excluded.background,
            goals_json=excluded.goals_json,
            tags_json=excluded.tags_json,
            notes=excluded.notes,
            updated_at=excluded.updated_at
    `, c.ID, c.Name, c.Summary, string(traitsJS), c.Background, string(goalsJS), string(tagsJS), c.Notes, c.CreatedAt, c.UpdatedAt)
    return err
}

func (s *RPStore) DeleteCharacter(id string) error {
    _, err := s.db.DB.Exec("DELETE FROM rp_characters WHERE id=?", id)
    return err
}

// ---- Story Cards ----
type RPStoryCard struct {
    ID        string   `json:"id"`
    Title     string   `json:"title"`
    Kind      string   `json:"kind"`
    Content   string   `json:"content"`
    Tags      []string `json:"tags"`
    Links     []string `json:"links"`
    CreatedAt string   `json:"created_at"`
    UpdatedAt string   `json:"updated_at"`
}

func (s *RPStore) ListStoryCards(query, kind string, limit, offset int) ([]RPStoryCard, error) {
    if limit <= 0 { limit = 100 }
    if offset < 0 { offset = 0 }
    where := []string{}
    args := []any{}
    if strings.TrimSpace(query) != "" {
        where = append(where, "(title LIKE ? OR content LIKE ?)")
        like := "%" + strings.ReplaceAll(query, "%", "") + "%"
        args = append(args, like, like)
    }
    if strings.TrimSpace(kind) != "" { where = append(where, "kind=?"); args = append(args, kind) }
    whereSQL := ""
    if len(where) > 0 { whereSQL = "WHERE " + strings.Join(where, " AND ") }
    rows, err := s.db.DB.Query(
        fmt.Sprintf("SELECT id,title,kind,content,tags_json,links_json,created_at,updated_at FROM rp_story_cards %s ORDER BY updated_at DESC LIMIT ? OFFSET ?", whereSQL),
        append(args, limit, offset)...,
    )
    if err != nil { return nil, err }
    defer rows.Close()
    out := []RPStoryCard{}
    for rows.Next() {
        var sc RPStoryCard
        var tagsJS, linksJS string
        if err := rows.Scan(&sc.ID, &sc.Title, &sc.Kind, &sc.Content, &tagsJS, &linksJS, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
            continue
        }
        _ = json.Unmarshal([]byte(emptyJSON(tagsJS, "[]")), &sc.Tags)
        _ = json.Unmarshal([]byte(emptyJSON(linksJS, "[]")), &sc.Links)
        out = append(out, sc)
    }
    return out, rows.Err()
}

func (s *RPStore) GetStoryCard(id string) (*RPStoryCard, error) {
    row := s.db.DB.QueryRow("SELECT id,title,kind,content,tags_json,links_json,created_at,updated_at FROM rp_story_cards WHERE id=?", id)
    var sc RPStoryCard
    var tagsJS, linksJS string
    if err := row.Scan(&sc.ID, &sc.Title, &sc.Kind, &sc.Content, &tagsJS, &linksJS, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
        if err == sql.ErrNoRows { return nil, nil }
        return nil, err
    }
    _ = json.Unmarshal([]byte(emptyJSON(tagsJS, "[]")), &sc.Tags)
    _ = json.Unmarshal([]byte(emptyJSON(linksJS, "[]")), &sc.Links)
    return &sc, nil
}

func (s *RPStore) SaveStoryCard(sc *RPStoryCard) error {
    now := time.Now().UTC().Format(time.RFC3339)
    if sc.CreatedAt == "" { sc.CreatedAt = now }
    sc.UpdatedAt = now
    tagsJS, _ := json.Marshal(sc.Tags)
    linksJS, _ := json.Marshal(sc.Links)
    _, err := s.db.DB.Exec(`
        INSERT INTO rp_story_cards(id,title,kind,content,tags_json,links_json,created_at,updated_at)
        VALUES(?,?,?,?,?,?,?,?)
        ON CONFLICT(id) DO UPDATE SET
            title=excluded.title,
            kind=excluded.kind,
            content=excluded.content,
            tags_json=excluded.tags_json,
            links_json=excluded.links_json,
            updated_at=excluded.updated_at
    `, sc.ID, sc.Title, sc.Kind, sc.Content, string(tagsJS), string(linksJS), sc.CreatedAt, sc.UpdatedAt)
    return err
}

func (s *RPStore) DeleteStoryCard(id string) error {
    _, err := s.db.DB.Exec("DELETE FROM rp_story_cards WHERE id=?", id)
    return err
}

// small helper
func emptyJSON(s, def string) string {
    if strings.TrimSpace(s) == "" { return def }
    return s
}
