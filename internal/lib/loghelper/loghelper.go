package loghelper

import (
	"context"
	"fmt"
	"log/slog"
)

// OpLogger структурв для loghelper
type OpLogger struct {
	logger *slog.Logger
	op     string
}

// New конструктор OpLogger
func New(logger *slog.Logger, op string) *OpLogger {
	return &OpLogger{
		logger: logger.With(slog.String("op", op)),
		op:     op,
	}
}

// LogDebug логирование LevelDebug
func (o *OpLogger) LogDebug(ctx context.Context, msg string, attrs ...slog.Attr) {
	o.logger.LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
}

// LogInfo логирование LevelInfo
func (o *OpLogger) LogInfo(ctx context.Context, msg string, attrs ...slog.Attr) {
	o.logger.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

// LogAndReturnError логирование LevelError и возврат ошибки
func (o *OpLogger) LogAndReturnError(ctx context.Context, msg string, err error, attrs ...slog.Attr) error {
	o.logger.LogAttrs(ctx, slog.LevelError, msg, append(attrs, slog.String("error", err.Error()))...)
	return fmt.Errorf("%s: %s: %w", o.op, msg, err)
}

// LogError логирование LevelError
func (o *OpLogger) LogError(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	o.logger.LogAttrs(ctx, slog.LevelError, msg, append(attrs, slog.String("error", err.Error()))...)
}
