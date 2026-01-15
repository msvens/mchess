package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/msvens/mchess/internal/model"
)

// PlayerRepository handles player cache database operations
type PlayerRepository struct {
	db *sql.DB
}

// NewPlayerRepository creates a new player repository
func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

// Get retrieves a cached player by member ID and rating date
func (r *PlayerRepository) Get(ctx context.Context, memberID int, ratingDate time.Time) (*model.PlayerInfo, error) {
	query := `
		SELECT data FROM player_cache
		WHERE member_id = $1 AND rating_date = $2
		AND (expires_at IS NULL OR expires_at > NOW())`

	var data []byte
	err := r.db.QueryRowContext(ctx, query, memberID, ratingDate).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query player cache: %w", err)
	}

	var player model.PlayerInfo
	if err := json.Unmarshal(data, &player); err != nil {
		return nil, fmt.Errorf("unmarshal player data: %w", err)
	}

	return &player, nil
}

// GetByFideID retrieves a cached player by FIDE ID and rating date
func (r *PlayerRepository) GetByFideID(ctx context.Context, fideID int, ratingDate time.Time) (*model.PlayerInfo, error) {
	query := `
		SELECT data FROM player_cache
		WHERE fide_id = $1 AND rating_date = $2
		AND (expires_at IS NULL OR expires_at > NOW())`

	var data []byte
	err := r.db.QueryRowContext(ctx, query, fideID, ratingDate).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query player cache by fide id: %w", err)
	}

	var player model.PlayerInfo
	if err := json.Unmarshal(data, &player); err != nil {
		return nil, fmt.Errorf("unmarshal player data: %w", err)
	}

	return &player, nil
}

// GetBatch retrieves multiple cached players by member IDs and rating date
func (r *PlayerRepository) GetBatch(ctx context.Context, memberIDs []int, ratingDate time.Time) (map[int]*model.PlayerInfo, error) {
	if len(memberIDs) == 0 {
		return make(map[int]*model.PlayerInfo), nil
	}

	query := `
		SELECT member_id, data FROM player_cache
		WHERE member_id = ANY($1) AND rating_date = $2
		AND (expires_at IS NULL OR expires_at > NOW())`

	rows, err := r.db.QueryContext(ctx, query, memberIDs, ratingDate)
	if err != nil {
		return nil, fmt.Errorf("query player cache batch: %w", err)
	}
	defer rows.Close()

	result := make(map[int]*model.PlayerInfo)
	for rows.Next() {
		var memberID int
		var data []byte
		if err := rows.Scan(&memberID, &data); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		var player model.PlayerInfo
		if err := json.Unmarshal(data, &player); err != nil {
			return nil, fmt.Errorf("unmarshal player data: %w", err)
		}

		result[memberID] = &player
	}

	return result, rows.Err()
}

// GetRatingHistory retrieves all cached player info for a player (all dates)
func (r *PlayerRepository) GetRatingHistory(ctx context.Context, memberID int) ([]model.PlayerInfo, error) {
	query := `
		SELECT data FROM player_cache
		WHERE member_id = $1
		AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY rating_date DESC`

	rows, err := r.db.QueryContext(ctx, query, memberID)
	if err != nil {
		return nil, fmt.Errorf("query rating history: %w", err)
	}
	defer rows.Close()

	var result []model.PlayerInfo
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		var player model.PlayerInfo
		if err := json.Unmarshal(data, &player); err != nil {
			return nil, fmt.Errorf("unmarshal player data: %w", err)
		}

		result = append(result, player)
	}

	return result, rows.Err()
}

// GetRatingHistoryForDates retrieves cached player info for specific dates
func (r *PlayerRepository) GetRatingHistoryForDates(ctx context.Context, memberID int, dates []time.Time) (map[string]*model.PlayerInfo, error) {
	if len(dates) == 0 {
		return make(map[string]*model.PlayerInfo), nil
	}

	query := `
		SELECT rating_date, data FROM player_cache
		WHERE member_id = $1 AND rating_date = ANY($2)
		AND (expires_at IS NULL OR expires_at > NOW())`

	rows, err := r.db.QueryContext(ctx, query, memberID, dates)
	if err != nil {
		return nil, fmt.Errorf("query rating history for dates: %w", err)
	}
	defer rows.Close()

	result := make(map[string]*model.PlayerInfo)
	for rows.Next() {
		var ratingDate time.Time
		var data []byte
		if err := rows.Scan(&ratingDate, &data); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		var player model.PlayerInfo
		if err := json.Unmarshal(data, &player); err != nil {
			return nil, fmt.Errorf("unmarshal player data: %w", err)
		}

		dateStr := ratingDate.Format("2006-01-02")
		result[dateStr] = &player
	}

	return result, rows.Err()
}

// Save stores a player in the cache
func (r *PlayerRepository) Save(ctx context.Context, player *model.PlayerInfo, ratingDate time.Time, expiresAt *time.Time) error {
	data, err := json.Marshal(player)
	if err != nil {
		return fmt.Errorf("marshal player data: %w", err)
	}

	var eloStandard, eloRapid, eloBlitz, laskRating *int
	if player.Elo != nil {
		eloStandard = &player.Elo.Rating
		eloRapid = &player.Elo.RapidRating
		eloBlitz = &player.Elo.BlitzRating
	}
	if player.Lask != nil {
		laskRating = &player.Lask.Rating
	}

	var clubID, fideID *int
	if player.ClubID != 0 {
		clubID = &player.ClubID
	}
	if player.FideID != 0 {
		fideID = &player.FideID
	}

	query := `
		INSERT INTO player_cache (
			member_id, rating_date, first_name, last_name, club, club_id, fide_id,
			elo_standard, elo_rapid, elo_blitz, lask_rating, data, fetched_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
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
			expires_at = EXCLUDED.expires_at`

	_, err = r.db.ExecContext(ctx, query,
		player.ID, ratingDate, player.FirstName, player.LastName,
		player.Club, clubID, fideID,
		eloStandard, eloRapid, eloBlitz, laskRating,
		data, time.Now(), expiresAt)

	if err != nil {
		return fmt.Errorf("insert player cache: %w", err)
	}

	return nil
}

// DeleteExpired removes expired cache entries
func (r *PlayerRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM player_cache WHERE expires_at IS NOT NULL AND expires_at < NOW()`)
	if err != nil {
		return 0, fmt.Errorf("delete expired cache: %w", err)
	}
	return result.RowsAffected()
}
