-- +goose Up
-- +goose StatementBegin
CREATE TABLE pizza_game_win
(
    update_id       Int64,
    game_start_date DateTime
) ENGINE = MergeTree()
      ORDER BY (update_id, game_start_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pizza_game_win;
-- +goose StatementEnd
