package handlers_test

import (
	"net/http"
	"testing"

	"github.com/msvens/mchess/internal/api/handlers"
)

func TestRatingListHandler(t *testing.T) {
	client := NewTestClient(t)
	handler := handlers.NewRatingListHandler(client)

	t.Run("GetFederationRatingList", func(t *testing.T) {
		t.Run("ValidParams_ReturnsSuccess", func(t *testing.T) {
			// Scenario 1: Correct input
			// ratingtype: 1=Standard, 6=Rapid, 7=Blitz
			// category: 0=All, 1=Juniors, etc.
			rr := MakeRequest(t, handler.GetFederationRatingList, http.MethodGet,
				"/ratinglist/federation/date/2024-01-01/ratingtype/1/category/0",
				map[string]string{
					"ratingdate": "2024-01-01",
					"ratingtype": "1",
					"category":   "0",
				})

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
		})

		t.Run("MalformedRatingType_Returns400", func(t *testing.T) {
			// Scenario 4: Invalid input format
			rr := MakeRequest(t, handler.GetFederationRatingList, http.MethodGet,
				"/ratinglist/federation/date/2024-01-01/ratingtype/abc/category/0",
				map[string]string{
					"ratingdate": "2024-01-01",
					"ratingtype": "abc",
					"category":   "0",
				})

			AssertStatus(t, rr, http.StatusBadRequest)
		})

		t.Run("MalformedCategory_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetFederationRatingList, http.MethodGet,
				"/ratinglist/federation/date/2024-01-01/ratingtype/1/category/abc",
				map[string]string{
					"ratingdate": "2024-01-01",
					"ratingtype": "1",
					"category":   "abc",
				})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetDistrictRatingList", func(t *testing.T) {
		t.Run("ValidParams_ReturnsSuccess", func(t *testing.T) {
			// District 1 = Stockholm
			rr := MakeRequest(t, handler.GetDistrictRatingList, http.MethodGet,
				"/ratinglist/district/1/date/2024-01-01/ratingtype/1/category/0",
				map[string]string{
					"id":         "1",
					"ratingdate": "2024-01-01",
					"ratingtype": "1",
					"category":   "0",
				})

			AssertStatus(t, rr, http.StatusOK)
			AssertContentType(t, rr, "application/json")
		})

		t.Run("MalformedDistrictID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetDistrictRatingList, http.MethodGet,
				"/ratinglist/district/abc/date/2024-01-01/ratingtype/1/category/0",
				map[string]string{
					"id":         "abc",
					"ratingdate": "2024-01-01",
					"ratingtype": "1",
					"category":   "0",
				})

			AssertStatus(t, rr, http.StatusBadRequest)
		})

		t.Run("MalformedRatingType_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetDistrictRatingList, http.MethodGet,
				"/ratinglist/district/1/date/2024-01-01/ratingtype/abc/category/0",
				map[string]string{
					"id":         "1",
					"ratingdate": "2024-01-01",
					"ratingtype": "abc",
					"category":   "0",
				})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})

	t.Run("GetClubRatingList", func(t *testing.T) {
		t.Run("ValidParams_ReturnsSuccess", func(t *testing.T) {
			// TODO: Find a valid club ID to test with
			t.Skip("TODO: Implement with valid club ID")
		})

		t.Run("MalformedClubID_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetClubRatingList, http.MethodGet,
				"/ratinglist/club/abc/date/2024-01-01/ratingtype/1/category/0",
				map[string]string{
					"id":         "abc",
					"ratingdate": "2024-01-01",
					"ratingtype": "1",
					"category":   "0",
				})

			AssertStatus(t, rr, http.StatusBadRequest)
		})

		t.Run("MalformedRatingType_Returns400", func(t *testing.T) {
			rr := MakeRequest(t, handler.GetClubRatingList, http.MethodGet,
				"/ratinglist/club/100/date/2024-01-01/ratingtype/abc/category/0",
				map[string]string{
					"id":         "100",
					"ratingdate": "2024-01-01",
					"ratingtype": "abc",
					"category":   "0",
				})

			AssertStatus(t, rr, http.StatusBadRequest)
		})
	})
}
