package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/msvens/mchess/internal/upstream"
)

// TournamentHandler handles tournament-related requests (pass-through)
type TournamentHandler struct {
	client *upstream.Client
}

// NewTournamentHandler creates a new tournament handler
func NewTournamentHandler(client *upstream.Client) *TournamentHandler {
	return &TournamentHandler{client: client}
}

// GetTournament returns tournament by ID
// @Summary Get tournament by ID
// @Description Get tournament information by tournament ID
// @Tags tournament
// @Produce json
// @Param id path int true "Tournament ID"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournament/tournament/id/{id} [get]
func (h *TournamentHandler) GetTournament(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid tournament id")
		return
	}

	path := fmt.Sprintf("/tournament/tournament/id/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetTournamentFromGroup returns tournament by group ID
// @Summary Get tournament by group ID
// @Description Get tournament information from a group ID
// @Tags tournament
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournament/group/id/{id} [get]
func (h *TournamentHandler) GetTournamentFromGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	path := fmt.Sprintf("/tournament/group/id/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetTournamentFromClass returns tournament by class ID
// @Summary Get tournament by class ID
// @Description Get tournament information from a class/division ID
// @Tags tournament
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournament/class/id/{id} [get]
func (h *TournamentHandler) GetTournamentFromClass(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid class id")
		return
	}

	path := fmt.Sprintf("/tournament/class/id/%d", id)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// SearchTournamentGroups searches for tournament groups
// @Summary Search tournament groups
// @Description Search for tournament groups by name or location
// @Tags tournament
// @Produce json
// @Param searchWord path string true "Search word"
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
// @Router /tournament/group/search/{searchWord} [get]
func (h *TournamentHandler) SearchTournamentGroups(w http.ResponseWriter, r *http.Request) {
	searchWord := chi.URLParam(r, "searchWord")

	path := fmt.Sprintf("/tournament/group/search/%s", url.PathEscape(searchWord))
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetComingTournaments returns upcoming tournaments
// @Summary Get upcoming tournaments
// @Description Get all upcoming tournaments
// @Tags tournament
// @Produce json
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
// @Router /tournament/group/coming [get]
func (h *TournamentHandler) GetComingTournaments(w http.ResponseWriter, r *http.Request) {
	data, err := h.client.GetRaw(r.Context(), "/tournament/group/coming")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetComingTournamentsByDistrict returns upcoming tournaments for a district
// @Summary Get upcoming tournaments by district
// @Description Get upcoming tournaments filtered by district
// @Tags tournament
// @Produce json
// @Param districtid path int true "District ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournament/group/coming/{districtid} [get]
func (h *TournamentHandler) GetComingTournamentsByDistrict(w http.ResponseWriter, r *http.Request) {
	districtIDStr := chi.URLParam(r, "districtid")
	districtID, err := strconv.Atoi(districtIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid district id")
		return
	}

	path := fmt.Sprintf("/tournament/group/coming/%d", districtID)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// SearchUpdatedTournaments returns tournaments updated between dates
// @Summary Search updated tournaments
// @Description Search for tournaments with results updated between dates
// @Tags tournament
// @Produce json
// @Param startdate path string true "Start date (ISO 8601)"
// @Param enddate path string true "End date (ISO 8601)"
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
// @Router /tournament/tournament/updated/{startdate}/{enddate} [get]
func (h *TournamentHandler) SearchUpdatedTournaments(w http.ResponseWriter, r *http.Request) {
	startDate := chi.URLParam(r, "startdate")
	endDate := chi.URLParam(r, "enddate")

	path := fmt.Sprintf("/tournament/tournament/updated/%s/%s", startDate, endDate)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// SearchUpdatedTournamentsByDistrict returns tournaments updated between dates for a district
// @Summary Search updated tournaments by district
// @Description Search for tournaments with results updated between dates, filtered by district
// @Tags tournament
// @Produce json
// @Param startdate path string true "Start date (ISO 8601)"
// @Param enddate path string true "End date (ISO 8601)"
// @Param districtid path int true "District ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournament/tournament/updated/{startdate}/{enddate}/{districtid} [get]
func (h *TournamentHandler) SearchUpdatedTournamentsByDistrict(w http.ResponseWriter, r *http.Request) {
	startDate := chi.URLParam(r, "startdate")
	endDate := chi.URLParam(r, "enddate")
	districtIDStr := chi.URLParam(r, "districtid")

	districtID, err := strconv.Atoi(districtIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid district id")
		return
	}

	path := fmt.Sprintf("/tournament/tournament/updated/%s/%s/%d", startDate, endDate, districtID)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// SearchUpdatedGroups returns groups updated between dates
// @Summary Search updated groups
// @Description Search for groups with results updated between dates
// @Tags tournament
// @Produce json
// @Param startdate path string true "Start date (ISO 8601)"
// @Param enddate path string true "End date (ISO 8601)"
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
// @Router /tournament/group/updated/{startdate}/{enddate} [get]
func (h *TournamentHandler) SearchUpdatedGroups(w http.ResponseWriter, r *http.Request) {
	startDate := chi.URLParam(r, "startdate")
	endDate := chi.URLParam(r, "enddate")

	path := fmt.Sprintf("/tournament/group/updated/%s/%s", startDate, endDate)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// SearchUpdatedGroupsByDistrict returns groups updated between dates for a district
// @Summary Search updated groups by district
// @Description Search for groups with results updated between dates, filtered by district
// @Tags tournament
// @Produce json
// @Param startdate path string true "Start date (ISO 8601)"
// @Param enddate path string true "End date (ISO 8601)"
// @Param districtid path int true "District ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tournament/group/updated/{startdate}/{enddate}/{districtid} [get]
func (h *TournamentHandler) SearchUpdatedGroupsByDistrict(w http.ResponseWriter, r *http.Request) {
	startDate := chi.URLParam(r, "startdate")
	endDate := chi.URLParam(r, "enddate")
	districtIDStr := chi.URLParam(r, "districtid")

	districtID, err := strconv.Atoi(districtIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid district id")
		return
	}

	path := fmt.Sprintf("/tournament/group/updated/%s/%s/%d", startDate, endDate, districtID)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}
