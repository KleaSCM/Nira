package memory

import (
    "crypto/sha1"
    "database/sql"
    "encoding/hex"
    "fmt"
    "path/filepath"
    "strings"
)

// RagIndex provides minimal text indexing and search over small local files.
type RagIndex struct{ db *Database }

func NewRagIndex(db *Database) *RagIndex { return &RagIndex{db: db} }

func (ri *RagIndex) Upsert(path, name, modTime string, size int64, content string) error {
    abs, _ := filepath.Abs(path)
    h := sha1.Sum([]byte(content))
    hash := hex.EncodeToString(h[:])
    _, err := ri.db.DB.Exec(`
        INSERT INTO rag_index(path, name, mod_time, size, hash, content)
        VALUES(?, ?, ?, ?, ?, ?)
        ON CONFLICT(path) DO UPDATE SET
            name=excluded.name,
            mod_time=excluded.mod_time,
            size=excluded.size,
            hash=excluded.hash,
            content=excluded.content
    `, abs, name, modTime, size, hash, content)
    return err
}

// Search performs a simple LIKE-based search on name and content with optional path prefix.
func (ri *RagIndex) Search(query string, limit int, pathPrefix string) ([]map[string]interface{}, error) {
    if limit <= 0 { limit = 10 }
    like := "%" + strings.ReplaceAll(query, "%", "") + "%"
    where := "WHERE (name LIKE ? OR content LIKE ?)"
    args := []interface{}{like, like}
    if pathPrefix != "" {
        where += " AND path LIKE ?"
        pp := strings.TrimRight(pathPrefix, "\\/") + "%"
        args = append(args, pp)
    }
    rows, err := ri.db.DB.Query(
        fmt.Sprintf("SELECT path, name, mod_time, size, content FROM rag_index %s ORDER BY mod_time DESC LIMIT ?", where),
        append(args, limit)...,
    )
    if err != nil { return nil, err }
    defer rows.Close()
    out := []map[string]interface{}{}
    for rows.Next() {
        var path, name, mod, content string
        var size int64
        if err := rows.Scan(&path, &name, &mod, &size, &content); err == nil {
            snippet := makeSnippet(content, query, 160)
            out = append(out, map[string]interface{}{
                "path": path,
                "name": name,
                "mod_time": mod,
                "size": size,
                "snippet": snippet,
                "score": scoreSimple(name, content, query),
            })
        }
    }
    return out, rows.Err()
}

// Helpers
func makeSnippet(content, query string, max int) string {
    if content == "" { return "" }
    lc := strings.ToLower(content)
    lq := strings.ToLower(query)
    idx := strings.Index(lc, lq)
    if idx < 0 {
        if len(content) <= max { return content }
        return content[:max]
    }
    start := idx - max/4
    if start < 0 { start = 0 }
    end := start + max
    if end > len(content) { end = len(content) }
    return content[start:end]
}

func scoreSimple(name, content, query string) float64 {
    n := strings.Count(strings.ToLower(name), strings.ToLower(query))
    c := strings.Count(strings.ToLower(content), strings.ToLower(query))
    return float64(n*5 + c)
}

// For cleanup when a directory is removed from allowed list
func (ri *RagIndex) DeleteByPathPrefix(tx *sql.Tx, prefix string) error {
    p := strings.TrimRight(prefix, "\\/") + "%"
    if tx != nil {
        _, err := tx.Exec("DELETE FROM rag_index WHERE path LIKE ?", p)
        return err
    }
    _, err := ri.db.DB.Exec("DELETE FROM rag_index WHERE path LIKE ?", p)
    return err
}
