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

func (s *Storage) AddEvent(_ *storage.Event) error {
	return nil
}
func (s *Storage) EditEvent(_ *storage.Event) error {
	return nil
}
func (s *Storage) DeleteEvent(_ int) error {
	return nil
}
func (s *Storage) ListEvents(_ time.Time, _ time.Time) ([]storage.Event, error) {
	return nil, nil
}
