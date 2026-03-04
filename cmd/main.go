package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Banchee501/RossWeatherBot/internal/config"
	"github.com/Banchee501/RossWeatherBot/internal/handler"
	"github.com/Banchee501/RossWeatherBot/internal/telegram"
	"github.com/Banchee501/RossWeatherBot/internal/utils"
	"github.com/Banchee501/RossWeatherBot/internal/weather"
)

func main() {
	cfg := config.Load()

	botClient := telegram.NewClient(cfg.TelegramToken)
	weatherClient := weather.NewClient(cfg.WeatherAPIKey)

	h := handler.New(botClient, weatherClient)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var lastUpdateID int

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
				if update.Message.Text == "" || update.Message.Chat.ID == 0 {
					continue
				}

				err := utils.Retry(ctx, 3, 500*time.Millisecond, func() error {
					h.Handle(ctx, update.Message.Chat.ID, update.Message.Text)
					return nil
				})
				if err != nil {
					log.Println("Handler failed:", err)
				}

				lastUpdateID = update.UpdateID + 1
			}
		}
	}
}
