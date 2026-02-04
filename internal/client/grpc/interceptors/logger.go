package interceptors

import (
	"context"
	"log/slog"
	"time"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger *slog.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {

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
				slog.String("code", inf.Code().String()),
				slog.String("duration", duration.String()),
				slog.String("error", errMsg),
				slog.String("statusCode", statusCode))
		}()

		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}
