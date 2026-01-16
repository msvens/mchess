package handlers_test

import (
	"net/http"
	"testing"

	"github.com/msvens/mchess/internal/api/handlers"
)

func TestOrganisationHandler(t *testing.T) {
	client := NewTestClient(t)
	handler := handlers.NewOrganisationHandler(client)

	t.Run("GetFederation", func(t *testing.T) {
		t.Run("ReturnsSuccess", func(t *testing.T) {
			// Scenario 1: Correct input - federation info
			rr := MakeRequest(t, handler.GetFederation, http.MethodGet,
				"/organisation/federation",
				nil)

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
			AssertBodyNotEmpty(t, rr)
		})
	})

	t.Run("GetDistricts", func(t *testing.T) {
		t.Run("ReturnsSuccess", func(t *testing.T) {
			// Scenario 1: Returns list of districts
			rr := MakeRequest(t, handler.GetDistricts, http.MethodGet,
				"/organisation/districts",
				nil)

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
			AssertBodyNotEmpty(t, rr)
		})
	})

	t.Run("GetClubsInDistrict", func(t *testing.T) {
		t.Run("ValidDistrictID_ReturnsSuccess", func(t *testing.T) {
			// Scenario 1: Correct input - district 1 (Stockholm)
			rr := MakeRequest(t, handler.GetClubsInDistrict, http.MethodGet,
				"/organisation/district/clubs/1",
				map[string]string{"districtid": "1"})

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
			AssertBodyNotEmpty(t, rr)
		})

		t.Run("InvalidDistrictID_ReturnsError", func(t *testing.T) {
			// Scenario 2: Wrong input - non-existent district ID
			// TODO: Check what upstream returns for invalid district
			t.Skip("TODO: Implement - check upstream behavior")
		})

		t.Run("MalformedDistrictID_Returns400", func(t *testing.T) {
			// Scenario 4: Invalid input format
			rr := MakeRequest(t, handler.GetClubsInDistrict, http.MethodGet,
				"/organisation/district/clubs/abc",
				map[string]string{"districtid": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
			AssertBodyContains(t, rr, "invalid")
		})
	})

	t.Run("GetClub", func(t *testing.T) {
		t.Run("ValidClubID_ReturnsSuccess", func(t *testing.T) {
			// Scenario 1: Correct input - use a known club ID
			// TODO: Find a valid club ID to test with
			t.Skip("TODO: Implement with valid club ID")
		})

		t.Run("InvalidClubID_ReturnsError", func(t *testing.T) {
			// Scenario 2: Wrong input - non-existent club ID
			t.Skip("TODO: Implement - check upstream behavior")
		})

		t.Run("MalformedClubID_Returns400", func(t *testing.T) {
			// Scenario 4: Invalid input format
			rr := MakeRequest(t, handler.GetClub, http.MethodGet,
				"/organisation/club/abc",
				map[string]string{"clubid": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("ClubNameExists", func(t *testing.T) {
		t.Run("ValidParams_ReturnsSuccess", func(t *testing.T) {
			// TODO: Find valid test params
			t.Skip("TODO: Implement")
		})

		t.Run("MalformedClubID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.ClubNameExists, http.MethodGet,
				"/organisation/club/exists/TestClub/abc",
				map[string]string{"name": "TestClub", "id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})
}
