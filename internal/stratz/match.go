package stratz

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akionka/akionkabot/internal/data"
)

type MatchQueryData struct {
	Match Match `json:"match"`
}

func (c *Client) GetMatchByID(ctx context.Context, matchID int64) (*data.Match, error) {
	body, err := c.query(context.Background(), fmt.Sprintf(`
		{
			match(id: %d) {
				id
				didRadiantWin
				actualRank
				startDateTime
				players {
					heroId
					item0Id
					item1Id
					item2Id
					item3Id
					item4Id
					item5Id
					isRadiant
					position
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
		return nil, fmt.Errorf("graphql query failed: %w", err)
	}
	defer body.Close()

	var qr queryResponse[MatchQueryData]
	if err = json.NewDecoder(body).Decode(&qr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if qr.Error != nil {
		return nil, fmt.Errorf("graphql error: %s", qr.Error)
	}

	if qr.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}

	m := qr.Data.Match.toDomain()

	return &m, nil
}
