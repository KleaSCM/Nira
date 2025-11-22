package tests

import (
	"nira/memory"
	"sync"
	"testing"
)

// TestIntegration_CompleteFlow verifies the end-to-end functionality
// of the memory management system including conversation handling,
// message persistence, and concurrent access patterns.
func TestIntegration_CompleteFlow(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	m, err := memory.NewManager(db)
	if err != nil {
		t.Fatalf("Manager initialization failed: %v", err)
	}

	t.Run("Conversation Lifecycle", func(t *testing.T) {
		// Initialize test conversation with a unique context
		convID, err := m.Conversations.CreateConversation("integration-test")
		if err != nil {
			t.Fatalf("Conversation creation failed: %v", err)
		}

		// Define test conversation flow to verify message ordering and persistence
		testDialogue := []struct {
			role    string
			content string
		}{
			{"user", "What's the weather like?"},
			{"assistant", "I'm sorry, I don't have weather information."},
			{"user", "Can you help with something else?"},
		}

		// Persist conversation messages
		for _, msg := range testDialogue {
			if err := m.Conversations.AddMessage(convID, msg.role, msg.content, ""); err != nil {
				t.Fatalf("Message persistence failed: %v", err)
			}
		}

		// Verify conversation integrity
		messages, err := m.Conversations.GetMessages(convID)
		if err != nil {
			t.Fatalf("Failed to retrieve conversation: %v", err)
		}

		if len(messages) != len(testDialogue) {
			t.Fatalf("Message count mismatch. Got %d, expected %d",
				len(messages), len(testDialogue))
		}

		// Validate each message's content and order
		for i, expected := range testDialogue {
			if messages[i].Role != expected.role || messages[i].Content != expected.content {
				t.Errorf("Message %d validation failed. Got: %+v, Expected: %+v",
					i, messages[i], expected)
			}
		}

		t.Run("Memory Operations", func(t *testing.T) {
			// Store memory with test data
			memoryKey := "integration:test:key"
			memoryContent := "Test memory for integration testing"
			if err := m.Memories.StoreMemory(memoryKey, memoryContent, "test", 75); err != nil {
				t.Fatalf("Memory storage failed: %v", err)
			}

			// Verify memory retrieval
			mem, err := m.Memories.GetMemory(memoryKey)
			if err != nil {
				t.Fatalf("Memory retrieval failed: %v", err)
			}

			if mem.Content != memoryContent {
				t.Fatalf("Memory content corrupted. Got: %s, Expected: %s",
					mem.Content, memoryContent)
			}

			// Test memory search functionality
			memories, err := m.Memories.SearchMemories("test", 50)
			if err != nil {
				t.Fatalf("Memory search failed: %v", err)
			}

			if len(memories) == 0 {
				t.Fatal("Expected to find stored memories")
			}

			// Clean up test memory
			if err := m.Memories.DeleteMemory(memoryKey); err != nil {
				t.Fatalf("Memory cleanup failed: %v", err)
			}

			// Verify deletion
			if _, err := m.Memories.GetMemory(memoryKey); err == nil {
				t.Error("Expected error when retrieving deleted memory")
			}
		})
	})

	t.Run("Concurrent Access Patterns", func(t *testing.T) {
		// Test concurrent write operations to verify thread safety
		const (
			numRoutines       = 5
			messagesPerThread = 10
		)

		// Initialize conversation for concurrent testing
		convID, err := m.Conversations.CreateConversation("concurrent-test")
		if err != nil {
			t.Fatalf("Failed to create test conversation: %v", err)
		}

		var wg sync.WaitGroup
		errCh := make(chan error, numRoutines*messagesPerThread)

		// Launch multiple goroutines to add messages concurrently
		for i := 0; i < numRoutines; i++ {
			wg.Add(1)
			go func(threadID int) {
				defer wg.Done()
				for j := 0; j < messagesPerThread; j++ {
					msg := memory.Message{
						ConversationID: convID,
						Role:           "user",
						Content:        string(rune('A'+threadID)) + "-msg-" + string(rune('0'+j)),
					}

					if err := m.Conversations.AddMessage(convID, msg.Role, msg.Content, ""); err != nil {
						errCh <- err
						return
					}
				}
			}(i)
		}

		// Wait for all goroutines to complete
		wg.Wait()
		close(errCh)

		// Check for any errors during concurrent execution
		for err := range errCh {
			t.Fatalf("Concurrent write error: %v", err)
		}

		// Verify all messages were persisted correctly
		messages, err := m.Conversations.GetMessages(convID)
		if err != nil {
			t.Fatalf("Failed to retrieve messages: %v", err)
		}

		expectedCount := numRoutines * messagesPerThread
		if len(messages) != expectedCount {
			t.Fatalf("Message count mismatch. Got %d, expected %d",
				len(messages), expectedCount)
		}

		// Verify message uniqueness and content
		messageSet := make(map[string]bool, len(messages))
		for _, msg := range messages {
			messageSet[msg.Content] = true
		}

		for i := 0; i < numRoutines; i++ {
			for j := 0; j < messagesPerThread; j++ {
				expected := string(rune('A'+i)) + "-msg-" + string(rune('0'+j))
				if !messageSet[expected] {
					t.Errorf("Expected message not found: %s", expected)
				}
			}
		}
	})
}

func TestIntegration_ConcurrentAccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create manager
	m, err := memory.NewManager(db)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Number of concurrent operations
	numOps := 10
	var wg sync.WaitGroup
	errChan := make(chan error, numOps*2)

	// Test concurrent message saving
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := m.SaveMessage("user", "Concurrent message", "")
			if err != nil {
				errChan <- err
			}
		}(i)
	}

	// Test concurrent memory storage
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := m.Memories.StoreMemory(
				"concurrent:test:"+string(rune('A'+id)),
				"Concurrent memory",
				"test",
				id+1,
			)
			if err != nil {
				errChan <- err
			}
		}(i)
	}

	// Wait for all operations to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for errors
	for err := range errChan {
		if err != nil {
			t.Errorf("Concurrent operation failed: %v", err)
		}
	}

	// Verify all messages were saved
	messages, err := m.LoadRecentMessages(numOps * 2)
	if err != nil {
		t.Fatalf("Failed to load messages: %v", err)
	}

	// We should have at least numOps messages (some might be from other tests)
	if len(messages) < numOps {
		t.Errorf("Expected at least %d messages, got %d", numOps, len(messages))
	}

	// Verify all memories were stored
	memories, err := m.Memories.SearchMemories("test", 0)
	if err != nil {
		t.Fatalf("Failed to search memories: %v", err)
	}

	// We should have at least numOps memories (some might be from other tests)
	if len(memories) < numOps {
		t.Errorf("Expected at least %d memories, got %d", numOps, len(memories))
	}
}
