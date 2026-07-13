package logger

import (
	"context"
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() error {
	logger, err := zap.NewProduction()

	if err != nil {
		return err
	}

	Log = logger

	return nil
}

func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Log
	}

	requestID := ""
	if reqID, ok := ctx.Value("request_id").(string); ok {
		requestID = reqID
	}

	traceID := ""
	if trace, ok := ctx.Value("trace_id").(string); ok {
		traceID = trace
	}

	fields := []zap.Field{}
	if requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if len(fields) == 0 {
		return Log
	}

	return Log.With(fields...)
}
