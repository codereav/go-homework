package app

import (
	"context"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

type App struct {
	Logger  logger.Log
	Storage Storage
}

var ErrNotExists = errors.New("record is not exists")

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	AddEvent(event *storage.Event) error
	EditEvent(event *storage.Event) error
	DeleteEvent(eventID int64) error
	ListEvents(from time.Time, to time.Time) ([]*storage.Event, error)
}

func New(logger logger.Log, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(
	ctx context.Context,
	title string,
	descr string,
	ownerID int64,
	startDate *time.Time,
	endDate *time.Time,
	remindFor *time.Time,
) (*int64, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("createEvent - context is done")
	default:
		event := &storage.Event{
			Title:     title,
			Descr:     descr,
			OwnerID:   ownerID,
			StartDate: startDate,
			EndDate:   endDate,
			RemindFor: remindFor,
		}
		err := a.Storage.AddEvent(event)
		if err != nil {
			return nil, errors.Wrap(err, "create event")
		}

		return &event.ID, nil
	}
}

func (a *App) EditEvent(
	ctx context.Context,
	id int64,
	title string,
	descr string,
	ownerID int64,
	startDate *time.Time,
	endDate *time.Time,
	remindFor *time.Time,
) error {
	select {
	case <-ctx.Done():
		return errors.New("editEvent - context is done")
	default:
		event := storage.Event{
			ID:        id,
			Title:     title,
			Descr:     descr,
			OwnerID:   ownerID,
			StartDate: startDate,
			EndDate:   endDate,
			RemindFor: remindFor,
		}
		err := a.Storage.EditEvent(&event)
		if err != nil {
			return errors.Wrap(err, "edit event")
		}

		return nil
	}
}

func (a *App) DeleteEvent(
	ctx context.Context,
	id int64,
) error {
	select {
	case <-ctx.Done():
		return errors.New("deleteEvent - context is done")
	default:
		err := a.Storage.DeleteEvent(id)
		if err != nil {
			return errors.Wrap(err, "delete event")
		}

		return nil
	}
}

func (a *App) ListEvents(
	ctx context.Context,
	dateFrom time.Time,
	dateTo time.Time,
) ([]*storage.Event, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("listEvents - context is done")
	default:
		events, err := a.Storage.ListEvents(dateFrom, dateTo)
		if err != nil {
			return nil, errors.Wrap(err, "list events")
		}

		return events, nil
	}
}
