package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/msvens/mchess/docs"

	"github.com/msvens/mchess/internal/api/handlers"
	"github.com/msvens/mchess/internal/config"
	"github.com/msvens/mchess/internal/db"
	"github.com/msvens/mchess/internal/repository"
	"github.com/msvens/mchess/internal/service"
	"github.com/msvens/mchess/internal/upstream"
)

// @title mchess API
// @version 1.0
// @description A caching proxy for the Swedish Chess Federation (schack.se) API
// @BasePath /api

// Server represents the API server
type Server struct {
	router              *chi.Mux
	cfg                 *config.Config
	playerHandler       *handlers.PlayerHandler
	organisationHandler *handlers.OrganisationHandler
	ratingListHandler   *handlers.RatingListHandler
	tournamentHandler   *handlers.TournamentHandler
	resultsHandler      *handlers.ResultsHandler
	registrationHandler *handlers.RegistrationHandler
	db                  *db.DB
}

// NewServer creates a new API server
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize database
	database, err := db.New(cfg.DBConnectionString())
	if err != nil {
		return nil, err
	}

	// Initialize upstream client
	upstreamClient := upstream.NewClient(
		cfg.Upstream.BaseURL,
		cfg.Upstream.Timeout,
		cfg.Upstream.RateLimit,
	)

	// Initialize repository
	playerRepo := repository.NewPlayerRepository(database.DB)

	// Initialize service
	playerService := service.NewPlayerService(playerRepo, upstreamClient, cfg)

	// Initialize handlers
	playerHandler := handlers.NewPlayerHandler(playerService, upstreamClient)
	organisationHandler := handlers.NewOrganisationHandler(upstreamClient)
	ratingListHandler := handlers.NewRatingListHandler(upstreamClient)
	tournamentHandler := handlers.NewTournamentHandler(upstreamClient)
	resultsHandler := handlers.NewResultsHandler(upstreamClient)
	registrationHandler := handlers.NewRegistrationHandler(upstreamClient)

	s := &Server{
		router:              chi.NewRouter(),
		cfg:                 cfg,
		playerHandler:       playerHandler,
		organisationHandler: organisationHandler,
		ratingListHandler:   ratingListHandler,
		tournamentHandler:   tournamentHandler,
		resultsHandler:      resultsHandler,
		registrationHandler: registrationHandler,
		db:                  database,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(30 * time.Second))

	// CORS for development
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}

func (s *Server) setupRoutes() {
	// Health endpoints (no prefix)
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/ready", s.handleReady)

	// Override Swagger BasePath at runtime to match configured prefix
	docs.SwaggerInfo.BasePath = s.cfg.Server.Prefix

	// API routes - schack.se compatible (under configured prefix, e.g., /api)
	s.router.Route(s.cfg.Server.Prefix, func(r chi.Router) {
		// Swagger UI - accessible at {prefix}/swagger/*
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(s.cfg.Server.Prefix+"/swagger/doc.json"),
		))
		// === schack.se compatible routes ===
		// These match the upstream API exactly for drop-in replacement

		// Player endpoints (schack.se compatible + mchess enhancements)
		r.Get("/player/{id}/date/{date}", s.playerHandler.GetPlayer)
		r.Get("/player/fideid/{id}/date/{date}", s.playerHandler.GetPlayerByFideID)
		r.Get("/player/fornamn/{fornamn}/efternamn/{efternamn}", s.playerHandler.SearchPlayers)
		r.Get("/player/batch", s.playerHandler.GetPlayers)           // mchess: Batch fetch ?ids=1,2,3&date=...
		r.Get("/player/{id}/ratings", s.playerHandler.GetPlayerRatings) // mchess: Rating history

		// Organisation endpoints
		r.Get("/organisation/federation", s.organisationHandler.GetFederation)
		r.Get("/organisation/districts", s.organisationHandler.GetDistricts)
		r.Get("/organisation/district/clubs/{districtid}", s.organisationHandler.GetClubsInDistrict)
		r.Get("/organisation/club/{clubid}", s.organisationHandler.GetClub)
		r.Get("/organisation/club/exists/{name}/{id}", s.organisationHandler.ClubNameExists)

		// Rating list endpoints
		r.Get("/ratinglist/federation/date/{ratingdate}/ratingtype/{ratingtype}/category/{category}", s.ratingListHandler.GetFederationRatingList)
		r.Get("/ratinglist/district/{id}/date/{ratingdate}/ratingtype/{ratingtype}/category/{category}", s.ratingListHandler.GetDistrictRatingList)
		r.Get("/ratinglist/club/{id}/date/{ratingdate}/ratingtype/{ratingtype}/category/{category}", s.ratingListHandler.GetClubRatingList)

		// Tournament structure endpoints
		r.Get("/tournament/tournament/id/{id}", s.tournamentHandler.GetTournament)
		r.Get("/tournament/tournament/updated/{startdate}/{enddate}", s.tournamentHandler.SearchUpdatedTournaments)
		r.Get("/tournament/tournament/updated/{startdate}/{enddate}/{districtid}", s.tournamentHandler.SearchUpdatedTournamentsByDistrict)
		r.Get("/tournament/group/id/{id}", s.tournamentHandler.GetTournamentFromGroup)
		r.Get("/tournament/group/search/{searchWord}", s.tournamentHandler.SearchTournamentGroups)
		r.Get("/tournament/group/coming", s.tournamentHandler.GetComingTournaments)
		r.Get("/tournament/group/coming/{districtid}", s.tournamentHandler.GetComingTournamentsByDistrict)
		r.Get("/tournament/group/updated/{startdate}/{enddate}", s.tournamentHandler.SearchUpdatedGroups)
		r.Get("/tournament/group/updated/{startdate}/{enddate}/{districtid}", s.tournamentHandler.SearchUpdatedGroupsByDistrict)
		r.Get("/tournament/class/id/{id}", s.tournamentHandler.GetTournamentFromClass)

		// Tournament results endpoints
		r.Get("/tournamentresults/table/id/{id}", s.resultsHandler.GetResultTable)
		r.Get("/tournamentresults/table/memberid/{id}", s.resultsHandler.GetMemberTableResults)
		r.Get("/tournamentresults/roundresults/id/{id}", s.resultsHandler.GetRoundResults)
		r.Get("/tournamentresults/team/table/id/{id}", s.resultsHandler.GetTeamResultTable)
		r.Get("/tournamentresults/team/roundresults/id/{id}", s.resultsHandler.GetTeamRoundResults)
		r.Get("/tournamentresults/team/roundresults/id/{id}/memberid/{memberid}", s.resultsHandler.GetTeamRoundResultsForMember)
		r.Get("/tournamentresults/game/memberid/{id}", s.resultsHandler.GetMemberGames)

		// Team registration endpoint
		r.Get("/tournamentteamregistration/tournament/{id}/club/{clubid}", s.registrationHandler.GetTeamRegistration)

	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	if err := s.db.PingContext(r.Context()); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"not ready","reason":"database unavailable"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ready"}`))
}

// Close closes the server and its resources
func (s *Server) Close() error {
	return s.db.Close()
}

// StartServer starts the API server with graceful shutdown
func StartServer(cfg *config.Config) {
	server, err := NewServer(cfg)
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:    cfg.ServerAddr(),
		Handler: server.router,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Starting server", "addr", cfg.ServerAddr())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	if err := server.Close(); err != nil {
		slog.Error("Error closing server resources", "error", err)
	}

	slog.Info("Server stopped")
}
