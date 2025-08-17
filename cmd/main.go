package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"spying_adelina/internal/app"
	telegramservice "spying_adelina/internal/telegram/service"
)

const (
	TelegramBotApiToken = "TELEGRAM_BOT_API_TOKEN"
	ErrosLogFilePath    = "/app/log/errors.log"
)

func main() {

	// Открываем файл для логирования (если не существует - создаётся)
	file, err := os.OpenFile(ErrosLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
		go telegramservice.MonitorPfp(bot, appConfig, chatMember)
	}

	go telegramservice.PizzaGame(bot, appConfig)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	ginError := r.Run()
	if ginError != nil {
		log.Fatal(ginError.Error())
	}
}
