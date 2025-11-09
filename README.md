# Athena

A secure bookmark management API server built with Go and Echo framework with JWT authentication.

## Overview

Athena is a lightweight RESTful API service for managing bookmarks. It provides functionality to store, retrieve, and organize bookmarks with support for archiving and user-specific collections. All bookmark endpoints are protected with JWT authentication to ensure users can only access their own bookmarks.

## Features

- ✅ JWT-based authentication and authorization
- ✅ User registration and login
- ✅ RESTful API for bookmark management
- ✅ Create, retrieve, and archive bookmarks
- ✅ User-specific bookmark collections with access control
- ✅ Auto-generated UUID for bookmark and user IDs
- ✅ Password hashing with bcrypt
- ✅ Structured logging with Uber's Zap (development and production modes)
- ✅ Automatic website title fetching for bookmarks
- ✅ Pluggable storage backends: In-memory and PostgreSQL
- ✅ Database migrations with automatic schema creation
- ✅ Clean architecture with separation of concerns (handler, service, repository layers)
- ✅ Built-in logging, CORS, and recovery middleware
- ✅ Docker support with multi-stage builds
- ✅ Comprehensive test coverage (95.9% handler, 100% repository, 98.4% service)

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: [Echo v4](https://echo.labstack.com/)
- **Authentication**: JWT with [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) & [echo-jwt/v4](https://github.com/labstack/echo-jwt)
- **Database**: [PostgreSQL](https://www.postgresql.org/) with [lib/pq](https://github.com/lib/pq) driver
- **Logging**: [Uber Zap](https://github.com/uber-go/zap)
- **Password Hashing**: bcrypt
- **ID Generation**: [Google UUID](https://github.com/google/uuid)
- **HTML Parsing**: [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html)
- **Testing**: [testify](https://github.com/stretchr/testify)

## Project Structure

```
athena/
├── cmd/
│   └── api-server/
│       └── main.go                      # Application entry point, route setup
├── internal/
│   ├── handler/                         # HTTP request handlers
│   │   ├── auth.go                      # Authentication handlers (login, create user)
│   │   ├── auth_test.go                 # Authentication handler tests (14 tests)
│   │   ├── bookmark.go                  # Bookmark handlers (CRUD operations)
│   │   ├── bookmark_test.go             # Bookmark handler tests (38 tests)
│   │   ├── jwt_helper.go                # JWT generation, validation, extraction
│   │   └── service.go                   # Service interfaces for handlers
│   ├── service/                         # Business logic layer
│   │   ├── bookmark_service.go          # Bookmark business logic
│   │   ├── bookmark_service_test.go     # Bookmark service tests
│   │   ├── user_service.go              # User authentication & management
│   │   ├── user_service_test.go         # User service tests
│   │   └── repository.go                # Repository interfaces (UserRepository, BookmarkRepository)
│   ├── repository/                      # Data access layer
│   │   ├── bookmark_inmem_repo.go       # In-memory bookmark storage
│   │   ├── bookmark_inmem_repo_test.go  # Bookmark repository tests
│   │   ├── user_inmem_repo.go           # In-memory user storage
│   │   └── user_inmem_repo_test.go      # User repository tests
│   ├── transport/                       # HTTP transport layer (DTOs)
│   │   ├── bookmark_transport.go        # Bookmark request/response DTOs
│   │   └── user_transport.go            # User request/response DTOs
│   └── model/                           # Domain models
│       ├── bookmark.go                  # Bookmark domain model
│       └── user.go                      # User domain model
├── go.mod                               # Go module dependencies
└── README.md                            # This file
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

3. Configure environment variables (all optional):
```bash
# JWT secret (defaults to development secret)
export JWT_SECRET="your-super-secret-key-change-this-in-production"

# Storage configuration (defaults to in-memory)
export STORAGE_TYPE="postgres"  # Options: memory (default), postgres

# PostgreSQL configuration (only needed if STORAGE_TYPE=postgres)
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="your_password"
export DB_NAME="athena"
export DB_SSLMODE="disable"  # Use "require" in production

# Logging configuration
export APP_ENV="production"  # Use "production" for JSON logs, default is development
export LOG_LEVEL="info"      # Options: debug, info, warn, error, fatal
```

### Running the Server

#### Option 1: Run Locally with Go

```bash
go run cmd/api-server/main.go
```

The server will start on `http://localhost:1323`

#### Option 2: Run with Docker

```bash
# Build and run with docker-compose
docker-compose up -d

# Or build and run manually
docker build -t athena:latest .
docker run -p 1323:1323 -e JWT_SECRET="your-secret-key" athena:latest
```

The server will be available at `http://localhost:1323`

To stop the Docker container:
```bash
docker-compose down
```

### Building

#### Build Native Binary

```bash
go build -o athena cmd/api-server/main.go
./athena
```

Or build from the cmd directory:
```bash
cd cmd/api-server
go build -o ../../athena
cd ../..
./athena
```

#### Build Docker Image

```bash
docker build -t athena:latest .
```

## API Endpoints

### Public Endpoints

#### Public Endpoints

####HealthCheck
- **GET** `/ping`
  - Response: `pong` (200 OK)
  - Purpose: Verify server is running#### UserUser RegistrationRegistration
- **POST** `/usersusers`
  - Request body:
    ```json
    {
      "namename": "John Doe",
      "email: "John Doe",
      "email": "john@example.com",
      "passwordpassword": "securepassword123securepassword123"
    }
    ```
  - Response: `201 Created`
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2025-11-02T14:00:00Z",
      "updated_at": "2025-11-02T14:00:00Z"
    }
    ```
  - Errors:
    - `400` - Name, email, or password missing
    - `409` - Email already exists

#### User Login
- `POST /login`
  - Request body:
    ```json
    {
      "email": "john@example.com",
      "password": "securepassword123"
    }
    ```
  - Response: `200 OK`
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "token_type": "Bearer",
      "expires_in": 86400,
      "user": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": "2025-11-02T14:00:00Z",
        "updated_at": "2025-11-02T14:00:00Z"
      }
    }
    ```

### Protected Endpoints (Require JWT Authentication)

All bookmark endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

#### Create Bookmark
- `POST /bookmarks`
  - Headers: `Authorization: Bearer <token>`
  - Request body:
    ```json
    {
      "url": "https://example.com"
    }
    ```
  - Response: `201 Created`
#### User Login
- **POST** `/login`
  - Request body:
    ```json
    {
      "email": "john@example.com",
      "password": "securepassword123"
    }
    ```
  - Response: `200 OK`
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "token_type": "Bearer",
      "expires_in": 86400,
      "user": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": "2025-11-02T14:00:00Z",
        "updated_at": "2025-11-02T14:00:00Z"
      }
    }
    ```
  - Errors:
    - `400` - Email or password missing
    - `401` - Invalid credentials

### Protected Endpoints (Require JWT Authentication)

All bookmark endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

#### Create Bookmark
- **POST** `/bookmarks`
  - Headers: `Authorization: Bearer <token>`
  - Request body:
    ```json
    {
      "url": "https://example.com"
    }
    ```
  - Response: `201 Created`
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "url": "https://example.com",
      "title": "",
      "user_id": "user-id-from-jwtid-from-jwt",
      "is_archived": false
    }
    ```
  - Note: `user_id` is automatically extracted from the JWT token
  - Errors:
    - `400` - URL is missing
    - `401` - Invalid or missing JWT token

#### Get Single Bookmark
- `GET /bookmarks/:id`
  - Headers: `Authorization: Bearer <token>`
  - Response: `200 OK`
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "url": "https://example.com",
      "title": "",
      "user_id": "user-id-from-jwt",
      "is_archived": false
    }
    ```
  - Returns `403 Forbidden` if the bookmark belongs to a different user

#### Get All Bookmarks
- `GET /bookmarks?archived=false`
  - Headers: `Authorization: Bearer <token>`
  - Query parameters:
    - `archived` (optional): `true` or `false` (default: `false`)
  - Response: `200 OK`
#### Get Single Bookmark
- **GET** `/bookmarks/:id`
  - Headers: `Authorization: Bearer <token>`
  - URL Parameters: `id` - Bookmark UUID
  - Response: `200 OK`
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "url": "https://example.com",
      "title": "",
      "user_id": "user-id-from-jwt",
      "is_archived": false
    }
    ```
  - Errors:
    - `400` - ID is missing
    - `401` - Invalid or missing JWT token
    - `403` - Bookmark belongs to a different user
    - `404` - Bookmark not found

#### Get All Bookmarks
- **GET** `/bookmarks?archived=false`
  - Headers: `Authorization: Bearer <token>`
  - Query parameters:
    - `archived` (optional): `true` or `false` (default: `false`)
  - Response: `200 OK`
    ```json
    [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "url": "https://example.com",
        "title": "",
        "user_id": "user-id-from-jwtid-from-jwt",
        "created_at": "2025-1111-02T1402T14:00:00Z",
        "is_archived": false
      }
    ]
    ```
  - Note: Only returns bookmarks for the authenticated user
  - Errors:
    - `401` - Invalid or missing JWT token

#### Archive Bookmark
- **POST** `/bookmarks/:id/archive`
  - Headers: `Authorization: Bearer <token>`
  - URL Parameters: `id` - Bookmark UUID
  - Headers: `Authorization: Bearer <token>`
  - Response: `204 No Content`
  - Errors:
    - `400` - ID is missing
    - `401` - Invalid or missing JWT token
    - `403` - Bookmark belongs to a different user
    - `404` - Bookmark not found

## Quick Start Example

```bash
# 1. Start the server
./athena

# 2. Create a user
curl -X POST http://localhost:1323/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# 3. Login to get JWT token
TOKEN=$(curl -s -X POST http://localhost:1323/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.token')

# 4. Create a bookmark
curl -X POST http://localhost:1323/bookmarks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url":"https://example.com"}'

# 5. Get all bookmarks
curl -X GET http://localhost:1323/bookmarks \
  -H "Authorization: Bearer $TOKEN"
```

## Security Features

### Authentication
- JWT-based stateless authentication
- Token expiration: 24 hours
- Secure password hashing with bcrypt (cost factor 10)
- Maximum password length: 72 bytes (bcrypt limitation)

### Authorization
- Users can only access their own bookmarks
- User ID is extracted from JWT token (not from request parameters)
- Authorization checks on all bookmark operations:
  - Get bookmark: Verifies ownership before returning
  - Archive bookmark: Verifies ownership before archiving
  - List bookmarks: Automatically filtered by authenticated user

### Password Requirements
- Minimum length: Not enforced (consider adding)
- Maximum length: 72 bytes
- Stored as bcrypt hash (never plaintext)

### Environment Variables
- `JWT_SECRET`: Secret key for signing JWT tokens
  - Default: `your-secret-key-change-this-in-production` (⚠️ CHANGE IN PRODUCTION!)
  - Recommended: Use a strong random string (at least 32 characters)

## Data Models

### User
```go
type User struct {
    ID        string    // Auto-generated UUID
    Name      string    // User's full name
    Email     string    // User's email (unique)
    Password  string    // Bcrypt hashed password
    CreatedAt time.Time // Registration timestamp
    UpdatedAt time.Time // Last update timestamp
}
```
## Quick Start Example

```bash
# 1. Start the server
./athena

# 2. Create a user
curl -X POST http://localhost:1323/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# 3. Login to get JWT token
TOKEN=$(curl -s -X POST http://localhost:1323/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.token')

echo "Token: $TOKEN"

# 4. Create a bookmark
curl -X POST http://localhost:1323/bookmarks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url":"https://example.com"}'

# 5. Get all bookmarks
curl -X GET http://localhost:1323/bookmarks \
  -H "Authorization: Bearer $TOKEN"

# 6. Get specific bookmark (replace BOOKMARK_ID)
curl -X GET http://localhost:1323/bookmarks/BOOKMARK_ID \
  -H "Authorization: Bearer $TOKEN"

# 7. Archive a bookmark (replace BOOKMARK_ID)
curl -X POST http://localhost:1323/bookmarks/BOOKMARK_ID/archive \
  -H "Authorization: Bearer $TOKEN"

# 8. Get archived bookmarks
curl -X GET http://localhost:1323/bookmarks?archived=true \
  -H "Authorization: Bearer $TOKEN"
```

## Security Features

### Authentication
- **JWT-based stateless authentication**
- Token expiration: 24 hours (86400 seconds)
- Secure password hashing with bcrypt (cost factor 10)
- Maximum password length: 72 bytes (bcrypt limitation)
- Token signed with HMAC-SHA256

### Authorization
- Users can only access their own bookmarks
- User ID is extracted from JWT token (not from request parameters)
- Authorization checks on all bookmark operations:
  - **Get bookmark**: Verifies ownership before returning (403 if unauthorized)
  - **Archive bookmark**: Verifies ownership before archiving (403 if unauthorized)
  - **List bookmarks**: Automatically filtered by authenticated user
  - **Create bookmark**: User ID automatically set from JWT token

### JWT Token Structure
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "name": "John Doe",
  "exp": 1730649600,
  "iat": 1730563200,
  "nbf": 1730563200,
  "iss": "athena",
  "sub": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Environment Variables
- `JWT_SECRET`: Secret key for signing JWT tokens
  - Default: `your-secret-key-change-this-in-production` (⚠️ **MUST CHANGE IN PRODUCTION!**)
  - Recommended: Use a strong random string (at least 32 characters)
  - Example: `export JWT_SECRET="$(openssl rand -base64 32)"`

## Data Models

### User
```go
type User struct {
    ID        string    // Auto-generated UUID
    Name      string    // User's full name
    Email     string    // User's email (unique)
    Password  string    // Bcrypt hashed password
    CreatedAt time.Time // Registration timestamp
    UpdatedAt time.Time // Last update timestamp
}
```

### Bookmark
```go
type Bookmark struct {
    ID         string    // Auto-generated UUID
    UserID     string    // User identifier (owner)
    URL        string    // Bookmark URL
    Title      string    // Bookmark title (reserved for future use)
    IsArchived bool      // Archive status
    CreatedAt  time.Time // Creation timestamp
}
```

### JWT Claims
```go
type JWTClaims struct {
    UserID string // User ID
    Email  string // User email
    Name   string // User name
    jwt.RegisteredClaims
}
```

### JWT Claims
```go
type JWTClaims struct {
    UserID string `json:"user_id"` // User ID
    Email  string `json:"email"`   // User email
    Name   string `json:"name"`    // User name
    jwt.RegisteredClaims             // Standard JWT claims (exp, iat, nbf, iss, sub)
}
```

## Docker

The application includes Docker support for easy deployment with PostgreSQL database.

### Quick Start with Docker Compose

```bash
# Start both API and PostgreSQL
docker-compose up -d

# Check logs
docker-compose logs -f

# Check specific service logs
docker-compose logs -f athena
docker-compose logs -f postgres

# Stop all services
docker-compose down

# Stop and remove volumes (deletes database data)
docker-compose down -v
```

### Docker Configuration

#### Environment Variables

Set environment variables via `.env` file:

```bash
# Copy example file
cp .env.example .env

# Edit .env and configure your settings
nano .env
```

**Important variables**:
- `JWT_SECRET` - Change in production!
- `DB_PASSWORD` - Database password
- `STORAGE_TYPE` - Set to `postgres` (default)
- `LOG_LEVEL` - Logging verbosity

#### Docker Compose Services

The `docker-compose.yml` includes:

**PostgreSQL Service**:
- Image: `postgres:16-alpine`
- Port: 5432
- Persistent storage with named volume
- Health checks for reliable startup
- Configurable via environment variables

**Athena API Service**:
- Built from Dockerfile
- Port: 1323
- Automatic database migrations on startup
- Waits for PostgreSQL to be healthy
- Health checks via `/ping` endpoint
- Auto-restart unless stopped
- JSON logging in production mode

**Network**:
- Custom bridge network for service communication
- Services communicate via service names

#### Alternative: In-Memory Storage

For development without PostgreSQL:

```bash
# Use in-memory configuration
docker-compose -f docker-compose.memory.yml up -d
```

### Building the Docker Image

```bash
# Build image
docker build -t athena:latest .

# Build with custom tag
docker build -t athena:v1.0.0 .

# Build with no cache
docker build --no-cache -t athena:latest .
```

### Running the Docker Container

```bash
# Run in foreground
docker run -p 1323:1323 athena:latest

# Run in background (detached)
docker run -d -p 1323:1323 --name athena-api athena:latest

# Run with custom JWT secret
docker run -d -p 1323:1323 \
  -e JWT_SECRET="$(openssl rand -base64 32)" \
  --name athena-api \
  athena:latest

# Run with volume mount (for future file storage)
docker run -d -p 1323:1323 \
  -v $(pwd)/data:/app/data \
  --name athena-api \
  athena:latest
```

### Docker Container Management

```bash
# View running containers
docker ps

# View all containers (including stopped)
docker ps -a

# View logs
docker-compose logs athena
docker-compose logs postgres
docker-compose logs -f  # Follow all logs

# Stop services
docker-compose stop

# Start services
docker-compose start

# Restart services
docker-compose restart

# Remove containers
docker-compose down

# Remove containers and volumes (deletes database)
docker-compose down -v
```

### Database Management

#### Access PostgreSQL

```bash
# Connect to PostgreSQL container
docker exec -it athena-postgres psql -U athena_user -d athena

# Run SQL query
docker exec -it athena-postgres psql -U athena_user -d athena -c "SELECT COUNT(*) FROM bookmarks;"

# Dump database
docker exec athena-postgres pg_dump -U athena_user athena > backup.sql

# Restore database
docker exec -i athena-postgres psql -U athena_user athena < backup.sql
```

#### View Database Data

```bash
# List all bookmarks
docker exec -it athena-postgres psql -U athena_user -d athena -c "SELECT id, title, url FROM bookmarks;"

# Check table structure
docker exec -it athena-postgres psql -U athena_user -d athena -c "\d bookmarks"
```

#### Database Volume Management

```bash
# View volumes
docker volume ls

# Inspect volume
docker volume inspect athena_postgres_data

# Backup volume
docker run --rm -v athena_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz -C /data .

# Restore volume
docker run --rm -v athena_postgres_data:/data -v $(pwd):/backup alpine sh -c "cd /data && tar xzf /backup/postgres_backup.tar.gz"
```

### Health Check

The Docker image includes a built-in health check:

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' athena-api

# View health check logs
docker inspect --format='{{json .State.Health}}' athena-api | jq
```

Health check endpoint: `GET /ping`
- Interval: 30 seconds
- Timeout: 3 seconds
- Retries: 3
- Start period: 5 seconds

### Multi-stage Build

The Dockerfile uses a multi-stage build for optimal image size:

1. **Builder stage**: Uses `golang:1.25.1-alpine` to compile the binary
2. **Runtime stage**: Uses `alpine:latest` with only the compiled binary

Benefits:
- Small image size (~20MB vs ~800MB with full Go image)
- No Go toolchain in final image (security)
- Static binary with no external dependencies
- Runs as non-root user for security

### Docker Best Practices Implemented

- ✅ Multi-stage build for minimal image size
- ✅ Non-root user (user `athena`, UID 1000)
- ✅ Health checks for container monitoring
- ✅ `.dockerignore` to exclude unnecessary files
- ✅ Static binary (CGO_ENABLED=0)
- ✅ Security: CA certificates for HTTPS, timezone data
- ✅ Optimized layer caching (dependencies before source)

### Troubleshooting

#### Container won't start
```bash
# Check logs
docker logs athena-api

# Check if port is already in use
lsof -i :1323
```

#### Health check failing
```bash
# Test health endpoint manually
curl http://localhost:1323/ping

# Check container logs
docker logs athena-api
```

#### Permission issues
```bash
# Ensure volumes have correct permissions
chmod -R 755 ./data
```

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

Expected output:
```
ok      github.com/tsongpon/athena/internal/handler     0.271s  coverage: 95.9% of statements
ok      github.com/tsongpon/athena/internal/repository  0.003s  coverage: 100.0% of statements
ok      github.com/tsongpon/athena/internal/service     0.050s  coverage: 98.4% of statements
```

### Run Tests Verbosely
```bash
go test -v ./...
```

### Run Specific Package Tests
```bash
# Handler tests (65 tests total)
go test -v ./internal/handler/

# Service tests
go test -v ./internal/service/

# Repository tests
go test -v ./internal/repository/
```

Run specific test:
```bash
go test -v -run TestAuthHandler_Login ./internal/handler/
go test -v -run TestBookmarkHandler_CreateBookmark ./internal/handler/
```

### Run Specific Test
```bash
# Auth handler tests
go test -v -run TestAuthHandler_Login ./internal/handler/

# Bookmark handler tests
go test -v -run TestBookmarkHandler_CreateBookmark ./internal/handler/

# JWT extraction tests
go test -v -run TestBookmarkHandler_GetBookmark_Forbidden ./internal/handler/
```

### Test Coverage Breakdown
- **Handler Layer**: 95.9% (65 tests total)
  - `auth.go`: Covered by 14 tests
  - `bookmark.go`: Covered by 38 tests
  - `jwt_helper.go`: Covered by handler tests
  - Other tests: 13 tests (marshaling, benchmarks, helpers)
- **Repository Layer**: 100.0%
  - `bookmark_inmem_repo.go`: Full coverage
  - `user_inmem_repo.go`: Full coverage
- **Service Layer**: 98.4%
  - `bookmark_service.go`: Nearly complete coverage
  - `user_service.go`: Nearly complete coverage
- **Overall**: ~89%

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

### Layers

1. **Handler Layer** (`internal/handler/`)
   - **Purpose**: HTTP request/response handling
   - **Responsibilities**:
     - Validate HTTP request parameters
     - Extract JWT claims from context
     - Perform authorization checks
     - Transform between transport objects and domain models
     - Return HTTP responses
   - **Files**:
     - `auth.go`: Login, user creation
     - `bookmark.go`: Bookmark CRUD operations
     - `jwt_helper.go`: JWT generation, validation, extraction helper
     - `service.go`: Service interfaces used by handlers

2. **Service Layer** (`internal/service/`)
   - **Purpose**: Business logic and orchestration
   - **Responsibilities**:
     - Enforce business rules
     - Password hashing and authentication
     - Bookmark validation (e.g., ID must be empty on creation)
     - Orchestrate repository operations
   - **Files**:
     - `user_service.go`: User authentication and management
     - `bookmark_service.go`: Bookmark business logic
     - `repository.go`: UserRepository interface
     - `bookmark_repository.go`: BookmarkRepository interface

3. **Repository Layer** (`internal/repository/`)
   - **Purpose**: Data persistence abstraction
   - **Responsibilities**:
     - CRUD operations on data store
     - Currently in-memory implementation
     - Interface-based for easy database migration
   - **Files**:
     - `user_inmem_repo.go`: In-memory user storage
     - `bookmark_inmem_repo.go`: In-memory bookmark storage

4. **Transport Layer** (`internal/transport/`)
   - **Purpose**: HTTP API contracts (DTOs)
   - **Responsibilities**:
     - Define request/response structures
     - Separate external API from internal domain
   - **Files**:
     - `user_transport.go`: LoginRequest, CreateUserRequest, UserResponse, LoginResponse
     - `bookmark_transport.go`: BookmarkTransport

5. **Model Layer** (`internal/model/`)
   - **Purpose**: Domain models
   - **Responsibilities**:
     - Pure data structures
     - Minimal business logic
   - **Files**:
     - `user.go`: User domain model
     - `bookmark.go`: Bookmark domain model, BookmarkQuery

### Data Flow

```
HTTP Request
    ↓
[Handler] ← Validates request, extracts JWT claims
    ↓
[Service] ← Business logic, password hashing
    ↓
[Repository] ← Data persistence
    ↓
[In-Memory Store]
```

### Design Patterns

- **Repository Pattern**: Abstracts data access logic
- **Dependency Injection**: Services and handlers receive dependencies via constructors
- **Interface Segregation**: Small, focused interfaces (e.g., `BookmarkRepository`, `UserRepository`)
- **Middleware Pattern**: JWT authentication, CORS, logging, recovery
- **Helper Functions**: Reusable JWT claim extraction logic

### Security Architecture

- **Stateless Authentication**: JWT tokens contain all necessary user information
- **Authorization at Handler Level**: Each protected endpoint verifies user ownership
- **Secure Password Storage**: bcrypt hashing with salt (never store plaintext)
- **Defense in Depth**: Multiple layers of validation (handler, service, repository)

## Error Handling

The API returns standard HTTP status codes:

- `200 OK` - Successful GET request
- `201 Created` - Successful POST request (resource created)
- `204 No Content` - Successful POST request (no content returned)
- `400 Bad Request` - Invalid request parameters or body
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - Valid token but insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists (e.g., duplicate email)
- `500 Internal Server Error` - Server-side error

Error response format:
```json
{
  "message": "Error description"
}
```
- **Interface Segregation**: Small, focused interfaces (BookmarkRepository, UserRepository, UserService, BookmarkService)
- **Middleware Pattern**: JWT authentication, CORS, logging, recovery
- **Helper Functions**: Reusable JWT claim extraction logic (`getAuthenticatedUser`)

### Key Design Decisions

1. **JWT in Context**: Echo JWT middleware v4 stores `*jwt.Token` in context. The `getAuthenticatedUser` helper extracts claims safely.

2. **User ID from JWT**: User ID is never accepted from request parameters for protected endpoints, always extracted from authenticated token.

3. **Authorization at Handler**: Authorization checks happen at the handler layer before calling services.

4. **Service Interfaces in Handler Package**: Service interfaces are defined in `internal/handler/service.go` to avoid circular dependencies.

5. **Repository Interfaces in Service Package**: Repository interfaces are defined in the service package where they're consumed.

## Error Handling

The API returns standard HTTP status codes:

- `200 OK` - Successful GET request
- `201 Created` - Successful POST request (resource created)
- `204 No Content` - Successful POST request (no content returned, e.g., archive)
- `400 Bad Request` - Invalid request parameters or body
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - Valid token but insufficient permissions (accessing other user's resources)
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists (e.g., duplicate email)
- `500 Internal Server Error` - Server-side error

Error response format:
```json
{
  "message": "Error description"
}
```

## PostgreSQL Storage

Athena supports PostgreSQL as a persistent storage backend. See [docs/postgresql.md](docs/postgresql.md) for detailed setup instructions.

### Quick Setup

```bash
# Create database
createdb athena

# Configure environment
export STORAGE_TYPE=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=athena

# Run server (migrations run automatically)
go run ./cmd/api-server
```

### Features

- ✅ Automatic database migrations
- ✅ Connection pooling
- ✅ Indexed queries for performance
- ✅ Ordered results (newest first)
- ✅ Transaction support

For production setup, Docker Compose configuration, and troubleshooting, see the [PostgreSQL documentation](docs/postgresql.md).

## Logging

Athena uses [Uber Zap](https://github.com/uber-go/zap) for structured, high-performance logging.

### Configuration

- **Development Mode** (default): Console output with colors, debug level
- **Production Mode** (`APP_ENV=production`): JSON formatted logs, info level
- **Custom Log Level**: Set `LOG_LEVEL` environment variable (debug, info, warn, error, fatal)

### Example Logs

**Development:**
```
2024-01-15T10:30:45.123+0700    INFO    service/bookmark_service.go:42    Created bookmark    {"id": "abc123", "user_id": "user1"}
```

**Production:**
```json
{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","msg":"Created bookmark","id":"abc123","user_id":"user1"}
```

For detailed logging documentation, see [docs/logging.md](docs/logging.md).

## CI/CD Pipeline

Athena uses GitHub Actions for continuous integration and deployment. The workflow automatically tests, builds, and deploys your application to Google Cloud Run.

### Workflow Overview

The CI/CD pipeline consists of three jobs:

1. **Test** - Runs on all pushes and pull requests
   - Runs unit tests with race detection
   - Generates code coverage reports
   - Uploads coverage to Codecov
   - Enforces 80% minimum coverage threshold

2. **Build** - Runs only on pushes to main/master branch
   - Builds Docker image
   - Pushes to GCP Artifact Registry
   - Tags images with branch name, SHA, and `latest`

3. **Deploy** - Runs only on pushes to main/master branch
   - Deploys to Google Cloud Run
   - Configures auto-scaling and environment variables
   - Outputs the deployed service URL

### Required GitHub Secrets

Configure these secrets in your GitHub repository settings (Settings → Secrets and variables → Actions):

- **`GCP_SA_KEY`** - GCP Service Account JSON key with the following permissions:
  - `roles/artifactregistry.writer` - Push Docker images
  - `roles/run.admin` - Deploy to Cloud Run
  - `roles/iam.serviceAccountUser` - Act as service account
  
- **`GCP_PROJECT_ID`** - Your GCP project ID (e.g., `my-project-123`)

- **`GCP_REGION`** - GCP region for deployment (e.g., `us-central1`, `asia-southeast1`)

- **`GCP_ARTIFACT_REGISTRY_REPO`** - Artifact Registry repository name (e.g., `athena-docker`)

- **`JWT_SECRET`** - JWT secret key for your application (generate with `openssl rand -base64 32`)

### Setting Up GCP Service Account

```bash
# Create service account
gcloud iam service-accounts create athena-ci \
  --display-name="Athena CI/CD Service Account"

# Grant permissions
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:athena-ci@PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:athena-ci@PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:athena-ci@PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

# Create and download key
gcloud iam service-accounts keys create key.json \
  --iam-account=athena-ci@PROJECT_ID.iam.gserviceaccount.com

# Copy the contents of key.json to GCP_SA_KEY secret
cat key.json
```

### Setting Up Artifact Registry

```bash
# Create Artifact Registry repository
gcloud artifacts repositories create athena-docker \
  --repository-format=docker \
  --location=us-central1 \
  --description="Athena Docker images"
```

### Workflow Triggers

- **Pull Requests** to main/master/develop: Run tests only
- **Push** to main/master: Run tests → build → deploy
- **Push** to develop: Run tests only

### Manual Deployment

You can also manually trigger deployment from the GitHub Actions tab.

### Monitoring Deployments

```bash
# View Cloud Run service details
gcloud run services describe athena \
  --platform=managed \
  --region=REGION

# View service logs
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=athena" \
  --limit=50 \
  --format=json

# Get service URL
gcloud run services describe athena \
  --platform=managed \
  --region=REGION \
  --format='value(status.url)'
```

### Cloud Run Configuration

The deployment configures:
- **Port**: 1323 (matches Dockerfile EXPOSE)
- **Auto-scaling**: 0-10 instances
- **Resources**: 512Mi memory, 1 CPU
- **Access**: Public (unauthenticated)
- **Environment**: JWT_SECRET from GitHub secrets

To customize these settings, edit the `gcloud run deploy` command in `.github/workflows/test.yml:deploy` job.

## Future Enhancements

- [x] Database persistence (PostgreSQL)
- [ ] Additional database support (MySQL, MongoDB)
- [x] Bookmark title fetching from URL metadata
- [ ] Full-text search across bookmarks
- [ ] Tagging/categorization system
- [ ] Bookmark collections/folders
- [ ] Refresh tokens for extended sessions
- [ ] Email verification for new users
- [ ] Password reset functionality
- [ ] Rate limiting per user/IP
- [ ] API documentation (Swagger/OpenAPI)
- [x] Docker support with docker-compose
- [x] CI/CD pipeline (GitHub Actions)
- [ ] Metrics and monitoring (Prometheus)
- [ ] Graceful shutdown
- [ ] Health check endpoint with dependencies
- [ ] Soft delete for bookmarks
- [ ] Bookmark sharing between users
- [ ] Import/export bookmarks (HTML, JSON)
- [ ] Bookmark duplicates detection
- [ ] Browser extension integration

## Known Limitations

- **In-memory storage**: Data is lost when server restarts
- **Single instance**: Not suitable for horizontal scaling without external session store
- **No refresh tokens**: Users must re-login after 24 hours
- **No password complexity requirements**: Any non-empty password accepted
- **No rate limiting**: Vulnerable to brute force attacks
- **Default JWT secret**: Must be changed in production
- **No email verification**: Anyone can register with any email
- **No pagination**: Large bookmark collections may cause performance issues
- **No bookmark deduplication**: Same URL can be bookmarked multiple times

## Production Deployment Checklist

Before deploying to production:

- [ ] **Set strong `JWT_SECRET`** environment variable (min 32 characters)
- [ ] Replace in-memory repositories with database implementations (PostgreSQL recommended)
- [ ] Add password complexity requirements (min length, special chars, etc.)
- [ ] Implement rate limiting (per IP and per user)
- [ ] Add logging to external service (CloudWatch, ELK, Datadog)
- [ ] Set up monitoring and alerts (server health, error rates)
- [ ] Enable HTTPS/TLS (use Let's Encrypt or load balancer)
- [ ] Review and harden CORS settings (whitelist specific origins)
- [ ] Add input sanitization and validation
- [ ] Implement refresh token mechanism
- [ ] Add database migrations tooling
- [ ] Set up automated database backups
- [ ] Configure proper error handling (don't leak stack traces)
- [ ] Add health check endpoint
- [ ] Implement graceful shutdown
- [ ] Add request timeout limits
- [ ] Set up CDN for static assets (if any)
- [ ] Configure firewall rules
- [ ] Enable audit logging
- [ ] Set up container orchestration (Kubernetes, ECS)

## Development

### Running in Development Mode

```bash
# With hot reload (using air or similar)
go install github.com/cosmtrek/air@latest
air

# Or manually
go run cmd/api-server/main.go
```

### Adding a New Endpoint

1. Define request/response DTOs in `internal/transport/`
2. Add handler method in `internal/handler/`
3. Add tests in `internal/handler/*_test.go`
4. Register route in `cmd/api-server/main.go`
5. Update this README

### Extending to Database

1. Create new repository implementations (e.g., `user_postgres_repo.go`)
2. Implement existing repository interfaces
3. Update `cmd/api-server/main.go` to use new repositories
4. Add database migrations
5. Update configuration for database connection
- [ ] Metrics and monitoring
- [ ] Soft delete for bookmarks
- [ ] Bookmark sharing between users
- [ ] Import/export bookmarks

## Known Limitations

- **In-memory storage**: Data is lost when server restarts
- **Single instance**: Not suitable for horizontal scaling
- **No refresh tokens**: Users must re-login after 24 hours
- **No password complexity requirements**: Consider adding validation
- **No rate limiting**: Vulnerable to brute force attacks
- **Default JWT secret**: Must be changed in production

## Production Deployment Checklist

Before deploying to production:

### Security
- [ ] Set a strong `JWT_SECRET` environment variable (min 32 characters)
- [ ] Add password complexity requirements
- [ ] Implement rate limiting
- [ ] Review and harden CORS settings
- [ ] Add input sanitization
- [ ] Configure proper error handling (don't leak stack traces)
- [ ] Enable HTTPS/TLS (handled by Cloud Run)
- [ ] Implement refresh token mechanism

### Infrastructure
- [x] Replace in-memory repositories with database implementations (PostgreSQL)
- [x] Set up CI/CD pipeline (GitHub Actions)
- [ ] Configure GCP service account with minimal permissions
- [ ] Set up automated database backups
- [ ] Configure database connection pooling
- [ ] Set up monitoring and alerts (Cloud Monitoring, Prometheus)
- [ ] Add logging to external service (Cloud Logging, ELK)
- [ ] Configure firewall rules and VPC if needed

### GitHub Actions Setup
- [ ] Add all required GitHub secrets (GCP_SA_KEY, GCP_PROJECT_ID, etc.)
- [ ] Create GCP Artifact Registry repository
- [ ] Set up GCP service account with proper IAM roles
- [ ] Test deployment pipeline on staging branch first
- [ ] Configure branch protection rules for main branch

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]

## Contact

[Add contact information here]
