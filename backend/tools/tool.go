/**
 * Tool framework module.
 *
 * Defines the tool interface, registry, and execution engine for NIRA's
 * tool system. Tools are validated before execution and sandboxed based
 * on permissions.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: tool.go
 * Description: Tool interface and registry implementation.
 */

package tools

import (
	"encoding/json"
	"fmt"
)

type Tool interface {
	Name() string
	Description() string
	Execute(args map[string]interface{}) (interface{}, error)
	Schema() map[string]interface{}
}

type Registry struct {
	Tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		Tools: make(map[string]Tool),
	}
}

func (tr *Registry) Register(tool Tool) {
	tr.Tools[tool.Name()] = tool
}

func (tr *Registry) Get(name string) (Tool, bool) {
	tool, exists := tr.Tools[name]
	return tool, exists
}

func (tr *Registry) ListTools() []map[string]interface{} {
	var tools []map[string]interface{}
	for _, tool := range tr.Tools {
		tools = append(tools, tool.Schema())
	}
	return tools
}

type Call struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func ParseCall(jsonStr string) (*Call, error) {
	var call Call
	if err := json.Unmarshal([]byte(jsonStr), &call); err != nil {
		return nil, fmt.Errorf("failed to parse tool call: %w", err)
	}
	return &call, nil
}
