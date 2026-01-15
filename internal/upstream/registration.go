package upstream

import (
	"context"
	"fmt"

	"github.com/msvens/mchess/internal/model"
)

// GetTeamRegistration fetches team registration for a tournament and club
func (c *Client) GetTeamRegistration(ctx context.Context, tournamentID, clubID int) (*model.TeamRegistration, error) {
	path := fmt.Sprintf("/tournamentteamregistration/tournament/%d/club/%d", tournamentID, clubID)
	var registration model.TeamRegistration
	if err := c.get(ctx, path, &registration); err != nil {
		return nil, err
	}
	return &registration, nil
}
