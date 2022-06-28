package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
	"ui-project/api"
	"ui-project/auth"
	"ui-project/database"
	"ui-project/lib"
	"ui-project/logger"
	v1 "ui-project/v1"

	"github.com/gorilla/mux"
)

const (
	MAXIMUM_FINISH_TIME time.Duration = 10
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := logger.CreateLogger(logger.LoggerSetting{
		Skip:     3,
		LogLevel: 4,
	})
	lib.LoadConfig(ctx, logger)

	r := mux.NewRouter()

	database := database.NewDatabase(ctx, logger)
	database.Bootstrap(ctx)

	auth := auth.NewAuthUsecase()
	rs := v1.Register(ctx, logger, auth, r, database.GetClient())
	app := api.NewApiUsecase(ctx, logger, auth, r)
	app.RegisterRoute(ctx, rs)

	app.Bootstrap(ctx)

	// Graceful Shutdown
	<-ctx.Done()
	logger.LogInfo(ctx, "Graceful Shutdown start")

	c, cancel := context.WithTimeout(context.Background(), MAXIMUM_FINISH_TIME*time.Second)
	defer cancel()

	app.Shutdown(c)
	database.Shutdown(c)

	// <-c.Done()
	logger.LogInfo(ctx, "Graceful Shutdown finished")
}
