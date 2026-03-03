package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	WeatherAPIKey string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set in .env file")
	}

	weatherKey := os.Getenv("WEATHER_API_KEY")
	if weatherKey == "" {
		log.Fatal("WEATHER_API_KEY is not set in .env file")
	}

	return &Config{
		TelegramToken: token,
		WeatherAPIKey: weatherKey,
	}
}
