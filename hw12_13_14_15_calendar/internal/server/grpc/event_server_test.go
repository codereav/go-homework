package internalgrpc_test

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/generated/event"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server/grpc"
	memorystorage "github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEventServer_AddEvent(t *testing.T) {
	ctx := context.Background()
	server, err := getGRPCServer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = server.Stop(ctx)
		require.NoError(t, err)
	}()

	client := getGRPCClient(ctx)

	e, err := client.AddEvent(ctx, &event.AddEventRequest{
		Title:     "title for new event",
		Descr:     "description for new event",
		OwnerId:   123,
		StartDate: timestamppb.Now(),
		EndDate:   timestamppb.New(time.Now().Add(1 * time.Hour)),
		RemindFor: nil,
	})
	require.Equal(t, reflect.TypeOf(&event.AddEventSuccessResponse{}), reflect.TypeOf(e))
	require.NoError(t, err)

	res, err := client.ListEvents(ctx, &event.ListEventsRequest{
		DateFrom: timestamppb.New(time.Now().Add(-1 * time.Hour)),
		DateTo:   timestamppb.New(time.Now().Add(1 * time.Hour)),
	})
	require.Equal(t, reflect.TypeOf(&event.ListEventsSuccessResponse{}), reflect.TypeOf(res))
	require.NoError(t, err)
	events := res.GetEvents()
	require.Len(t, events, 1)
	require.Equal(t, events[0].Title, "title for new event")

	require.NoError(t, err)
}

func TestEventServer_DeleteEvent(t *testing.T) {
	ctx := context.Background()
	server, err := getGRPCServer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = server.Stop(ctx)
		require.NoError(t, err)
	}()

	client := getGRPCClient(ctx)

	e, err := client.AddEvent(ctx, &event.AddEventRequest{
		Title:     "title for new event",
		Descr:     "description for new event",
		OwnerId:   123,
		StartDate: timestamppb.Now(),
		EndDate:   timestamppb.New(time.Now().Add(1 * time.Hour)),
		RemindFor: nil,
	})
	require.Equal(t, reflect.TypeOf(&event.AddEventSuccessResponse{}), reflect.TypeOf(e))
	require.NoError(t, err)
	eventID := e.Id
	res, err := client.DeleteEvent(ctx, &event.DeleteEventRequest{Id: eventID})
	require.Equal(t, reflect.TypeOf(&event.DeleteEventSuccessResponse{}), reflect.TypeOf(res))
	require.NoError(t, err)

	err = server.Stop(ctx)
	require.NoError(t, err)
}

func TestEventServer_EditEvent(t *testing.T) {
	ctx := context.Background()
	server, err := getGRPCServer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = server.Stop(ctx)
		require.NoError(t, err)
	}()

	client := getGRPCClient(ctx)

	e, err := client.AddEvent(ctx, &event.AddEventRequest{
		Title:     "old title for new event",
		Descr:     "description for new event",
		OwnerId:   123,
		StartDate: timestamppb.Now(),
		EndDate:   timestamppb.New(time.Now().Add(1 * time.Hour)),
		RemindFor: nil,
	})
	require.Equal(t, reflect.TypeOf(&event.AddEventSuccessResponse{}), reflect.TypeOf(e))
	require.NoError(t, err)
	eventID := e.Id
	res, err := client.EditEvent(ctx, &event.EventRequest{
		Id:        eventID,
		Title:     "New title",
		Descr:     "New sescr",
		OwnerId:   456,
		StartDate: nil,
		EndDate:   nil,
		RemindFor: nil,
		DeletedAt: nil,
	})
	require.Equal(t, reflect.TypeOf(&event.EditEventSuccessResponse{}), reflect.TypeOf(res))
	require.NoError(t, err)

	err = server.Stop(ctx)
	require.NoError(t, err)
}

func TestEventServer_ListEvents(t *testing.T) {
	ctx := context.Background()
	server, err := getGRPCServer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = server.Stop(ctx)
		require.NoError(t, err)
	}()
	client := getGRPCClient(ctx)

	// Проверяем, что записи отсутствуют на начало теста
	res, err := client.ListEvents(ctx, &event.ListEventsRequest{
		DateFrom: timestamppb.New(time.Now().Add(-1 * time.Hour)),
		DateTo:   timestamppb.New(time.Now().Add(1 * time.Hour)),
	})
	require.Equal(t, reflect.TypeOf(&event.ListEventsSuccessResponse{}), reflect.TypeOf(res))
	require.NoError(t, err)
	events := res.GetEvents()
	require.Len(t, events, 0)

	// Добавляем событие
	e, err := client.AddEvent(ctx, &event.AddEventRequest{
		Title:     "new event",
		Descr:     "description",
		OwnerId:   123,
		StartDate: timestamppb.Now(),
		EndDate:   timestamppb.New(time.Now().Add(1 * time.Hour)),
		RemindFor: nil,
	})
	require.Equal(t, reflect.TypeOf(&event.AddEventSuccessResponse{}), reflect.TypeOf(e))
	require.NoError(t, err)
	eventID := e.Id

	// Проверяем, что событие есть в списке
	res, err = client.ListEvents(ctx, &event.ListEventsRequest{
		DateFrom: timestamppb.New(time.Now().Add(-1 * time.Hour)),
		DateTo:   timestamppb.New(time.Now().Add(1 * time.Hour)),
	})
	require.Equal(t, reflect.TypeOf(&event.ListEventsSuccessResponse{}), reflect.TypeOf(res))
	require.NoError(t, err)
	events = res.GetEvents()
	require.Len(t, events, 1)
	require.Equal(t, events[0].Id, eventID)

	require.NoError(t, err)
}

func getGRPCServer(ctx context.Context) (*internalgrpc.Server, error) {
	logs := logger.DummyLogger{}
	storage := memorystorage.New()
	err := storage.Connect(ctx)
	if err != nil {
		return nil, err
	}
	application := app.New(&logs, storage)
	server := internalgrpc.NewServer(&logs, application, "localhost:9191")
	go func() {
		if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	return server, nil
}

func getGRPCClient(ctx context.Context) event.EventServiceClient {
	// create gRpc client
	conn, err := grpc.DialContext(ctx, "localhost:9191",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return event.NewEventServiceClient(conn)
}
