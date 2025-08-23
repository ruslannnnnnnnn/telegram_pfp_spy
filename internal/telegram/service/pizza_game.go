package service

import (
	"fmt"
	"log"
	"math/rand"
	"spying_adelina/internal/common"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PizzaGame struct {
	bot       *tgbotapi.BotAPI
	appConfig *common.Config

	pizzaUpdatesChan chan tgbotapi.Update
}

func NewPizzaGame(bot *tgbotapi.BotAPI, appConfig *common.Config) *PizzaGame {
	return &PizzaGame{
		bot:              bot,
		appConfig:        appConfig,
		pizzaUpdatesChan: make(chan tgbotapi.Update),
	}
}

func (p *PizzaGame) HandleTelegramUpdate(update tgbotapi.Update) {
	p.pizzaUpdatesChan <- update
}

func (p *PizzaGame) Start() {
	for {
		randomDelayInSeconds := rand.Intn(p.appConfig.DelayBetweenPizzaGamesMax-p.appConfig.DelayBetweenPizzaGamesMin+1) + p.appConfig.DelayBetweenPizzaGamesMin
		sleepTime := time.Duration(randomDelayInSeconds) * time.Second
		time.Sleep(sleepTime)

		// Отправляем начальное сообщение и запоминаем время старта игры
		msg := tgbotapi.NewMessage(int64(p.appConfig.SpyingConfig.ChatId), "Кто первый напишет \"пицца\" тот победил")
		_, startPizzaGameErr := p.bot.Send(msg)
		if startPizzaGameErr != nil {
			log.Println("Не удалось отправить сообщение о начале игры в пиццу " + startPizzaGameErr.Error())
			continue
		}

		gameStartTime := time.Now() // Время старта игры
		gameTimeout := time.After(time.Duration(p.appConfig.TimeoutForPizzaGame) * time.Minute)
		gameOver := false

		for !gameOver {
			select {
			case update := <-p.pizzaUpdatesChan:
				if update.Message == nil {
					continue
				}

				// Игнорируем сообщения, отправленные ДО старта игры
				if update.Message.Time().Before(gameStartTime) {
					continue
				}

				// Проверяем, содержит ли сообщение слово "пицца"
				if strings.Contains(strings.ToLower(update.Message.Text), "пицца") {

					winner := fmt.Sprintf("@%s победил(а)!", update.Message.From.UserName)
					winnerMsg := tgbotapi.NewMessage(update.Message.Chat.ID, winner)
					_, winnerMsgSendErr := p.bot.Send(winnerMsg)
					if winnerMsgSendErr != nil {
						log.Println("Не удалось отправить сообщение о победе" + winnerMsgSendErr.Error())
						panic(winnerMsgSendErr)
					}

					gameOver = true
				}

			case <-gameTimeout:
				timeoutMsg := tgbotapi.NewMessage(int64(p.appConfig.SpyingConfig.ChatId), "Время вышло! Никто не успел написать \"пицца\"")
				_, sendTimeoutMessageError := p.bot.Send(timeoutMsg)
				if sendTimeoutMessageError != nil {
					log.Println("Не удалось отправить сообщение о таймауте в чат" + sendTimeoutMessageError.Error())
					panic(sendTimeoutMessageError)
				}
				gameOver = true
			}
		}
	}
}
