# Athena

A bookmark management API server built with Go and Echo framework.

## Overview

Athena is a lightweight RESTful API service for managing bookmarks. It provides functionality to store, retrieve, and organize bookmarks with support for archiving and user-specific collections.

## Features

- RESTful API for bookmark management
- In-memory storage (with repository pattern for easy extension)
- Clean architecture with separation of concerns (handler, service, repository layers)
- Built-in logging and recovery middleware

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: [Echo v4](https://echo.labstack.com/)

## Project Structure

```
athena/
├── cmd/
│   └── api-server/          # Application entry point
├── internal/
│   ├── handler/             # HTTP request handlers
│   ├── service/             # Business logic layer
│   ├── repository/          # Data access layer
│   └── model/               # Domain models
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
```

## API Endpoints

- `GET /ping` - Health check endpoint

## Testing

Run tests with:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Development

The project follows a clean architecture pattern with three main layers:

- **Handler Layer**: Handles HTTP requests and responses
- **Service Layer**: Contains business logic
- **Repository Layer**: Manages data persistence

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
