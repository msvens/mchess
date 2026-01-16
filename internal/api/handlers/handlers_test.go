package handlers_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/msvens/mchess/internal/db"
	"github.com/msvens/mchess/internal/upstream"
)

// Default upstream URL - uses the real schack.se API
const defaultUpstreamURL = "https://member.schack.se/public/api/v1"

// Default test database connection string
// Override with MCHESS_TEST_DB environment variable
func testDBConnectionString() string {
	if connStr := os.Getenv("MCHESS_TEST_DB"); connStr != "" {
		return connStr
	}
	return "postgres://mchess:mchess@localhost:5432/mchess_test?sslmode=disable"
}

// NewTestClient creates an upstream client for testing against the real API
func NewTestClient(t *testing.T) *upstream.Client {
	t.Helper()
	return upstream.NewClient(defaultUpstreamURL, 30*time.Second, 10)
}

// NewTestDB creates a database connection for testing
// Call ClearTestDB before running tests to ensure clean state
func NewTestDB(t *testing.T) *db.DB {
	t.Helper()

	database, err := db.New(testDBConnectionString())
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return database
}

// SetupTestDB runs migrations on the test database
func SetupTestDB(t *testing.T) {
	t.Helper()

	if err := db.MigrateUp(testDBConnectionString()); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}

// ClearTestDB clears all data from the test database (keeps schema)
func ClearTestDB(t *testing.T) {
	t.Helper()

	conn, err := sql.Open("pgx", testDBConnectionString())
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer conn.Close()

	// Clear tables in correct order (respecting foreign keys if any)
	tables := []string{
		"player_cache",
		"cache_stats",
	}

	for _, table := range tables {
		_, err := conn.Exec("TRUNCATE TABLE " + table + " CASCADE")
		if err != nil {
			// Table might not exist yet, that's ok
			t.Logf("Warning: Could not truncate %s: %v", table, err)
		}
	}
}

// MakeRequest creates an HTTP request with chi URL params and executes it against a handler
func MakeRequest(t *testing.T, handler http.HandlerFunc, method, path string, urlParams map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)

	// Add chi URL params to context
	rctx := chi.NewRouteContext()
	for k, v := range urlParams {
		rctx.URLParams.Add(k, v)
	}
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// MakeRequestWithQuery creates an HTTP request with query parameters
func MakeRequestWithQuery(t *testing.T, handler http.HandlerFunc, method, path string, urlParams map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)

	// chi context for URL params
	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// AssertStatus checks the response status code
func AssertStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rr.Code != want {
		t.Errorf("Status code: got %d, want %d. Body: %s", rr.Code, want, rr.Body.String())
	}
}

// AssertContentType checks the Content-Type header
func AssertContentType(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := rr.Header().Get("Content-Type")
	if got != want {
		t.Errorf("Content-Type: got %q, want %q", got, want)
	}
}

// AssertBodyContains checks if the response body contains a substring
func AssertBodyContains(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()
	body := rr.Body.String()
	if !containsString(body, want) {
		t.Errorf("Body should contain %q, got: %s", want, body)
	}
}

// AssertBodyNotEmpty checks that the response body is not empty
func AssertBodyNotEmpty(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	if rr.Body.Len() == 0 {
		t.Error("Body should not be empty")
	}
}

func containsString(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}