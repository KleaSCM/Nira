package tests

import (
	"nira/memory"
	"testing"
)

// TestManager_ConversationFlow verifies the complete conversation management workflow
// including message handling and conversation switching.
func TestManager_ConversationFlow(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Initialize manager with a default conversation
	m, err := memory.NewManager(db)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	if m.CurrentConvID == 0 {
		t.Error("Manager should initialize with a default conversation")
	}

	// Validate basic message operations
	err = m.SaveMessage("user", "Hello, Nira!", "")
	if err != nil {
		t.Fatalf("Failed to save message: %v", err)
	}

	messages, err := m.LoadRecentMessages(10)
	if err != nil {
		t.Fatalf("Failed to load messages: %v", err)
	}

	if len(messages) != 1 || messages[0].Content != "Hello, Nira!" {
		t.Errorf("Unexpected messages: %+v", messages)
	}

	// Test conversation management
	newConvID, err := m.StartNewConversation("test")
	if err != nil {
		t.Fatalf("Failed to start new conversation: %v", err)
	}

	if newConvID == m.CurrentConvID {
		t.Error("New conversation should have a different ID than the current one")
	}

	// Verify conversation isolation
	oldConvID := m.CurrentConvID
	m.CurrentConvID = newConvID

	err = m.SaveMessage("user", "New conversation", "")
	if err != nil {
		t.Fatalf("Failed to save message in new conversation: %v", err)
	}

	// Ensure conversations remain isolated
	m.CurrentConvID = oldConvID
	messages, err = m.LoadRecentMessages(10)
	if err != nil {
		t.Fatalf("Failed to load messages from old conversation: %v", err)
	}

	if len(messages) != 1 || messages[0].Content != "Hello, Nira!" {
		t.Error("Messages from different conversations should not interfere")
	}
}

// TestManager_MemoryOperations verifies memory storage and retrieval
// through the manager interface.
func TestManager_MemoryOperations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	m, err := memory.NewManager(db)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Verify empty state handling
	memories, err := m.GetContextMemories(5)
	if err != nil {
		t.Fatalf("Failed to get context memories: %v", err)
	}

	if len(memories) != 0 {
		t.Errorf("Expected no memories in new database, got %d", len(memories))
	}

	// Test memory storage and retrieval
	err = m.Memories.StoreMemory("test:1", "Memory 1", "test", 50)
	if err != nil {
		t.Fatalf("Failed to store memory: %v", err)
	}

	err = m.Memories.StoreMemory("test:2", "Memory 2", "test", 75)
	if err != nil {
		t.Fatalf("Failed to store memory: %v", err)
	}

	// Verify memory retrieval
	memories, err = m.GetContextMemories(5)
	if err != nil {
		t.Fatalf("Failed to get context memories: %v", err)
	}

	if len(memories) != 2 {
		t.Errorf("Expected 2 memories, got %d", len(memories))
	}
}
