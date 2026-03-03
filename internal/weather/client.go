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
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
}

func windDirection(deg float64) string {
	switch {
	case deg >= 337.5 || deg < 22.5:
		return "північний"
	case deg >= 22.5 && deg < 67.5:
		return "північно-східний"
	case deg >= 67.5 && deg < 112.5:
		return "східний"
	case deg >= 112.5 && deg < 157.5:
		return "південно-східний"
	case deg >= 157.5 && deg < 202.5:
		return "південний"
	case deg >= 202.5 && deg < 247.5:
		return "південно-західний"
	case deg >= 247.5 && deg < 292.5:
		return "західний"
	case deg >= 292.5 && deg < 337.5:
		return "північно-західний"
	default:
		return "невідомий"
	}
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

	windDir := windDirection(data.Wind.Deg)
	windSpeed := data.Wind.Speed

	result := fmt.Sprintf(
		"Погода в %s:\nТемпература: %.1f°C\nОпис: %s\n Вітер: %.1f м/с, %s",
		city,
		data.Main.Temp,
		data.Weather[0].Description,
		windSpeed,
		windDir,
	)

	return result, nil
}
