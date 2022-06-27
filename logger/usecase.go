package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
)

const ContextRequestID string = "requestID"

// log format
// %d{MM-dd-yyyy HH:mm:ss.SSS} %p [%t] [%traceId] [%X{REQUEST_ID}] %c{1.}: %m%n
// %d{MM-dd-yyyy HH:mm:ss.SSS}: time format
// %p: log level
// %t: current thread
// %traceId: trace ID (from tracing)
// %X{REQUEST_ID}: used for search
// %c{1.}: class name
// %m: message, use | if we need to output multiple messages
// %n: new line

type loggerTime struct {
	time.Time
}

func (t loggerTime) String() string {
	yyyy, MM, dd := t.Date()
	HH, mm, ss, SSS := t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%03d", yyyy, MM, dd, HH, mm, ss, SSS)
}

type Logger struct {
	Logger   *log.Logger
	LogLevel int
	Skip     int
}

func CreateLogger(settings LoggerSetting) LogUsecase {
	return &Logger{
		LogLevel: settings.LogLevel,
		Logger:   log.New(os.Stdout, "", 0),
		Skip:     settings.Skip,
	}
}

func (l *Logger) SetRequestId(ctx context.Context) context.Context {
	key := uuid.New()

	return context.WithValue(ctx, ContextRequestID, key.String())
}

func (l *Logger) GetRequestId(ctx context.Context) string {
	requestId := ctx.Value(ContextRequestID)

	if value, ok := requestId.(string); ok {
		return value
	}

	return ""
}

func (l *Logger) LogErr(ctx context.Context, msgs ...interface{}) {
	if l.LogLevel < LogLevelErr {
		return
	}

	l.logMessage(errLabel, l.GetRequestId(ctx), fmt.Sprintln(msgs...))
}

func (l *Logger) LogWarning(ctx context.Context, msgs ...interface{}) {
	if l.LogLevel < LogLevelWarning {
		return
	}

	l.logMessage(warnLabel, l.GetRequestId(ctx), fmt.Sprintln(msgs...))
}

func (l *Logger) LogInfo(ctx context.Context, msgs ...interface{}) {
	if l.LogLevel < LogLevelInfo {
		return
	}

	l.logMessage(infoLabel, l.GetRequestId(ctx), fmt.Sprintln(msgs...))
}

func (l *Logger) LogDebug(ctx context.Context, msgs ...interface{}) {
	if l.LogLevel < LogLevelDebug {
		return
	}

	l.logMessage(debugLabel, l.GetRequestId(ctx), fmt.Sprintln(msgs...))
}

func (l *Logger) logMessage(logLevel string, requestId string, msg string) {
	msg = strings.TrimRight(msg, "\n")
	now := &loggerTime{Time: time.Now()}

	l.Logger.Printf(
		"%s %-6s [%d] [%s] %s %s\n",
		now,
		logLevel,
		syscall.Gettid(),
		requestId,
		msg,
		l.getFunctionPosition(l.Skip),
	)
}

func (l *Logger) getFunctionPosition(skip int) (result string) {
	var (
		function uintptr
		fileName string
	)

	function, file, lineNumber, _ := runtime.Caller(skip)

	functionName := strings.Split(runtime.FuncForPC(function).Name(), ".")
	fileFullPath := strings.Split(file, "/")

	filePathSkip := 3
	count := 0
	for i := len(fileFullPath); i > 0; i-- {
		fileName = fileFullPath[i-1] + "/" + fileName
		count++
		if count == filePathSkip {
			break
		}
	}

	result = fmt.Sprintf("%s:%d-%s", strings.TrimRight(fileName, "/"), lineNumber, functionName[len(functionName)-1])
	return
}
