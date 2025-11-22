# Phase 2 — Memory System Design

## Overview

The NIRA memory system provides persistent storage for conversations, long-term context, and knowledge fragments. It uses SQLite as the storage backend, allowing for offline operation and easy portability.

## Design Goals

1. **Persistent Conversations**: Store all chat messages across application restarts
2. **Long-term Context**: Remember important facts, user preferences, and context
3. **Efficient Retrieval**: Fast lookup of past conversations and memories
4. **Scalability**: Support for future features like embeddings and vector search
5. **Isolation**: Separate storage for normal chat vs RP mode

## Database Schema

### Tables

#### `conversations`
Stores conversation sessions with metadata.

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER PRIMARY KEY | Unique conversation ID |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |
| title | TEXT | Optional conversation title |
| mode | TEXT | 'normal' or 'rp' |
| metadata | TEXT | JSON metadata (user preferences, etc.) |

#### `messages`
Stores individual chat messages within conversations.

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER PRIMARY KEY | Unique message ID |
| conversation_id | INTEGER | Foreign key to conversations |
| role | TEXT | 'user', 'assistant', 'system', 'tool' |
| content | TEXT | Message content |
| timestamp | TEXT | ISO 8601 timestamp |
| metadata | TEXT | JSON metadata (tool calls, etc.) |

#### `memories`
Stores long-term knowledge fragments and facts.

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER PRIMARY KEY | Unique memory ID |
| key | TEXT UNIQUE | Memory key/identifier |
| content | TEXT | Memory content |
| category | TEXT | Memory category (fact, preference, etc.) |
| created_at | TEXT | ISO 8601 timestamp |
| updated_at | TEXT | ISO 8601 timestamp |
| importance | INTEGER | Importance score (0-100) |

#### `embeddings`
Stores vector embeddings for semantic search (future use).

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER PRIMARY KEY | Unique embedding ID |
| source_type | TEXT | 'message', 'memory', 'file' |
| source_id | INTEGER | ID of source entity |
| embedding | BLOB | Vector embedding data |
| created_at | TEXT | ISO 8601 timestamp |

## Architecture

### Components

1. **Database Module** (`backend/memory/database.go`)
   - Database initialization
   - Schema creation and migrations
   - Connection management

2. **Conversation Store** (`backend/memory/conversation.go`)
   - Save/load conversations
   - Message persistence
   - Conversation management

3. **Memory Store** (`backend/memory/memory.go`)
   - Long-term memory CRUD operations
   - Memory retrieval by key/category
   - Importance-based filtering

4. **Memory Manager** (`backend/memory/manager.go`)
   - High-level memory operations
   - Integration with conversation flow
   - Memory summarization (future)

## Implementation Details

### Database Location
- Default: `./nira.db` (configurable via config)
- Created automatically on first run
- SQLite WAL mode for better concurrency

### Conversation Flow Integration
1. User sends message → Save to database
2. Assistant responds → Save to database
3. On startup → Load recent conversations
4. Memory extraction → Store important facts

### Memory Categories
- `fact`: General knowledge facts
- `preference`: User preferences
- `context`: Contextual information
- `tool_result`: Important tool execution results

## Future Enhancements

1. **Embeddings**: Vector search for semantic memory retrieval
2. **Compression**: Summarize old conversations to save space
3. **Memory Pruning**: Remove low-importance memories
4. **RP Isolation**: Separate tables for RP mode memories
5. **Memory Indexing**: Full-text search on memories

## API Design

### Conversation Operations
```go
CreateConversation() (int64, error)
GetConversation(id int64) (*Conversation, error)
ListConversations(limit, offset int) ([]*Conversation, error)
AddMessage(conversationID int64, role, content string) error
GetMessages(conversationID int64) ([]*Message, error)
```

### Memory Operations
```go
StoreMemory(key, content, category string, importance int) error
GetMemory(key string) (*Memory, error)
SearchMemories(query string, category string) ([]*Memory, error)
DeleteMemory(key string) error
```

## Migration Strategy

- Version 1: Initial schema (Phase 2)
- Future: Add embedding tables when needed
- Future: Add compression tables for summarized conversations

## Author

Author: KleaSCM  
Email: KleaSCM@gmail.com  
Date: 2025-11-22

