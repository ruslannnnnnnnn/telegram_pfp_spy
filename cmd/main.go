package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"spying_adelina/internal/app"
	telegram_service "spying_adelina/internal/telegram/service"
)

const (
	TelegramBotApiToken = "TELEGRAM_BOT_API_TOKEN"
)

func main() {

	appConfig := app.LoadConfig()

	// создаем бота епта
	bot := telegram_service.MakeBotByToken(os.Getenv(TelegramBotApiToken))

	// для каждого пользователя из чата запускаем слежку за автаркой
	for _, chatMember := range appConfig.SpyingConfig.ChatMembers {
		go telegram_service.MonitorPfp(bot, appConfig, chatMember)
	}

	go telegram_service.PizzaGame(bot, appConfig)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	ginError := r.Run()
	if ginError != nil {
		log.Fatal(ginError.Error())
	}
}
