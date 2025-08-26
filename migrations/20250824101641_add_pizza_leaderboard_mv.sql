-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW pizza_game_leaderboard_mv
    TO pizza_game_leaderboard
AS
SELECT
    JSONExtractInt(r.message_json, 'message', 'from', 'id') AS user_id,
    JSONExtractString(r.message_json, 'message', 'from', 'username') AS username,
    count() AS wins
FROM pizza_game_win w
         JOIN telegram_message_raw r
              ON r.update_id = w.update_id
GROUP BY user_id, username;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pizza_game_leaderboard_mv;
-- +goose StatementEnd
