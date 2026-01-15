package model

// TeamRegistration represents team registration for a tournament
type TeamRegistration struct {
	TournamentID int                      `json:"tournamentid"`
	ClubID       int                      `json:"clubid"`
	Players      []TeamRegistrationPlayer `json:"players,omitempty"`
}

// TeamRegistrationPlayer represents a registered player in a team
type TeamRegistrationPlayer struct {
	Registered     *Date       `json:"registered,omitempty"`
	Available      *Date       `json:"available,omitempty"`
	SwedishCitizen bool        `json:"swedishCitizen,omitempty"`
	PlayerInfoDTO  *PlayerInfo `json:"playerInfoDto,omitempty"` // Keep json tag to match upstream
}
