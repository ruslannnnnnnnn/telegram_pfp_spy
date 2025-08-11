package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"time"
)

func GetBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	return bot
}

type SpyingConfigJson struct {
	UserId      int    `json:"user_id"` // whom's pfp we're going to monitor
	ChatId      int    `json:"chat_id"` // where we will send messages about pfp updates
	MessageText string `json:"message_text"`
}

type SpyingConfig struct {
	SpyingConfigJson
	PollInterval int
	MinDelay     int
}

func MonitorPfp(config SpyingConfig) {
	bot := GetBot(os.Getenv("TELEGRAM_BOT_API_TOKEN"))
	var lastPhotoID string

	for {
		photos, err := getPhoto(bot, int64(config.UserId))

		pollInterval := time.Duration(config.PollInterval) * time.Second
		if err != nil {
			log.Printf("Failed to get profile photos: %v", err)
			time.Sleep(pollInterval)
			continue
		}

		if len(photos.Photos) > 0 && len(photos.Photos[0]) > 0 {
			currentPhotoID := photos.Photos[0][0].FileID
			if currentPhotoID != lastPhotoID && lastPhotoID != "" {
				log.Println("New profile photo detected")

				msg := tgbotapi.NewMessage(int64(config.ChatId), config.MessageText)
				delay := time.Duration(config.MinDelay) * time.Second

				if _, err := bot.Send(msg); err != nil {
					log.Printf("Failed to send message: %v", err)
					time.Sleep(delay)
					continue
				}

				lastPhotoID = currentPhotoID

				time.Sleep(delay)
				continue
			} else if lastPhotoID == "" {
				lastPhotoID = currentPhotoID
			}
		}

		time.Sleep(pollInterval)
	}
}

func getPhoto(bot *tgbotapi.BotAPI, userId int64) (tgbotapi.UserProfilePhotos, error) {
	photos, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userId, Limit: 1})
	return photos, err
}
