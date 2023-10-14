package app

import (
	"context"
	"errors"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

var (
	ErrDateBusy  = errors.New("date is busy")
	ErrNotExists = errors.New("record is not exists")
)

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	AddEvent(event *storage.Event) error
	EditEvent(event *storage.Event) error
	DeleteEvent(eventID int) error
	ListEvents(from time.Time, to time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(
	ctx context.Context,
	title string,
	descr string,
	startDate time.Time,
	endDate time.Time,
	remindFor time.Duration,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			return a.Storage.AddEvent(
				&storage.Event{
					Title:     title,
					Descr:     descr,
					StartDate: startDate,
					EndDate:   endDate,
					RemindFor: remindFor,
				},
			)
		}
	}
}
