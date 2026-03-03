package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

func (c *Client) GetUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=30",
		c.Token,
		offset,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateResp UpdateResponse
	err = json.NewDecoder(resp.Body).Decode(&updateResp)
	if err != nil {
		return nil, err
	}

	return updateResp.Result, nil
}
