package service

import (
	"spying_adelina/internal/common"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramMessageAnalyser struct {
	appConfig *common.Config
	bot       *tgbotapi.BotAPI
	storage   *common.IAnalyticsStorage
}

func NewTelegramMessageAnalyser(bot *tgbotapi.BotAPI, appConfig *common.Config, storage *common.IAnalyticsStorage) *TelegramMessageAnalyser {
	return &TelegramMessageAnalyser{
		appConfig: appConfig,
		bot:       bot,
		storage:   storage,
	}
}

func (t *TelegramMessageAnalyser) HandleTelegramUpdate(update tgbotapi.Update) {

}
