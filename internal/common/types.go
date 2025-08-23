package common

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type SpyingConfig struct {
	ChatId      int64                        `json:"chat_id"`
	ChatMembers map[int64]TelegramChatMember `json:"users"`
}

type TelegramChatMember struct {
	Name          string `json:"name"`
	UserId        int64  `json:"user_id"`
	PfpUpdateText string `json:"pfp_update_text"`
	PizzaWinText  string `json:"pizza_win_text"`
}

type Config struct {
	SpyingConfig              SpyingConfig
	PollInterval              int
	MinDelay                  int
	DelayBetweenPizzaGamesMin int
	DelayBetweenPizzaGamesMax int
}

type ITelegramUpdateHandler interface {
	HandleTelegramUpdate(update tgbotapi.Update)
}

type IAnalyticsStorage interface {
}
