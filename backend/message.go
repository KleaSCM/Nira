/**
 * Message handling module.
 *
 * Defines message types for WebSocket communication between frontend
 * and backend, including chat messages and tool call requests.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: message.go
 * Description: Message type definitions and serialization.
 */

package main

import "encoding/json"

type MessageType string

const (
	MessageTypeUser      MessageType = "user"
	MessageTypeAssistant MessageType = "assistant"
	MessageTypeSystem    MessageType = "system"
	MessageTypeError    MessageType = "error"
	MessageTypeChunk    MessageType = "chunk"
)

type WSMessage struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
	ID      string      `json:"id,omitempty"`
}

func (m *WSMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func ParseWSMessage(data []byte) (*WSMessage, error) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

