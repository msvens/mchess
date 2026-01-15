package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/msvens/mchess/internal/service"
	"github.com/msvens/mchess/internal/upstream"
)

// ErrorResponse represents an API error response
// @Description API error response
// @name ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"invalid player id"`
	Code  int    `json:"code" example:"400"`
}

// PlayerHandler handles player-related HTTP requests
type PlayerHandler struct {
	service *service.PlayerService
	client  *upstream.Client // For pass-through when no caching needed
}

// NewPlayerHandler creates a new player handler
func NewPlayerHandler(service *service.PlayerService, client *upstream.Client) *PlayerHandler {
	return &PlayerHandler{service: service, client: client}
}

// GetPlayer handles GET /player/{id}/date/{date}
// @Summary Get player by member ID
// @Description Get player information by their Swedish Chess Federation member ID (cached)
// @Tags player
// @Produce json
// @Param id path int true "Member ID (Swedish Chess Federation ID)"
// @Param date path string true "Rating date (YYYY-MM-DD)"
// @Success 200 {object} model.PlayerInfo "Player information"
// @Failure 400 {object} ErrorResponse "Invalid player ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /player/{id}/date/{date} [get]
func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid player id")
		return
	}

	dateStr := chi.URLParam(r, "date")
	date := parseDate(dateStr)

	player, err := h.service.GetPlayer(r.Context(), id, date)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, player)
}

// GetPlayerByFideID handles GET /player/fideid/{id}/date/{date}
// @Summary Get player by FIDE ID
// @Description Get player information using their FIDE ID (cached)
// @Tags player
// @Produce json
// @Param id path int true "FIDE ID"
// @Param date path string true "Rating date (YYYY-MM-DD)"
// @Success 200 {object} model.PlayerInfo "Player information"
// @Failure 400 {object} ErrorResponse "Invalid FIDE ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /player/fideid/{id}/date/{date} [get]
func (h *PlayerHandler) GetPlayerByFideID(w http.ResponseWriter, r *http.Request) {
	fideIDStr := chi.URLParam(r, "id")
	fideID, err := strconv.Atoi(fideIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid fide id")
		return
	}

	dateStr := chi.URLParam(r, "date")
	date := parseDate(dateStr)

	player, err := h.service.GetPlayerByFideID(r.Context(), fideID, date)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, player)
}

// SearchPlayers handles GET /player/fornamn/{fornamn}/efternamn/{efternamn}
// @Summary Search players by name
// @Description Search for players by first and last name
// @Tags player
// @Produce json
// @Param fornamn path string true "First name"
// @Param efternamn path string true "Last name"
// @Success 200 {array} model.PlayerInfo "List of matching players"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /player/fornamn/{fornamn}/efternamn/{efternamn} [get]
func (h *PlayerHandler) SearchPlayers(w http.ResponseWriter, r *http.Request) {
	firstName := chi.URLParam(r, "fornamn")
	lastName := chi.URLParam(r, "efternamn")

	players, err := h.service.SearchPlayers(r.Context(), firstName, lastName)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, players)
}

// GetPlayers handles GET /player/batch?ids=1,2,3&date=2024-06-01
// @Summary Batch fetch multiple players
// @Description Fetch multiple players in a single request by providing comma-separated member IDs. Maximum 100 IDs per request. All lookups are cached.
// @Tags player
// @Produce json
// @Param ids query string true "Comma-separated member IDs (max 100)" example:"12345,67890,11111"
// @Param date query string false "Rating date (YYYY-MM-DD), defaults to current date"
// @Success 200 {object} model.PlayersResponse "Players with any errors for failed lookups"
// @Failure 400 {object} ErrorResponse "Invalid request (missing/invalid IDs, too many IDs)"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /player/batch [get]
func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	idsStr := r.URL.Query().Get("ids")
	if idsStr == "" {
		WriteError(w, http.StatusBadRequest, "ids parameter is required")
		return
	}

	ids, err := parseIDs(idsStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid ids format")
		return
	}

	if len(ids) == 0 {
		WriteError(w, http.StatusBadRequest, "at least one id is required")
		return
	}

	if len(ids) > 100 {
		WriteError(w, http.StatusBadRequest, "maximum 100 ids allowed")
		return
	}

	date := parseDate(r.URL.Query().Get("date"))

	response, err := h.service.GetPlayers(r.Context(), ids, date)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

// GetPlayerRatings handles GET /player/{id}/ratings
// @Summary Get player rating history
// @Description Get historical rating data for a player within a date range. Use either from/to dates or the months parameter. Defaults to last 12 months if no range specified. All lookups are cached.
// @Tags player
// @Produce json
// @Param id path int true "Member ID"
// @Param from query string false "Start date (YYYY-MM-DD or YYYY-MM)"
// @Param to query string false "End date (YYYY-MM-DD or YYYY-MM)"
// @Param months query int false "Number of months back from today (alternative to from/to)"
// @Success 200 {object} model.RatingHistoryResponse "Rating history sorted newest first"
// @Failure 400 {object} ErrorResponse "Invalid player ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /player/{id}/ratings [get]
func (h *PlayerHandler) GetPlayerRatings(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid player id")
		return
	}

	// Parse date range
	fromDate, toDate := parseDateRange(r)

	response, err := h.service.GetPlayerRatings(r.Context(), id, fromDate, toDate)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

// Internal helper functions

func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}

	// Try YYYY-MM-DD format
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t
	}

	// Try YYYY-MM format
	if t, err := time.Parse("2006-01", dateStr); err == nil {
		return t
	}

	return time.Now()
}

func parseIDs(idsStr string) ([]int, error) {
	parts := strings.Split(idsStr, ",")
	ids := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func parseDateRange(r *http.Request) (from, to time.Time) {
	now := time.Now()

	// Check for "months" parameter first
	if monthsStr := r.URL.Query().Get("months"); monthsStr != "" {
		months, err := strconv.Atoi(monthsStr)
		if err == nil && months > 0 {
			from = now.AddDate(0, -months, 0)
			to = now
			return
		}
	}

	// Parse from/to parameters
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr != "" {
		from = parseDate(fromStr)
	} else {
		from = now.AddDate(0, -12, 0) // Default: 12 months back
	}

	if toStr != "" {
		to = parseDate(toStr)
	} else {
		to = now
	}

	return
}
