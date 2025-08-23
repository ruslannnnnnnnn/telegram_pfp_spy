-- +goose Up
-- +goose StatementBegin
CREATE TABLE telegram_message_raw
(
    update_id    Int64,
    message_json String,
    created_at   DateTime DEFAULT now()
) ENGINE = MergeTree()
      ORDER BY (update_id)
      SETTINGS allow_nullable_key = 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE telegram_message_raw;
-- +goose StatementEnd
