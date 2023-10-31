package sqlstorage

import (
	"context"
	"strconv"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Storage struct {
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
	query := `insert into events(title, descr, start_date, end_date, remind_for)
 values($1, $2, $3, $4, $5) returning id`
	var eventID int64
	err := s.conn.QueryRow(
		context.Background(),
		query,
		event.Title,
		event.Descr,
		strconv.Quote(event.StartDate.Format(time.DateTime)),
		strconv.Quote(event.EndDate.Format(time.DateTime)),
		strconv.Quote(event.RemindFor.Format(time.DateTime)),
	).Scan(&eventID)
	if err != nil {
		return err
	}
	event.ID = eventID

	return nil
}

func (s *Storage) EditEvent(event *storage.Event) error {
	query := `update events SET title=$1, descr=$2, start_date=$3, end_date=$4, remind_for=$5
 WHERE id=$6 AND deleted_at IS NULL`
	_, err := s.conn.Exec(
		context.Background(),
		query,
		event.Title,
		event.Descr,
		strconv.Quote(event.StartDate.Format(time.DateTime)),
		strconv.Quote(event.EndDate.Format(time.DateTime)),
		strconv.Quote(event.RemindFor.Format(time.DateTime)),
		event.ID,
	)

	return err
}

func (s *Storage) DeleteEvent(id int64) error {
	query := `UPDATE events SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NOT NULL`
	_, err := s.conn.Exec(
		context.Background(),
		query,
		id,
	)

	return err
}

func (s *Storage) ListEvents(from time.Time, to time.Time) ([]*storage.Event, error) {
	query := `
 select id, title, descr, start_date, end_date, remind_for
 from events
 where start_date BETWEEN $1 AND $2
 AND deleted_at IS NULL
`
	rows, err := s.conn.Query(context.Background(), query, strconv.Quote(from.Format(time.DateTime)),
		strconv.Quote(to.Format(time.DateTime)),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*storage.Event

	for rows.Next() {
		event := &storage.Event{}
		if err := rows.Scan(&event.ID, &event.Title, &event.Descr,
			&event.StartDate, &event.EndDate, &event.RemindFor); err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) ListNotSentEvents(from time.Time) ([]*storage.Event, error) {
	query := `
 select id, title, descr, start_date, end_date
 from events
 where start_date >= $1 
   AND deleted_at IS NULL AND remind_for IS NOT NULL AND notify_sent = false
   AND CURRENT_TIMESTAMP < start_date
   AND CURRENT_TIMESTAMP > remind_for
`
	rows, err := s.conn.Query(context.Background(), query, strconv.Quote(from.Format(time.DateTime)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*storage.Event

	for rows.Next() {
		event := &storage.Event{}
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

func (s *Storage) MarkEventsAsSent(eventIDs []int64) (int64, error) {
	query := `
 UPDATE events SET notify_sent=true
 WHERE id = any($1)
`
	c, err := s.conn.Exec(context.Background(), query, eventIDs)
	if err != nil {
		return 0, err
	}
	rowsCount := c.RowsAffected()

	return rowsCount, nil
}

func (s *Storage) DeleteOldEvents(oldDate time.Time) (int64, error) {
	query := `
 UPDATE events SET deleted_at=CURRENT_TIMESTAMP
 WHERE start_date < $1 AND deleted_at IS NULL;
`
	c, err := s.conn.Exec(context.Background(), query, strconv.Quote(oldDate.Format(time.DateTime)))
	if err != nil {
		return 0, err
	}
	rowsCount := c.RowsAffected()

	return rowsCount, nil
}
