package interceptors

import (
	"context"
	"log/slog"
	"time"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		log := loghelper.New(logger, "LoggerInterceptor")
		start := time.Now()
		defer func() {
			duration := time.Since(start)

			inf := status.Convert(err)
			var statusCode, errMsg string

			if err != nil {
				statusCode = inf.Code().String()
				errMsg = inf.Message()
			} else {
				statusCode = "OK"
			}

			log.LogInfo(ctx, "request gRPC",
				slog.String("method", info.FullMethod),
				slog.String("code", inf.Code().String()),
				slog.String("duration", duration.String()),
				slog.String("error", errMsg),
				slog.String("statusCode", statusCode))
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}
