package stratz

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akionka/akionkabot/internal/data"
)

type PlayerQueryData struct {
	Player *Player `json:"player"`
}

func (c *Client) GetPlayerByID(ctx context.Context, steamID int64) (*data.SteamAccount, error) {
	body, err := c.query(context.Background(), fmt.Sprintf(`
		{
			player (steamAccountId: %d) {
				steamAccount {
					id
					name
					proSteamAccount {
						isPro
						name
					}
				}
			}
		}`, steamID))
	if err != nil {
		return nil, fmt.Errorf("graphql query failed: %w", err)
	}
	defer body.Close()

	var qr queryResponse[PlayerQueryData]
	if err = json.NewDecoder(body).Decode(&qr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if qr.Error != nil {
		return nil, fmt.Errorf("graphql error: %s", qr.Error)
	}

	if qr.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}

	if qr.Data.Player.SteamAccount == nil {
		return nil, data.ErrNotFound
	}

	p := qr.Data.Player.SteamAccount.toDomain()

	return &p, nil
}
