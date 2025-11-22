/**
 * Memory manager module.
 *
 * Provides high-level memory operations and integrates memory storage
 * with the conversation flow. Handles automatic memory extraction and
 * context injection.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: manager.go
 * Description: High-level memory management and integration.
 */

package memory

type Manager struct {
	Conversations *ConversationStore
	Memories      *MemoryStore
	CurrentConvID int64
}

func NewManager(db *Database) (*Manager, error) {
	convStore := NewConversationStore(db)
	memStore := NewMemoryStore(db)

	manager := &Manager{
		Conversations: convStore,
		Memories:      memStore,
	}

	currentID, err := convStore.GetCurrentConversation()
	if err != nil {
		return nil, err
	}

	if currentID == 0 {
		currentID, err = convStore.CreateConversation("normal")
		if err != nil {
			return nil, err
		}
	}

	manager.CurrentConvID = currentID
	return manager, nil
}

func (m *Manager) SaveMessage(role, content, metadata string) error {
	return m.Conversations.AddMessage(m.CurrentConvID, role, content, metadata)
}

func (m *Manager) LoadRecentMessages(limit int) ([]*Message, error) {
	return m.Conversations.GetMessages(m.CurrentConvID)
}

func (m *Manager) GetContextMemories(limit int) ([]*Memory, error) {
	return m.Memories.SearchMemories("", 30)
}

func (m *Manager) StartNewConversation(mode string) (int64, error) {
	id, err := m.Conversations.CreateConversation(mode)
	if err != nil {
		return 0, err
	}
	m.CurrentConvID = id
	return id, nil
}

