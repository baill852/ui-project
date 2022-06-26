package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type TplinkLogger struct {
	enableColor bool
	logger      *log.Logger
	LogLevel    int
	skip        int
}

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

func (tl *Logger) LogErr(requestId string, msgs ...interface{}) {
	if tl.LogLevel < LogLevelErr {
		return
	}

	tl.logMessage(errLabel, requestId, fmt.Sprintln(msgs...))
}

func (tl *Logger) LogWarning(requestId string, msgs ...interface{}) {
	if tl.LogLevel < LogLevelWarning {
		return
	}

	tl.logMessage(warnLabel, requestId, fmt.Sprintln(msgs...))
}

func (tl *Logger) LogInfo(requestId string, msgs ...interface{}) {
	if tl.LogLevel < LogLevelInfo {
		return
	}

	tl.logMessage(infoLabel, requestId, fmt.Sprintln(msgs...))
}

func (tl *Logger) LogDebug(requestId string, msgs ...interface{}) {
	if tl.LogLevel < LogLevelDebug {
		return
	}

	tl.logMessage(debugLabel, requestId, fmt.Sprintln(msgs...))
}

func (tl *Logger) logMessage(logLevel string, requestId string, msg string) {
	msg = strings.TrimRight(msg, "\n")
	now := &loggerTime{Time: time.Now()}

	tl.Logger.Printf(
		"%s %-6s [%d] [%s] %s %s\n",
		now,
		logLevel,
		syscall.Gettid(),
		requestId,
		msg,
		getFunctionPosition(tl.Skip),
	)
}

func getFunctionPosition(skip int) (result string) {
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
