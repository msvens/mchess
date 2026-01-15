# mchess

A Go-based caching proxy for the Swedish Chess Federation (schack.se) API.

## Overview

mchess sits between your application and the schack.se API, providing:

- **Smart Caching**: Reduces redundant API calls by caching player data with intelligent TTL (historical data never expires, current month data has configurable TTL)
- **Batch Operations**: Fetch multiple players in a single request (the upstream API requires individual calls)
- **Rate Limiting**: Built-in request throttling to respect upstream API limits
- **Drop-in Replacement**: API routes match schack.se exactly, allowing seamless migration

## Why?

The schack.se API is excellent but has limitations for frontend applications:

1. **No batch operations** - fetching 10 player profiles requires 10 separate API calls
2. **Static data** - player ratings update monthly, tournament results are immutable once published
3. **Rate limits** - repeated fetches of the same data waste quota

mchess solves these by caching intelligently and adding batch endpoints.

## Architecture

```
┌─────────────────┐
│   Your App      │
│   (Frontend)    │
└────────┬────────┘
         │ HTTP/REST
         ▼
┌─────────────────┐
│   mchess        │
│   - Caching     │
│   - Batching    │
│   - Rate Limit  │
└────────┬────────┘
         │ HTTP/REST
         ▼
┌─────────────────┐
│  schack.se API  │
└─────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 14+

### Installation

```bash
git clone https://github.com/msvens/mchess.git
cd mchess
go build -o mchess .
```

### Setup

1. Create the database:
   ```bash
   createdb mchess
   ```

2. Copy and configure:
   ```bash
   cp config.example.yaml config.yaml
   # Edit config.yaml with your database credentials
   ```

3. Run migrations:
   ```bash
   ./mchess db create
   ```

4. Start the server:
   ```bash
   ./mchess serve
   ```

5. Visit `http://localhost:8080/swagger/` for interactive API documentation.

## CLI Commands

```bash
mchess                  # Start the server (default action)
mchess serve            # Start the API server
mchess db create        # Create database tables (run migrations)
mchess db upgrade       # Upgrade database schema (pending migrations)
mchess db delete        # Delete all database tables (WARNING: destroys data)
mchess db version       # Show current database schema version
mchess version          # Show mchess version
mchess --help           # Show help
```

### Global Flags

```bash
--config string    Config file path (default: ./config.yaml)
```

## Configuration

Configuration is loaded from `config.yaml` or environment variables:

```yaml
server:
  port: 8080
  # host: localhost      # default: localhost
  # prefix: /api         # default: /api

db:
  host: localhost
  port: 5432
  user: mchess
  password: mchess
  name: mchess

upstream:
  baseUrl: https://member.schack.se/public/api/v1
  # timeout: 30s         # default: 30s
  # rateLimit: 10        # default: 10 requests/second

cache:
  # TTL for "current" data (current month or no date specified)
  # Historical data (past months) never expires - it's immutable
  ttl: 24h

log:
  level: info    # debug, info, warn, error
  # format: text  # text (dev, default) or json (prod)
```

## API Documentation

### Swagger UI

Interactive API documentation is available at `/swagger/` when the server is running.

### Health Endpoints

| Endpoint  | Description |
|-----------|-------------|
| `/health` | Basic health check (always returns OK) |
| `/ready`  | Readiness check (verifies database connection) |

### API Endpoints

All API endpoints are under the configured prefix (default: `/api`).

#### Player Endpoints (with caching)

| Endpoint | Description |
|----------|-------------|
| `GET /api/player/{id}/date/{date}` | Get player by member ID |
| `GET /api/player/fideid/{id}/date/{date}` | Get player by FIDE ID |
| `GET /api/player/fornamn/{fornamn}/efternamn/{efternamn}` | Search players by name |
| `GET /api/player/batch?ids=1,2,3&date=...` | **mchess**: Batch fetch multiple players |
| `GET /api/player/{id}/ratings?from=...&to=...` | **mchess**: Get rating history |

#### Organisation Endpoints (pass-through)

| Endpoint | Description |
|----------|-------------|
| `GET /api/organisation/federation` | Get federation info |
| `GET /api/organisation/districts` | Get all districts |
| `GET /api/organisation/district/clubs/{districtid}` | Get clubs in district |
| `GET /api/organisation/club/{clubid}` | Get club by ID |

#### Rating List Endpoints (pass-through)

| Endpoint | Description |
|----------|-------------|
| `GET /api/ratinglist/federation/date/{date}/ratingtype/{type}/category/{cat}` | Federation rating list |
| `GET /api/ratinglist/district/{id}/date/{date}/ratingtype/{type}/category/{cat}` | District rating list |
| `GET /api/ratinglist/club/{id}/date/{date}/ratingtype/{type}/category/{cat}` | Club rating list |

#### Tournament Endpoints (pass-through)

| Endpoint | Description |
|----------|-------------|
| `GET /api/tournament/tournament/id/{id}` | Get tournament by ID |
| `GET /api/tournament/group/id/{id}` | Get tournament from group |
| `GET /api/tournament/group/coming` | Get upcoming tournaments |
| `GET /api/tournament/group/search/{searchWord}` | Search tournaments |

#### Tournament Results Endpoints (pass-through)

| Endpoint | Description |
|----------|-------------|
| `GET /api/tournamentresults/table/id/{id}` | Get tournament standings |
| `GET /api/tournamentresults/roundresults/id/{id}` | Get round results |
| `GET /api/tournamentresults/game/memberid/{id}` | Get games for member |

## Cache Strategy

mchess uses intelligent caching based on data immutability:

- **Historical data** (any month before current): **Never expires** - ratings are immutable after the month ends
- **Current month data**: Configurable TTL (default 24h)

```
Example: Today is 2024-06-15, TTL is 24h

Request for 2024-05-01 → Cache forever (historical)
Request for 2024-06-01 → Cache for 24h (current month)
Request with no date   → Cache for 24h
```

## Project Structure

```
mchess/
├── cmd/                    # CLI commands (serve, db, version)
├── internal/
│   ├── api/               # HTTP server and routes
│   │   └── handlers/      # Request handlers
│   ├── config/            # Configuration loading
│   ├── db/                # Database connection and migrations
│   ├── model/             # Domain types
│   ├── repository/        # Database access layer
│   ├── service/           # Business logic
│   └── upstream/          # schack.se API client
├── migrations/            # SQL migration files
├── api-specs/             # OpenAPI spec from schack.se
├── docs/                  # Generated Swagger documentation
├── config.example.yaml    # Example configuration
└── main.go
```

## Development

### Build

```bash
go build -o mchess .
```

### Run in development

```bash
go run . serve
```

### Generate Swagger docs

```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init
```

### Run tests

```bash
go test ./...
```

## API Compatibility

mchess aims to be a drop-in replacement for the schack.se API. All upstream endpoints are available under `/api/` with identical request/response formats.

**mchess-only endpoints:**
- `GET /api/player/batch` - Fetch multiple players in one request
- `GET /api/player/{id}/ratings` - Get rating history for a player

## Status

### Implemented
- Player endpoints with PostgreSQL caching
- Batch player fetch
- Rating history
- All organisation, tournament, and results endpoints (pass-through)
- Swagger UI documentation
- Database migrations
- Graceful shutdown

### Planned
- Caching for organisation, tournament, and results endpoints
- Background cache warming
- Redis support for distributed caching
- Cache statistics and metrics

## License

MIT