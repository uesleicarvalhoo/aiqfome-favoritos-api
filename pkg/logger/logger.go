package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.SugaredLogger

func Configure(opts Options) {
	zl := buildZap(opts)

	zap.ReplaceGlobals(zl)
	defaultLogger = zl.Sugar()
}

func init() {
	Configure(Options{
		Level:          "INFO",
		ServiceName:    "vision-manager",
		ServiceVersion: "0.0.0",
		Env:            "dev",
	})
}

func buildZap(opts Options) *zap.Logger {
	var lvl zapcore.Level
	if err := lvl.Set(opts.Level); err != nil {
		panic(err)
	}

	enc := zap.NewProductionEncoderConfig()
	enc.MessageKey = "message"
	enc.TimeKey = "time"
	enc.LevelKey = "level"
	enc.EncodeTime = zapcore.ISO8601TimeEncoder
	enc.EncodeCaller = zapcore.ShortCallerEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(lvl),
		Development:       opts.Env != "prd",
		DisableCaller:     false,
		DisableStacktrace: true,
		Encoding:          "json",
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig:     enc,
		InitialFields: map[string]any{
			"logVersion": "2.0.0",
			"app": map[string]string{
				"name":    opts.ServiceName,
				"version": opts.ServiceVersion,
				"env":     opts.Env,
			},
		},
	}

	return zap.Must(cfg.Build()).WithOptions(zap.AddCallerSkip(1))
}

func Info(ctx context.Context, msg string, args ...any) {
	withContext(ctx).Infof(msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	withContext(ctx).Debugf(msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	withContext(ctx).Warnf(msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	withContext(ctx).Errorf(msg, args...)
}

func Fatal(ctx context.Context, msg string, args ...any) {
	withContext(ctx).Fatalf(msg, args...)
}

func InfoF(ctx context.Context, msg string, f Fields, args ...any) {
	withContextAndFields(ctx, f).Infof(msg, args...)
}

func DebugF(ctx context.Context, msg string, f Fields, args ...any) {
	withContextAndFields(ctx, f).Debugf(msg, args...)
}

func WarnF(ctx context.Context, msg string, f Fields, args ...any) {
	withContextAndFields(ctx, f).Warnf(msg, args...)
}

func ErrorF(ctx context.Context, msg string, f Fields, args ...any) {
	withContextAndFields(ctx, f).Errorf(msg, args...)
}

func FatalF(ctx context.Context, msg string, f Fields, args ...any) {
	withContextAndFields(ctx, f).Fatalf(msg, args...)
}
