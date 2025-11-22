package tests

import (
	"nira/memory"
	"testing"
	"time"
)

// TestConversationStore_CRUD verifies the complete lifecycle of conversation operations.
// It ensures data consistency and proper error handling across all CRUD operations.
func TestConversationStore_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := memory.NewConversationStore(db)

	t.Run("Create and Get Conversation", func(t *testing.T) {
		// Verifies conversation creation returns a valid ID and can be retrieved
		convID, err := store.CreateConversation("test")
		if err != nil {
			t.Fatalf("Failed to create conversation: %v", err)
		}

		if convID <= 0 {
			t.Error("Expected positive conversation ID")
		}

		// Validates conversation data integrity after creation
		conv, err := store.GetConversation(convID)
		if err != nil {
			t.Fatalf("Failed to get conversation: %v", err)
		}

		if conv == nil || conv.ID != convID || conv.Mode != "test" {
			t.Errorf("Unexpected conversation data: %+v", conv)
		}

		// Ensures conversation timestamps are updated on message addition
		oldUpdatedAt := conv.UpdatedAt
		err = store.AddMessage(convID, "user", "Test message", "")
		if err != nil {
			t.Fatalf("Failed to add message: %v", err)
		}

		conv, err = store.GetConversation(convID)
		if err != nil {
			t.Fatalf("Failed to get updated conversation: %v", err)
		}

		if !conv.UpdatedAt.After(oldUpdatedAt) {
			t.Error("Expected UpdatedAt to be updated after adding message")
		}

		// Validates message storage and retrieval
		messages, err := store.GetMessages(convID)
		if err != nil {
			t.Fatalf("Failed to get messages: %v", err)
		}

		if len(messages) != 1 || messages[0].Content != "Test message" {
			t.Errorf("Unexpected messages: %+v", messages)
		}
	})

	t.Run("List Conversations", func(t *testing.T) {
		// Creates multiple conversations to test list and pagination
		convIDs := make([]int64, 3)
		for i := 0; i < 3; i++ {
			id, err := store.CreateConversation("test")
			if err != nil {
				t.Fatalf("Failed to create conversation %d: %v", i, err)
			}
			convIDs[i] = id

			err = store.AddMessage(id, "user", "Test message", "")
			if err != nil {
				t.Fatalf("Failed to add message: %v", err)
			}

			time.Sleep(10 * time.Millisecond)
		}

		// Verifies conversation listing returns expected results
		conversations, err := store.ListConversations(10, 0)
		if err != nil {
			t.Fatalf("Failed to list conversations: %v", err)
		}

		if len(conversations) < 3 {
			t.Fatalf("Expected at least 3 conversations, got %d", len(conversations))
		}

		// Ensures conversations are sorted by update time
		for i := 0; i < len(conversations)-1; i++ {
			if conversations[i].UpdatedAt.Before(conversations[i+1].UpdatedAt) {
				t.Error("Conversations not ordered by updated_at desc")
			}
		}

		// Tests pagination functionality
		firstPage, err := store.ListConversations(2, 0)
		if err != nil {
			t.Fatalf("Failed to get first page: %v", err)
		}

		if len(firstPage) != 2 {
			t.Errorf("Expected 2 conversations on first page, got %d", len(firstPage))
		}

		secondPage, err := store.ListConversations(2, 2)
		if err != nil {
			t.Fatalf("Failed to get second page: %v", err)
		}

		if len(secondPage) == 0 {
			t.Error("Expected at least one conversation on second page")
		}

		t.Run("Delete Conversation", func(t *testing.T) {
			// Tests conversation deletion and data consistency
			convID := convIDs[0]

			// Pre-deletion validation
			conv, err := store.GetConversation(convID)
			if err != nil {
				t.Fatalf("Failed to get conversation before deletion: %v", err)
			}
			if conv == nil {
				t.Fatal("Expected conversation to exist before deletion")
			}

			// Validates conversation deletion
			err = store.DeleteConversation(convID)
			if err != nil {
				t.Fatalf("Failed to delete conversation: %v", err)
			}

			// Verifies conversation is no longer accessible
			_, err = store.GetConversation(convID)
			if err == nil {
				t.Error("Expected error when getting deleted conversation, got nil")
			}

			// Ensures related messages are also removed
			messages, err := store.GetMessages(convID)
			if err != nil {
				t.Fatalf("Failed to get messages: %v", err)
			}
			if len(messages) > 0 {
				t.Errorf("Expected no messages after conversation deletion, got %d", len(messages))
			}

			// Tests error handling for non-existent conversation
			err = store.DeleteConversation(999999)
			if err == nil {
				t.Error("Expected error when deleting non-existent conversation, got nil")
			}

			// Cleanup remaining test data
			for _, id := range convIDs[1:] {
				err := store.DeleteConversation(id)
				if err != nil {
					t.Errorf("Failed to clean up conversation %d: %v", id, err)
				}
			}
		})
	})
}
