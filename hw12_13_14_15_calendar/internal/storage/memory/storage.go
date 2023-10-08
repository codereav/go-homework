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
	events      map[int]*storage.Event
	lastEventID int
	mu          sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(_ context.Context) error {
	s.events = make(map[int]*storage.Event)
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
	_, ok := s.events[event.ID]
	if !ok {
		return app.ErrNotExists
	}
	s.events[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(eventID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, eventID)

	return nil
}

func (s *Storage) ListEvents(from time.Time, to time.Time) ([]storage.Event, error) {
	var result []storage.Event
	for _, event := range s.events {
		event := event
		if event.StartDate.After(from) && event.StartDate.Before(to) {
			result = append(result, *event)
		}
	}

	return result, nil
}
