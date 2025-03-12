package d2pt

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

// IntBool is a custom boolean type that unmarshals from JSON integers (0 or 1).
// It is used to represent boolean values in the Dota2ProTracker API.
type IntBool bool

func (ib *IntBool) UnmarshalJSON(data []byte) error {
	if len(data) != 1 {
		return errors.New("ib length is not equal 1")
	}
	switch data[0] {
	case '0':
		*ib = false
	case '1':
		*ib = true
	default:
		return errors.New("invalid value for IntBool")
	}

	return nil
}

func (ib IntBool) MarshalJSON() ([]byte, error) {
	v := 0
	if ib {
		v = 1
	}

	return json.Marshal(v)
}
