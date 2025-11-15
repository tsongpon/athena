[![CI/CD Pipeline](https://github.com/tsongpon/athena/actions/workflows/athena.yml/badge.svg)](https://github.com/tsongpon/athena/actions/workflows/athena.yml)

# Athena

A secure bookmark management API server built with Go and Echo framework with JWT authentication, automatic web metadata extraction, and AI-powered content summarization.

## Overview

Athena is a production-ready RESTful API service for managing bookmarks with intelligent features. It automatically fetches website metadata (titles, images) and generates AI-powered content summaries for bookmarks. All bookmark endpoints are protected with JWT authentication to ensure users can only access their own bookmarks.

## Features

- ✅ JWT-based authentication and authorization
- ✅ User registration and login with tier support (free/paid)
- ✅ RESTful API for bookmark management
- ✅ Create, retrieve, archive, and delete bookmarks
- ✅ **Automatic website metadata extraction:**
  - Page titles from HTML `<title>` tags
  - Images from Open Graph meta tags
  - AI-powered content summaries (for paid tier)
- ✅ **Concurrent metadata fetching** for optimal performance
- ✅ Pluggable storage backends: In-memory and Cloud Firestore
- ✅ Clean architecture with separation of concerns (handler, service, repository layers)
- ✅ Docker support with multi-stage builds
- ✅ Comprehensive test coverage
- ✅ CI/CD pipeline with GitHub Actions
- ✅ Pagination support for bookmark lists
- ✅ LLM integration (Anthropic Claude, OpenAI, Google Gemini)

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: [Echo v4](https://echo.labstack.com/)
- **Authentication**: JWT with [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) & [echo-jwt/v4](https://github.com/labstack/echo-jwt)
- **Cloud Storage**: [Cloud Firestore](https://cloud.google.com/firestore) for persistent storage
- **Logging**: [Uber Zap](https://github.com/uber-go/zap)
- **Password Hashing**: bcrypt
- **ID Generation**: [Google UUID](https://github.com/google/uuid)
- **HTML Parsing**: [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html)
- **LLM Integration**: [LangChain Go](https://github.com/tmc/langchaingo) with Anthropic/OpenAI/Gemini support
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
│   │   ├── bookmark_service.go          # Bookmark business logic with metadata fetching
│   │   ├── bookmark_service_test.go     # Bookmark service tests
│   │   ├── user_service.go              # User authentication & management
│   │   ├── user_service_test.go         # User service tests
│   │   └── repository.go                # Repository interfaces
│   ├── repository/                      # Data access layer
│   │   ├── bookmark_inmem_repo.go       # In-memory bookmark storage
│   │   ├── bookmark_inmem_repo_test.go  # Bookmark repository tests
│   │   ├── bookmark_firestore_repo.go   # Cloud Firestore bookmark storage
│   │   ├── user_inmem_repo.go           # In-memory user storage
│   │   ├── user_inmem_repo_test.go      # User repository tests
│   │   ├── user_firestore_repo.go       # Cloud Firestore user storage
│   │   ├── web_repo.go                  # Web metadata fetching & LLM integration
│   │   └── web_repo_test.go             # Web repository tests
│   ├── transport/                       # HTTP transport layer (DTOs)
│   │   ├── bookmark_transport.go        # Bookmark request/response DTOs
│   │   └── user_transport.go            # User request/response DTOs
│   ├── model/                           # Domain models
│   │   ├── bookmark.go                  # Bookmark domain model
│   │   └── user.go                      # User domain model with tier support
│   └── logger/                          # Structured logging
│       ├── logger.go                    # Zap logger configuration
│       └── logger_test.go               # Logger tests
├── .github/workflows/
│   └── athena.yml                       # CI/CD pipeline
├── docker-compose.yml                   # Docker services configuration
├── Dockerfile                           # Multi-stage Docker build
├── go.mod                               # Go module dependencies
└── README.md                            # This file
```

## Getting Started

### Prerequisites

- Go 1.25.1 or higher
- (Optional) GCP account with Firestore for persistent storage
- (Optional) Docker

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
export STORAGE_TYPE="firestore"  # Options: memory (default), firestore

# Cloud Firestore configuration (only needed if STORAGE_TYPE=firestore)
export GCP_PROJECT_ID="your-project-id"
export GCP_FIRESTORE_DATABASE_ID="athena"

# Logging configuration
export APP_ENV="production"  # Use "production" for JSON logs, default is development
export LOG_LEVEL="info"      # Options: debug, info, warn, error, fatal

# LLM features for content summarization (optional, for paid tier)
export LLM_SUMMARY_CONTENT="true"        # Enable AI content summaries
export LLM_MODEL="anthropic"             # Options: anthropic, openai, gemini
export ANTHROPIC_API_KEY="your-key"      # If using Anthropic Claude
export OPENAI_API_KEY="your-key"         # If using OpenAI
export GEMINI_API_KEY="your-key"         # If using Google Gemini
```

### Running the Server

#### Option 1: Run Locally with Go

```bash
go run cmd/api-server/main.go
```

The server will start on `http://localhost:1323`

#### Option 2: Run with Docker

```bash
# Build and run with docker
docker build -t athena:latest .
docker run -p 1323:1323 -e JWT_SECRET="your-secret-key" athena:latest
```

The server will be available at `http://localhost:1323`



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

#### Health Check
- **GET** `/ping`
  - Response: `pong` (200 OK)
  - Purpose: Verify server is running

#### User Registration
- **POST** `/users`
  - Request body:
    ```json
    {
      "name": "John Doe",
      "email": "john@example.com",
      "password": "securepassword123"
    }
    ```
  - Response: `201 Created`
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "tier": "free",
      "created_at": "2025-11-02T14:00:00Z",
      "updated_at": "2025-11-02T14:00:00Z"
    }
    ```
  - Errors:
    - `400` - Name, email, or password missing
    - `409` - Email already exists

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
        "tier": "free",
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
      "title": "Example Domain",
      "main_image_url": "https://example.com/og-image.png",
      "content_summary": "AI-generated summary...",
      "user_id": "user-id-from-jwt",
      "is_archived": false,
      "created_at": "2025-11-15T10:30:45.123Z"
    }
    ```
  - Note: 
    - `user_id` is automatically extracted from the JWT token
    - Metadata (title, image, summary) is fetched automatically and concurrently
    - `content_summary` is only populated for paid tier users
  - Errors:
    - `400` - URL is missing
    - `401` - Invalid or missing JWT token

#### Get Single Bookmark
- **GET** `/bookmarks/:id`
  - Headers: `Authorization: Bearer <token>`
  - URL Parameters: `id` - Bookmark UUID
  - Response: `200 OK`
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "url": "https://example.com",
      "title": "Example Domain",
      "main_image_url": "https://example.com/og-image.png",
      "content_summary": "AI-generated summary...",
      "user_id": "user-id-from-jwt",
      "is_archived": false,
      "created_at": "2025-11-15T10:30:45.123Z"
    }
    ```
  - Errors:
    - `400` - ID is missing
    - `401` - Invalid or missing JWT token
    - `403` - Bookmark belongs to a different user
    - `404` - Bookmark not found

#### Get All Bookmarks
- **GET** `/bookmarks?archived=false&page=1&page_size=20`
  - Headers: `Authorization: Bearer <token>`
  - Query parameters:
    - `archived` (optional): `true` or `false` (default: `false`)
    - `page` (optional): Page number (default: `1`)
    - `page_size` (optional): Items per page (default: `20`, max: `100`)
  - Response: `200 OK`
    ```json
    {
      "bookmarks": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "url": "https://example.com",
          "title": "Example Domain",
          "main_image_url": "https://example.com/og-image.png",
          "content_summary": "AI-generated summary...",
          "user_id": "user-id-from-jwt",
          "created_at": "2025-11-02T14:00:00Z",
          "is_archived": false
        }
      ],
      "total_count": 150,
      "page": 1,
      "page_size": 20,
      "total_pages": 8
    }
    ```
  - Note: Only returns bookmarks for the authenticated user
  - Errors:
    - `401` - Invalid or missing JWT token

#### Archive Bookmark
- **POST** `/bookmarks/:id/archive`
  - Headers: `Authorization: Bearer <token>`
  - URL Parameters: `id` - Bookmark UUID
  - Response: `204 No Content`
  - Errors:
    - `400` - ID is missing
    - `401` - Invalid or missing JWT token
    - `403` - Bookmark belongs to a different user
    - `404` - Bookmark not found

#### Delete Bookmark
- **DELETE** `/bookmarks/:id`
  - Headers: `Authorization: Bearer <token>`
  - URL Parameters: `id` - Bookmark UUID
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

echo "Token: $TOKEN"

# 4. Create a bookmark (automatically fetches metadata)
curl -X POST http://localhost:1323/bookmarks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url":"https://example.com"}'

# 5. Get all bookmarks
curl -X GET http://localhost:1323/bookmarks \
  -H "Authorization: Bearer $TOKEN"

# 6. Get bookmarks with pagination
curl -X GET "http://localhost:1323/bookmarks?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 7. Get specific bookmark (replace BOOKMARK_ID)
curl -X GET http://localhost:1323/bookmarks/BOOKMARK_ID \
  -H "Authorization: Bearer $TOKEN"

# 8. Archive a bookmark (replace BOOKMARK_ID)
curl -X POST http://localhost:1323/bookmarks/BOOKMARK_ID/archive \
  -H "Authorization: Bearer $TOKEN"

# 9. Get archived bookmarks
curl -X GET http://localhost:1323/bookmarks?archived=true \
  -H "Authorization: Bearer $TOKEN"

# 10. Delete a bookmark (replace BOOKMARK_ID)
curl -X DELETE http://localhost:1323/bookmarks/BOOKMARK_ID \
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
  - **Delete bookmark**: Verifies ownership before deleting (403 if unauthorized)
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
    Tier      string    // User tier: "free" or "paid"
    CreatedAt time.Time // Registration timestamp
    UpdatedAt time.Time // Last update timestamp
}
```

### Bookmark
```go
type Bookmark struct {
    ID             string    // Auto-generated UUID
    UserID         string    // User identifier (owner)
    URL            string    // Bookmark URL
    Title          string    // Page title (auto-fetched from HTML)
    MainImageURL   string    // OpenGraph image (auto-fetched)
    ContentSummary string    // AI-generated summary (paid tier only)
    IsArchived     bool      // Archive status
    CreatedAt      time.Time // Creation timestamp
    UpdatedAt      time.Time // Last update timestamp
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

## Intelligent Features

### Automatic Metadata Extraction

When creating a bookmark, Athena automatically fetches metadata from the target URL using **concurrent goroutines** for optimal performance:

1. **Page Title**: Extracted from HTML `<title>` tag
2. **Main Image**: Extracted from OpenGraph `<meta property="og:image">` tag
3. **Content Summary**: AI-generated summary (paid tier only)

All three operations run in parallel and gracefully handle failures. If metadata fetching fails, the bookmark is still created with the URL.

### AI-Powered Content Summarization

For paid tier users, Athena can generate intelligent content summaries using Large Language Models:

**Supported LLM Providers:**
- **Anthropic Claude** (recommended)
- **OpenAI GPT**
- **Google Gemini**

**Configuration:**
```bash
export LLM_SUMMARY_CONTENT="true"
export LLM_MODEL="anthropic"  # or "openai" or "gemini"
export ANTHROPIC_API_KEY="your-api-key"
```

The content summary is generated asynchronously and won't block bookmark creation if it fails.

### User Tier System

Athena supports a tiered user system:

- **Free Tier**: Basic bookmark management (no AI summaries)
- **Paid Tier**: Full features including AI-powered content summarization

The tier is stored in the user model and checked before generating LLM summaries.

## Storage Backends

### In-Memory Storage (Development)

Default storage backend using Go maps. Data is lost when the server restarts.

```bash
export STORAGE_TYPE="memory"
go run cmd/api-server/main.go
```

### Cloud Firestore (Production)

Google Cloud Platform integration for serverless deployment.

```bash
# Configure Firestore
export STORAGE_TYPE="firestore"
export GCP_PROJECT_ID="your-project-id"
export GCP_FIRESTORE_DATABASE_ID="athena"

# Run server
go run cmd/api-server/main.go
```

**Features:**
- Serverless, auto-scaling storage
- Built-in replication and backups
- Native GCP integration

## Docker

The application includes Docker support for easy deployment.

### Docker Configuration

Configure via environment variables:

```bash
# Core settings
JWT_SECRET=your-super-secret-key-change-in-production
STORAGE_TYPE=firestore
APP_ENV=production
LOG_LEVEL=info

# Firestore (for persistent storage)
GCP_PROJECT_ID=your-project-id
GCP_FIRESTORE_DATABASE_ID=athena

# LLM features (optional)
LLM_SUMMARY_CONTENT=true
LLM_MODEL=anthropic
ANTHROPIC_API_KEY=your-api-key
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

# Run with LLM support
docker run -d -p 1323:1323 \
  -e JWT_SECRET="your-secret" \
  -e LLM_SUMMARY_CONTENT="true" \
  -e LLM_MODEL="anthropic" \
  -e ANTHROPIC_API_KEY="your-api-key" \
  --name athena-api \
  athena:latest
```

### Docker Container Management

```bash
# View running containers
docker ps

# View logs
docker logs athena-api
docker logs -f athena-api  # Follow logs

# Stop container
docker stop athena-api

# Start container
docker start athena-api

# Restart container
docker restart athena-api

# Remove container
docker rm athena-api
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

### Run Specific Test
```bash
# Auth handler tests
go test -v -run TestAuthHandler_Login ./internal/handler/

# Bookmark handler tests
go test -v -run TestBookmarkHandler_CreateBookmark ./internal/handler/

# Web repository tests
go test -v -run TestWebRepository ./internal/repository/
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
  - `web_repo.go`: Covered by integration tests
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
     - Bookmark validation
     - Orchestrate repository operations
     - Coordinate concurrent metadata fetching
   - **Files**:
     - `user_service.go`: User authentication and management
     - `bookmark_service.go`: Bookmark business logic with metadata fetching
     - `repository.go`: Repository interfaces

3. **Repository Layer** (`internal/repository/`)
   - **Purpose**: Data persistence abstraction
   - **Responsibilities**:
     - CRUD operations on data store
     - Multiple implementations (in-memory, Firestore)
     - Web metadata extraction and LLM integration
   - **Files**:
     - `user_inmem_repo.go`: In-memory user storage
     - `bookmark_inmem_repo.go`: In-memory bookmark storage
     - `user_firestore_repo.go`: Firestore user storage
     - `bookmark_firestore_repo.go`: Firestore bookmark storage
     - `web_repo.go`: Web scraping and LLM integration

4. **Transport Layer** (`internal/transport/`)
   - **Purpose**: HTTP API contracts (DTOs)
   - **Responsibilities**:
     - Define request/response structures
     - Separate external API from internal domain
   - **Files**:
     - `user_transport.go`: LoginRequest, CreateUserRequest, UserResponse, LoginResponse
     - `bookmark_transport.go`: BookmarkTransport, BookmarkListResponse

5. **Model Layer** (`internal/model/`)
   - **Purpose**: Domain models
   - **Responsibilities**:
     - Pure data structures
     - Minimal business logic
   - **Files**:
     - `user.go`: User domain model with tier support
     - `bookmark.go`: Bookmark domain model with metadata fields

6. **Logger Layer** (`internal/logger/`)
   - **Purpose**: Structured logging
   - **Responsibilities**:
     - Configure Zap logger for development/production
     - Provide consistent logging interface
   - **Files**:
     - `logger.go`: Zap logger configuration

### Data Flow

```
HTTP Request
    ↓
[Handler] ← Validates request, extracts JWT claims
    ↓
[Service] ← Business logic, concurrent metadata fetching
    ↓ ↓ ↓
[BookmarkRepo] [UserRepo] [WebRepo] ← Data persistence & external fetching
    ↓
[Storage: Firestore/Memory]
```

### Design Patterns

- **Repository Pattern**: Abstracts data access logic
- **Dependency Injection**: Services and handlers receive dependencies via constructors
- **Interface Segregation**: Small, focused interfaces (BookmarkRepository, UserRepository, WebRepository)
- **Middleware Pattern**: JWT authentication, CORS, logging, recovery
- **Concurrent Processing**: Goroutines for parallel metadata fetching
- **Helper Functions**: Reusable JWT claim extraction logic (`getAuthenticatedUser`)

### Security Architecture

- **Stateless Authentication**: JWT tokens contain all necessary user information
- **Authorization at Handler Level**: Each protected endpoint verifies user ownership
- **Secure Password Storage**: bcrypt hashing with salt (never store plaintext)
- **Defense in Depth**: Multiple layers of validation (handler, service, repository)

### Key Design Decisions

1. **JWT in Context**: Echo JWT middleware v4 stores `*jwt.Token` in context. The `getAuthenticatedUser` helper extracts claims safely.

2. **User ID from JWT**: User ID is never accepted from request parameters for protected endpoints, always extracted from authenticated token.

3. **Authorization at Handler**: Authorization checks happen at the handler layer before calling services.

4. **Concurrent Metadata Fetching**: Title, image, and LLM summary are fetched concurrently using goroutines for optimal performance.

5. **Graceful Degradation**: If metadata fetching fails, bookmarks are still created with the URL.

## CI/CD Pipeline

Athena uses GitHub Actions for continuous integration and deployment to Google Cloud Run.

### Workflow Overview

The CI/CD pipeline consists of three jobs:

1. **Test** - Runs on all pushes and pull requests
   - Runs unit tests with race detection
   - Generates code coverage reports
   - Enforces 50% minimum coverage threshold

2. **Build** - Runs only on pushes to main branch
   - Builds Docker image
   - Pushes to GCP Artifact Registry
   - Tags images with branch name, SHA, and `latest`

3. **Deploy** - Runs only on pushes to main branch
   - Deploys to Google Cloud Run
   - Configures auto-scaling and environment variables
   - Outputs the deployed service URL

### Required GitHub Secrets

Configure these secrets in your GitHub repository settings (Settings → Secrets and variables → Actions):

- **`GCP_SA_KEY`** - GCP Service Account JSON key
- **`GCP_PROJECT_ID`** - Your GCP project ID
- **`GCP_REGION`** - GCP region for deployment (e.g., `us-central1`)
- **`GCP_ARTIFACT_REGISTRY_REPO`** - Artifact Registry repository name
- **`JWT_SECRET`** - JWT secret key for your application

### Workflow Triggers

- **Pull Requests** to main/develop: Run tests only
- **Push** to main: Run tests → build → deploy
- **Push** to other branches: Run tests only

## Logging

Athena uses [Uber Zap](https://github.com/uber-go/zap) for structured, high-performance logging.

### Configuration

- **Development Mode** (default): Console output with colors, debug level
- **Production Mode** (`APP_ENV=production`): JSON formatted logs, info level
- **Custom Log Level**: Set `LOG_LEVEL` environment variable (debug, info, warn, error, fatal)

### Example Logs

**Development:**
```
2025-11-15T10:30:45.123+0700    INFO    service/bookmark_service.go:42    Created bookmark    {"id": "abc123", "user_id": "user1"}
```

**Production:**
```json
{"level":"info","timestamp":"2025-11-15T10:30:45.123Z","msg":"Created bookmark","id":"abc123","user_id":"user1"}
```

## Error Handling

The API returns standard HTTP status codes:

- `200 OK` - Successful GET request
- `201 Created` - Successful POST request (resource created)
- `204 No Content` - Successful POST/DELETE request (no content returned)
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

## Future Enhancements

- [x] Database persistence (Firestore)
- [x] Bookmark metadata fetching from URL (title, image)
- [x] AI-powered content summarization
- [x] Pagination support
- [x] User tier system
- [x] Docker support
- [x] CI/CD pipeline (GitHub Actions)
- [ ] PostgreSQL support (repository implementation needed)
- [ ] Full-text search across bookmarks
- [ ] Tagging/categorization system
- [ ] Bookmark collections/folders
- [ ] Refresh tokens for extended sessions
- [ ] Email verification for new users
- [ ] Password reset functionality
- [ ] Rate limiting per user/IP
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Metrics and monitoring (Prometheus)
- [ ] Graceful shutdown
- [ ] Bookmark sharing between users
- [ ] Import/export bookmarks (HTML, JSON)
- [ ] Bookmark duplicate detection
- [ ] Browser extension integration

## Known Limitations

- **In-memory storage**: Data is lost when server restarts (use Firestore for production persistence)
- **No PostgreSQL support**: PostgreSQL repository implementation not yet available
- **No refresh tokens**: Users must re-login after 24 hours
- **No password complexity requirements**: Consider adding validation
- **No rate limiting**: Vulnerable to brute force attacks
- **Default JWT secret**: Must be changed in production
- **No email verification**: Anyone can register with any email
- **LLM cost**: Content summarization may incur API costs for paid users

## Production Deployment Checklist

Before deploying to production:

### Security
- [ ] Set a strong `JWT_SECRET` environment variable (min 32 characters)
- [ ] Add password complexity requirements
- [ ] Implement rate limiting (per IP/user)
- [ ] Review and harden CORS settings
- [ ] Add input sanitization
- [ ] Configure proper error handling (don't leak stack traces)
- [ ] Enable HTTPS/TLS (handled by Cloud Run)
- [ ] Implement refresh token mechanism
- [ ] Add email verification

### Infrastructure
- [x] Replace in-memory repositories with persistent storage (Firestore)
- [x] Set up CI/CD pipeline (GitHub Actions)
- [ ] Configure GCP service account with minimal permissions
- [ ] Set up automated database backups
- [ ] Configure database connection pooling
- [ ] Set up monitoring and alerts (Cloud Monitoring)
- [ ] Add logging to external service (Cloud Logging)
- [ ] Configure firewall rules and VPC if needed
- [ ] Set up rate limiting and CDN (Cloud Armor, Cloud CDN)

### LLM Features
- [ ] Monitor LLM API costs and set budget alerts
- [ ] Implement caching for duplicate content summaries
- [ ] Add retry logic with exponential backoff for LLM calls
- [ ] Set up fallback if LLM provider is down

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
