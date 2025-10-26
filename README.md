# Athena

A bookmark management API server built with Go and Echo framework.

## Overview

Athena is a lightweight RESTful API service for managing bookmarks. It provides functionality to store, retrieve, and organize bookmarks with support for archiving and user-specific collections.

## Features

- ✅ RESTful API for bookmark management
- ✅ Create, retrieve, and archive bookmarks
- ✅ User-specific bookmark collections
- ✅ Auto-generated UUID for bookmark IDs
- ✅ In-memory storage with repository pattern (easy to extend to database)
- ✅ Clean architecture with separation of concerns (handler, service, repository layers)
- ✅ Built-in logging and recovery middleware
- ✅ Comprehensive test coverage

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: [Echo v4](https://echo.labstack.com/)
- **ID Generation**: [Google UUID](https://github.com/google/uuid)

## Project Structure

```
athena/
├── cmd/
│   └── api-server/          # Application entry point
│       └── main.go
├── internal/
│   ├── handler/             # HTTP request handlers
│   │   └── http_handler.go
│   ├── service/             # Business logic layer
│   │   ├── bookmark_service.go
│   │   └── bookmark_service_test.go
│   ├── repository/          # Data access layer
│   │   ├── bookmark_inmem_repo.go
│   │   └── bookmark_inmem_repo_test.go
│   ├── transport/           # HTTP transport layer (DTOs)
│   │   └── bookmark_transport.go
│   └── model/               # Domain models
│       └── bookmark.go
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.25.1 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/tsongpon/athena.git
cd athena
```

2. Install dependencies:
```bash
go mod download
```

### Running the Server

```bash
go run cmd/api-server/main.go
```

The server will start on `http://localhost:1323`

### Building

```bash
go build -o athena cmd/api-server/main.go
./athena
```

## API Endpoints

### Health Check
- `GET /ping` - Health check endpoint
  - Response: `pong`

### Bookmark Management

#### Create Bookmark
- `POST /bookmarks`
  - Request body:
    ```json
    {
      "url": "https://example.com",
      "user_id": "user-123"
    }
    ```
  - Response: `201 Created`

#### Get Single Bookmark
- `GET /bookmarks/:id`
  - Response:
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "url": "https://example.com",
      "title": "",
      "user_id": "user-123",
      "created_at": "2025-10-26T14:00:00Z"
    }
    ```

#### Get All Bookmarks (by User)
- `GET /bookmarks?userid=user-123`
  - Query parameter: `userid` - Filter bookmarks by user ID
  - Response:
    ```json
    [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "url": "https://example.com",
        "title": "",
        "user_id": "user-123",
        "created_at": "2025-10-26T14:00:00Z"
      }
    ]
    ```

#### Archive Bookmark
- `POST /bookmarks/:id/archive`
  - Response: `204 No Content`

## Data Model

### Bookmark

```go
type Bookmark struct {
    ID         string    // Auto-generated UUID
    UserID     string    // User identifier
    URL        string    // Bookmark URL
    Title      string    // Bookmark title (reserved for future use)
    IsArchived bool      // Archive status
    CreatedAt  time.Time // Creation timestamp
}
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests verbosely:
```bash
go test -v ./...
```

Run tests for a specific package:
```bash
go test ./internal/service/
go test ./internal/repository/
```

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

### Layers

1. **Handler Layer** (`internal/handler/`)
   - Handles HTTP requests and responses
   - Validates request parameters
   - Transforms between transport objects and domain models

2. **Service Layer** (`internal/service/`)
   - Contains business logic
   - Validates business rules (e.g., bookmark ID must be empty on creation)
   - Orchestrates repository operations

3. **Repository Layer** (`internal/repository/`)
   - Manages data persistence
   - Currently implements in-memory storage
   - Interface-based design for easy extension to databases

4. **Transport Layer** (`internal/transport/`)
   - Defines DTOs (Data Transfer Objects) for HTTP API
   - Separates internal domain models from external API contracts

5. **Model Layer** (`internal/model/`)
   - Contains domain models
   - Pure data structures with no business logic

### Design Patterns

- **Repository Pattern**: Abstracts data access logic
- **Dependency Injection**: Services and handlers receive dependencies via constructors
- **Interface Segregation**: Small, focused interfaces (e.g., `BookmarkRepository`)

## Future Enhancements

- [ ] Database persistence (PostgreSQL, MySQL)
- [ ] Bookmark title fetching from URL
- [ ] Full-text search
- [ ] Tagging system
- [ ] Authentication and authorization
- [ ] Rate limiting
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Docker support
- [ ] CI/CD pipeline

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
