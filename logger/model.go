package logger

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
	LogErr(requestId string, msg ...interface{})
	LogWarning(requestId string, msg ...interface{})
	LogInfo(requestId string, msg ...interface{})
	LogDebug(requestId string, msg ...interface{})
}
