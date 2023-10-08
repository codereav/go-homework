package internalhttp

import (
	"fmt"
	"net/http"
)

func loggingMiddleware(log Logger, next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(
			fmt.Sprintf(
				"%s %s %s %s %s %s",
				r.Host, r.Method, r.URL.Path, r.Proto, r.Response.StatusCode, r.UserAgent()),
		)
		next.ServeHTTP(w, r)
	})
}
