package common

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type SpyingConfig struct {
	ChatId      int                        `json:"chat_id"`
	ChatMembers map[int]TelegramChatMember `json:"users"`
}

type TelegramChatMember struct {
	Name          string `json:"name"`
	UserId        int    `json:"user_id"`
	PfpUpdateText string `json:"pfp_update_text"`
	PizzaWinText  string `json:"pizza_win_text"`
}

type Config struct {
	SpyingConfig              SpyingConfig
	PollInterval              int
	MinDelay                  int
	DelayBetweenPizzaGamesMin int
	DelayBetweenPizzaGamesMax int
	TimeoutForPizzaGame       int
}

type TelegramUpdateHandler interface {
	HandleTelegramUpdate(update tgbotapi.Update)
}
