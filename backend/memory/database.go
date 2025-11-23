/**
 * Database initialization and management module.
 *
 * Handles SQLite database creation, schema initialization, and connection
 * management for the NIRA memory system.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: database.go
 * Description: Database setup and connection management.
 */

package memory

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}
	if err := database.InitializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

func (d *Database) InitializeSchema() error {
    schema := `
	CREATE TABLE IF NOT EXISTS conversations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		title TEXT,
		mode TEXT NOT NULL DEFAULT 'normal',
		metadata TEXT
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id INTEGER NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		metadata TEXT,
		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS memories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT UNIQUE NOT NULL,
		content TEXT NOT NULL,
		category TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		importance INTEGER NOT NULL DEFAULT 50
	);

	CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id);
	CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);
	CREATE INDEX IF NOT EXISTS idx_memories_key ON memories(key);
	CREATE INDEX IF NOT EXISTS idx_memories_category ON memories(category);

	-- Allowed directories for sandboxed filesystem access
	CREATE TABLE IF NOT EXISTS allowed_directories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT UNIQUE NOT NULL,
		added_at TEXT NOT NULL
	);

	-- Lightweight RAG text index (basic, non-embedding)
	CREATE TABLE IF NOT EXISTS rag_index (
		path TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		mod_time TEXT NOT NULL,
		size INTEGER NOT NULL,
		hash TEXT,
		content TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_rag_index_name ON rag_index(name);
	CREATE INDEX IF NOT EXISTS idx_rag_index_mod ON rag_index(mod_time);
    `

	if _, err := d.DB.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
