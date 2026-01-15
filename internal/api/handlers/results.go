package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/msvens/mchess/internal/upstream"
)

// ResultsHandler handles tournament results requests (pass-through)
type ResultsHandler struct {
	client *upstream.Client
}

// NewResultsHandler creates a new results handler
func NewResultsHandler(client *upstream.Client) *ResultsHandler {
	return &ResultsHandler{client: client}
}

// GetResultTable returns individual tournament table
// @Summary Get tournament table
// @Description Get individual tournament standings by group ID
// @Tags tournamentresults
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournamentresults/table/id/{id} [get]
func (h *ResultsHandler) GetResultTable(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	path := fmt.Sprintf("/tournamentresults/table/id/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetMemberTableResults returns member's tournament results
// @Summary Get member tournament results
// @Description Get all tournament results for a specific member
// @Tags tournamentresults
// @Produce json
// @Param id path int true "Member ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournamentresults/table/memberid/{id} [get]
func (h *ResultsHandler) GetMemberTableResults(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid member id")
		return
	}

	path := fmt.Sprintf("/tournamentresults/table/memberid/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetRoundResults returns round results
// @Summary Get round results
// @Description Get round-by-round results for a group
// @Tags tournamentresults
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournamentresults/roundresults/id/{id} [get]
func (h *ResultsHandler) GetRoundResults(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	path := fmt.Sprintf("/tournamentresults/roundresults/id/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetTeamResultTable returns team tournament table
// @Summary Get team tournament table
// @Description Get team tournament standings by group ID
// @Tags tournamentresults
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournamentresults/team/table/id/{id} [get]
func (h *ResultsHandler) GetTeamResultTable(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	path := fmt.Sprintf("/tournamentresults/team/table/id/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetTeamRoundResults returns team round results
// @Summary Get team round results
// @Description Get team round-by-round results for a group
// @Tags tournamentresults
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournamentresults/team/roundresults/id/{id} [get]
func (h *ResultsHandler) GetTeamRoundResults(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	path := fmt.Sprintf("/tournamentresults/team/roundresults/id/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetTeamRoundResultsForMember returns team round results for a member
// @Summary Get team round results for member
// @Description Get team round-by-round results for a specific member
// @Tags tournamentresults
// @Produce json
// @Param id path int true "Group ID"
// @Param memberid path int true "Member ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournamentresults/team/roundresults/id/{id}/memberid/{memberid} [get]
func (h *ResultsHandler) GetTeamRoundResultsForMember(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	memberIDStr := chi.URLParam(r, "memberid")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid member id")
		return
	}

	path := fmt.Sprintf("/tournamentresults/team/roundresults/id/%d/memberid/%d", id, memberID)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetMemberGames returns games for a member
// @Summary Get member games
// @Description Get all games for a specific member
// @Tags tournamentresults
// @Produce json
// @Param id path int true "Member ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournamentresults/game/memberid/{id} [get]
func (h *ResultsHandler) GetMemberGames(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid member id")
		return
	}

	path := fmt.Sprintf("/tournamentresults/game/memberid/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}