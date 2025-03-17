package d2pt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func (c *Client) FetchQuestion(ctx context.Context, isWon bool) (*QuestionResponse, error) {
	qr := new(QuestionResponse)
	isWonStr := "0"
	if isWon {
		isWonStr = "1"
	}
	requestURL := url.URL{
		Scheme: "https",
		Host:   "dota2protracker.com",
		Path:   "api/gth/random-question",
		RawQuery: url.Values{
			"won":         []string{isWonStr},
			"p_min_score": []string{"3"},
			"TCACHE":      []string{"0"},
		}.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Referer", "https://dota2protracker.com/")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(qr); err != nil {
		return nil, err
	}

	return qr, nil
}

// QuestionResponse represents a question fetched from the Dota2ProTracker API.
type QuestionResponse struct {
	ID              int     `json:"id"`
	MatchID         int64   `json:"match_id"`
	HeroID          int     `json:"hero_id"`
	SteamID         int64   `json:"account_id"`
	IsPro           IntBool `json:"is_pro"`
	Name            string  `json:"name"`
	Pos             string  `json:"pos"`
	MMR             int     `json:"mmr"`
	ActivateTime    int64   `json:"activate_time"`
	Score1          string  `json:"score1"`
	Score2          string  `json:"score2"`
	Items           []int   `json:"items"`
	ItemsWithScores []struct {
		P      float64 `json:"p"`
		Score1 float64 `json:"score1"`
		Score2 int     `json:"score2"`
		ItemID int     `json:"item_id"`
	}
	IsWon IntBool `json:"won"`
	// CreatedAt    *time.Time `json:"created_at"`
	WrongOptions []int `json:"wrong_options"`
}
