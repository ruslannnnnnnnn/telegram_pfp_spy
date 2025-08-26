-- +goose Up
-- +goose StatementBegin
INSERT INTO pizza_game_leaderboard
SELECT
    JSONExtractInt(r.message_json, 'message', 'from', 'id') AS user_id,
    JSONExtractString(r.message_json, 'message', 'from', 'username') AS username,
    1 AS wins
FROM pizza_game_win w
         JOIN telegram_message_raw r
              ON r.update_id = w.update_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
