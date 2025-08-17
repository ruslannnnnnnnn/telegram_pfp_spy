package common

type SpyingConfig struct {
	ChatId      int                  `json:"chat_id"`
	ChatMembers []TelegramChatMember `json:"users"`
}

type TelegramChatMember struct {
	UserId      int    `json:"user_id"`
	MessageText string `json:"message_text"`
}

type Config struct {
	SpyingConfig              SpyingConfig
	PollInterval              int
	MinDelay                  int
	DelayBetweenPizzaGamesMin int
	DelayBetweenPizzaGamesMax int
	TimeoutForPizzaGame       int
}
