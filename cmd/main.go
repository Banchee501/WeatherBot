package main

import (
	"log"

	"github.com/Banchee501/RossWeatherBot/internal/config"
	"github.com/Banchee501/RossWeatherBot/internal/telegram"
)

func main() {
	cfg := config.Load()

	bot := telegram.NewClient(cfg.TelegramToken)

	var lastUpdateID int

	for {
		updates, err := bot.GetUpdates(lastUpdateID)
		if err != nil {
			log.Println("Error getting updates:", err)
			continue
		}

		for _, update := range updates {
			if update.Message.Text == "" {
				continue
			}

			err := bot.SendMessage(
				update.Message.Chat.ID,
				"Echo: "+update.Message.Text,
			)
			if err != nil {
				log.Println("SendMessage error:", err)
			}

			lastUpdateID = update.UpdateID + 1
		}
	}
}
