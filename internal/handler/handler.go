package handler

import (
	"context"
	"strings"
)

type Handler struct {
	weather WeatherService
	bot     TelegramAPI
}

func New(weather WeatherService, bot TelegramAPI) *Handler {
	return &Handler{
		weather: weather,
		bot:     bot,
	}
}

func (h *Handler) Handle(ctx context.Context, chatID int64, text string) {

	text = strings.TrimSpace(text)

	if text == "/start" {
		_ = h.bot.SendMessage(ctx, chatID, "Напишіть назву міста")
		return
	}

	weatherText, err := h.weather.GetWeather(ctx, text)
	if err != nil {
		_ = h.bot.SendMessage(ctx, chatID, "Місто не знайдено")
		return
	}

	_ = h.bot.SendMessage(ctx, chatID, weatherText)
}
