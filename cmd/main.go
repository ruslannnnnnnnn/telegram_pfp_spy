package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"spying_adelina/internal/telegram/service"
	"strconv"
)

const AppEnvPath = "/app/.env"
const ConfigFilePath = "/app/config.json"

func main() {
	err := godotenv.Load(AppEnvPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var appConfig service.SpyingConfig

	configJson, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatal("Error reading appConfig file")
	}

	err = json.Unmarshal(configJson, &appConfig)
	if err != nil {
		log.Fatal("Error parsing appConfig file")
	}

	pollInterval, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err != nil {
		log.Fatal("Error parsing POLL_INTERVAL")
	}
	minDelay, err := strconv.Atoi(os.Getenv("MIN_DELAY"))
	if err != nil {
		log.Fatal("Error parsing MIN_DELAY")
	}

	for _, chatMember := range appConfig.ChatMembers {

		spyingConfig := service.AppConfig{
			SpyingConfig: appConfig,
			PollInterval: pollInterval,
			MinDelay:     minDelay,
		}

		go service.MonitorPfp(spyingConfig, chatMember)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	ginError := r.Run()
	if ginError != nil {
		log.Fatal(ginError.Error())
	}
}
