-- +goose Up
-- +goose StatementBegin
CREATE TABLE pizza_game_leaderboard
(
    user_id   Int64,
    username  String,
    wins      UInt64
) ENGINE = SummingMergeTree()
      ORDER BY (user_id, username);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pizza_game_leaderboard;
-- +goose StatementEnd
