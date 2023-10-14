package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	logger  Logger
	app     Application
	address string
	timeout time.Duration
	server  *http.Server
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface{}

func NewServer(logger Logger, app Application, address string) *Server {
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

		s.server = &http.Server{
			Addr:        s.address,
			Handler:     loggingMiddleware(s.logger, mux),
			ReadTimeout: s.timeout,
		}

		s.logger.Info(fmt.Sprintf("server start on address %s", s.address))
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
