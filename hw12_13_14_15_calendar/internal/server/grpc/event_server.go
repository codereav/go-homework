package internalgrpc

import (
	"context"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/generated/event"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventServer struct {
	App    server.Application
	Logger server.Logger
	event.UnimplementedEventServiceServer
}

func (s *EventServer) AddEvent(ctx context.Context, r *event.AddEventRequest) (*event.AddEventSuccessResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("context done triggered when AddEvent call - operation failed")
	default:
		s.Logger.Info("AddEvent method called")
		id, err := s.App.CreateEvent(
			ctx,
			r.Title,
			r.Descr,
			r.OwnerId,
			r.StartDate.AsTime(),
			r.EndDate.AsTime(),
			r.RemindFor.AsDuration(),
		)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, status.Errorf(codes.Aborted, err.Error())
		}

		response := event.AddEventSuccessResponse{
			Status: "success",
			Id:     *id,
		}

		return &response, nil
	}
}

func (s *EventServer) EditEvent(ctx context.Context, r *event.EventRequest) (*event.EditEventSuccessResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("context done triggered when EditEvent call - operation failed")
	default:
		s.Logger.Info("EditEvent method called")
		err := s.App.EditEvent(
			ctx,
			r.GetId(),
			r.GetTitle(),
			r.GetDescr(),
			r.GetOwnerId(),
			r.StartDate.AsTime(),
			r.EndDate.AsTime(),
			r.RemindFor.AsDuration(),
		)
		if err != nil {
			if errors.Is(err, app.ErrNotExists) {
				return nil, status.Errorf(codes.NotFound, err.Error())
			}
			s.Logger.Error(err.Error())
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		response := event.EditEventSuccessResponse{
			Status: "success",
		}

		return &response, nil
	}
}

func (s *EventServer) DeleteEvent(
	ctx context.Context, r *event.DeleteEventRequest,
) (*event.DeleteEventSuccessResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("context done triggered when EditEvent call - operation failed")
	default:
		s.Logger.Info("DeleteEvent method called")
		err := s.App.DeleteEvent(
			ctx,
			r.GetId(),
		)
		if err != nil {
			if errors.Is(err, app.ErrNotExists) {
				return nil, status.Errorf(codes.NotFound, err.Error())
			}
			s.Logger.Error(err.Error())

			return nil, status.Errorf(codes.Internal, err.Error())
		}

		response := event.DeleteEventSuccessResponse{
			Status: "success",
		}

		return &response, nil
	}
}

func (s *EventServer) ListEvents(
	ctx context.Context, r *event.ListEventsRequest,
) (*event.ListEventsSuccessResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("context done triggered when ListEvents call - operation failed")
	default:
		s.Logger.Info("ListEvents method called")
		res, err := s.App.ListEvents(
			ctx,
			r.DateFrom.AsTime(),
			r.DateTo.AsTime(),
		)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		var events []*event.EventResponse

		for _, e := range res {
			var deletedAt string
			if e.DeletedAt != nil {
				deletedAt = e.DeletedAt.String()
			}
			events = append(events, &event.EventResponse{
				Id:        e.ID,
				Title:     e.Title,
				Descr:     e.Descr,
				OwnerId:   e.OwnerID,
				StartDate: e.StartDate.String(),
				EndDate:   e.EndDate.String(),
				RemindFor: e.RemindFor.String(),
				DeletedAt: deletedAt,
			})
		}

		response := event.ListEventsSuccessResponse{
			Events: events,
		}

		return &response, nil
	}
}
