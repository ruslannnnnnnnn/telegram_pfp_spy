package service

import (
	"log"
	"spying_adelina/internal/common"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramMessageAnalyser struct {
	appConfig *common.Config
	bot       *tgbotapi.BotAPI
	storage   common.IAnalyticsStorage
}

func NewTelegramMessageAnalyser(bot *tgbotapi.BotAPI, appConfig *common.Config, storage common.IAnalyticsStorage) *TelegramMessageAnalyser {
	return &TelegramMessageAnalyser{
		appConfig: appConfig,
		bot:       bot,
		storage:   storage,
	}
}

func (t *TelegramMessageAnalyser) HandleTelegramUpdate(update tgbotapi.Update) {
	err := t.storage.SaveTelegramUpdate(update)
	if err != nil {
		log.Println("Не удалось сохранить сообщение в базу " + err.Error())
	}

	if update.Message.IsCommand() && update.Message.Command() == "stats" {
		pizzaPlayers, leaderBoardErr := t.storage.GetPizzaWinnersLeaderBoard()
		if leaderBoardErr != nil {
			log.Fatal("Не удалось выполнить запрос на получение лидербордов по пицце" + leaderBoardErr.Error())
		}
		messageText := ""

		for _, pizzaPlayer := range pizzaPlayers {
			messageText += "@" + pizzaPlayer.Username + " - " + strconv.Itoa(pizzaPlayer.AmountOfWins) + "\n"
		}
		telegramMessage := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
		_, statsMessageError := t.bot.Send(telegramMessage)
		if statsMessageError != nil {
			log.Println("Не удалось отправить сообщение с лидербордом пиццы" + statsMessageError.Error())
		}
	}

}
