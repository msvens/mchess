package upstream

import (
	"context"
	"fmt"

	"github.com/msvens/mchess/internal/model"
)

// GetResultTable fetches individual tournament table by group ID
func (c *Client) GetResultTable(ctx context.Context, groupID int) ([]model.TournamentEndResult, error) {
	path := fmt.Sprintf("/tournamentresults/table/id/%d", groupID)
	var results []model.TournamentEndResult
	if err := c.get(ctx, path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetMemberTableResults fetches member's tournament results
func (c *Client) GetMemberTableResults(ctx context.Context, memberID int) ([]model.TournamentEndResult, error) {
	path := fmt.Sprintf("/tournamentresults/table/memberid/%d", memberID)
	var results []model.TournamentEndResult
	if err := c.get(ctx, path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetRoundResults fetches round results for a group
func (c *Client) GetRoundResults(ctx context.Context, groupID int) ([]model.TournamentRoundResult, error) {
	path := fmt.Sprintf("/tournamentresults/roundresults/id/%d", groupID)
	var results []model.TournamentRoundResult
	if err := c.get(ctx, path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetTeamResultTable fetches team tournament table by group ID
func (c *Client) GetTeamResultTable(ctx context.Context, groupID int) ([]model.TeamTournamentEndResult, error) {
	path := fmt.Sprintf("/tournamentresults/team/table/id/%d", groupID)
	var results []model.TeamTournamentEndResult
	if err := c.get(ctx, path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetTeamRoundResults fetches team round results for a group
func (c *Client) GetTeamRoundResults(ctx context.Context, groupID int) ([]model.TournamentRoundResult, error) {
	path := fmt.Sprintf("/tournamentresults/team/roundresults/id/%d", groupID)
	var results []model.TournamentRoundResult
	if err := c.get(ctx, path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetTeamRoundResultsForMember fetches team round results for a specific member
func (c *Client) GetTeamRoundResultsForMember(ctx context.Context, groupID, memberID int) ([]model.TournamentRoundResult, error) {
	path := fmt.Sprintf("/tournamentresults/team/roundresults/id/%d/memberid/%d", groupID, memberID)
	var results []model.TournamentRoundResult
	if err := c.get(ctx, path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetMemberGames fetches all games for a member
func (c *Client) GetMemberGames(ctx context.Context, memberID int) ([]model.Game, error) {
	path := fmt.Sprintf("/tournamentresults/game/memberid/%d", memberID)
	var games []model.Game
	if err := c.get(ctx, path, &games); err != nil {
		return nil, err
	}
	return games, nil
}
