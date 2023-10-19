-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ALTER COLUMN deleted_at DROP NOT NULL;
ALTER TABLE events ALTER COLUMN deleted_at SET DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE events SET deleted_at=0 WHERE deleted_at IS NULL;
ALTER TABLE events ALTER COLUMN deleted_at SET NOT NULL;
-- +goose StatementEnd
