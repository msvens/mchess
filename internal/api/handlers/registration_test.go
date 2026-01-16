package handlers_test

import (
	"net/http"
	"testing"

	"github.com/msvens/mchess/internal/api/handlers"
)

func TestRegistrationHandler(t *testing.T) {
	client := NewTestClient(t)
	handler := handlers.NewRegistrationHandler(client)

	t.Run("GetTeamRegistration", func(t *testing.T) {
		t.Run("ValidIDs_ReturnsSuccess", func(t *testing.T) {
			// Scenario 1: Correct input - existing tournament and club IDs
			// TODO: Find valid tournament and club IDs to test with
			t.Skip("TODO: Implement with valid tournament and club IDs")
		})

		t.Run("InvalidTournamentID_ReturnsError", func(t *testing.T) {
			// Scenario 2: Wrong input - non-existent tournament ID
			t.Skip("TODO: Implement - check upstream behavior")
		})

		t.Run("InvalidClubID_ReturnsError", func(t *testing.T) {
			// Scenario 2: Wrong input - non-existent club ID
			t.Skip("TODO: Implement - check upstream behavior")
		})

		t.Run("MalformedTournamentID_Returns400", func(t *testing.T) {
			// Scenario 4: Invalid input format
			rr := MakeRequest(t, handler.GetTeamRegistration, http.MethodGet,
				"/tournamentteamregistration/tournament/abc/club/100",
				map[string]string{"id": "abc", "clubid": "100"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})

		t.Run("MalformedClubID_Returns400", func(t *testing.T) {
			// Scenario 4: Invalid input format
			rr := MakeRequest(t, handler.GetTeamRegistration, http.MethodGet,
				"/tournamentteamregistration/tournament/1/club/abc",
				map[string]string{"id": "1", "clubid": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})
}
