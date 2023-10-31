-- +goose Up
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN start_date;
ALTER TABLE events ADD COLUMN start_date timestamp;

ALTER TABLE events DROP COLUMN start_time;

ALTER TABLE events DROP COLUMN end_date;
ALTER TABLE events ADD COLUMN end_date timestamp;

ALTER TABLE events DROP COLUMN end_time;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN start_date;
ALTER TABLE events ADD COLUMN start_date date;
ALTER TABLE events ADD COLUMN start_time time;

ALTER TABLE events DROP COLUMN end_date;
ALTER TABLE events ADD COLUMN end_date date;
ALTER TABLE events ADD COLUMN end_time time;
-- +goose StatementEnd
