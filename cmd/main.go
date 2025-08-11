package main

import (
	"encoding/json"
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

	var configs []service.SpyingConfigJson

	configJson, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatal("Error reading config file")
	}

	err = json.Unmarshal(configJson, &configs)
	if err != nil {
		log.Fatal("Error parsing config file")
	}

	pollInterval, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err != nil {
		log.Fatal("Error parsing POLL_INTERVAL")
	}
	minDelay, err := strconv.Atoi(os.Getenv("MIN_DELAY"))
	if err != nil {
		log.Fatal("Error parsing MIN_DELAY")
	}

	for _, jsonConfig := range configs {

		config := service.SpyingConfig{
			SpyingConfigJson: jsonConfig,
			PollInterval:     pollInterval,
			MinDelay:         minDelay,
		}

		go service.MonitorPfp(config)
	}

	select {}
}
