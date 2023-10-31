package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) { //nolint:funlen
	forAdd := []struct {
		in          storage.Event
		expectedErr error
	}{
		{
			in: storage.Event{
				Title:     "Test title1",
				Descr:     "descr",
				StartDate: time.Now().Add(3 * time.Hour),
				EndDate:   time.Now().Add(6 * time.Hour),
				RemindFor: 1 * time.Hour,
			},
			expectedErr: nil,
		},
		{
			in: storage.Event{
				Title:     "Test title2",
				Descr:     "descr",
				StartDate: time.Now().Add(3 * time.Hour),
				EndDate:   time.Now().Add(6 * time.Hour),
				RemindFor: 1 * time.Hour,
			},
			expectedErr: nil,
		},
		{
			in: storage.Event{
				Title:     "Test title3",
				Descr:     "descr",
				StartDate: time.Now().Add(3 * time.Hour),
				EndDate:   time.Now().Add(6 * time.Hour),
				RemindFor: 1 * time.Hour,
			},
			expectedErr: nil,
		},
		{
			in: storage.Event{
				Title:     "Test title4",
				Descr:     "descr",
				StartDate: time.Now().Add(3 * time.Hour),
				EndDate:   time.Now().Add(6 * time.Hour),
				RemindFor: 1 * time.Hour,
			},
			expectedErr: nil,
		},
		{
			in: storage.Event{
				Title:     "Test title5",
				Descr:     "descr",
				StartDate: time.Now().Add(3 * time.Hour),
				EndDate:   time.Now().Add(6 * time.Hour),
				RemindFor: 1 * time.Hour,
			},
			expectedErr: nil,
		},
		{
			in: storage.Event{
				Title:     "Test title6",
				Descr:     "descr",
				StartDate: time.Now().Add(3 * time.Hour),
				EndDate:   time.Now().Add(6 * time.Hour),
				RemindFor: 1 * time.Hour,
			},
			expectedErr: nil,
		},
		{
			in: storage.Event{
				Title:     "Test title7",
				Descr:     "descr",
				StartDate: time.Now().Add(3 * time.Hour),
				EndDate:   time.Now().Add(6 * time.Hour),
				RemindFor: 1 * time.Hour,
			},
			expectedErr: nil,
		},
	}
	s := New()
	ctx := context.Background()
	err := s.Connect(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("%w", err))
		return
	}
	defer func(s *Storage, ctx context.Context) {
		err := s.Close(ctx)
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
			return
		}
	}(s, ctx)

	for i, tt := range forAdd {
		t.Run(fmt.Sprintf("Test AddEvent %d", i), func(t *testing.T) {
			tt := tt
			err := s.AddEvent(&tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}

	t.Run("Test ListEvents 1", func(t *testing.T) {
		events, err := s.ListEvents(time.Now().Add(1*time.Hour), time.Now().Add(3*time.Hour))
		require.NoError(t, err)
		require.Equal(t, len(forAdd), len(events))
	})

	t.Run("Test ListEvents 2", func(t *testing.T) {
		var err error
		err = s.AddEvent(&storage.Event{
			Title:     "listEvents 2",
			StartDate: time.Now().Add(10 * time.Hour),
			EndDate:   time.Now().Add(12 * time.Hour),
		})
		require.NoError(t, err)
		events, err := s.ListEvents(time.Now().Add(9*time.Hour), time.Now().Add(11*time.Hour))
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
	})

	t.Run("Test DeleteEvent 1", func(t *testing.T) {
		eventsBefore, err := s.ListEvents(time.Now(), time.Now().Add(10*time.Hour))
		require.NoError(t, err)
		err = s.DeleteEvent(1)
		require.NoError(t, err)
		err = s.DeleteEvent(2)
		require.NoError(t, err)
		eventsAfter, err := s.ListEvents(time.Now(), time.Now().Add(11*time.Hour))
		require.NoError(t, err)
		require.Equal(t, len(eventsBefore)-2, len(eventsAfter))
	})

	t.Run("Test DeleteEvent 2", func(t *testing.T) {
		err = s.DeleteEvent(1)
		require.Error(t, err)
	})

	t.Run("Test EditEvent 1", func(t *testing.T) {
		var err error
		event := storage.Event{
			Title:     "editEvent 1",
			StartDate: time.Now().Add(10 * time.Hour),
			EndDate:   time.Now().Add(12 * time.Hour),
		}
		err = s.AddEvent(&event)
		require.NoError(t, err)
		require.Condition(t, func() (success bool) {
			return event.ID > 0
		})
		newEvent := storage.Event{
			ID:        event.ID,
			Title:     "editEvent 2",
			Descr:     "Test description 2",
			StartDate: time.Now().Add(10 * time.Hour),
			EndDate:   time.Now().Add(12 * time.Hour),
		}
		err = s.EditEvent(&newEvent)
		require.NoError(t, err)
		events, err := s.ListEvents(time.Now(), time.Now().Add(12*time.Hour))
		require.NoError(t, err)
		eventFound := false
		for _, val := range events {
			if val.ID == event.ID {
				eventFound = true
				require.Equal(t, val.Title, newEvent.Title)
				require.Equal(t, val.Descr, newEvent.Descr)
				break
			}
		}
		require.Equal(t, true, eventFound)
	})
}
