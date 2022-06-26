package main

import (
	"ui-project/lib"
	"ui-project/logger"
)

func main() {
	logger := logger.CreateLogger(logger.LoggerSetting{
		Skip:     3,
		LogLevel: 4,
	})
	lib.LoadConfig(logger)
}
