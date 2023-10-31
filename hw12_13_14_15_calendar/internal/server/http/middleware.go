package internalhttp

import (
	"fmt"
	"net/http"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rww *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rww.ResponseWriter.WriteHeader(statusCode)
	rww.statusCode = statusCode
}

func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func loggingMiddleware(log server.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rww := NewResponseWriterWrapper(w)
		next.ServeHTTP(rww, r)
		log.Info(
			fmt.Sprintf(
				"%s %s %s %s %d %s",
				r.Host, r.Method, r.URL.Path, r.Proto, rww.statusCode, r.UserAgent()),
		)
	})
}
