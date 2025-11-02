# PostgreSQL Setup

Athena supports PostgreSQL as a persistent storage backend for bookmarks. This guide will help you configure and use PostgreSQL with Athena.

## Prerequisites

- PostgreSQL 12 or higher installed
- Database user with CREATE DATABASE privileges

## Quick Start

### 1. Create Database

```bash
# Connect to PostgreSQL as superuser
psql -U postgres

# Create database
CREATE DATABASE athena;

# Create user (optional)
CREATE USER athena_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE athena TO athena_user;
```

### 2. Configure Environment Variables

Set the following environment variables to use PostgreSQL:

```bash
# Storage configuration
export STORAGE_TYPE=postgres

# Database connection
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=athena
export DB_SSLMODE=disable  # Use 'require' in production
```

### 3. Start the Server

```bash
go run ./cmd/api-server
```

The server will automatically:
- Connect to PostgreSQL
- Run database migrations
- Create the bookmarks table and indexes

## Configuration Options

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `STORAGE_TYPE` | Storage backend (`memory` or `postgres`) | `memory` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | (empty) |
| `DB_NAME` | Database name | `athena` |
| `DB_SSLMODE` | SSL mode (`disable`, `require`, `verify-ca`, `verify-full`) | `disable` |

## Database Schema

### Bookmarks Table

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

### Indexes

- `idx_bookmarks_user_id` - Index on `user_id` for efficient user queries
- `idx_bookmarks_user_id_archived` - Composite index on `user_id` and `is_archived`
- `idx_bookmarks_created_at` - Index on `created_at` for sorting (descending)

## Docker Compose with PostgreSQL

Use the provided docker-compose configuration with PostgreSQL:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: athena-postgres
    environment:
      POSTGRES_DB: athena
      POSTGRES_USER: athena_user
      POSTGRES_PASSWORD: athena_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U athena_user -d athena"]
      interval: 10s
      timeout: 5s
      retries: 5

  api:
    build: .
    container_name: athena-api
    environment:
      STORAGE_TYPE: postgres
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: athena_user
      DB_PASSWORD: athena_password
      DB_NAME: athena
      DB_SSLMODE: disable
      JWT_SECRET: your-secret-key-change-this
    ports:
      - "1323:1323"
    depends_on:
      postgres:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:1323/ping"]
      interval: 30s
      timeout: 3s
      start_period: 10s
      retries: 3

volumes:
  postgres_data:
```

Start with:

```bash
docker-compose up -d
```

## Migrations

Migrations run automatically on server startup. The migration system creates:
- Bookmarks table
- All necessary indexes

To run migrations manually:

```bash
# Connect to database
psql -h localhost -U athena_user -d athena

# Run migration from migrations/001_create_bookmarks_table.sql
\i migrations/001_create_bookmarks_table.sql
```

## Connection Pooling

The PostgreSQL repository is configured with the following connection pool settings:

- **Max Open Connections**: 25
- **Max Idle Connections**: 5
- **Connection Max Lifetime**: 5 minutes

These can be adjusted in `internal/database/postgres.go` if needed.

## Testing

### Running Tests with PostgreSQL

Tests require a test database. Set the `TEST_DATABASE_URL` environment variable:

```bash
export TEST_DATABASE_URL="postgres://athena_user:athena_password@localhost/athena_test?sslmode=disable"
go test ./internal/repository -v -run TestBookmarkPostgres
```

If `TEST_DATABASE_URL` is not set, PostgreSQL tests will be skipped.

### Create Test Database

```sql
CREATE DATABASE athena_test;
GRANT ALL PRIVILEGES ON DATABASE athena_test TO athena_user;
```

## Production Recommendations

### 1. Use SSL/TLS

```bash
export DB_SSLMODE=require
```

### 2. Use Strong Passwords

Generate a strong password for the database user:

```bash
openssl rand -base64 32
```

### 3. Connection String Security

Store sensitive database credentials securely:
- Use environment variables (not hardcoded)
- Use secrets management (AWS Secrets Manager, HashiCorp Vault)
- Restrict database user permissions

### 4. Database Backups

Set up regular backups:

```bash
# Backup
pg_dump -U athena_user athena > backup.sql

# Restore
psql -U athena_user athena < backup.sql
```

### 5. Monitoring

Monitor database performance:
- Connection pool metrics
- Query execution time
- Database size and growth
- Index usage

### 6. Index Maintenance

Regularly analyze and vacuum:

```sql
ANALYZE bookmarks;
VACUUM bookmarks;
```

## Troubleshooting

### Connection Refused

```
Error: failed to connect to PostgreSQL: connection refused
```

**Solution**: Ensure PostgreSQL is running and accepting connections:

```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Check if listening on correct port
netstat -an | grep 5432
```

### Authentication Failed

```
Error: password authentication failed for user "athena_user"
```

**Solution**: Check credentials and pg_hba.conf configuration:

```bash
# Edit pg_hba.conf
sudo vi /etc/postgresql/*/main/pg_hba.conf

# Add line for local connections
host    athena    athena_user    127.0.0.1/32    md5

# Restart PostgreSQL
sudo systemctl restart postgresql
```

### Permission Denied

```
Error: permission denied for table bookmarks
```

**Solution**: Grant proper permissions:

```sql
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO athena_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO athena_user;
```

### Migration Errors

If migrations fail, check:
1. Database user has CREATE TABLE privileges
2. Table doesn't already exist with different schema
3. Database logs for specific error messages

```bash
# View PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-*.log
```

## Performance Tips

### 1. Use Prepared Statements

The repository uses parameterized queries ($1, $2) which are automatically prepared by the PostgreSQL driver.

### 2. Monitor Slow Queries

Enable slow query logging in postgresql.conf:

```
log_min_duration_statement = 1000  # Log queries slower than 1 second
```

### 3. Optimize Indexes

Check index usage:

```sql
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan;
```

### 4. Connection Pool Tuning

Adjust based on your workload in `internal/database/postgres.go`:

```go
db.SetMaxOpenConns(50)    // Increase for high concurrency
db.SetMaxIdleConns(10)    // Keep more idle connections
db.SetConnMaxLifetime(10 * time.Minute)
```

## Switching Between Storage Types

You can switch between in-memory and PostgreSQL storage by changing the `STORAGE_TYPE` environment variable:

```bash
# Use in-memory storage (default)
export STORAGE_TYPE=memory
go run ./cmd/api-server

# Use PostgreSQL storage
export STORAGE_TYPE=postgres
go run ./cmd/api-server
```

**Note**: Data is not automatically migrated between storage types.
