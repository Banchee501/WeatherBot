package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			return
		default:
			updates, err := botClient.GetUpdates(ctx, lastUpdateID)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Println("Error getting updates:", err)
				time.Sleep(1 * time.Second)
				continue
			}

			for _, update := range updates {
				if update.Message.Text == "" {
					continue
				}

				text := strings.TrimSpace(update.Message.Text)

				if text == "/start" {
					botClient.SendMessage(update.Message.Chat.ID, "Напишіть назву міста")
				} else {
					weatherText, err := weatherClient.GetWeather(ctx, text)
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
}
