package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 35 * time.Second,
		},
	}
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func (c *Client) GetWeather(city string) (string, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ua",
		city,
		c.apiKey,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("weather api returned %d", resp.StatusCode)
	}

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Weather) == 0 {
		return "", fmt.Errorf("no weather data")
	}

	result := fmt.Sprintf(
		"🌤 Погода в %s:\nТемпература: %.1f°C\nОпис: %s",
		city,
		data.Main.Temp,
		data.Weather[0].Description,
	)

	return result, nil
}
