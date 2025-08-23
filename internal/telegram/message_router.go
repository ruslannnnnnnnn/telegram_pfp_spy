package telegram

import (
	"spying_adelina/internal/common"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ListenToUpdates(bot *tgbotapi.BotAPI, handlers ...common.TelegramUpdateHandler) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:

			if update.Message == nil {
				continue
			}

			for _, handler := range handlers {
				go handler.HandleTelegramUpdate(update)
			}
		}
	}
}
