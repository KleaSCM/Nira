/**
 * Ollama client module.
 *
 * Handles communication with the Ollama API for model inference, including
 * streaming chat completions and system prompt injection.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: ollama.go
 * Description: Ollama API client implementation.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaClient struct {
	Endpoint string
	Model    string
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatResponseChunk struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func NewOllamaClient(endpoint string, model string) *OllamaClient {
	return &OllamaClient{
		Endpoint: endpoint,
		Model:    model,
	}
}

func (c *OllamaClient) Chat(messages []ChatMessage, onChunk func(string) error) error {
	url := fmt.Sprintf("%s/api/chat", c.Endpoint)

	reqBody := ChatRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ollama API error: %s", string(bodyBytes))
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk ChatResponseChunk
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode chunk: %w", err)
		}

		if chunk.Message.Content != "" {
			if err := onChunk(chunk.Message.Content); err != nil {
				return err
			}
		}

		if chunk.Done {
			break
		}
	}

	return nil
}
