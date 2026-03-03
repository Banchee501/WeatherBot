package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//читаємо .env

type Config struct {
	TelegramToken string
	WeatherAPIKey string
}

//зберігаємо токени

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set in .env file")
	}

	APIKey := os.Getenv("WEATHER_API_KEY")
	if APIKey == "" {
		log.Fatal("WEATHER_API_KEY is not set in .env file")
	}

	//повертаємо struct Config

	return &Config{
		TelegramToken: token,
		WeatherAPIKey: APIKey,
	}
}
