package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	app.Storage
	events      map[int64]*storage.Event
	lastEventID int64
	mu          sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(_ context.Context) error {
	s.events = make(map[int64]*storage.Event)

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = nil

	return nil
}

func (s *Storage) AddEvent(event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	event.ID = s.lastEventID + 1
	s.events[event.ID] = event
	s.lastEventID = event.ID

	return nil
}

func (s *Storage) EditEvent(event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.events[event.ID]
	if !ok || e.DeletedAt != nil {
		return app.ErrNotExists
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(eventID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.events[eventID]
	if !ok || e.DeletedAt != nil {
		return app.ErrNotExists
	}
	now := time.Now()
	s.events[eventID].DeletedAt = &now

	return nil
}

func (s *Storage) ListEvents(from time.Time, to time.Time) ([]*storage.Event, error) {
	var result []*storage.Event

	for _, event := range s.events {
		if event.StartDate.After(from) && event.StartDate.Before(to) && event.DeletedAt == nil {
			result = append(result, event)
		}
	}

	return result, nil
}
