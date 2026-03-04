package telegram

import (
	"bytes"
	"context"
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
			Timeout: 35 * time.Second,
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

func (c *Client) GetUpdates(ctx context.Context, offset int) ([]Update, error) {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=30",
		c.Token,
		offset,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram api returned %d", resp.StatusCode)
	}

	var updateResp UpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return nil, err
	}

	return updateResp.Result, nil
}

func (c *Client) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		c.Token,
	)

	body := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	return nil
}
