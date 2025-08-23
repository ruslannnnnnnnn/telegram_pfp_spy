package main

import (
	"log"
	"os"
	"spying_adelina/internal/app"
	"spying_adelina/internal/telegram"
	telegramservice "spying_adelina/internal/telegram/service"
)

const (
	TelegramBotApiToken = "TELEGRAM_BOT_API_TOKEN"
	ErrorsLogFilePath   = "/app/log/errors.log"
)

func main() {

	// Открываем файл для логирования (если не существует - создаётся)
	file, err := os.OpenFile(ErrorsLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка открытия файла логов: ", err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Не удалось закрыть файл с логами" + err.Error())
		}
	}(file)

	// Устанавливаем вывод логов в файл
	log.SetOutput(file)

	appConfig := app.LoadConfig()

	// создаем бота епта
	bot := telegramservice.MakeBotByToken(os.Getenv(TelegramBotApiToken))

	// для каждого пользователя из чата запускаем слежку за автаркой
	for _, chatMember := range appConfig.SpyingConfig.ChatMembers {
		go telegramservice.MonitorPfp(bot, &appConfig, chatMember)
	}

	pizzaGame := telegramservice.NewPizzaGame(bot, &appConfig)

	go pizzaGame.Start()

	telegram.ListenToUpdates(bot, pizzaGame)
}
