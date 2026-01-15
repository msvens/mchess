package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/msvens/mchess/internal/upstream"
)

// RegistrationHandler handles team registration requests (pass-through)
type RegistrationHandler struct {
	client *upstream.Client
}

// NewRegistrationHandler creates a new registration handler
func NewRegistrationHandler(client *upstream.Client) *RegistrationHandler {
	return &RegistrationHandler{client: client}
}

// GetTeamRegistration returns team registration for a tournament and club
// @Summary Get team registration
// @Description Get player registrations for a tournament and club
// @Tags tournamentteamregistration
// @Produce json
// @Param id path int true "Tournament ID"
// @Param clubid path int true "Club ID"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournamentteamregistration/tournament/{id}/club/{clubid} [get]
func (h *RegistrationHandler) GetTeamRegistration(w http.ResponseWriter, r *http.Request) {
	tournamentIDStr := chi.URLParam(r, "id")
	tournamentID, err := strconv.Atoi(tournamentIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid tournament id")
		return
	}

	clubIDStr := chi.URLParam(r, "clubid")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid club id")
		return
	}

	path := fmt.Sprintf("/tournamentteamregistration/tournament/%d/club/%d", tournamentID, clubID)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}