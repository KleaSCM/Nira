/**
 * Conversation storage module.
 *
 * Handles persistence and retrieval of conversation sessions and messages.
 * Provides CRUD operations for conversations and their associated messages.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: conversation.go
 * Description: Conversation and message persistence.
 */

package memory

import (
	"database/sql"
	"fmt"
	"time"
)

type Conversation struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Mode      string
	Metadata  string
}

type Message struct {
	ID             int64
	ConversationID int64
	Role           string
	Content        string
	Timestamp      time.Time
	Metadata       string
}

type ConversationStore struct {
	DB *Database
}

func NewConversationStore(db *Database) *ConversationStore {
	return &ConversationStore{DB: db}
}

func (cs *ConversationStore) CreateConversation(mode string) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := cs.DB.DB.Exec(
		"INSERT INTO conversations (created_at, updated_at, mode) VALUES (?, ?, ?)",
		now, now, mode,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create conversation: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get conversation ID: %w", err)
	}

	return id, nil
}

func (cs *ConversationStore) GetConversation(id int64) (*Conversation, error) {
	var conv Conversation
	var createdAt, updatedAt string

	err := cs.DB.DB.QueryRow(
		"SELECT id, created_at, updated_at, title, mode, metadata FROM conversations WHERE id = ?",
		id,
	).Scan(&conv.ID, &createdAt, &updatedAt, &conv.Title, &conv.Mode, &conv.Metadata)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	conv.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	conv.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &conv, nil
}

func (cs *ConversationStore) ListConversations(limit, offset int) ([]*Conversation, error) {
	rows, err := cs.DB.DB.Query(
		"SELECT id, created_at, updated_at, title, mode, metadata FROM conversations ORDER BY updated_at DESC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		var conv Conversation
		var createdAt, updatedAt string

		if err := rows.Scan(&conv.ID, &createdAt, &updatedAt, &conv.Title, &conv.Mode, &conv.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}

		conv.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		conv.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		conversations = append(conversations, &conv)
	}

	return conversations, nil
}

func (cs *ConversationStore) AddMessage(conversationID int64, role, content, metadata string) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	_, err := cs.DB.DB.Exec(
		"INSERT INTO messages (conversation_id, role, content, timestamp, metadata) VALUES (?, ?, ?, ?, ?)",
		conversationID, role, content, timestamp, metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to add message: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	_, err = cs.DB.DB.Exec(
		"UPDATE conversations SET updated_at = ? WHERE id = ?",
		now, conversationID,
	)
	if err != nil {
		return fmt.Errorf("failed to update conversation timestamp: %w", err)
	}

	return nil
}

func (cs *ConversationStore) GetMessages(conversationID int64) ([]*Message, error) {
	rows, err := cs.DB.DB.Query(
		"SELECT id, conversation_id, role, content, timestamp, metadata FROM messages WHERE conversation_id = ? ORDER BY timestamp ASC",
		conversationID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		var timestamp string

		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &timestamp, &msg.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		msg.Timestamp, _ = time.Parse(time.RFC3339, timestamp)
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (cs *ConversationStore) DeleteConversation(id int64) error {
	tx, err := cs.DB.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete all messages in the conversation first
	_, err = tx.Exec("DELETE FROM messages WHERE conversation_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	// Delete the conversation
	result, err := tx.Exec("DELETE FROM conversations WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation with ID %d not found", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (cs *ConversationStore) GetCurrentConversation() (int64, error) {
	var id int64
	err := cs.DB.DB.QueryRow(
		"SELECT id FROM conversations ORDER BY updated_at DESC LIMIT 1",
	).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get current conversation: %w", err)
	}

	return id, nil
}
