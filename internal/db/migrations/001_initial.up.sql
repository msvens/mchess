-- Schema version tracking
CREATE TABLE IF NOT EXISTS schema_version (
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

-- Cache statistics for monitoring
CREATE TABLE cache_stats (
    cache_type      TEXT PRIMARY KEY,           -- 'player', 'tournament', etc.
    hits            BIGINT DEFAULT 0,
    misses          BIGINT DEFAULT 0,
    upstream_calls  BIGINT DEFAULT 0,
    last_reset      TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO cache_stats (cache_type) VALUES ('player');
