package handlers_test

import (
	"net/http"
	"testing"

	"github.com/msvens/mchess/internal/api/handlers"
)

func TestTournamentHandler(t *testing.T) {
	client := NewTestClient(t)
	handler := handlers.NewTournamentHandler(client)

	t.Run("GetTournament", func(t *testing.T) {
		t.Run("ValidID_ReturnsSuccess", func(t *testing.T) {
			// Scenario 1: Correct input - existing tournament ID
			// Use a known tournament ID from schack.se
			rr := MakeRequest(t, handler.GetTournament, http.MethodGet,
				"/tournament/tournament/id/6058",
				map[string]string{"id": "6058"})

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
			AssertBodyNotEmpty(t, rr)
		})

		t.Run("InvalidID_ReturnsError", func(t *testing.T) {
			// Scenario 2: Wrong input - non-existent tournament ID
			// TODO: Determine how upstream handles this and assert accordingly
			t.Skip("TODO: Implement - check upstream behavior for non-existent ID")
		})

		t.Run("MalformedID_Returns400", func(t *testing.T) {
			// Scenario 4: Invalid input format
			rr := MakeRequest(t, handler.GetTournament, http.MethodGet,
				"/tournament/tournament/id/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
			AssertBodyContains(t, rr, "invalid")
		})
	})

	t.Run("GetTournamentFromGroup", func(t *testing.T) {
		t.Run("ValidID_ReturnsSuccess", func(t *testing.T) {
			// TODO: Find a valid group ID to test with
			t.Skip("TODO: Implement with valid group ID")
		})

		t.Run("MalformedID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetTournamentFromGroup, http.MethodGet,
				"/tournament/group/id/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetTournamentFromClass", func(t *testing.T) {
		t.Run("ValidID_ReturnsSuccess", func(t *testing.T) {
			// TODO: Find a valid class ID to test with
			t.Skip("TODO: Implement with valid class ID")
		})

		t.Run("MalformedID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetTournamentFromClass, http.MethodGet,
				"/tournament/class/id/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("SearchTournamentGroups", func(t *testing.T) {
		t.Run("ValidSearch_ReturnsResults", func(t *testing.T) {
			rr := MakeRequest(t, handler.SearchTournamentGroups, http.MethodGet,
				"/tournament/group/search/stockholm",
				map[string]string{"searchWord": "stockholm"})

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
		})
	})

	t.Run("GetComingTournaments", func(t *testing.T) {
		t.Run("ReturnsSuccess", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetComingTournaments, http.MethodGet,
				"/tournament/group/coming",
				nil)

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
		})
	})

	t.Run("GetComingTournamentsByDistrict", func(t *testing.T) {
		t.Run("ValidDistrictID_ReturnsSuccess", func(t *testing.T) {
			// District 1 = Stockholms Schackförbund
			rr := MakeRequest(t, handler.GetComingTournamentsByDistrict, http.MethodGet,
				"/tournament/group/coming/1",
				map[string]string{"districtid": "1"})

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
		})

		t.Run("MalformedDistrictID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetComingTournamentsByDistrict, http.MethodGet,
				"/tournament/group/coming/abc",
				map[string]string{"districtid": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("SearchUpdatedTournaments", func(t *testing.T) {
		t.Run("ValidDateRange_ReturnsSuccess", func(t *testing.T) {
			rr := MakeRequest(t, handler.SearchUpdatedTournaments, http.MethodGet,
				"/tournament/tournament/updated/2024-01-01/2024-12-31",
				map[string]string{"startdate": "2024-01-01", "enddate": "2024-12-31"})

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
		})
	})

	t.Run("SearchUpdatedGroups", func(t *testing.T) {
		t.Run("ValidDateRange_ReturnsSuccess", func(t *testing.T) {
			rr := MakeRequest(t, handler.SearchUpdatedGroups, http.MethodGet,
				"/tournament/group/updated/2024-01-01/2024-12-31",
				map[string]string{"startdate": "2024-01-01", "enddate": "2024-12-31"})

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
		})
	})
}
