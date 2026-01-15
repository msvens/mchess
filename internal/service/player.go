package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/msvens/mchess/internal/config"
	"github.com/msvens/mchess/internal/model"
	"github.com/msvens/mchess/internal/repository"
	"github.com/msvens/mchess/internal/upstream"
)

// PlayerService handles player-related business logic
type PlayerService struct {
	repo     *repository.PlayerRepository
	upstream *upstream.Client
	cacheTTL time.Duration
}

// NewPlayerService creates a new player service
func NewPlayerService(repo *repository.PlayerRepository, client *upstream.Client, cfg *config.Config) *PlayerService {
	return &PlayerService{
		repo:     repo,
		upstream: client,
		cacheTTL: cfg.Cache.TTL,
	}
}

// GetPlayer retrieves a player, checking cache first then upstream
func (s *PlayerService) GetPlayer(ctx context.Context, memberID int, date time.Time) (*model.PlayerInfo, error) {
	ratingDate := normalizeToMonthStart(date)

	// Check cache first
	cached, err := s.repo.Get(ctx, memberID, ratingDate)
	if err != nil {
		slog.Error("Cache lookup failed", "error", err, "memberID", memberID)
		// Continue to upstream on cache error
	}
	if cached != nil {
		slog.Debug("Cache hit", "memberID", memberID, "date", ratingDate)
		return cached, nil
	}

	slog.Debug("Cache miss, fetching from upstream", "memberID", memberID, "date", ratingDate)

	// Fetch from upstream
	player, err := s.upstream.GetPlayer(ctx, memberID, date.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("upstream fetch: %w", err)
	}

	// Determine TTL
	expiresAt := s.determineTTL(ratingDate)

	// Store in cache
	if err := s.repo.Save(ctx, player, ratingDate, expiresAt); err != nil {
		slog.Error("Failed to cache player", "error", err, "memberID", memberID)
		// Continue even if caching fails
	}

	return player, nil
}

// GetPlayers retrieves multiple players in batch
func (s *PlayerService) GetPlayers(ctx context.Context, memberIDs []int, date time.Time) (*model.PlayersResponse, error) {
	ratingDate := normalizeToMonthStart(date)

	// Check cache for all IDs
	cached, err := s.repo.GetBatch(ctx, memberIDs, ratingDate)
	if err != nil {
		slog.Error("Batch cache lookup failed", "error", err)
		cached = make(map[int]*model.PlayerInfo)
	}

	// Find missing IDs
	var missingIDs []int
	for _, id := range memberIDs {
		if _, found := cached[id]; !found {
			missingIDs = append(missingIDs, id)
		}
	}

	slog.Debug("Batch lookup", "total", len(memberIDs), "cached", len(cached), "missing", len(missingIDs))

	// Fetch missing from upstream in parallel
	if len(missingIDs) > 0 {
		fetched, errors := s.fetchPlayersParallel(ctx, missingIDs, date, ratingDate)
		for id, player := range fetched {
			cached[id] = player
		}

		// Build response with results and errors
		response := &model.PlayersResponse{
			Players: make([]model.PlayerInfo, 0, len(memberIDs)),
			Errors:  errors,
		}

		// Maintain order from input
		for _, id := range memberIDs {
			if player, found := cached[id]; found {
				response.Players = append(response.Players, *player)
			}
		}

		return response, nil
	}

	// All were cached
	response := &model.PlayersResponse{
		Players: make([]model.PlayerInfo, 0, len(memberIDs)),
	}
	for _, id := range memberIDs {
		if player, found := cached[id]; found {
			response.Players = append(response.Players, *player)
		}
	}

	return response, nil
}

// fetchPlayersParallel fetches multiple players from upstream concurrently
func (s *PlayerService) fetchPlayersParallel(ctx context.Context, memberIDs []int, date time.Time, ratingDate time.Time) (map[int]*model.PlayerInfo, []model.PlayerError) {
	results := make(map[int]*model.PlayerInfo)
	var errors []model.PlayerError
	var mu sync.Mutex
	var wg sync.WaitGroup

	dateStr := date.Format("2006-01-02")
	expiresAt := s.determineTTL(ratingDate)

	for _, id := range memberIDs {
		wg.Add(1)
		go func(memberID int) {
			defer wg.Done()

			player, err := s.upstream.GetPlayer(ctx, memberID, dateStr)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				slog.Warn("Failed to fetch player", "memberID", memberID, "error", err)
				errors = append(errors, model.PlayerError{
					ID:    memberID,
					Error: err.Error(),
				})
				return
			}

			// Cache the player
			if saveErr := s.repo.Save(ctx, player, ratingDate, expiresAt); saveErr != nil {
				slog.Error("Failed to cache player", "error", saveErr, "memberID", memberID)
			}

			results[memberID] = player
		}(id)
	}

	wg.Wait()
	return results, errors
}

// GetPlayerRatings retrieves rating history for a player
func (s *PlayerService) GetPlayerRatings(ctx context.Context, memberID int, fromDate, toDate time.Time) (*model.RatingHistoryResponse, error) {
	// Generate list of months in range
	dates := generateMonthRange(fromDate, toDate)

	// Check cache for all dates
	cached, err := s.repo.GetRatingHistoryForDates(ctx, memberID, dates)
	if err != nil {
		slog.Error("Rating history cache lookup failed", "error", err)
		cached = make(map[string]*model.PlayerInfo)
	}

	// Find missing dates
	var missingDates []time.Time
	for _, d := range dates {
		dateStr := d.Format("2006-01-02")
		if _, found := cached[dateStr]; !found {
			missingDates = append(missingDates, d)
		}
	}

	slog.Debug("Rating history lookup", "total", len(dates), "cached", len(cached), "missing", len(missingDates))

	// Fetch missing from upstream
	for _, d := range missingDates {
		player, err := s.upstream.GetPlayer(ctx, memberID, d.Format("2006-01-02"))
		if err != nil {
			slog.Warn("Failed to fetch player for date", "memberID", memberID, "date", d, "error", err)
			continue
		}

		expiresAt := s.determineTTL(d)
		if saveErr := s.repo.Save(ctx, player, d, expiresAt); saveErr != nil {
			slog.Error("Failed to cache player", "error", saveErr, "memberID", memberID)
		}

		dateStr := d.Format("2006-01-02")
		cached[dateStr] = player
	}

	// Build response in chronological order (newest first)
	response := &model.RatingHistoryResponse{
		PlayerID: memberID,
		Ratings:  make([]model.PlayerInfo, 0, len(dates)),
	}

	// Reverse order (newest first)
	for i := len(dates) - 1; i >= 0; i-- {
		dateStr := dates[i].Format("2006-01-02")
		if player, found := cached[dateStr]; found {
			response.Ratings = append(response.Ratings, *player)
		}
	}

	return response, nil
}

// SearchPlayers proxies search to upstream (no caching)
func (s *PlayerService) SearchPlayers(ctx context.Context, firstName, lastName string) ([]model.PlayerInfo, error) {
	return s.upstream.SearchPlayers(ctx, firstName, lastName)
}

// GetPlayerByFideID retrieves a player by their FIDE ID
func (s *PlayerService) GetPlayerByFideID(ctx context.Context, fideID int, date time.Time) (*model.PlayerInfo, error) {
	ratingDate := normalizeToMonthStart(date)

	// Check cache first (by FIDE ID)
	cached, err := s.repo.GetByFideID(ctx, fideID, ratingDate)
	if err != nil {
		slog.Error("Cache lookup by FIDE ID failed", "error", err, "fideID", fideID)
		// Continue to upstream on cache error
	}
	if cached != nil {
		slog.Debug("Cache hit by FIDE ID", "fideID", fideID, "date", ratingDate)
		return cached, nil
	}

	slog.Debug("Cache miss by FIDE ID, fetching from upstream", "fideID", fideID, "date", ratingDate)

	// Fetch from upstream
	player, err := s.upstream.GetPlayerByFideID(ctx, fideID, date.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("upstream fetch by FIDE ID: %w", err)
	}

	// Determine TTL
	expiresAt := s.determineTTL(ratingDate)

	// Store in cache (using member ID from response)
	if err := s.repo.Save(ctx, player, ratingDate, expiresAt); err != nil {
		slog.Error("Failed to cache player", "error", err, "fideID", fideID)
		// Continue even if caching fails
	}

	return player, nil
}

// determineTTL calculates the cache expiration time based on the rating date
func (s *PlayerService) determineTTL(ratingDate time.Time) *time.Time {
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Historical data: never expires
	if ratingDate.Before(currentMonth) {
		return nil
	}

	// Current month: use configured TTL
	expires := now.Add(s.cacheTTL)
	return &expires
}

// normalizeToMonthStart returns the first day of the month for the given date
func normalizeToMonthStart(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// generateMonthRange generates a slice of first-of-month dates between from and to
func generateMonthRange(from, to time.Time) []time.Time {
	from = normalizeToMonthStart(from)
	to = normalizeToMonthStart(to)

	var dates []time.Time
	for d := from; !d.After(to); d = d.AddDate(0, 1, 0) {
		dates = append(dates, d)
	}
	return dates
}
