-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD COLUMN notify_sent bool NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN notify_sent;
-- +goose StatementEnd
