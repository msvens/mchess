package upstream

import (
	"context"
	"fmt"

	"github.com/msvens/mchess/internal/model"
)

// GetFederationRatingList fetches federation-wide rating list
func (c *Client) GetFederationRatingList(ctx context.Context, date string, ratingType, category int) ([]model.PlayerInfo, error) {
	path := fmt.Sprintf("/ratinglist/federation/date/%s/ratingtype/%d/category/%d", date, ratingType, category)
	var players []model.PlayerInfo
	if err := c.get(ctx, path, &players); err != nil {
		return nil, err
	}
	return players, nil
}

// GetDistrictRatingList fetches district rating list
func (c *Client) GetDistrictRatingList(ctx context.Context, districtID int, date string, ratingType, category int) ([]model.PlayerInfo, error) {
	path := fmt.Sprintf("/ratinglist/district/%d/date/%s/ratingtype/%d/category/%d", districtID, date, ratingType, category)
	var players []model.PlayerInfo
	if err := c.get(ctx, path, &players); err != nil {
		return nil, err
	}
	return players, nil
}

// GetClubRatingList fetches club rating list
func (c *Client) GetClubRatingList(ctx context.Context, clubID int, date string, ratingType, category int) ([]model.PlayerInfo, error) {
	path := fmt.Sprintf("/ratinglist/club/%d/date/%s/ratingtype/%d/category/%d", clubID, date, ratingType, category)
	var players []model.PlayerInfo
	if err := c.get(ctx, path, &players); err != nil {
		return nil, err
	}
	return players, nil
}
