package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Banchee501/RossWeatherBot/internal/utils"
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

	var updates []Update

	err := utils.Retry(ctx, 3, 500*time.Millisecond, func() error {

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			return fmt.Errorf("temporary telegram error %d", resp.StatusCode)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("telegram api returned %d", resp.StatusCode)
		}

		var updateResp UpdateResponse

		if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
			return err
		}

		updates = updateResp.Result
		return nil
	})

	return updates, err
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {

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

	return utils.Retry(ctx, 3, 500*time.Millisecond, func() error {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			return fmt.Errorf("temporary telegram error %d", resp.StatusCode)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("telegram api returned %d", resp.StatusCode)
		}

		return nil
	})
}
