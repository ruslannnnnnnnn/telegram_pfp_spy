package main

import (
	"log"
	"net/http"
	"os"
	"spying_adelina/internal/app"
	telegramservice "spying_adelina/internal/telegram/service"

	"github.com/gin-gonic/gin"
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
		go telegramservice.MonitorPfp(bot, appConfig, chatMember)
	}

	go telegramservice.PizzaGame(bot, appConfig)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://youtu.be/dQw4w9WgXcQ?si=lDs5Dg8PRgZLTM6T")
	})

	ginError := r.Run()
	if ginError != nil {
		log.Fatal(ginError.Error())
	}
}
