/**
 * Long-term memory storage module.
 *
 * Handles persistence and retrieval of long-term knowledge fragments,
 * facts, preferences, and contextual information.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: memory.go
 * Description: Long-term memory persistence.
 */

package memory

import (
	"database/sql"
	"fmt"
	"time"
)

type Memory struct {
	ID         int64
	Key        string
	Content    string
	Category   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Importance int
}

type MemoryStore struct {
	DB *Database
}

func NewMemoryStore(db *Database) *MemoryStore {
	return &MemoryStore{DB: db}
}

func (ms *MemoryStore) StoreMemory(key, content, category string, importance int) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := ms.DB.DB.Exec(
		`INSERT INTO memories (key, content, category, created_at, updated_at, importance)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET
		 content = excluded.content,
		 category = excluded.category,
		 updated_at = excluded.updated_at,
		 importance = excluded.importance`,
		key, content, category, now, now, importance,
	)
	if err != nil {
		return fmt.Errorf("failed to store memory: %w", err)
	}

	return nil
}

func (ms *MemoryStore) GetMemory(key string) (*Memory, error) {
	var mem Memory
	var createdAt, updatedAt string

	err := ms.DB.DB.QueryRow(
		"SELECT id, key, content, category, created_at, updated_at, importance FROM memories WHERE key = ?",
		key,
	).Scan(&mem.ID, &mem.Key, &mem.Content, &mem.Category, &createdAt, &updatedAt, &mem.Importance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("memory not found")
		}
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	mem.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	mem.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &mem, nil
}

func (ms *MemoryStore) SearchMemories(category string, minImportance int) ([]*Memory, error) {
	query := "SELECT id, key, content, category, created_at, updated_at, importance FROM memories WHERE importance >= ?"
	args := []interface{}{minImportance}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY importance DESC, updated_at DESC"

	rows, err := ms.DB.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories: %w", err)
	}
	defer rows.Close()

	var memories []*Memory
	for rows.Next() {
		var mem Memory
		var createdAt, updatedAt string

		if err := rows.Scan(&mem.ID, &mem.Key, &mem.Content, &mem.Category, &createdAt, &updatedAt, &mem.Importance); err != nil {
			return nil, fmt.Errorf("failed to scan memory: %w", err)
		}

		mem.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		mem.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		memories = append(memories, &mem)
	}

	return memories, nil
}

func (ms *MemoryStore) DeleteMemory(key string) error {
	_, err := ms.DB.DB.Exec("DELETE FROM memories WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}
	return nil
}

