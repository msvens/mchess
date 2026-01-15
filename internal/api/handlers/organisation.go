package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/msvens/mchess/internal/upstream"
)

// OrganisationHandler handles organisation-related requests (pass-through)
type OrganisationHandler struct {
	client *upstream.Client
}

// NewOrganisationHandler creates a new organisation handler
func NewOrganisationHandler(client *upstream.Client) *OrganisationHandler {
	return &OrganisationHandler{client: client}
}

// GetFederation returns Swedish Chess Federation info
// @Summary Get federation info
// @Description Get Swedish Chess Federation (SSF) information
// @Tags organisation
// @Produce json
// @Success 200 {object} object
// @Failure 500 {object} ErrorResponse
// @Router /organisation/federation [get]
func (h *OrganisationHandler) GetFederation(w http.ResponseWriter, r *http.Request) {
	data, err := h.client.GetRaw(r.Context(), "/organisation/federation")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetDistricts returns all districts
// @Summary Get all districts
// @Description Get all chess districts
// @Tags organisation
// @Produce json
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
// @Router /organisation/districts [get]
func (h *OrganisationHandler) GetDistricts(w http.ResponseWriter, r *http.Request) {
	data, err := h.client.GetRaw(r.Context(), "/organisation/districts")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetClubsInDistrict returns clubs in a district
// @Summary Get clubs in district
// @Description Get all clubs in a specific district
// @Tags organisation
// @Produce json
// @Param districtid path int true "District ID"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organisation/district/clubs/{districtid} [get]
func (h *OrganisationHandler) GetClubsInDistrict(w http.ResponseWriter, r *http.Request) {
	districtIDStr := chi.URLParam(r, "districtid")
	districtID, err := strconv.Atoi(districtIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid district id")
		return
	}

	path := fmt.Sprintf("/organisation/district/clubs/%d", districtID)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetClub returns a specific club
// @Summary Get club
// @Description Get a specific club by ID
// @Tags organisation
// @Produce json
// @Param clubid path int true "Club ID"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organisation/club/{clubid} [get]
func (h *OrganisationHandler) GetClub(w http.ResponseWriter, r *http.Request) {
	clubIDStr := chi.URLParam(r, "clubid")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid club id")
		return
	}

	path := fmt.Sprintf("/organisation/club/%d", clubID)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// ClubNameExists checks if a club name exists
// @Summary Check if club name exists
// @Description Check if a club name exists (other than for the given club ID)
// @Tags organisation
// @Produce json
// @Param name path string true "Club name"
// @Param id path int true "Club ID to exclude"
// @Success 200 {boolean} bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organisation/club/exists/{name}/{id} [get]
func (h *OrganisationHandler) ClubNameExists(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	clubIDStr := chi.URLParam(r, "id")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid club id")
		return
	}

	path := fmt.Sprintf("/organisation/club/exists/%s/%d", url.PathEscape(name), clubID)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}