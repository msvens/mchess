# mchess - Backend API Design

A Go-based caching proxy for the Swedish Chess Federation (schack.se) API.

## Project Structure

```
mchess/
├── cmd/
│   ├── root.go              # Root command, config initialization
│   ├── serve.go             # Start the API server
│   ├── migrate.go           # Run database migrations
│   └── version.go           # Print version info
├── internal/
│   ├── config/
│   │   └── config.go        # Viper-based config with typed accessors
│   ├── api/
│   │   ├── server.go        # Chi router setup, server lifecycle
│   │   ├── routes.go        # All route definitions in one place
│   │   ├── middleware.go    # Rate limiting, logging, recovery
│   │   ├── response.go      # JSON response helpers
│   │   ├── errors.go        # Structured API errors
│   │   └── handlers/
│   │       ├── players.go   # Player endpoints
│   │       └── health.go    # Health/ready checks
│   ├── db/
│   │   ├── db.go            # Database connection, migrations
│   │   └── queries.sql      # sqlc query definitions
│   ├── repository/
│   │   ├── repository.go    # Repository interfaces
│   │   └── player.go        # Player cache repository
│   ├── upstream/
│   │   ├── client.go        # HTTP client for schack.se
│   │   ├── ratelimit.go     # Rate limiter for upstream calls
│   │   └── players.go       # Player API methods
│   ├── service/
│   │   └── player.go        # Business logic: cache → upstream → store
│   └── model/
│       └── player.go        # Domain types (separate from upstream DTOs)
├── migrations/
│   ├── 001_initial.up.sql
│   └── 001_initial.down.sql
├── main.go
├── go.mod
├── config.yaml
├── config.example.yaml
└── Makefile
```

## Configuration

Simplified config - sensible defaults are built into the code. Only configure what you need to change.

```yaml
# config.yaml
server:
  port: 8080

db:
  host: localhost
  port: 5432
  user: mchess
  password: mchess
  name: mchess

upstream:
  baseUrl: https://member.schack.se/public/api/v1

cache:
  # TTL for "current" data (current month or no date specified)
  # Historical data (past months) never expires - it's immutable
  ttl: 24h

log:
  level: info    # debug, info, warn, error
```

**Defaults (in code, not config):**
- `server.host`: `localhost`
- `server.prefix`: `/api`
- `upstream.timeout`: `30s`
- `upstream.rateLimit`: `10` requests/second
- `cache.ttl`: `24h`
- `log.format`: `text` (dev) / `json` (when deployed)

## Database Schema

```sql
-- migrations/001_initial.up.sql

-- Schema version tracking
CREATE TABLE schema_version (
    version     INTEGER PRIMARY KEY,
    applied_at  TIMESTAMPTZ DEFAULT NOW(),
    description TEXT
);

INSERT INTO schema_version (version, description)
VALUES (1, 'Initial schema with player cache');

-- Player cache with temporal awareness
-- Key insight: historical data never changes, only "current" data needs TTL
CREATE TABLE player_cache (
    member_id       INTEGER NOT NULL,
    rating_date     DATE NOT NULL,              -- First of month: 2024-06-01
    first_name      TEXT,
    last_name       TEXT,
    club            TEXT,
    club_id         INTEGER,
    fide_id         INTEGER,
    -- Denormalized ratings for quick access
    elo_standard    INTEGER,
    elo_rapid       INTEGER,
    elo_blitz       INTEGER,
    lask_rating     INTEGER,
    -- Full response stored as JSONB for complete data
    data            JSONB NOT NULL,
    -- Cache metadata
    fetched_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ,                -- NULL = never expires

    PRIMARY KEY (member_id, rating_date)
);

-- Indexes
CREATE INDEX idx_player_cache_expires
    ON player_cache(expires_at)
    WHERE expires_at IS NOT NULL;

CREATE INDEX idx_player_cache_member
    ON player_cache(member_id);

CREATE INDEX idx_player_cache_name
    ON player_cache(last_name, first_name);

CREATE INDEX idx_player_cache_club
    ON player_cache(club_id)
    WHERE club_id IS NOT NULL;

-- Cache statistics (optional, for monitoring)
CREATE TABLE cache_stats (
    cache_type      TEXT PRIMARY KEY,           -- 'player', 'tournament', etc.
    hits            BIGINT DEFAULT 0,
    misses          BIGINT DEFAULT 0,
    upstream_calls  BIGINT DEFAULT 0,
    last_reset      TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO cache_stats (cache_type) VALUES ('player');
```

## API Routes

```go
// internal/api/routes.go

func (s *Server) routes() {
    r := s.router

    // Global middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(s.requestLogger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(30 * time.Second))

    // Health endpoints (no prefix)
    r.Get("/health", s.handleHealth)
    r.Get("/ready", s.handleReady)

    // API routes
    r.Route(s.prefix, func(r chi.Router) {
        // Players
        r.Route("/players", func(r chi.Router) {
            r.Get("/", s.handleGetPlayers)           // Batch: ?ids=1,2,3&date=2024-06-01
            r.Get("/search", s.handleSearchPlayers)  // ?firstName=&lastName=
            r.Get("/{id}", s.handleGetPlayer)        // Single player, optional ?date=
            r.Get("/{id}/ratings", s.handleGetPlayerRatings) // ?from=&to= or ?dates=
        })

        // Future: Tournaments, Organizations, etc.
        // r.Route("/tournaments", func(r chi.Router) { ... })

        // Future: Cache management (admin)
        // r.Route("/admin/cache", func(r chi.Router) {
        //     r.Delete("/players/{id}", s.handleInvalidatePlayer)
        //     r.Post("/warm", s.handleWarmCache)
        // })
    })
}
```

## API Compatibility

**Important**: mchess does NOT change the upstream schack.se API response format.
We return the exact same `PlayerInfoDto` structure, just with optional metadata fields added.

This means:
- Frontend code that works with schack.se works with mchess unchanged
- We only ADD fields (like `_cached`, `_cachedAt`), never modify existing ones
- New endpoints (like `/players` batch, `/players/{id}/ratings`) are additions

## API Endpoints Detail

### GET /api/docs
Returns OpenAPI 3.0 specification (JSON).

### GET /api/swagger/
Interactive Swagger UI for testing endpoints.

---

### GET /api/players
Batch fetch multiple players. **NEW - not in upstream API.**

**Query Parameters:**
- `ids` (required): Comma-separated member IDs
- `date` (optional): Rating date (YYYY-MM-DD), defaults to current

**Response:**
```json
{
  "players": [
    {
      "id": 12345,
      "firstName": "Magnus",
      "lastName": "Carlsen",
      "birthdate": "1990-11-30T00:00:00Z",
      "sex": 1,
      "fideid": 1503014,
      "country": "SWE",
      "club": "SK Stockholm",
      "clubId": 100,
      "elo": {
        "rating": 2830,
        "rapidRating": 2800,
        "blitzRating": 2850,
        "title": "GM",
        "date": "2024-06-01T00:00:00Z",
        "k": 10
      },
      "lask": {
        "rating": 2400,
        "date": "2024-06-01T00:00:00Z"
      },
      "_cached": true,
      "_cachedAt": "2024-06-15T10:30:00Z"
    }
  ],
  "errors": [
    { "id": 99999, "error": "player not found" }
  ]
}
```

Note: The player object matches `PlayerInfoDto` from schack.se exactly.
Only `_cached` and `_cachedAt` are added (prefixed with `_` to avoid conflicts).

---

### GET /api/players/{id}
Get single player. **Matches upstream:** `/player/{id}/date/{date}`

**Query Parameters:**
- `date` (optional): Rating date (YYYY-MM-DD), defaults to current

**Response:** Same as upstream `PlayerInfoDto`, plus `_cached` and `_cachedAt`.

---

### GET /api/players/{id}/ratings
Get rating history for a player. **NEW - not in upstream API.**

**Query Parameters (one of):**
- `from` + `to`: Date range (YYYY-MM-DD or YYYY-MM)
- `dates`: Comma-separated specific dates
- `months`: Number of months back from today (e.g., `12`)

**Response:**
```json
{
  "playerId": 12345,
  "ratings": [
    {
      "date": "2024-06-01",
      "elo": {
        "rating": 2830,
        "rapidRating": 2800,
        "blitzRating": 2850,
        "title": "GM"
      },
      "lask": {
        "rating": 2400
      },
      "_cached": true
    },
    {
      "date": "2024-05-01",
      "elo": {
        "rating": 2825,
        "rapidRating": 2795,
        "blitzRating": 2845,
        "title": "GM"
      },
      "lask": {
        "rating": 2395
      },
      "_cached": true
    }
  ]
}
```

---

### GET /api/players/search
Search players by name. **Proxies to upstream:** `/player/fornamn/{firstName}/efternamn/{lastName}`

**Query Parameters:**
- `firstName` (required)
- `lastName` (required)

**Response:** Array of `PlayerInfoDto` (same as upstream, no caching).

## Key Dependencies

```go
// go.mod
module github.com/msvens/mchess

go 1.22

require (
    // Router
    github.com/go-chi/chi/v5

    // Database
    github.com/jackc/pgx/v5
    github.com/jackc/pgx/v5/stdlib  // database/sql driver

    // Migrations
    github.com/golang-migrate/migrate/v4

    // Configuration
    github.com/spf13/cobra
    github.com/spf13/viper

    // Utilities
    golang.org/x/time/rate          // Rate limiting
    golang.org/x/sync/errgroup      // Parallel fetching

    // OpenAPI/Swagger
    github.com/swaggo/swag          // Generates OpenAPI spec from annotations
    github.com/swaggo/http-swagger  // Serves Swagger UI
)
```

## Cache Strategy Details

### TTL Decision Logic

Simple rule:
- **Historical data** (any month before current month): **Never expires** - ratings are immutable
- **Current month data**: Expires after configured TTL (default 24h)

```
Example: Today is 2024-06-15, TTL is 24h

Request for player rating on 2024-05-01 → Cache forever (historical)
Request for player rating on 2024-04-01 → Cache forever (historical)
Request for player rating on 2024-06-01 → Cache for 24h (current month)
Request with no date (latest)           → Cache for 24h
```

This is simple and effective. If you later want smarter logic (e.g., "ratings update on the 5th"), we can add it.

### Batch Fetch Flow

```
Request: GET /api/players?ids=1,2,3&date=2024-06-01
                    │
                    ▼
            ┌───────────────┐
            │ Parse request │
            └───────┬───────┘
                    │
                    ▼
            ┌───────────────┐
            │ Check cache   │ ─── Cache hits → collect results
            │ for each ID   │
            └───────┬───────┘
                    │ Cache misses
                    ▼
            ┌───────────────┐
            │ Rate-limited  │ ─── Parallel fetch with concurrency limit
            │ upstream fetch│
            └───────┬───────┘
                    │
                    ▼
            ┌───────────────┐
            │ Store in cache│ ─── Calculate TTL per item
            └───────┬───────┘
                    │
                    ▼
            ┌───────────────┐
            │ Return merged │
            │ results       │
            └───────────────┘
```

## Makefile

```makefile
.PHONY: build run test migrate-up migrate-down generate

# Build
build:
	go build -o bin/mchess .

# Run development server
run:
	go run . serve

# Run with config file
run-config:
	go run . serve --config ./config.yaml

# Database migrations
migrate-up:
	go run . migrate up

migrate-down:
	go run . migrate down

# Generate sqlc (if using)
generate:
	sqlc generate

# Tests
test:
	go test ./...

# Lint
lint:
	golangci-lint run

# Docker
docker-build:
	docker build -t mchess .

docker-run:
	docker run -p 8080:8080 mchess
```

## Future Considerations

### Phase 2: Additional Endpoints
- Tournaments (results are immutable - cache forever)
- Organizations (districts, clubs - long TTL)
- Rating lists (monthly snapshots - cache forever)

### Phase 3: Advanced Features
- Cache warming jobs (scheduled)
- Admin endpoints for cache management
- Metrics/Prometheus endpoint
- Background job for cleaning expired cache entries

### Phase 4: Frontend Migration
- Generate TypeScript types from OpenAPI spec
- Update Next.js to call mchess instead of schack.se directly

---

## Decisions Made

1. **sqlc** for type-safe database queries (generates Go from SQL)
2. **Search**: Proxy through to upstream, no caching for now
3. **Batch errors**: Return partial results (players found + errors for failures)
4. **Swagger**: Use swaggo to generate OpenAPI spec from code annotations
5. **Router**: chi (modern, actively maintained, built-in middleware)
6. **Logging**: slog (Go 1.21+ stdlib, structured logging)