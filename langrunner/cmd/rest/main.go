package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Marattttt/coderunner/langrunner/internal/config"
	"github.com/Marattttt/coderunner/langrunner/internal/runner"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.TODO()

	conf, err := config.Config(ctx)
	checkFail(err, "Creating config")

	logger := conf.CreateLogger()
	runtimeManager := runner.NewRuntimeManager(&conf.RunnerConig, logger)

	e := echo.New()
	e.Use(
		handlerParentContextMiddeware(ctx),
		middleware.RequestID(),
		requestLogMiddleware(logger),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentType},
			AllowMethods:     []string{http.MethodGet, http.MethodPost},
			AllowCredentials: true,
		}),
	)

	applyRoutes(e, conf, logger, runtimeManager)

	checkFail(
		e.Start(fmt.Sprintf(":%d", conf.Port)),
		"Running server",
	)
}

func checkFail(err error, msg string) {
	if err == nil {
		return
	}

	slog.Error(msg, slog.String("err", err.Error()))
	os.Exit(1)
}
