package tests

import (
	"nira/memory"
	"testing"
)

// TestMemoryStore_CRUD verifies the complete lifecycle of memory operations
// including storage, retrieval, updating, deletion, and searching.
func TestMemoryStore_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := memory.NewMemoryStore(db)

	t.Run("Memory Lifecycle", func(t *testing.T) {
		// Test data setup with clear purpose for each field
		key := "test:key:1"
		content := "Test memory content"
		category := "test"
		importance := 75

		// Verify memory creation and retrieval
		err := store.StoreMemory(key, content, category, importance)
		if err != nil {
			t.Fatalf("Memory storage failed: %v", err)
		}

		mem, err := store.GetMemory(key)
		if err != nil {
			t.Fatalf("Failed to retrieve stored memory: %v", err)
		}

		// Validate data integrity after storage
		if mem.Content != content || mem.Category != category || mem.Importance != importance {
			t.Fatalf("Memory data corrupted. Got: %+v, Expected: {Content:%s Category:%s Importance:%d}",
				mem, content, category, importance)
		}

		// Test memory update functionality
		updatedContent := "Updated memory content"
		err = store.StoreMemory(key, updatedContent, category, importance+1)
		if err != nil {
			t.Fatalf("Memory update operation failed: %v", err)
		}

		mem, err = store.GetMemory(key)
		if err != nil {
			t.Fatalf("Failed to retrieve updated memory: %v", err)
		}

		// Verify update was applied correctly
		if mem.Content != updatedContent || mem.Importance != importance+1 {
			t.Errorf("Memory update verification failed. Got: %+v, Expected content: %s, importance: %d",
				mem, updatedContent, importance+1)
		}

		// Test memory deletion
		err = store.DeleteMemory(key)
		if err != nil {
			t.Fatalf("Memory deletion failed: %v", err)
		}

		// Verify deletion was successful
		_, err = store.GetMemory(key)
		if err == nil {
			t.Error("Expected error when retrieving deleted memory, got nil")
		}
	})

	t.Run("Memory Search Functionality", func(t *testing.T) {
		// Setup diverse test data to verify search capabilities
		testMemories := []struct {
			key        string
			content    string
			category   string
			importance int
		}{
			{"search:1", "Test memory 1", "test", 50},   // Medium importance test
			{"search:2", "Test memory 2", "test", 75},   // High importance test
			{"search:3", "Another memory", "other", 25}, // Different category
		}

		for _, m := range testMemories {
			err := store.StoreMemory(m.key, m.content, m.category, m.importance)
			if err != nil {
				t.Fatalf("Failed to initialize test memory %s: %v", m.key, err)
			}
		}

		// Define test cases to verify different search scenarios
		testCases := []struct {
			description   string
			category      string
			minImportance int
			expectedCount int
		}{
			{"Retrieve all memories", "", 0, 3},
			{"Filter by minimum importance", "", 50, 2},
			{"Filter by category", "test", 0, 2},
			{"Combined category and importance filter", "test", 50, 2},
			{"Non-existent category filter", "nonexistent", 0, 0},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				memories, err := store.SearchMemories(tc.category, tc.minImportance)
				if err != nil {
					t.Fatalf("SearchMemories failed: %v", err)
				}
				if len(memories) != tc.expectedCount {
					t.Errorf("Expected %d memories, got %d", tc.expectedCount, len(memories))
				}
			})
		}

		// Cleanup test data
		for _, m := range testMemories {
			if err := store.DeleteMemory(m.key); err != nil {
				t.Logf("Warning: Failed to clean up test memory %s: %v", m.key, err)
			}
		}
	})
}
