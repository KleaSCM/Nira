package tests

import (
	"database/sql"
	"nira/memory"
	"os"
	"path/filepath"
	"testing"
)

// TestDatabase_Initialization verifies database creation, schema initialization,
// and basic read/write operations for both in-memory and file-based databases.
func TestDatabase_Initialization(t *testing.T) {
	t.Run("InMemory Database", func(t *testing.T) {
		// In-memory databases are used for testing to ensure isolation between tests
		db, err := memory.NewDatabase(":memory:")
		if err != nil {
			t.Fatalf("Failed to create in-memory database: %v", err)
		}
		defer db.Close()

		// Verify all required tables exist after initialization
		tables := []string{"conversations", "messages", "memories"}
		for _, table := range tables {
			var exists bool
			err := db.DB.QueryRow(
				"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
				table,
			).Scan(&exists)

			if err == sql.ErrNoRows {
				t.Errorf("Required table %s was not created during initialization", table)
			} else if err != nil {
				t.Fatalf("Error verifying table %s: %v", table, err)
			}
		}
	})

	t.Run("File-based Database", func(t *testing.T) {
		// File-based databases are used in production for persistence
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "test.db")

		db, err := memory.NewDatabase(dbPath)
		if err != nil {
			t.Fatalf("Failed to initialize file-based database: %v", err)
		}
		defer db.Close()

		// Verify database file creation
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Fatalf("Database file not created at expected location: %s", dbPath)
		}

		// Test data persistence by writing and reading back
		testKey := "test:key"
		testContent := "test content"
		testCategory := "test"
		testImportance := 50

		_, err = db.DB.Exec(
			"INSERT INTO memories (key, content, category, importance) VALUES (?, ?, ?, ?)",
			testKey, testContent, testCategory, testImportance,
		)
		if err != nil {
			t.Fatalf("Failed to write test data: %v", err)
		}

		var content string
		err = db.DB.QueryRow(
			"SELECT content FROM memories WHERE key = ?",
			testKey,
		).Scan(&content)

		if err != nil {
			t.Fatalf("Failed to verify data persistence: %v", err)
		}

		if content != testContent {
			t.Errorf("Data integrity check failed: expected '%s', got '%s'", testContent, content)
		}
	})

	t.Run("Invalid Path", func(t *testing.T) {
		// Verify proper error handling for invalid database paths
		_, err := memory.NewDatabase("/non/existent/path/test.db")
		if err == nil {
			t.Fatal("Expected error for invalid database path, got nil")
		}
	})
}
