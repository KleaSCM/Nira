// WebSearchTool implementation for Nira
//
// Allows the AI to perform web searches and return summarized results.
// Uses DuckDuckGo Instant Answer API (no API key required, privacy-friendly).
//
// Author: KleaSCM
// Email: KleaSCM@gmail.com

package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// WebSearchResult represents a single search result
type WebSearchResult struct {
	Title   string
	Snippet string
	URL     string
	Source  string
}

// WebSearchTool implements the Tool interface for web search
type WebSearchTool struct{}

// Schema returns the tool's metadata as a map for registry listing
func (t *WebSearchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.Name(),
		"description": t.Description(),
		"permission":  t.PermissionLevel(),
		"input":       t.InputSchema(),
		"output":      t.OutputSchema(),
	}
}

func (t *WebSearchTool) Name() string {
	return "web_search"
}

func (t *WebSearchTool) Description() string {
	return "Searches the web for a query and returns summarized results."
}

func (t *WebSearchTool) PermissionLevel() string {
	return "internet"
}

func (t *WebSearchTool) InputSchema() string {
	return `{ "query": "string" }`
}

func (t *WebSearchTool) OutputSchema() string {
	return `[{ "Title": "string", "Snippet": "string", "URL": "string", "Source": "string" }]`
}

// Execute performs the web search
func (t *WebSearchTool) Execute(input map[string]interface{}) (interface{}, error) {
	query, ok := input["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("missing or invalid query")
	}

	// URL encode the query
	encodedQuery := url.QueryEscape(query)

	// DuckDuckGo Instant Answer API
	apiURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_redirect=1&no_html=1&skip_disambig=1", encodedQuery)

	// Create request with proper headers
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check if we got HTML instead of JSON (common DuckDuckGo issue)
	bodyStr := string(body)
	if strings.HasPrefix(strings.TrimSpace(bodyStr), "<") {
		return nil, fmt.Errorf("DuckDuckGo returned HTML instead of JSON. The Instant Answer API may not have results for '%s'. Try a more specific query about well-known topics (e.g., famous people, places, or Wikipedia subjects)", query)
	}

	// Parse JSON response
	var ddg struct {
		AbstractText   string `json:"AbstractText"`
		AbstractURL    string `json:"AbstractURL"`
		AbstractSource string `json:"AbstractSource"`
		Heading        string `json:"Heading"`
		RelatedTopics  []struct {
			Text     string `json:"Text"`
			FirstURL string `json:"FirstURL"`
		} `json:"RelatedTopics"`
		Results []struct {
			Text     string `json:"Text"`
			FirstURL string `json:"FirstURL"`
		} `json:"Results"`
	}

	if err := json.Unmarshal(body, &ddg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w (response: %s)", err, truncate(bodyStr, 200))
	}

	results := []WebSearchResult{}

	// Add abstract/main result if available
	if ddg.AbstractText != "" && ddg.AbstractURL != "" {
		results = append(results, WebSearchResult{
			Title:   ddg.Heading,
			Snippet: ddg.AbstractText,
			URL:     ddg.AbstractURL,
			Source:  ddg.AbstractSource,
		})
	}

	// Add related topics
	for _, topic := range ddg.RelatedTopics {
		if topic.Text != "" && topic.FirstURL != "" {
			results = append(results, WebSearchResult{
				Title:   extractTitle(topic.Text),
				Snippet: topic.Text,
				URL:     topic.FirstURL,
				Source:  "DuckDuckGo",
			})
		}
		if len(results) >= 5 {
			break
		}
	}

	// Add direct results
	for _, result := range ddg.Results {
		if result.Text != "" && result.FirstURL != "" {
			results = append(results, WebSearchResult{
				Title:   extractTitle(result.Text),
				Snippet: result.Text,
				URL:     result.FirstURL,
				Source:  "DuckDuckGo",
			})
		}
		if len(results) >= 5 {
			break
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for '%s'. DuckDuckGo's Instant Answer API works best for:\n- Famous people (e.g., 'Albert Einstein')\n- Well-known places (e.g., 'Eiffel Tower')\n- Wikipedia topics (e.g., 'Artificial Intelligence')\n\nFor better results, consider using a proper search API like Brave Search or Google Custom Search", query)
	}

	return results, nil
}

// extractTitle extracts the first sentence or up to 60 chars as title
func extractTitle(text string) string {
	if len(text) == 0 {
		return "Result"
	}

	// Try to get first sentence
	if idx := strings.Index(text, "."); idx > 0 && idx < 60 {
		return text[:idx]
	}

	// Otherwise truncate to 60 chars
	if len(text) > 60 {
		return text[:57] + "..."
	}

	return text
}

// truncate helper for error messages
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// RegisterWebSearchTool adds the tool to the registry
func RegisterWebSearchTool(registry map[string]Tool) {
	registry["web_search"] = &WebSearchTool{}
}
