package main

import (
	"log"
	"time"

	"github.com/Banchee501/RossWeatherBot/internal/config"
	"github.com/Banchee501/RossWeatherBot/internal/telegram"
	"github.com/Banchee501/RossWeatherBot/internal/weather"
)

func main() {
	cfg := config.Load()

	botClient := telegram.NewClient(cfg.TelegramToken)
	weatherClient := weather.NewClient(cfg.WeatherAPIKey)

	var lastUpdateID int

	for {
		updates, err := botClient.GetUpdates(lastUpdateID)
		if err != nil {
			log.Println("Error getting updates:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if len(updates) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		for _, update := range updates {
			if update.Message.Text == "" {
				continue
			}

			text := update.Message.Text

			if text == "/start" {
				botClient.SendMessage(update.Message.Chat.ID, "Напишіть назву міста")
			} else {
				weatherText, err := weatherClient.GetWeather(text)
				if err != nil {
					botClient.SendMessage(update.Message.Chat.ID, "Місто не знайдено")
				} else {
					botClient.SendMessage(update.Message.Chat.ID, weatherText)
				}
			}

			lastUpdateID = update.UpdateID + 1
		}
	}
}
