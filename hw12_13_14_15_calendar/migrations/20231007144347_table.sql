-- +goose Up
-- +goose StatementBegin
create table events (
                        id serial primary key,
                        owner bigint,
                        title text,
                        descr text,
                        start_date date not null,
                        start_time time,
                        end_date date not null,
                        end_time time,
                        remind_for timestamp,
                        deleted_at timestamp
);
create index owner_idx on events (owner);
create index start_idx on events using btree (start_date, start_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events;
-- +goose StatementEnd
