package server

import (
	"context"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface {
	CreateEvent(
		ctx context.Context,
		title string,
		descr string,
		ownerID int64,
		startDate *time.Time,
		endDate *time.Time,
		remindFor *time.Time,
	) (*int64, error)
	EditEvent(
		ctx context.Context,
		id int64,
		title string,
		descr string,
		ownerID int64,
		startDate *time.Time,
		endDate *time.Time,
		remindFor *time.Time,
	) error
	DeleteEvent(
		ctx context.Context,
		id int64,
	) error
	ListEvents(
		ctx context.Context,
		dateFrom time.Time,
		dateTo time.Time,
	) ([]*storage.Event, error)
}
