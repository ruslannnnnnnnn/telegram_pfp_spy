package app

import (
	"encoding/json"
	"log"
	"os"
	"spying_adelina/internal/common"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	EnvFilePath    = "/app/.env"
	ConfigFilePath = "/app/config.json"
	// ENV var names
	PollIntervalENV             = "POLL_INTERVAL"
	MinDelayEnv                 = "MIN_DELAY"
	DelayBetweenPizzaGameMinEnv = "DELAY_BETWEEN_PIZZA_GAMES_MIN"
	DelayBetweenPizzaGameMaxEnv = "DELAY_BETWEEN_PIZZA_GAMES_MAX"
	TimeoutFormPizzaGameEnv     = "TIMEOUT_FOR_PIZZA_GAME"
)

func LoadConfig() common.Config {
	err := godotenv.Load(EnvFilePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var spyingConfig common.SpyingConfig

	configJson, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatal("Error reading config.json file")
	}

	err = json.Unmarshal(configJson, &spyingConfig)
	if err != nil {
		log.Fatal("Error parsing config.json file")
	}

	pollInterval, err := strconv.Atoi(os.Getenv(PollIntervalENV))
	if err != nil {
		log.Fatal("Error parsing " + PollIntervalENV)
	}
	minDelay, err := strconv.Atoi(os.Getenv(MinDelayEnv))
	if err != nil {
		log.Fatal("Error parsing " + MinDelayEnv)
	}

	delayBetweenPizzaGamesMin, err := strconv.Atoi(os.Getenv(DelayBetweenPizzaGameMinEnv))
	if err != nil {
		log.Fatal("Error parsing " + DelayBetweenPizzaGameMinEnv)
	}

	delayBetweenPizzaGamesMax, err := strconv.Atoi(os.Getenv(DelayBetweenPizzaGameMaxEnv))
	if err != nil {
		log.Fatal("Error parsing " + DelayBetweenPizzaGameMaxEnv)
	}

	return common.Config{
		SpyingConfig:              spyingConfig,
		PollInterval:              pollInterval,
		MinDelay:                  minDelay,
		DelayBetweenPizzaGamesMin: delayBetweenPizzaGamesMin,
		DelayBetweenPizzaGamesMax: delayBetweenPizzaGamesMax,
	}
}
