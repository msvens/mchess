package upstream

import (
	"context"
	"fmt"

	"github.com/msvens/mchess/internal/model"
)

// GetFederation fetches Swedish Chess Federation info
func (c *Client) GetFederation(ctx context.Context) (*model.Federation, error) {
	var federation model.Federation
	if err := c.get(ctx, "/organisation/federation", &federation); err != nil {
		return nil, err
	}
	return &federation, nil
}

// GetDistricts fetches all districts
func (c *Client) GetDistricts(ctx context.Context) ([]model.District, error) {
	var districts []model.District
	if err := c.get(ctx, "/organisation/districts", &districts); err != nil {
		return nil, err
	}
	return districts, nil
}

// GetClubsInDistrict fetches clubs in a district
func (c *Client) GetClubsInDistrict(ctx context.Context, districtID int) ([]model.Club, error) {
	path := fmt.Sprintf("/organisation/district/clubs/%d", districtID)
	var clubs []model.Club
	if err := c.get(ctx, path, &clubs); err != nil {
		return nil, err
	}
	return clubs, nil
}

// GetClub fetches a specific club
func (c *Client) GetClub(ctx context.Context, clubID int) (*model.Club, error) {
	path := fmt.Sprintf("/organisation/club/%d", clubID)
	var club model.Club
	if err := c.get(ctx, path, &club); err != nil {
		return nil, err
	}
	return &club, nil
}

// ClubNameExists checks if a club name exists (other than for the given club ID)
func (c *Client) ClubNameExists(ctx context.Context, name string, clubID int) (bool, error) {
	path := fmt.Sprintf("/organisation/club/exists/%s/%d", name, clubID)
	var exists bool
	if err := c.get(ctx, path, &exists); err != nil {
		return false, err
	}
	return exists, nil
}
