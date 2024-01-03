package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

type Scheduler struct {
	Logger    logger.Log
	Storage   Storage
	Client    rabbitmq.Client
	PeriodSec int16
	OldDate   time.Time
	Exchange  string
	Key       string
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	ListNotSentEvents(from time.Time) ([]*storage.Event, error)
	MarkEventsAsSent(eventIDs []int64) (int64, error)
	DeleteOldEvents(oldDate time.Time) (int64, error)
}

func New(logger logger.Log, storage Storage, client rabbitmq.Client,
	periodSec int16, oldDate time.Time, exchange, key string,
) *Scheduler {
	return &Scheduler{
		Logger:    logger,
		Storage:   storage,
		Client:    client,
		PeriodSec: periodSec,
		OldDate:   oldDate,
		Exchange:  exchange,
		Key:       key,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(s.PeriodSec) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return errors.New("scheduler context is done")
		case <-ticker.C:
			s.Logger.Info("try to delete old events...")

			countOfDeleted, err := s.Storage.DeleteOldEvents(s.OldDate)
			if err != nil {
				return errors.Wrap(err, "scheduler fails on deleteOldEvents")
			}
			if countOfDeleted > 0 {
				s.Logger.Info(fmt.Sprintf("Old events was deleted. Count of deleted events: %d", countOfDeleted))
			} else {
				s.Logger.Info("no old events found")
			}
			s.Logger.Info("try to find events for notification...")

			var events []*storage.Event

			events, err = s.Storage.ListNotSentEvents(s.OldDate)
			if err != nil {
				return errors.Wrap(err, "scheduler fails on listNotSentEvents")
			}

			var sentEventIDs []int64
			for _, e := range events {
				data := &rabbitmq.NotifyMessage{
					ID:        e.ID,
					Title:     e.Title,
					OwnerID:   e.OwnerID,
					StartDate: *e.StartDate,
				}
				msg, err := json.Marshal(data)
				if err != nil {
					return errors.Wrap(err, "scheduler fails on message preparing")
				}
				if err = s.Client.Publish(s.Exchange, s.Key, msg); err != nil {
					return errors.Wrap(err, "scheduler fails on message publishing")
				}
				sentEventIDs = append(sentEventIDs, e.ID)
			}
			if len(sentEventIDs) == 0 {
				s.Logger.Info("no new events found")
				continue
			}
			s.Logger.Info("events were published.")
			s.Logger.Info("try to update events notify sent flag...")
			var cnt int64
			cnt, err = s.Storage.MarkEventsAsSent(sentEventIDs)
			if err != nil {
				return errors.Wrap(err, "scheduler fails on mark events as sent")
			}
			s.Logger.Info(fmt.Sprintf("events marked as sent successfully. Total: %d", cnt))
		}
	}
}
