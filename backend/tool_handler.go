/**
 * Tool call handler module.
 *
 * Handles detection, parsing, validation, and execution of tool calls
 * from the model, including result injection back into the conversation.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: tool_handler.go
 * Description: Tool call processing and execution.
 */

package main

import (
	"encoding/json"
	"fmt"
	"nira/tools"
	"regexp"
	"strings"
)

type ToolHandler struct {
	Registry *tools.Registry
	Logger   *Logger
}

func NewToolHandler(registry *tools.Registry, logger *Logger) *ToolHandler {
	return &ToolHandler{
		Registry: registry,
		Logger:   logger,
	}
}

func (th *ToolHandler) DetectToolCall(content string) (*tools.Call, bool) {
	// Look for JSON tool call patterns in the response
	// Common patterns: {"tool": "...", "args": {...}} or <tool_call>...</tool_call>

	// Pattern 1: JSON object with tool/function fields
	jsonPattern := regexp.MustCompile(`\{[^{}]*"(?:tool|function|name)"[^{}]*\}`)
	matches := jsonPattern.FindString(content)
	if matches != "" {
		var call tools.Call
		if err := json.Unmarshal([]byte(matches), &call); err == nil && call.Name != "" {
			return &call, true
		}
	}

	// Pattern 2: XML-like tool call tags
	xmlPattern := regexp.MustCompile(`<tool_call[^>]*>([^<]+)</tool_call>`)
	xmlMatches := xmlPattern.FindStringSubmatch(content)
	if len(xmlMatches) > 1 {
		var call tools.Call
		if err := json.Unmarshal([]byte(xmlMatches[1]), &call); err == nil && call.Name != "" {
			return &call, true
		}
	}

	// Pattern 3: Simple function call format: tool_name(arg1="value1", arg2="value2")
	simplePattern := regexp.MustCompile(`(\w+)\s*\(([^)]*)\)`)
	simpleMatches := simplePattern.FindStringSubmatch(content)
	if len(simpleMatches) > 1 {
		call := &tools.Call{
			Name:      simpleMatches[1],
			Arguments: th.parseSimpleArgs(simpleMatches[2]),
		}
		if call.Name != "" {
			return call, true
		}
	}

	return nil, false
}

func (th *ToolHandler) parseSimpleArgs(argsStr string) map[string]interface{} {
	args := make(map[string]interface{})
	parts := strings.Split(argsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "="); idx > 0 {
			key := strings.TrimSpace(part[:idx])
			value := strings.TrimSpace(part[idx+1:])
			value = strings.Trim(value, `"'`)
			args[key] = value
		}
	}
	return args
}

func (th *ToolHandler) ExecuteTool(call *tools.Call) (interface{}, error) {
	tool, exists := th.Registry.Get(call.Name)
	if !exists {
		return nil, fmt.Errorf("tool '%s' not found", call.Name)
	}

	th.Logger.LogToolCall(call.Name, call.Arguments)

	result, err := tool.Execute(call.Arguments)
	th.Logger.LogToolResult(call.Name, result, err)

	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	return result, nil
}

func (th *ToolHandler) FormatToolResult(toolName string, result interface{}) string {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("Tool %s returned result (unable to serialize)", toolName)
	}
	return fmt.Sprintf("Tool %s result: %s", toolName, string(resultJSON))
}
