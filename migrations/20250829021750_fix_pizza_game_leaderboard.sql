-- +goose Up
-- +goose StatementBegin
CREATE TABLE pizza_game_leaderboard_temp
(
    user_id Int64,
    wins    UInt64
) ENGINE = SummingMergeTree()
      ORDER BY user_id;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO pizza_game_leaderboard_temp
SELECT user_id, wins
FROM pizza_game_leaderboard FINAL;
-- +goose StatementEnd


-- +goose StatementBegin
DROP TABLE pizza_game_leaderboard;
-- +goose StatementEnd

-- +goose StatementBegin
rename table pizza_game_leaderboard_temp TO pizza_game_leaderboard;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE pizza_game_leaderboard_mv;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW pizza_game_leaderboard_mv
            TO pizza_game_leaderboard
AS
SELECT
    JSONExtractInt(r.message_json, 'message', 'from', 'id') AS user_id,
    count() AS wins
FROM pizza_game_win w
         JOIN telegram_message_raw r
              ON r.update_id = w.update_id
GROUP BY user_id;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE user_profiles
(
    user_id    Int64,
    username   String,
    updated_at DateTime
) ENGINE = ReplacingMergeTree(updated_at)
      ORDER BY user_id;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW user_profiles_mv
            TO user_profiles
AS
SELECT
    JSONExtractInt(message_json, 'message', 'from', 'id') AS user_id,
    JSONExtractString(message_json, 'message', 'from', 'username') AS username,
    now() AS updated_at
FROM telegram_message_raw;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'чет лень писать обратную миграцию';
-- +goose StatementEnd
