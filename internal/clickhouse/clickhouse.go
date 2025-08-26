package clickhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"spying_adelina/internal/common"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ClickHouse struct {
	conn      driver.Conn
	appConfig *common.Config
}

func NewClickHouse(config *common.Config) (*ClickHouse, error) {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"clickhouse:9000"},
		Auth: clickhouse.Auth{
			Database: config.ClickHouseDb,
			Username: config.ClickHouseUser,
			Password: config.ClickHousePassword,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			//dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "telegram_pfp_spy", Version: "0.1"},
			},
		},
	})

	if err != nil {
		return &ClickHouse{}, err
	}

	return &ClickHouse{conn: conn, appConfig: config}, nil
}

func (c *ClickHouse) SaveTelegramUpdate(update tgbotapi.Update) error {

	jsonMessage, err := json.Marshal(update)
	if err != nil {
		return err
	}

	ctx := context.Background()

	err = c.conn.Exec(ctx, `
		INSERT INTO telegram_message_raw
		VALUES (?, ?, now())
	`, update.UpdateID, string(jsonMessage))

	if err != nil {
		return err
	}

	return nil
}

func (c *ClickHouse) SavePizzaWin(update tgbotapi.Update, gameStartTime time.Time) error {

	ctx := context.Background()

	err := c.conn.Exec(ctx, `
		INSERT INTO pizza_game_win
		VALUES (?, ?)
	`, update.UpdateID, gameStartTime.Format("2006-01-02T15:04:05"))

	if err != nil {
		return err
	}

	return nil
}

func (c *ClickHouse) GetPizzaWinnersLeaderBoard() ([]common.PizzaPlayer, error) {
	ctx := context.Background()
	rows, err := c.conn.Query(ctx, "SELECT username, wins FROM pizza_game_leaderboard FINAL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []common.PizzaPlayer

	for rows.Next() {
		var userName string
		var amountOfWins *uint64

		if err := rows.Scan(&userName, &amountOfWins); err != nil {
			return nil, err
		}
		result = append(result, common.PizzaPlayer{Username: userName, AmountOfWins: int(*amountOfWins)})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
