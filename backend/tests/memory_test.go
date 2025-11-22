// memory_test.go
package tests

import (
	"nira/memory"
	"testing"
)

func setupTestDB(t *testing.T) *memory.Database {
	db, err := memory.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func TestMemoryManager(t *testing.T) {
	// Initialize test database
	db := setupTestDB(t)
	defer db.Close()

	// Create memory manager
	m, err := memory.NewManager(db)
	if err != nil {
		t.Fatalf("Failed to create memory manager: %v", err)
	}

	t.Run("Conversation Management", func(t *testing.T) {
		// Test creating a new conversation
		convID, err := m.Conversations.CreateConversation("test")
		if err != nil {
			t.Fatalf("Failed to create conversation: %v", err)
		}

		// Test adding a message
		err = m.Conversations.AddMessage(convID, "user", "Test message", "")
		if err != nil {
			t.Fatalf("Failed to add message: %v", err)
		}

		// Test retrieving messages
		messages, err := m.Conversations.GetMessages(convID)
		if err != nil {
			t.Fatalf("Failed to get messages: %v", err)
		}

		if len(messages) == 0 {
			t.Fatal("Expected at least one message")
		}

		// Test message content
		if messages[0].Content != "Test message" {
			t.Errorf("Expected message content 'Test message', got '%s'", messages[0].Content)
		}
	})

	t.Run("Manager Methods", func(t *testing.T) {
		// Test saving a message through manager
		err := m.SaveMessage("user", "Manager test message", "")
		if err != nil {
			t.Fatalf("Failed to save message through manager: %v", err)
		}

		// Test getting context memories
		memories, err := m.GetContextMemories(5)
		if err != nil {
			t.Fatalf("Failed to get context memories: %v", err)
		}
		if len(memories) == 0 {
			t.Error("Expected to get some context memories, got none")
		}

		// Test starting a new conversation
		newConvID, err := m.StartNewConversation("test")
		if err != nil {
			t.Fatalf("Failed to start new conversation: %v", err)
		}

		if newConvID == 0 {
			t.Error("Expected non-zero conversation ID")
		}

		// Test getting current conversation
		conv, err := m.Conversations.GetConversation(m.CurrentConvID)
		if err != nil {
			t.Fatalf("Failed to get current conversation: %v", err)
		}

		if conv == nil {
			t.Fatal("Expected to get current conversation, got nil")
		}

		// Test listing conversations
		conversations, err := m.Conversations.ListConversations(10, 0)
		if err != nil {
			t.Fatalf("Failed to list conversations: %v", err)
		}

		if len(conversations) == 0 {
			t.Error("Expected at least one conversation")
		}

		// Test storing and retrieving memory
		memoryKey := "test:key"
		memoryContent := "Test memory content"
		err = m.Memories.StoreMemory(memoryKey, memoryContent, "test", 50)
		if err != nil {
			t.Fatalf("Failed to store memory: %v", err)
		}

		mem, err := m.Memories.GetMemory(memoryKey)
		if err != nil {
			t.Fatalf("Failed to get memory: %v", err)
		}

		if mem.Content != memoryContent {
			t.Errorf("Expected memory content '%s', got '%s'", memoryContent, mem.Content)
		}

		// Test searching memories
		memories, err = m.Memories.SearchMemories("test", 0)
		if err != nil {
			t.Fatalf("Failed to search memories: %v", err)
		}

		if len(memories) == 0 {
			t.Error("Expected to find at least one memory")
		}

		// Test deleting memory
		err = m.Memories.DeleteMemory(memoryKey)
		if err != nil {
			t.Fatalf("Failed to delete memory: %v", err)
		}

		// Verify memory was deleted
		_, err = m.Memories.GetMemory(memoryKey)
		if err == nil {
			t.Error("Expected error when getting deleted memory, got nil")
		}
	})
}

// TestConversationFlow verifies the complete message exchange pattern
// between user and assistant within a single conversation.
func TestConversationFlow(t *testing.T) {
	// Initialize test database
	db := setupTestDB(t)
	defer db.Close()

	// Initialize manager
	m, err := memory.NewManager(db)
	if err != nil {
		t.Fatalf("Manager initialization failed: %v", err)
	}

	// Simulate user interaction
	userMessage := "Hello, how are you?"
	if err := m.SaveMessage("user", userMessage, ""); err != nil {
		t.Fatalf("Failed to persist user message: %v", err)
	}

	// Simulate AI response
	assistantResponse := "I'm doing well, thank you!"
	if err := m.SaveMessage("assistant", assistantResponse, ""); err != nil {
		t.Fatalf("Failed to persist assistant response: %v", err)
	}

	// Verify conversation history integrity
	messages, err := m.LoadRecentMessages(10)
	if err != nil {
		t.Fatalf("Conversation history retrieval failed: %v", err)
	}

	// Verify message count and order
	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	// Verify message content and role preservation
	if messages[0].Role != "user" || messages[0].Content != userMessage {
		t.Errorf("User message corrupted. Got: %+v", messages[0])
	}

	if messages[1].Role != "assistant" || messages[1].Content != assistantResponse {
		t.Errorf("Assistant response corrupted. Got: %+v", messages[1])
	}
}

// TestMultipleConversations verifies isolation between different conversation contexts
// to ensure messages don't leak between conversations.
func TestMultipleConversations(t *testing.T) {
	// Initialize test database
	db := setupTestDB(t)
	defer db.Close()

	// Initialize manager
	m, err := memory.NewManager(db)
	if err != nil {
		t.Fatalf("Manager initialization failed: %v", err)
	}

	// Initialize first conversation with a message
	firstMessage := "First conversation message"
	if err := m.SaveMessage("user", firstMessage, ""); err != nil {
		t.Fatalf("Failed to save first conversation message: %v", err)
	}

	// Create and switch to second conversation
	_, err = m.StartNewConversation("test")
	if err != nil {
		t.Fatalf("Failed to create second conversation: %v", err)
	}

	// Add message to second conversation
	secondMessage := "Second conversation message"
	if err := m.SaveMessage("user", secondMessage, ""); err != nil {
		t.Fatalf("Failed to save second conversation message: %v", err)
	}

	// Verify second conversation content
	messages, err := m.LoadRecentMessages(10)
	if err != nil {
		t.Fatalf("Failed to load second conversation: %v", err)
	}

	if len(messages) != 1 || messages[0].Content != secondMessage {
		t.Fatalf("Second conversation verification failed. Got: %+v", messages)
	}

	// Switch back to first conversation and verify isolation
	currentConvID, err := m.Conversations.GetCurrentConversation()
	if err != nil {
		t.Fatalf("Failed to get current conversation: %v", err)
	}
	m.CurrentConvID = currentConvID
	messages, err = m.LoadRecentMessages(10)
	if err != nil {
		t.Fatalf("Failed to reload first conversation: %v", err)
	}

	if len(messages) != 1 || messages[0].Content != firstMessage {
		t.Fatalf("First conversation isolation failed. Got: %+v", messages)
	}
}
