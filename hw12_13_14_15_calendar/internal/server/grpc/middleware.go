package internalgrpc

import (
	"context"
	"fmt"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func loggingMiddleware(log server.Logger) grpc.ServerOption {
	return grpc.UnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (interface{}, error) {
			resp, err := handler(ctx, req)
			// Доступ к метаданным запроса
			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				host := md.Get(":host")
				method := md.Get(":method")
				scheme := md.Get(":scheme")
				status := md.Get(":status")
				path := md.Get(":path")
				ua := md.Get("user-agent")
				log.Info(
					fmt.Sprintf(
						"%s %s %s %s %s %s",
						host, method, path, scheme, status, ua),
				)
			}

			return resp, err
		})
}
