-- name: GetPlayerCache :one
SELECT * FROM player_cache
WHERE member_id = $1 AND rating_date = $2
AND (expires_at IS NULL OR expires_at > NOW());

-- name: GetPlayerCacheByMember :many
SELECT * FROM player_cache
WHERE member_id = $1
AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY rating_date DESC;

-- name: GetPlayersCacheBatch :many
SELECT * FROM player_cache
WHERE member_id = ANY($1::int[]) AND rating_date = $2
AND (expires_at IS NULL OR expires_at > NOW());

-- name: UpsertPlayerCache :exec
INSERT INTO player_cache (
    member_id, rating_date, first_name, last_name, club, club_id, fide_id,
    elo_standard, elo_rapid, elo_blitz, lask_rating, data, fetched_at, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
ON CONFLICT (member_id, rating_date) DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    club = EXCLUDED.club,
    club_id = EXCLUDED.club_id,
    fide_id = EXCLUDED.fide_id,
    elo_standard = EXCLUDED.elo_standard,
    elo_rapid = EXCLUDED.elo_rapid,
    elo_blitz = EXCLUDED.elo_blitz,
    lask_rating = EXCLUDED.lask_rating,
    data = EXCLUDED.data,
    fetched_at = EXCLUDED.fetched_at,
    expires_at = EXCLUDED.expires_at;

-- name: DeleteExpiredCache :execrows
DELETE FROM player_cache
WHERE expires_at IS NOT NULL AND expires_at < NOW();

-- name: GetCacheStats :one
SELECT * FROM cache_stats WHERE cache_type = $1;

-- name: IncrementCacheHit :exec
UPDATE cache_stats SET hits = hits + 1 WHERE cache_type = $1;

-- name: IncrementCacheMiss :exec
UPDATE cache_stats SET misses = misses + 1 WHERE cache_type = $1;

-- name: IncrementUpstreamCalls :exec
UPDATE cache_stats SET upstream_calls = upstream_calls + $2 WHERE cache_type = $1;
