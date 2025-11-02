# Athena Architecture

This document describes the architecture and design decisions of the Athena bookmark management API.

## Overview

Athena follows a clean architecture pattern with clear separation of concerns across multiple layers. The application is designed to be modular, testable, and easy to extend.

## Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                     HTTP Layer (Echo)                    │
│                    Port: 1323                            │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Handler Layer                         │
│  • BookmarkHandler (CRUD operations)                    │
│  • AuthHandler (login, registration)                    │
│  • JWT middleware integration                           │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Service Layer                         │
│  • BookmarkService (business logic)                     │
│  • UserService (authentication, user management)        │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                        │
│  • BookmarkInMemRepository                              │
│  • BookmarkPostgresRepository                           │
│  • UserInMemRepository                                  │
│  • WebRepository (title fetching)                       │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   Data Storage                           │
│  • In-Memory (Map)                                      │
│  • PostgreSQL Database                                  │
│  • External HTTP (websites)                             │
└─────────────────────────────────────────────────────────┘
```

## Components

### 1. Handler Layer (`internal/handler`)

**Responsibility**: HTTP request/response handling

**Components**:
- `BookmarkHandler`: Handles bookmark CRUD endpoints
- `AuthHandler`: Handles authentication and user registration
- `jwt_helper.go`: JWT token generation and validation utilities

**Key Features**:
- Request validation
- Response formatting
- JWT authentication middleware integration
- Error handling with appropriate HTTP status codes
- Structured logging for all operations

### 2. Service Layer (`internal/service`)

**Responsibility**: Business logic and orchestration

**Components**:
- `BookmarkService`: Bookmark business logic
- `UserService`: User management and authentication logic

**Key Features**:
- Input validation
- Business rule enforcement
- Repository coordination
- Password hashing (bcrypt)
- Authorization checks

### 3. Repository Layer (`internal/repository`)

**Responsibility**: Data access abstraction

**Components**:
- `BookmarkInMemRepository`: In-memory bookmark storage
- `BookmarkPostgresRepository`: PostgreSQL bookmark storage
- `UserInMemRepository`: In-memory user storage
- `WebRepository`: External HTTP title fetching

**Key Features**:
- Interface-based design for easy swapping
- Concurrent access safety (mutex for in-memory)
- Connection pooling (PostgreSQL)
- Automatic sorting (newest first)
- Indexed queries (PostgreSQL)

### 4. Model Layer (`internal/model`)

**Responsibility**: Domain models and data structures

**Components**:
- `Bookmark`: Bookmark entity
- `User`: User entity
- `BookmarkQuery`: Query parameters for filtering

### 5. Transport Layer (`internal/transport`)

**Responsibility**: API request/response DTOs

**Components**:
- `BookmarkTransport`: Bookmark API representation
- `LoginRequest/LoginResponse`: Authentication payloads
- `CreateUserRequest`: User registration payload
- `UserResponse`: User API representation

### 6. Database Layer (`internal/database`)

**Responsibility**: Database connection management

**Components**:
- `postgres.go`: PostgreSQL connection and migration utilities

**Key Features**:
- Connection pooling configuration
- Automatic schema migrations
- Health checks

### 7. Logger (`internal/logger`)

**Responsibility**: Centralized logging

**Components**:
- `logger.go`: Zap logger wrapper

**Key Features**:
- Development and production modes
- Structured logging
- Configurable log levels
- Performance optimized

## Storage Backends

### In-Memory Storage

**Use Cases**:
- Development and testing
- Demonstrations
- Stateless deployments (data lost on restart)

**Implementation**: `BookmarkInMemRepository`

**Features**:
- Thread-safe with RWMutex
- Instant operations (no I/O)
- Automatic UUID generation
- Sorted results

### PostgreSQL Storage

**Use Cases**:
- Production deployments
- Multi-instance deployments
- Data persistence requirements

**Implementation**: `BookmarkPostgresRepository`

**Features**:
- ACID transactions
- Connection pooling (25 max, 5 idle)
- Indexed queries
- Automatic migrations
- Timestamp tracking (created_at, updated_at)

**Schema**:
```sql
CREATE TABLE bookmarks (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## Authentication & Authorization

### JWT Authentication

- **Algorithm**: HS256
- **Token Expiration**: 24 hours
- **Claims**: UserID, Email, Name
- **Secret**: Configurable via `JWT_SECRET` environment variable

### Authorization

- **Endpoint Protection**: All bookmark endpoints require JWT
- **User Isolation**: Users can only access their own bookmarks
- **Check Location**: Handler layer validates user ownership

### Security Features

- Password hashing with bcrypt (cost 10)
- Generic error messages (prevent user enumeration)
- CORS middleware
- Request recovery middleware

## Data Flow

### Create Bookmark Flow

```
1. Client sends POST /bookmarks with URL
2. Handler validates JWT token
3. Handler extracts user ID from token
4. Handler calls BookmarkService.CreateBookmark()
5. Service calls WebRepository.GetTitle() to fetch title
6. Service calls BookmarkRepository.CreateBookmark()
7. Repository stores bookmark (in-memory or PostgreSQL)
8. Response flows back with created bookmark
```

### Authentication Flow

```
1. Client sends POST /login with email and password
2. Handler calls UserService.AuthenticateUser()
3. Service retrieves user by email
4. Service compares password hash with bcrypt
5. Handler generates JWT token with user claims
6. Response includes token and user info
```

## Design Patterns

### Repository Pattern

Abstracts data access behind interfaces, allowing easy switching between storage implementations.

```go
type BookmarkRepository interface {
    CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error)
    GetBookmark(id string) (model.Bookmark, error)
    ListBookmarks(query model.BookmarkQuery) ([]model.Bookmark, error)
    UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error)
    DeleteBookmark(id string) error
}
```

### Dependency Injection

Components receive their dependencies through constructors, making them testable and decoupled.

```go
func NewBookmarkService(
    bookmarkRepo BookmarkRepository,
    webRepo WebRepository,
) *BookmarkService
```

### Interface Segregation

Small, focused interfaces for specific functionality.

```go
type WebRepository interface {
    GetTitle(url string) (string, error)
}
```

## Configuration

### Environment Variables

- `STORAGE_TYPE`: Storage backend selection (memory, postgres)
- `DB_*`: PostgreSQL connection parameters
- `JWT_SECRET`: JWT signing secret
- `APP_ENV`: Environment mode (development, production)
- `LOG_LEVEL`: Logging level

### Defaults

Sensible defaults allow quick development startup:
- In-memory storage
- Development logging
- Default JWT secret (warning displayed)

## Error Handling

### Approach

- **Handler Layer**: Returns HTTP errors with appropriate status codes
- **Service Layer**: Returns Go errors with context
- **Repository Layer**: Returns errors for not found, conflicts, etc.

### Logging

- **Info**: Successful operations
- **Warn**: Validation failures, expected errors
- **Error**: Unexpected errors, database failures
- **Debug**: Detailed diagnostic information

## Testing Strategy

### Unit Tests

- **Handler Tests**: Mock service layer, test HTTP handling
- **Service Tests**: Mock repository layer, test business logic
- **Repository Tests**: Test data access (in-memory and PostgreSQL)

### Test Coverage

- Handler: 95.9%
- Repository: 100%
- Service: 98.4%

### Test Database

PostgreSQL tests can run against a test database when `TEST_DATABASE_URL` is set. Tests are skipped if not configured.

## Performance Considerations

### Database

- Indexed columns: user_id, is_archived, created_at
- Connection pooling: 25 max connections, 5 minute lifetime
- Prepared statements via parameterized queries

### Logging

- Structured logging (no string concatenation)
- Conditional compilation in production
- Buffered writes

### Concurrency

- Thread-safe in-memory repository with RWMutex
- Database connection pool handles concurrent requests
- Stateless handlers for horizontal scaling

## Scalability

### Horizontal Scaling

- Stateless application (JWT tokens, no sessions)
- Database-backed storage supports multiple instances
- No in-memory shared state

### Vertical Scaling

- Connection pool tunable
- Configurable log levels to reduce I/O
- Efficient data structures

## Security Best Practices

1. **Passwords**: Never logged, always hashed
2. **Secrets**: Configurable via environment, not hardcoded
3. **Authorization**: Checked at handler layer
4. **Input Validation**: Multiple layers (handler, service)
5. **Error Messages**: Generic to prevent information leakage
6. **SQL Injection**: Prevented via parameterized queries
7. **CORS**: Configurable middleware

## Future Architecture Considerations

### Planned Enhancements

- [ ] Event-driven architecture (bookmark events)
- [ ] Caching layer (Redis)
- [ ] Full-text search (PostgreSQL FTS or Elasticsearch)
- [ ] Async job processing (title fetching)
- [ ] API rate limiting
- [ ] Metrics and observability (Prometheus)
- [ ] Distributed tracing (OpenTelemetry)

### Extension Points

- Additional storage backends (MySQL, MongoDB)
- Additional authentication methods (OAuth, SAML)
- Webhook notifications
- Import/export functionality
- Bookmark sharing and permissions
