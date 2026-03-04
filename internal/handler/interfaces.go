package handler

import "context"

type WeatherService interface {
	GetWeather(ctx context.Context, city string) (string, error)
}

type TelegramAPI interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}
