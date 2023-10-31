package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/generated/event"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
)

type Server struct {
	logger  server.Logger
	app     server.Application
	address string
	timeout time.Duration
	server  *http.Server
}

func NewServer(logger server.Logger, app server.Application, address string) *Server {
	return &Server{
		logger:  logger,
		app:     app,
		address: address,
	}
}

func (s *Server) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		var err error
		mux := http.NewServeMux()
		mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
			_, err = writer.Write([]byte("Hello world"))
			if err != nil {
				s.logger.Error(err.Error())
			}
		})

		runtimeMux := runtime.NewServeMux()
		err = s.registerGRPCGatewayAPIHandlers(ctx, runtimeMux)
		if err != nil {
			return errors.Wrap(err, "registering grpc-gateway api handlers")
		}
		mux.Handle("/api/", runtimeMux)

		s.server = &http.Server{
			Addr:        s.address,
			Handler:     loggingMiddleware(s.logger, mux),
			ReadTimeout: s.timeout,
		}

		s.logger.Info(fmt.Sprintf("http server start on address %s", s.address))

		return s.server.ListenAndServe()
	}
}

func (s *Server) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		return s.server.Close()
	}
}

func (s *Server) registerGRPCGatewayAPIHandlers(ctx context.Context, runtimeMux *runtime.ServeMux) error {
	err := event.RegisterEventServiceHandlerServer(
		ctx,
		runtimeMux,
		&internalgrpc.EventServer{App: s.app, Logger: s.logger})
	if err != nil {
		return errors.Wrap(err, "registering EventServer for grpc-gateway at http server")
	}

	return nil
}
