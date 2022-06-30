package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"ui-project/api"
	"ui-project/auth"
	"ui-project/database"
	"ui-project/lib"
	"ui-project/logger"
	v1 "ui-project/v1"
	ws "ui-project/websocket"

	_ "ui-project/docs"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	MAXIMUM_FINISH_TIME time.Duration = 10
)

// @title           UI Project API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:8888
// @BasePath  /v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type \"Bearer\" followed by a space and JWT token.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := logger.CreateLogger(logger.LoggerSetting{
		Skip:     3,
		LogLevel: 4,
	})
	lib.LoadConfig(ctx, logger)

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", viper.GetString("host"), viper.GetInt("port"))),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	database := database.NewDatabase(ctx, logger)
	database.Bootstrap(ctx)

	auth := auth.NewAuthUsecase()
	server := ws.NewServer()
	rs := v1.Register(ctx, logger, auth, r, database.GetClient(), server)
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
