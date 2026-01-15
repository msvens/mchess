package model

// Tournament represents a tournament
type Tournament struct {
	ID                           int               `json:"id"`
	Name                         string            `json:"name"`
	Start                        *Date             `json:"start,omitempty"`
	End                          *Date             `json:"end,omitempty"`
	City                         string            `json:"city,omitempty"`
	Arena                        string            `json:"arena,omitempty"`
	Type                         int               `json:"type,omitempty"`
	IA                           int               `json:"ia,omitempty"`
	SecJudges                    string            `json:"secjudges,omitempty"`
	ThinkingTime                 string            `json:"thinkingTime,omitempty"`
	State                        int               `json:"state,omitempty"`
	AllowForeignPlayers          int               `json:"allowForeignPlayers,omitempty"`
	TeamTournamentPlayerListType int               `json:"teamtournamentPlayerListType,omitempty"`
	AgeFilter                    int               `json:"ageFilter,omitempty"`
	NrOfPartLink                 string            `json:"nrOfPartLink,omitempty"`
	OrgType                      int               `json:"orgType,omitempty"`
	OrgNumber                    int               `json:"orgNumber,omitempty"`
	RatingRegDate                *Date             `json:"ratingRegDate,omitempty"`
	RatingRegDate2               *Date             `json:"ratingRegDate2,omitempty"`
	FIDERegged                   int               `json:"fideregged,omitempty"`
	Online                       int               `json:"online,omitempty"`
	Y2CRules                     int               `json:"y2cRules,omitempty"`
	TeamNrOfDaysRegged           int               `json:"teamNrOfDaysRegged,omitempty"`
	ShowPublic                   int               `json:"showPublic,omitempty"`
	InvitationURL                string            `json:"invitationurl,omitempty"`
	LatestUpdated                *Date             `json:"latestUpdated,omitempty"`
	SecParsedJudges              []int             `json:"secParsedJudges,omitempty"`
	RootClasses                  []TournamentClass `json:"rootClasses,omitempty"`
}

// TournamentClass represents a tournament class/division
type TournamentClass struct {
	ClassID                         int                    `json:"classID"`
	TournamentID                    int                    `json:"tournamentID,omitempty"`
	ParentClassID                   int                    `json:"parentClassID,omitempty"`
	Order                           int                    `json:"order,omitempty"`
	ClassName                       string                 `json:"className,omitempty"`
	NumberOfTeamsToGoUpClass        int                    `json:"numberOfTeamsToGoUpClass,omitempty"`
	NumberOfTeamsToGoDownClass      int                    `json:"numberOfTeamsToGoDownClass,omitempty"`
	PercentageOfPointsToGoUpClass   int                    `json:"percentageOfPointsToGoUpClass,omitempty"`
	PercentageOfPointsToGoDownClass int                    `json:"percentageOfPointsToGoDownClass,omitempty"`
	GamesURL                        string                 `json:"gamesUrl,omitempty"`
	Groups                          []TournamentClassGroup `json:"groups,omitempty"`
}

// TournamentClassGroup represents a group within a tournament class
type TournamentClassGroup struct {
	ID                              int             `json:"id"`
	ClassID                         int             `json:"classID,omitempty"`
	Name                            string          `json:"name,omitempty"`
	PairingSystemMember             int             `json:"pairingSystemMember,omitempty"`
	Order                           int             `json:"order,omitempty"`
	Start                           *Date           `json:"start,omitempty"`
	End                             *Date           `json:"end,omitempty"`
	NumberOfTeamsToGoUpClass        int             `json:"numberOfTeamsToGoUpClass,omitempty"`
	NumberOfTeamsToGoDownClass      int             `json:"numberOfTeamsToGoDownClass,omitempty"`
	PercentageOfPointsToGoUpClass   int             `json:"percentageOfPointsToGoUpClass,omitempty"`
	PercentageOfPointsToGoDownClass int             `json:"percentageOfPointsToGoDownClass,omitempty"`
	PlayersInTeam                   int             `json:"playersinteam,omitempty"`
	DoubleRounded                   int             `json:"doubleRounded,omitempty"`
	Cost                            int             `json:"cost,omitempty"`
	PossibleToRegister              int             `json:"possibleToRegister,omitempty"`
	TiebreakSystem                  int             `json:"tiebreakSystem,omitempty"`
	PointSystem                     int             `json:"pointSystem,omitempty"`
	PrintName                       string          `json:"printName,omitempty"`
	RankingAlgorithm                int             `json:"rankingAlgorithm,omitempty"`
	SplitGroupSize                  int             `json:"splitgroupSize,omitempty"`
	NrOfRounds                      int             `json:"nrofrounds,omitempty"`
	ArenaStart                      *LocalTime      `json:"arenaStart,omitempty"`
	ArenaEnd                        *LocalTime      `json:"arenaEnd,omitempty"`
	RegistrationCategories          []PrizeCategory `json:"registrationCategories,omitempty"`
	PrizeCategories                 []PrizeCategory `json:"prizeCategories,omitempty"`
	TournamentRounds                []Round         `json:"tournamentRounds,omitempty"`
}

// LocalTime represents a time without date
type LocalTime struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second,omitempty"`
	Nano   int `json:"nano,omitempty"`
}

// PrizeCategory represents a prize or registration category
type PrizeCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name,omitempty"`
	Start     int    `json:"start,omitempty"`
	End       int    `json:"end,omitempty"`
	Type      int    `json:"type,omitempty"`
	GroupID   int    `json:"groupid,omitempty"`
	Order     int    `json:"order,omitempty"`
	UsageType int    `json:"usagetype,omitempty"`
	AndLogic  int    `json:"andlogic,omitempty"`
}

// Round represents a tournament round
type Round struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"groupId,omitempty"`
	RoundNumber int    `json:"roundNumber,omitempty"`
	RoundDate   *Date  `json:"roundDate,omitempty"`
	Rated       int    `json:"rated,omitempty"`
	JudgeID     int    `json:"judgeId,omitempty"`
	LockTime    *Date  `json:"lockTime,omitempty"`
	PublishTime *Date  `json:"publishTime,omitempty"`
	MatchID     int    `json:"matchId,omitempty"`
}

// TournamentSearchAnswer represents a search result
type TournamentSearchAnswer struct {
	ID                int    `json:"id"`
	Name              string `json:"name,omitempty"`
	TournamentID      int    `json:"tournamentid,omitempty"`
	TournamentName    string `json:"tournamentname,omitempty"`
	LatestUpdatedGame *Date  `json:"latestUpdatedGame,omitempty"`
}
