package d2pt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (d2 Client) FetchQuestion(ctx context.Context, isWon bool) (*Question, error) {
	qr := new(Question)
	isWonStr := "0"
	if isWon {
		isWonStr = "1"
	}
	url := url.URL{
		Scheme: "https",
		Host:   "dota2protracker.com",
		Path:   "api/gth/random-question",
		RawQuery: url.Values{
			"won":         []string{isWonStr},
			"p_min_score": []string{"3"},
			"TCACHE":      []string{"0"},
		}.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("referer", "https://dota2protracker.com/")
	resp, err := d2.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(qr); err != nil {
		return nil, err
	}
	return qr, nil
}

type Question struct {
	ID              int     `json:"id"`
	MatchID         int64   `json:"match_id"`
	HeroID          int     `json:"hero_id"`
	AccountID       int64   `json:"account_id"`
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

type IntBool bool

func (ib *IntBool) UnmarshalJSON(data []byte) error {
	if len(data) != 1 {
		return errors.New("ib length is not equal 1")
	}
	*ib = data[0] == '1'
	return nil
}

func (ib IntBool) MarshalJSON() ([]byte, error) {
	v := 0
	if ib {
		v = 1
	}

	return json.Marshal(v)
}
