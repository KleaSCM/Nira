/**
 * Logging module.
 *
 * Provides structured logging functionality for the NIRA backend,
 * including different log levels and optional file output.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: logger.go
 * Description: Logging system implementation.
 */

package main

import (
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Logger struct {
	Level  LogLevel
	Logger *log.Logger
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{
		Level:  level,
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.Level <= LogLevelDebug {
		l.Logger.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.Level <= LogLevelInfo {
		l.Logger.Printf("[INFO] "+format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.Level <= LogLevelWarn {
		l.Logger.Printf("[WARN] "+format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.Level <= LogLevelError {
		l.Logger.Printf("[ERROR] "+format, v...)
	}
}

func (l *Logger) LogToolCall(toolName string, args map[string]interface{}) {
	l.Info("Tool call: %s with args: %v", toolName, args)
}

func (l *Logger) LogToolResult(toolName string, result interface{}, err error) {
	if err != nil {
		l.Error("Tool %s failed: %v", toolName, err)
	} else {
		l.Info("Tool %s completed successfully", toolName)
	}
}

func (l *Logger) LogWebSocketEvent(event string, details string) {
	l.Debug("WebSocket %s: %s", event, details)
}

func (l *Logger) LogOllamaRequest(model string, messageCount int) {
	l.Debug("Ollama request: model=%s, messages=%d", model, messageCount)
}

func (l *Logger) LogOllamaResponse(duration time.Duration, chunkCount int) {
	l.Debug("Ollama response: duration=%v, chunks=%d", duration, chunkCount)
}
