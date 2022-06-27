package logger

import (
	"context"
)

const (
	LogLevelErr int = iota
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

const (
	errLabel   = "ERROR"
	warnLabel  = "WARN"
	infoLabel  = "INFO"
	debugLabel = "DEBUG"
)

type LoggerSetting struct {
	LogLevel int
	Skip     int
}

type LogUsecase interface {
	SetRequestId(ctx context.Context) context.Context
	GetRequestId(ctx context.Context) string
	LogErr(ctx context.Context, msg ...interface{})
	LogWarning(ctx context.Context, msg ...interface{})
	LogInfo(ctx context.Context, msg ...interface{})
	LogDebug(ctx context.Context, msg ...interface{})
}
