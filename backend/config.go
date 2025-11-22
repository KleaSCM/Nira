/**
 * Configuration management module.
 *
 * Handles loading and validation of configuration settings including
 * Ollama endpoint, model selection, database paths, and tool permissions.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: config.go
 * Description: Configuration loading and validation.
 */

package main

type Config struct {
	OllamaEndpoint string
	DefaultModel   string
	DatabasePath   string
	WebSocketPort  int
	AllowedPaths   []string
}

func LoadConfig() (Config, error) {
	return Config{
		OllamaEndpoint: "http://localhost:11434",
		DefaultModel:   "HammerAI/mythomax-l2",
		DatabasePath:   "./nira.db",
		WebSocketPort:  8080,
		AllowedPaths:   []string{},
	}, nil
}
