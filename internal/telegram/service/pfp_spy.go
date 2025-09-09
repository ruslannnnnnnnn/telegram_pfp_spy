package service

import (
	"log"
	"spying_adelina/internal/common"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MonitorPfp Устанавливает слежку за аватаркой пользователя, при обновлении пишет сообщение в чат
func MonitorPfp(tgbot *tgbotapi.BotAPI, appConfig *common.Config, chatMember common.TelegramChatMember) {

	var lastPhotoID string

	for {
		photos, err := getPhoto(tgbot, int64(chatMember.UserId))

		pollInterval := time.Duration(appConfig.PollInterval) * time.Second
		if err != nil {
			log.Println("Не удалось получить фотки человека по имени " + chatMember.Name + ": " + err.Error())
			time.Sleep(pollInterval)
			continue
		}

		if len(photos.Photos) > 0 && len(photos.Photos[0]) > 0 {
			currentPhotoID := photos.Photos[0][0].FileID
			if currentPhotoID != lastPhotoID && lastPhotoID != "" {
				msg := tgbotapi.NewMessage(appConfig.SpyingConfig.ChatId, chatMember.PfpUpdateText)
				delay := time.Duration(appConfig.MinDelay) * time.Second

				if _, err := tgbot.Send(msg); err != nil {
					log.Println("Не удалось отправить сообщение об обновлении фотографии : " + err.Error())
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
