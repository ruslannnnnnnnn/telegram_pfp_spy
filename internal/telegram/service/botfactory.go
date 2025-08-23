package service

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MakeBotByToken(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	return bot
}
