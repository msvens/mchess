package model

import "time"

// PlayerInfo represents a player from the schack.se API
// This matches the upstream PlayerInfoDto exactly
// @Description Player information from the Swedish Chess Federation
// @name PlayerInfo
type PlayerInfo struct {
	ID        int         `json:"id" example:"12345"`
	FirstName string      `json:"firstName" example:"Magnus"`
	LastName  string      `json:"lastName" example:"Carlsen"`
	Birthdate string      `json:"birthdate,omitempty" example:"1990"`
	Sex       int         `json:"sex" example:"1"`
	FideID    int         `json:"fideid,omitempty" example:"1503014"`
	Country   string      `json:"country,omitempty" example:"SWE"`
	Club      string      `json:"club,omitempty" example:"Stockholms SS"`
	ClubID    int         `json:"clubId,omitempty" example:"101"`
	Elo       *EloRating  `json:"elo,omitempty"`
	Lask      *LaskRating `json:"lask,omitempty"`
}

// EloRating represents FIDE ELO ratings
// @Description FIDE ELO rating information
// @name EloRating
type EloRating struct {
	Rating      int    `json:"rating" example:"2830"`
	Title       string `json:"title,omitempty" example:"GM"`
	Date        *Date  `json:"date,omitempty"`
	K           int    `json:"k,omitempty" example:"10"`
	RapidRating int    `json:"rapidRating,omitempty" example:"2800"`
	RapidK      int    `json:"rapidk,omitempty" example:"10"`
	BlitzRating int    `json:"blitzRating,omitempty" example:"2850"`
	BlitzK      int    `json:"blitzK,omitempty" example:"10"`
}

// LaskRating represents Swedish national rating
// @Description Swedish national (LASK) rating
// @name LaskRating
type LaskRating struct {
	Rating int   `json:"rating" example:"2450"`
	Date   *Date `json:"date,omitempty"`
}

// PlayerCacheEntry represents a row in the player_cache table (internal)
type PlayerCacheEntry struct {
	MemberID    int
	RatingDate  time.Time
	FirstName   string
	LastName    string
	Club        string
	ClubID      *int
	FideID      *int
	EloStandard *int
	EloRapid    *int
	EloBlitz    *int
	LaskRating  *int
	Data        []byte // JSONB
	FetchedAt   time.Time
	ExpiresAt   *time.Time
}

// PlayersResponse is the response for batch player requests
// @Description Response containing multiple players and any errors
// @name PlayersResponse
type PlayersResponse struct {
	Players []PlayerInfo  `json:"players"`
	Errors  []PlayerError `json:"errors,omitempty"`
}

// PlayerError represents an error for a specific player ID
// @Description Error information for a failed player lookup
// @name PlayerError
type PlayerError struct {
	ID    int    `json:"id" example:"99999"`
	Error string `json:"error" example:"player not found"`
}

// RatingHistoryResponse is the response for player rating history
// @Description Historical rating data for a player (array of PlayerInfo for each date)
// @name RatingHistoryResponse
type RatingHistoryResponse struct {
	PlayerID int          `json:"playerId" example:"12345"`
	Ratings  []PlayerInfo `json:"ratings"`
}
