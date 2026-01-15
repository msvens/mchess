package model

// TournamentEndResult represents a player's final result in a tournament
type TournamentEndResult struct {
	Points      float32     `json:"points"`
	SecPoints   float64     `json:"secPoints,omitempty"`
	Place       int         `json:"place,omitempty"`
	ContenderID int         `json:"contenderId,omitempty"`
	TeamNumber  int         `json:"teamNumber,omitempty"`
	WonGames    int         `json:"wonGames,omitempty"`
	DrawGames   int         `json:"drawGames,omitempty"`
	LostGames   int         `json:"lostGames,omitempty"`
	GroupID     int         `json:"groupId,omitempty"`
	PlayerInfo  *PlayerInfo `json:"playerInfo,omitempty"`
}

// TeamTournamentEndResult represents a team's final result in a tournament
type TeamTournamentEndResult struct {
	Points      float32 `json:"points"`
	SecPoints   float64 `json:"secPoints,omitempty"`
	Place       int     `json:"place,omitempty"`
	ContenderID int     `json:"contenderId,omitempty"`
	TeamNumber  int     `json:"teamNumber,omitempty"`
	WonGames    int     `json:"wonGames,omitempty"`
	DrawGames   int     `json:"drawGames,omitempty"`
	LostGames   int     `json:"lostGames,omitempty"`
	Club        *Club   `json:"club,omitempty"`
}

// TournamentRoundResult represents results for a round
type TournamentRoundResult struct {
	ID             int     `json:"id"`
	GroupID        int     `json:"groupdId,omitempty"` // Note: typo in upstream API
	RoundNr        int     `json:"roundNr,omitempty"`
	Board          int     `json:"board,omitempty"`
	HomeID         int     `json:"homeId,omitempty"`
	HomeTeamNumber int     `json:"homeTeamNumber,omitempty"`
	AwayID         int     `json:"awayId,omitempty"`
	AwayTeamNumber int     `json:"awayTeamNumber,omitempty"`
	HomeResult     float32 `json:"homeResult,omitempty"`
	AwayResult     float32 `json:"awayResult,omitempty"`
	Date           *Date   `json:"date,omitempty"`
	Finalized      bool    `json:"finalized,omitempty"`
	Publisher      int     `json:"publisher,omitempty"`
	PublishDate    *Date   `json:"publishDate,omitempty"`
	PublishedNote  string  `json:"publishedNote,omitempty"`
	Games          []Game  `json:"games,omitempty"`
}

// Game represents an individual chess game
type Game struct {
	ID                 int    `json:"id"`
	TournamentResultID int    `json:"tournamentResultID,omitempty"`
	TableNr            int    `json:"tableNr,omitempty"`
	WhiteID            int    `json:"whiteId,omitempty"`
	BlackID            int    `json:"blackId,omitempty"`
	Result             int    `json:"result,omitempty"`
	PGN                string `json:"pgn,omitempty"`
	GroupID            int    `json:"groupiD,omitempty"` // Note: capital D in upstream API
}
