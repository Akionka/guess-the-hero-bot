package stratz

import (
	"context"
	"fmt"

	"github.com/akionka/akionkabot/internal/data"
)

type MatchQueryData struct {
	Match Match `json:"match"`
}

func (c *Client) GetMatch(ctx context.Context, matchID data.MatchID) (*data.Match, error) {
	qr, err := query[MatchQueryData](context.Background(), c, fmt.Sprintf(`
		{
			match(id: %d) {
				id
				didRadiantWin
				actualRank
				startDateTime
				players {
					hero {
						id
						shortName
						displayName
					}
					item0Id
					item1Id
					item2Id
					item3Id
					item4Id
					item5Id
					isRadiant
					position
					steamAccountId
				}
			}
		}`, matchID))
	if err != nil {
		return nil, err
	}
	m := qr.Match.toDomain()

	return &m, nil
}

func (c *Client) GetMatchSteamAccounts(ctx context.Context, matchID data.MatchID) ([]data.SteamAccount, error) {
	qr, err := query[MatchQueryData](context.Background(), c, fmt.Sprintf(`
	{
		match(id: %d) {
			players {
				steamAccount {
					id
					name
					proSteamAccount {
						name
					}
				}
			}
		}
	}`, matchID))
	if err != nil {
		return nil, err
	}

	accounts := make([]data.SteamAccount, 0, len(qr.Match.Players))
	for _, p := range qr.Match.Players {
		if p.SteamAccount == nil {
			continue
		}
		accounts = append(accounts, p.SteamAccount.toDomain())
	}

	return accounts, nil
}
