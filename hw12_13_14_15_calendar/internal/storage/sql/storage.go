package sqlstorage

import (
	"context"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Storage struct {
	app.Storage
	dsn  string
	conn *pgx.Conn
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		return err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return err
	}

	s.conn = conn

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.conn != nil {
		err := s.conn.Close(ctx)
		if err != nil {
			return errors.Wrap(err, "fail when close db connection")
		}
	}

	return nil
}

func (s *Storage) AddEvent(event *storage.Event) error {
	query := `insert into events(title, descr, start_date, end_date)
 values($1, $2, $3, $4) returning id`
	var eventID int
	err := s.conn.QueryRow(
		context.Background(),
		query,
		event.Title,
		event.Descr,
		event.StartDate.String(),
		event.EndDate.String(),
	).Scan(&eventID)
	if err != nil {
		return err
	}
	event.ID = eventID

	return nil
}

func (s *Storage) EditEvent(event *storage.Event) error {
	query := `update events SET title=$1, descr=$2, start_date=$3, end_date=$4
 WHERE id=$5`
	_, err := s.conn.Exec(
		context.Background(),
		query,
		event.Title,
		event.Descr,
		event.StartDate.String(),
		event.EndDate.String(),
	)

	return err
}

func (s *Storage) DeleteEvent(id int) error {
	query := `delete from eventsWHERE id=$1`
	_, err := s.conn.Exec(
		context.Background(),
		query,
		id,
	)

	return err
}

func (s *Storage) ListEvents(from time.Time, to time.Time) ([]storage.Event, error) {
	query := `
 select id, title, descr, start_date, end_date
 from events
 where start_date BETWEEN $1 AND $2
`
	rows, err := s.conn.Query(context.Background(), query, from.String(), to.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []storage.Event
	for rows.Next() {
		event := storage.Event{}
		if err := rows.Scan(&event.ID, &event.Title, &event.Descr, &event.StartDate, &event.EndDate); err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
