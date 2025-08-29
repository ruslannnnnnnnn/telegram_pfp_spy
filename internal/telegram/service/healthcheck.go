package service

import (
	"bufio"
	"log"
	"os"
	"spying_adelina/internal/app"
	"spying_adelina/internal/common"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HealthCheckService struct {
	appConfig *common.Config
	bot       *tgbotapi.BotAPI
	storage   common.IAnalyticsStorage
}

func NewHealthCheckService(bot *tgbotapi.BotAPI, appConfig *common.Config, storage common.IAnalyticsStorage) *HealthCheckService {
	return &HealthCheckService{
		appConfig: appConfig,
		bot:       bot,
		storage:   storage,
	}
}

// HandleTelegramUpdate выдает информацию о состоянии сервера: сколько ошибок за последние 3 часа, подключение к бд
func (t *HealthCheckService) HandleTelegramUpdate(update tgbotapi.Update) {
	if update.Message.IsCommand() && update.Message.Command() == "health" {
		healthCheckMessageText := "Дата: " + time.Now().Format(time.RFC3339) + "\n"
		statusOk := true
		responseTime, dbHealthErr := t.storage.Ping()
		if dbHealthErr != nil {
			healthCheckMessageText += "Ошибка при подключении к БД: " + dbHealthErr.Error() + "\n"
			statusOk = false
		} else {
			healthCheckMessageText += "Скорость ответа БД: " + strconv.FormatInt(responseTime.Milliseconds(), 10) + "мс\n"
		}

		errorCount, err := t.countRecentErrorsTail(app.ErrorsLogFilePath, 3*time.Hour, 128*1024) // читаем последние 128KB
		if err != nil {
			healthCheckMessageText += "Ошибка при чтении логов: " + err.Error() + "\n"
		} else {
			healthCheckMessageText += "Количество логов с ошибкой за последние 3 часа: " + strconv.Itoa(errorCount) + "\n"
		}

		if errorCount > 10 {
			statusOk = false
		}

		if statusOk {
			healthCheckMessageText = "Все ОК\n" + healthCheckMessageText
		} else {
			healthCheckMessageText = "Все плохо\n" + healthCheckMessageText
		}
		telegramMessage := tgbotapi.NewMessage(update.Message.Chat.ID, healthCheckMessageText)

		_, botSendMessageErr := t.bot.Send(telegramMessage)

		if botSendMessageErr != nil {
			log.Println("Не удалось отправить сообщение о HealthCheck в телеграм " + botSendMessageErr.Error())
		}
	}

}

// countRecentErrorsTail Подсчет количества логов за указанное время, читает из файла с логами последние tailSize байт
func (t *HealthCheckService) countRecentErrorsTail(
	logFilePath string,
	duration time.Duration,
	tailSize int64,
) (int, error) {
	f, err := os.Open(logFilePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size := fi.Size()
	start := size - tailSize
	if start < 0 {
		start = 0
	}

	_, err = f.Seek(start, 0)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(f)
	count := 0
	now := time.Now()

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 20 {
			continue
		}

		timestampStr := line[:19]
		timestamp, err := time.Parse("2006/01/02 15:04:05", timestampStr)
		if err != nil {
			continue
		}

		if now.Sub(timestamp) <= duration {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
