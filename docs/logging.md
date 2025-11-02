# Logging with Zap

This project uses [Uber's Zap](https://github.com/uber-go/zap) as the logging framework for structured, high-performance logging.

## Configuration

The logger is automatically initialized in `main.go`. Configuration is controlled via environment variables:

### Environment Variables

- **`APP_ENV`**: Sets the logging environment
  - `production` - JSON formatted logs, Info level
  - Any other value (default) - Console formatted logs with colors, Debug level

- **`LOG_LEVEL`**: Overrides the default log level
  - Values: `debug`, `info`, `warn`, `error`, `fatal`
  - Example: `LOG_LEVEL=warn`

## Usage

### Basic Logging

```go
import (
    "github.com/tsongpon/athena/internal/logger"
    "go.uber.org/zap"
)

// Simple log messages
logger.Info("Server started")
logger.Debug("Processing request")
logger.Warn("Configuration missing, using defaults")
logger.Error("Failed to connect to database")
```

### Structured Logging

Add context to your logs using fields:

```go
logger.Info("User created",
    zap.String("user_id", userID),
    zap.String("email", email),
    zap.Time("created_at", createdAt))

logger.Error("Database operation failed",
    zap.String("operation", "INSERT"),
    zap.String("table", "users"),
    zap.Error(err))
```

### Common Field Types

- `zap.String(key, value)` - String fields
- `zap.Int(key, value)` - Integer fields
- `zap.Bool(key, value)` - Boolean fields
- `zap.Error(err)` - Error fields
- `zap.Time(key, value)` - Timestamp fields
- `zap.Duration(key, value)` - Duration fields
- `zap.Any(key, value)` - Any type (uses reflection, slower)

### Creating Child Loggers

Create loggers with preset fields for specific components:

```go
// Create a child logger with service context
serviceLogger := logger.With(
    zap.String("service", "bookmark"),
    zap.String("version", "1.0.0"))

// All logs from this logger will include the preset fields
serviceLogger.Info("Processing bookmark",
    zap.String("bookmark_id", id))
```

### Accessing the Logger Directly

```go
// Get the underlying zap.Logger instance
zapLogger := logger.Get()
zapLogger.Info("Direct access", zap.String("key", "value"))
```

## Log Levels

From least to most severe:

1. **Debug** - Detailed information for debugging
2. **Info** - General informational messages
3. **Warn** - Warning messages for potentially harmful situations
4. **Error** - Error messages for failures that don't stop execution
5. **Fatal** - Critical errors that cause application termination

## Examples from the Codebase

### Service Layer

```go
// internal/service/bookmark_service.go
logger.Info("Created bookmark",
    zap.String("id", createdBookmark.ID),
    zap.String("user_id", createdBookmark.UserID),
    zap.String("url", createdBookmark.URL),
    zap.String("title", createdBookmark.Title))
```

### Handler Layer

```go
// internal/handler/bookmark.go
logger.Error("Failed to create bookmark",
    zap.String("user_id", authenticatedUser.UserID),
    zap.String("url", bt.URL),
    zap.Error(err))
```

### Repository Layer

```go
// internal/repository/bookmark_inmem_repo.go
logger.Debug("Getting bookmark", zap.String("id", id))
```

## Development vs Production

### Development Mode (default)
- Human-readable console output with colors
- Includes stack traces for warnings and errors
- Debug level enabled
- Easier for local development

Example output:
```
2024-01-15T10:30:45.123+0700    INFO    service/bookmark_service.go:42    Created bookmark    {"id": "abc123", "user_id": "user1", "url": "https://example.com"}
```

### Production Mode (`APP_ENV=production`)
- JSON formatted output for log aggregation
- Info level by default
- Optimized for performance
- Easy to parse by log management systems

Example output:
```json
{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","caller":"service/bookmark_service.go:42","msg":"Created bookmark","id":"abc123","user_id":"user1","url":"https://example.com"}
```

## Best Practices

1. **Use appropriate log levels**
   - Debug: Detailed diagnostics
   - Info: Normal operations
   - Warn: Recoverable issues
   - Error: Failures that need attention
   - Fatal: Unrecoverable errors

2. **Include context with fields**
   ```go
   // Good
   logger.Error("Failed to save bookmark",
       zap.String("bookmark_id", id),
       zap.Error(err))
   
   // Avoid
   logger.Error("Failed to save bookmark: " + id + " " + err.Error())
   ```

3. **Don't log sensitive information**
   - Avoid logging passwords, tokens, credit cards
   - Be careful with PII (Personally Identifiable Information)

4. **Use consistent field names**
   - `user_id` not `userId` or `userID`
   - `bookmark_id` not `id` or `bookmarkId`

5. **Log at boundaries**
   - API requests/responses
   - Database operations
   - External service calls
   - Important business operations

## Testing

Logger tests are available in `internal/logger/logger_test.go`. The logger is designed to work in test environments and will not interfere with test output.
