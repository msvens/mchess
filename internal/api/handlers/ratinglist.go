package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/msvens/mchess/internal/upstream"
)

// RatingListHandler handles rating list requests (pass-through)
type RatingListHandler struct {
	client *upstream.Client
}

// NewRatingListHandler creates a new rating list handler
func NewRatingListHandler(client *upstream.Client) *RatingListHandler {
	return &RatingListHandler{client: client}
}

// GetFederationRatingList returns federation-wide rating list
// @Summary Get federation rating list
// @Description Get rating list for the entire federation
// @Tags ratinglist
// @Produce json
// @Param ratingdate path string true "Rating date (YYYY-MM-DD)"
// @Param ratingtype path int true "Rating type: 1=Standard, 6=Rapid, 7=Blitz"
// @Param category path int true "Member category: 0=All, 1=Juniors, 2=Cadets, 4=Veterans, 5=Women, 6=Minors, 7=Kids"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /ratinglist/federation/date/{ratingdate}/ratingtype/{ratingtype}/category/{category} [get]
func (h *RatingListHandler) GetFederationRatingList(w http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "ratingdate")
	ratingTypeStr := chi.URLParam(r, "ratingtype")
	categoryStr := chi.URLParam(r, "category")

	ratingType, err := strconv.Atoi(ratingTypeStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid rating type")
		return
	}

	category, err := strconv.Atoi(categoryStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid category")
		return
	}

	path := fmt.Sprintf("/ratinglist/federation/date/%s/ratingtype/%d/category/%d", date, ratingType, category)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetDistrictRatingList returns district rating list
// @Summary Get district rating list
// @Description Get rating list for a specific district
// @Tags ratinglist
// @Produce json
// @Param id path int true "District ID"
// @Param ratingdate path string true "Rating date (YYYY-MM-DD)"
// @Param ratingtype path int true "Rating type: 1=Standard, 6=Rapid, 7=Blitz"
// @Param category path int true "Member category: 0=All, 1=Juniors, 2=Cadets, 4=Veterans, 5=Women, 6=Minors, 7=Kids"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /ratinglist/district/{id}/date/{ratingdate}/ratingtype/{ratingtype}/category/{category} [get]
func (h *RatingListHandler) GetDistrictRatingList(w http.ResponseWriter, r *http.Request) {
	districtIDStr := chi.URLParam(r, "id")
	date := chi.URLParam(r, "ratingdate")
	ratingTypeStr := chi.URLParam(r, "ratingtype")
	categoryStr := chi.URLParam(r, "category")

	districtID, err := strconv.Atoi(districtIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid district id")
		return
	}

	ratingType, err := strconv.Atoi(ratingTypeStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid rating type")
		return
	}

	category, err := strconv.Atoi(categoryStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid category")
		return
	}

	path := fmt.Sprintf("/ratinglist/district/%d/date/%s/ratingtype/%d/category/%d", districtID, date, ratingType, category)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}

// GetClubRatingList returns club rating list
// @Summary Get club rating list
// @Description Get rating list for a specific club
// @Tags ratinglist
// @Produce json
// @Param id path int true "Club ID"
// @Param ratingdate path string true "Rating date (YYYY-MM-DD)"
// @Param ratingtype path int true "Rating type: 1=Standard, 6=Rapid, 7=Blitz"
// @Param category path int true "Member category: 0=All, 1=Juniors, 2=Cadets, 4=Veterans, 5=Women, 6=Minors, 7=Kids"
// @Success 200 {array} object
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /ratinglist/club/{id}/date/{ratingdate}/ratingtype/{ratingtype}/category/{category} [get]
func (h *RatingListHandler) GetClubRatingList(w http.ResponseWriter, r *http.Request) {
	clubIDStr := chi.URLParam(r, "id")
	date := chi.URLParam(r, "ratingdate")
	ratingTypeStr := chi.URLParam(r, "ratingtype")
	categoryStr := chi.URLParam(r, "category")

	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid club id")
		return
	}

	ratingType, err := strconv.Atoi(ratingTypeStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid rating type")
		return
	}

	category, err := strconv.Atoi(categoryStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid category")
		return
	}

	path := fmt.Sprintf("/ratinglist/club/%d/date/%s/ratingtype/%d/category/%d", clubID, date, ratingType, category)
	data, err := h.client.GetRaw(r.Context(), path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteRawJSON(w, http.StatusOK, data)
}