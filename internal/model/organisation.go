package model

// Federation represents the Swedish Chess Federation
type Federation struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Street      string `json:"street,omitempty"`
	Zipcode     int    `json:"zipcode,omitempty"`
	City        string `json:"city,omitempty"`
	Started     *Date  `json:"started,omitempty"`
	StartSeason *Date  `json:"startseason,omitempty"`
	EndSeason   *Date  `json:"endseason,omitempty"`
	PhoneNr     string `json:"phonenr,omitempty"`
	PostGiro    string `json:"postgiro,omitempty"`
	Email       string `json:"email,omitempty"`
	OrgNumber   string `json:"orgnumber,omitempty"`
	URL         string `json:"url,omitempty"`
	SisuID      int    `json:"sisuid,omitempty"`
}

// District represents a chess district
type District struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	CoContactPerson    string `json:"co_ContantPerson,omitempty"`
	Street             string `json:"street,omitempty"`
	Zipcode            int    `json:"zipcode,omitempty"`
	City               string `json:"city,omitempty"`
	Started            *Date  `json:"started,omitempty"`
	StartSeason        *Date  `json:"startSeason,omitempty"`
	EndSeason          *Date  `json:"endSeason,omitempty"`
	PhoneNr            string `json:"phonenr,omitempty"`
	PostGiro           string `json:"postgiro,omitempty"`
	Active             int    `json:"active,omitempty"`
	Email              string `json:"email,omitempty"`
	OrgNumber          string `json:"orgnumber,omitempty"`
	AuthSchoolClub     string `json:"authschoolclub,omitempty"`
	URL                string `json:"url,omitempty"`
	SisuID             int    `json:"sisuid,omitempty"`
	JoinFederationDate *Date  `json:"joinFederationDate,omitempty"`
}

// Club represents a chess club
type Club struct {
	ID               int                  `json:"id"`
	Name             string               `json:"name"`
	Street           string               `json:"street,omitempty"`
	Zipcode          int                  `json:"zipcode,omitempty"`
	City             string               `json:"city,omitempty"`
	StartDate        *Date                `json:"startdate,omitempty"`
	StartSeason      *Date                `json:"startSeason,omitempty"`
	EndSeason        *Date                `json:"endSeason,omitempty"`
	PhoneNr          string               `json:"phonenr,omitempty"`
	PostGiro         string               `json:"postgiro,omitempty"`
	AlliansClub      int                  `json:"alliansclub,omitempty"`
	Email            string               `json:"email,omitempty"`
	OrgNumber        string               `json:"orgnumber,omitempty"`
	Districts        []DistrictMembership `json:"districts,omitempty"`
	RegYear          *RegistrationYear    `json:"regYear,omitempty"`
	URL              string               `json:"url,omitempty"`
	VBDescr          string               `json:"vbdescr,omitempty"`
	SchoolClub       int                  `json:"schoolClub,omitempty"`
	SchoolName       string               `json:"schoolName,omitempty"`
	SchoolID         int                  `json:"schoolid,omitempty"`
	Lan              int                  `json:"lan,omitempty"`
	NoEconomy        int                  `json:"noEconomy,omitempty"`
	SpecialType      int                  `json:"specialType,omitempty"`
	Yes2ChessReason  bool                 `json:"yes2chessreason,omitempty"`
	Schack4anReason  bool                 `json:"schack4anreason,omitempty"`
	OvrigtReason     bool                 `json:"ovrigtreason,omitempty"`
	HasRatingPlayers int                  `json:"hasRatingPlayers,omitempty"`
	SisuID           int                  `json:"sisuid,omitempty"`
}

// DistrictMembership represents a club's membership in a district
type DistrictMembership struct {
	Start      *Date `json:"start,omitempty"`
	End        *Date `json:"end,omitempty"`
	Year       int   `json:"year,omitempty"`
	DistrictID int   `json:"districtid,omitempty"`
	ClubID     int   `json:"clubid,omitempty"`
	Active     int   `json:"active,omitempty"`
	Benefit    int   `json:"benefit,omitempty"`
}

// RegistrationYear represents a registration year period
type RegistrationYear struct {
	StartDate *Date `json:"startDate,omitempty"`
	EndDate   *Date `json:"endDate,omitempty"`
}
