package upstream

import (
	"context"
	"fmt"
	"net/url"

	"github.com/msvens/mchess/internal/model"
)

// GetTournament fetches tournament by ID
func (c *Client) GetTournament(ctx context.Context, tournamentID int) (*model.Tournament, error) {
	path := fmt.Sprintf("/tournament/tournament/id/%d", tournamentID)
	var tournament model.Tournament
	if err := c.get(ctx, path, &tournament); err != nil {
		return nil, err
	}
	return &tournament, nil
}

// GetTournamentFromGroup fetches tournament info from a group ID
func (c *Client) GetTournamentFromGroup(ctx context.Context, groupID int) (*model.Tournament, error) {
	path := fmt.Sprintf("/tournament/group/id/%d", groupID)
	var tournament model.Tournament
	if err := c.get(ctx, path, &tournament); err != nil {
		return nil, err
	}
	return &tournament, nil
}

// GetTournamentFromClass fetches tournament info from a class ID
func (c *Client) GetTournamentFromClass(ctx context.Context, classID int) (*model.Tournament, error) {
	path := fmt.Sprintf("/tournament/class/id/%d", classID)
	var tournament model.Tournament
	if err := c.get(ctx, path, &tournament); err != nil {
		return nil, err
	}
	return &tournament, nil
}

// SearchTournamentGroups searches for tournament groups by name/location
func (c *Client) SearchTournamentGroups(ctx context.Context, searchWord string) ([]model.TournamentSearchAnswer, error) {
	path := fmt.Sprintf("/tournament/group/search/%s", url.PathEscape(searchWord))
	var results []model.TournamentSearchAnswer
	if err := c.get(ctx, path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetComingTournaments fetches upcoming tournaments
func (c *Client) GetComingTournaments(ctx context.Context) ([]model.Tournament, error) {
	var tournaments []model.Tournament
	if err := c.get(ctx, "/tournament/group/coming", &tournaments); err != nil {
		return nil, err
	}
	return tournaments, nil
}

// GetComingTournamentsByDistrict fetches upcoming tournaments for a district
func (c *Client) GetComingTournamentsByDistrict(ctx context.Context, districtID int) ([]model.Tournament, error) {
	path := fmt.Sprintf("/tournament/group/coming/%d", districtID)
	var tournaments []model.Tournament
	if err := c.get(ctx, path, &tournaments); err != nil {
		return nil, err
	}
	return tournaments, nil
}

// SearchUpdatedTournaments fetches tournaments updated between dates
func (c *Client) SearchUpdatedTournaments(ctx context.Context, startDate, endDate string) ([]model.Tournament, error) {
	path := fmt.Sprintf("/tournament/tournament/updated/%s/%s", startDate, endDate)
	var tournaments []model.Tournament
	if err := c.get(ctx, path, &tournaments); err != nil {
		return nil, err
	}
	return tournaments, nil
}

// SearchUpdatedTournamentsByDistrict fetches tournaments updated between dates for a district
func (c *Client) SearchUpdatedTournamentsByDistrict(ctx context.Context, startDate, endDate string, districtID int) ([]model.Tournament, error) {
	path := fmt.Sprintf("/tournament/tournament/updated/%s/%s/%d", startDate, endDate, districtID)
	var tournaments []model.Tournament
	if err := c.get(ctx, path, &tournaments); err != nil {
		return nil, err
	}
	return tournaments, nil
}

// SearchUpdatedGroups fetches groups updated between dates
func (c *Client) SearchUpdatedGroups(ctx context.Context, startDate, endDate string) ([]model.TournamentSearchAnswer, error) {
	path := fmt.Sprintf("/tournament/group/updated/%s/%s", startDate, endDate)
	var groups []model.TournamentSearchAnswer
	if err := c.get(ctx, path, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// SearchUpdatedGroupsByDistrict fetches groups updated between dates for a district
func (c *Client) SearchUpdatedGroupsByDistrict(ctx context.Context, startDate, endDate string, districtID int) ([]model.TournamentSearchAnswer, error) {
	path := fmt.Sprintf("/tournament/group/updated/%s/%s/%d", startDate, endDate, districtID)
	var groups []model.TournamentSearchAnswer
	if err := c.get(ctx, path, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}
