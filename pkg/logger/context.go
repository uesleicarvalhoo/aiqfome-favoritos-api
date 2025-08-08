package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

var loggerCtxKey ctxKey = "logger_ctx_key"

func withContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return defaultLogger
	}

	if v := ctx.Value(loggerCtxKey); v != nil {
		if sug, ok := v.(*zap.SugaredLogger); ok {
			return sug
		}
	}

	return defaultLogger
}

func withContextAndFields(ctx context.Context, f Fields) *zap.SugaredLogger {
	sug := withContext(ctx)
	if len(f) == 0 {
		return sug
	}

	return sug.With(f.toKeysAndValues()...)
}

func ContextWithFields(ctx context.Context, f Fields) context.Context {
	sug := defaultLogger
	if v := ctx.Value(loggerCtxKey); v != nil {
		if s2, ok := v.(*zap.SugaredLogger); ok {
			sug = s2
		}
	}

	child := sug.With(f.toKeysAndValues()...)

	return context.WithValue(ctx, loggerCtxKey, child)
}
