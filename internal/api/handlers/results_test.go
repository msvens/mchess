package handlers_test

import (
	"net/http"
	"testing"

	"github.com/msvens/mchess/internal/api/handlers"
)

func TestResultsHandler(t *testing.T) {
	client := NewTestClient(t)
	handler := handlers.NewResultsHandler(client)

	t.Run("GetResultTable", func(t *testing.T) {
		t.Run("ValidGroupID_ReturnsSuccess", func(t *testing.T) {
			// Scenario 1: Correct input - existing group ID
			// TODO: Find a valid group ID to test with
			t.Skip("TODO: Implement with valid group ID")
		})

		t.Run("InvalidGroupID_ReturnsError", func(t *testing.T) {
			// Scenario 2: Wrong input - non-existent group ID
			t.Skip("TODO: Implement - check upstream behavior")
		})

		t.Run("MalformedGroupID_Returns400", func(t *testing.T) {
			// Scenario 4: Invalid input format
			rr := MakeRequest(t, handler.GetResultTable, http.MethodGet,
				"/tournamentresults/table/id/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetMemberTableResults", func(t *testing.T) {
		t.Run("ValidMemberID_ReturnsSuccess", func(t *testing.T) {
			// TODO: Find a valid member ID to test with
			t.Skip("TODO: Implement with valid member ID")
		})

		t.Run("MalformedMemberID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetMemberTableResults, http.MethodGet,
				"/tournamentresults/table/memberid/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetRoundResults", func(t *testing.T) {
		t.Run("ValidGroupID_ReturnsSuccess", func(t *testing.T) {
			t.Skip("TODO: Implement with valid group ID")
		})

		t.Run("MalformedGroupID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetRoundResults, http.MethodGet,
				"/tournamentresults/roundresults/id/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetTeamResultTable", func(t *testing.T) {
		t.Run("ValidGroupID_ReturnsSuccess", func(t *testing.T) {
			t.Skip("TODO: Implement with valid group ID")
		})

		t.Run("MalformedGroupID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetTeamResultTable, http.MethodGet,
				"/tournamentresults/team/table/id/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetTeamRoundResults", func(t *testing.T) {
		t.Run("ValidGroupID_ReturnsSuccess", func(t *testing.T) {
			t.Skip("TODO: Implement with valid group ID")
		})

		t.Run("MalformedGroupID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetTeamRoundResults, http.MethodGet,
				"/tournamentresults/team/roundresults/id/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetTeamRoundResultsForMember", func(t *testing.T) {
		t.Run("ValidIDs_ReturnsSuccess", func(t *testing.T) {
			t.Skip("TODO: Implement with valid IDs")
		})

		t.Run("MalformedGroupID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetTeamRoundResultsForMember, http.MethodGet,
				"/tournamentresults/team/roundresults/id/abc/memberid/123",
				map[string]string{"id": "abc", "memberid": "123"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})

		t.Run("MalformedMemberID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetTeamRoundResultsForMember, http.MethodGet,
				"/tournamentresults/team/roundresults/id/123/memberid/abc",
				map[string]string{"id": "123", "memberid": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetMemberGames", func(t *testing.T) {
		t.Run("ValidMemberID_ReturnsSuccess", func(t *testing.T) {
			t.Skip("TODO: Implement with valid member ID")
		})

		t.Run("MalformedMemberID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetMemberGames, http.MethodGet,
				"/tournamentresults/game/memberid/abc",
				map[string]string{"id": "abc"})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})
}
