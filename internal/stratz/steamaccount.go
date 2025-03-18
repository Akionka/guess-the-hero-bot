package stratz

import (
	"context"
	"fmt"

	"github.com/akionka/akionkabot/internal/data"
)

type PlayerQueryData struct {
	Player *Player `json:"player"`
}

func (c *Client) GetSteamAccount(ctx context.Context, id data.SteamID) (*data.SteamAccount, error) {
	qr, err := query[PlayerQueryData](context.Background(), c, fmt.Sprintf(`
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
		}`, id))
	if err != nil {
		return nil, err
	}

	if qr.Player.SteamAccount == nil {
		return nil, data.ErrNotFound
	}

	p := qr.Player.SteamAccount.toDomain()

	return &p, nil
}
