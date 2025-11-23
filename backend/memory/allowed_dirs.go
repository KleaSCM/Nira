package memory

import (
    "database/sql"
    "path/filepath"
    "time"
)

// AllowedDirsStore manages the list of allowed root directories for file tools
type AllowedDirsStore struct {
    db *Database
    // in-memory cache of absolute, cleaned paths
    cache []string
}

func NewAllowedDirsStore(db *Database) (*AllowedDirsStore, error) {
    s := &AllowedDirsStore{db: db}
    if err := s.loadCache(); err != nil {
        return nil, err
    }
    return s, nil
}

func (s *AllowedDirsStore) loadCache() error {
    rows, err := s.db.DB.Query("SELECT path FROM allowed_directories ORDER BY id ASC")
    if err != nil {
        // If table is missing for some reason, try to init again
        if err2 := s.db.InitializeSchema(); err2 != nil {
            return err
        }
        rows, err = s.db.DB.Query("SELECT path FROM allowed_directories ORDER BY id ASC")
        if err != nil {
            return err
        }
    }
    defer rows.Close()
    s.cache = []string{}
    for rows.Next() {
        var p string
        if err := rows.Scan(&p); err == nil {
            s.cache = append(s.cache, p)
        }
    }
    return rows.Err()
}

// EnsureSeed inserts initial allowed paths if the table is empty.
func (s *AllowedDirsStore) EnsureSeed(paths []string) error {
    // if any rows exist, skip
    var count int
    _ = s.db.DB.QueryRow("SELECT COUNT(1) FROM allowed_directories").Scan(&count)
    if count > 0 {
        return s.loadCache()
    }
    for _, p := range paths {
        if p == "" { continue }
        _ = s.Add(p)
    }
    return s.loadCache()
}

// List returns the cached list of allowed directories (absolute paths)
func (s *AllowedDirsStore) List() []string { return append([]string{}, s.cache...) }

// Add inserts a directory into the table (normalized absolute path). No-op if exists.
func (s *AllowedDirsStore) Add(path string) error {
    if path == "" { return nil }
    abs, err := filepath.Abs(path)
    if err != nil { return err }
    abs = filepath.Clean(abs)
    _, err = s.db.DB.Exec(
        "INSERT OR IGNORE INTO allowed_directories(path, added_at) VALUES(?, ?)",
        abs, time.Now().UTC().Format(time.RFC3339),
    )
    if err != nil { return err }
    return s.loadCache()
}

// Remove deletes a directory row (by absolute normalized path).
func (s *AllowedDirsStore) Remove(path string) error {
    abs, err := filepath.Abs(path)
    if err != nil { return err }
    abs = filepath.Clean(abs)
    _, err = s.db.DB.Exec("DELETE FROM allowed_directories WHERE path = ?", abs)
    if err != nil { return err }
    return s.loadCache()
}

// IsAllowed checks whether the given path is within any allowed directory.
func (s *AllowedDirsStore) IsAllowed(path string) bool {
    if len(s.cache) == 0 { return false }
    absPath, err := filepath.Abs(path)
    if err != nil { return false }
    for _, allowed := range s.cache {
        rel, err := filepath.Rel(allowed, absPath)
        if err != nil { continue }
        if rel != ".." && !filepath.IsAbs(rel) {
            return true
        }
    }
    return false
}

// Helper used in tests to clear all rows
func (s *AllowedDirsStore) clearAll(tx *sql.Tx) error {
    if tx != nil {
        _, err := tx.Exec("DELETE FROM allowed_directories")
        return err
    }
    _, err := s.db.DB.Exec("DELETE FROM allowed_directories")
    return err
}
