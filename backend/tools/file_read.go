/**
 * File read tool implementation.
 *
 * Provides safe file reading capability with path validation and
 * permission checking. Only reads files within allowed directories.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: file_read.go
 * Description: File reading tool with sandboxing.
 */

package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileReadTool struct {
    AllowedPaths []string
    checker      PathChecker
}

func NewFileReadTool(allowedPaths []string) *FileReadTool {
    return &FileReadTool{AllowedPaths: allowedPaths}
}

// NewFileReadToolWithChecker uses a centralized PathChecker; falls back to AllowedPaths if nil.
func NewFileReadToolWithChecker(allowedPaths []string, checker PathChecker) *FileReadTool {
    return &FileReadTool{AllowedPaths: allowedPaths, checker: checker}
}

func (t *FileReadTool) Name() string {
	return "read_file"
}

func (t *FileReadTool) Description() string {
	return "Reads the contents of a text file. Requires a file path as input."
}

func (t *FileReadTool) Execute(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path argument is required and must be a string")
	}

 if !t.isPathAllowed(path) {
        return nil, fmt.Errorf("path not in allowed directories")
    }

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return map[string]interface{}{
		"content": string(content),
		"path":    path,
	}, nil
}

func (t *FileReadTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.Name(),
		"description": t.Description(),
		"parameters": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The file path to read",
				},
			},
			"required": []string{"path"},
		},
	}
}

func (t *FileReadTool) isPathAllowed(path string) bool {
    if t.checker != nil {
        return t.checker.IsAllowed(path)
    }
    if len(t.AllowedPaths) == 0 {
        return false
    }

	absPath, err := filepath.Abs(path)
	if err != nil {
    return false
}

	for _, allowed := range t.AllowedPaths {
		absAllowed, err := filepath.Abs(allowed)
		if err != nil {
			continue
		}

		rel, err := filepath.Rel(absAllowed, absPath)
		if err != nil {
			continue
		}

		if rel != ".." && !filepath.IsAbs(rel) {
			return true
		}
	}

	return false
}
