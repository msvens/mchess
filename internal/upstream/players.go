package upstream

import (
	"context"
	"fmt"
	"net/url"

	"github.com/msvens/mchess/internal/model"
)

// GetPlayer fetches a player by ID and date from the upstream API
func (c *Client) GetPlayer(ctx context.Context, memberID int, date string) (*model.PlayerInfo, error) {
	path := fmt.Sprintf("/player/%d/date/%s", memberID, date)

	var player model.PlayerInfo
	if err := c.get(ctx, path, &player); err != nil {
		return nil, err
	}

	return &player, nil
}

// SearchPlayers searches for players by first name and last name
func (c *Client) SearchPlayers(ctx context.Context, firstName, lastName string) ([]model.PlayerInfo, error) {
	path := fmt.Sprintf("/player/fornamn/%s/efternamn/%s",
		url.PathEscape(firstName),
		url.PathEscape(lastName))

	var players []model.PlayerInfo
	if err := c.get(ctx, path, &players); err != nil {
		return nil, err
	}

	return players, nil
}

// GetPlayerByFideID fetches a player by their FIDE ID
func (c *Client) GetPlayerByFideID(ctx context.Context, fideID int, date string) (*model.PlayerInfo, error) {
	path := fmt.Sprintf("/player/fideid/%d/date/%s", fideID, date)

	var player model.PlayerInfo
	if err := c.get(ctx, path, &player); err != nil {
		return nil, err
	}

	return &player, nil
}
