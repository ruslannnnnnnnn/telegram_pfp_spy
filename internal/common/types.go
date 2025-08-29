package common

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	ClickHouseUser            string
	ClickHousePassword        string
	ClickHouseDb              string
}

type ITelegramUpdateHandler interface {
	HandleTelegramUpdate(tgbotapi.Update)
}

type IAnalyticsStorage interface {
	SaveTelegramUpdate(tgbotapi.Update) error
	SavePizzaWin(tgbotapi.Update, time.Time) error
	GetPizzaWinnersLeaderBoard() ([]PizzaPlayer, error)
	Ping() (time.Duration, error)
}

type PizzaPlayer struct {
	Username     string `json:"username"`
	AmountOfWins int    `json:"amount_of_wins"`
}
