package main

import (
	"log"
	"os"
	"spying_adelina/internal/app"
	clickhouse "spying_adelina/internal/clickhouse"
	"spying_adelina/internal/telegram"
	telegramservice "spying_adelina/internal/telegram/service"
)

const (
	TelegramBotApiToken = "TELEGRAM_BOT_API_TOKEN"
)

func main() {

	// Открываем файл для логирования (если не существует - создаётся)
	file, err := os.OpenFile(app.ErrorsLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

	clickhouseInstance, err := clickhouse.NewClickHouse(&appConfig)
	if err != nil {
		log.Fatal("Could not connect to clickhouse db " + err.Error())
	}

	// создаем бота епта
	bot := telegramservice.MakeBotByToken(os.Getenv(TelegramBotApiToken))

	// сборщик статистики
	telegramMessageAnalyser := telegramservice.NewTelegramMessageAnalyser(bot, &appConfig, clickhouseInstance)
	healthCheckService := telegramservice.NewHealthCheckService(bot, &appConfig, clickhouseInstance)

	// для каждого пользователя из чата запускаем слежку за аватаркой
	for _, chatMember := range appConfig.SpyingConfig.ChatMembers {
		go telegramservice.MonitorPfp(bot, &appConfig, chatMember)
	}

	pizzaGame := telegramservice.NewPizzaGame(bot, &appConfig, clickhouseInstance)

	go pizzaGame.Start()

	telegram.ListenToUpdates(bot, pizzaGame, telegramMessageAnalyser, healthCheckService)
}
