package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/generated/event"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	logger  server.Logger
	app     server.Application
	address string
	server  *grpc.Server
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
		lsn, err := net.Listen("tcp", s.address)
		if err != nil {
			return errors.Wrap(err, "listening grpc")
		}

		logMiddleware := loggingMiddleware(s.logger)
		s.server = grpc.NewServer(logMiddleware)

		event.RegisterEventServiceServer(s.server, &EventServer{App: s.app, Logger: s.logger})

		s.logger.Info(fmt.Sprintf("grpc server start on address %s", s.address))

		return s.server.Serve(lsn)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		s.server.GracefulStop()
		return nil
	}
}
