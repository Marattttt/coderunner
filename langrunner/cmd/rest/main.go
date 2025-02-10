package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Marattttt/coderunner/langrunner/internal/config"
	"github.com/Marattttt/coderunner/langrunner/internal/models"
	"github.com/Marattttt/coderunner/langrunner/internal/runner"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.TODO()

	conf, err := config.Config(ctx)
	checkFail(err, "Creating config")

	logger := conf.CreateLogger()

	e := echo.New()
	e.Use(
		handlerParentContextMiddeware(ctx),
		middleware.RequestID(),
		requestLogMiddleware(logger),
	)

	runtimeManager := runner.NewRuntimeManager(&conf.RunnerConig, logger)
	e.POST("/run", runHandler(conf, logger, runtimeManager))

	checkFail(e.Start(":8080"), "Running server")
}

func runHandler(conf *config.AppConfig, logger *slog.Logger, man *runner.RuntimeManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		gorunner := runner.NewGoRunner(logger, *conf.Go, man)

		var req models.CodeReq

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				map[string]string{"error": "bad request"},
			)
		}

		res, err := gorunner.RunCode(ctx, []byte(req.Code))
		if err != nil {
			logger.Error("Failed to run code", slog.String("code", req.Code), slog.String("err", err.Error()))
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"error": "internal"},
			)
		}

		return c.JSON(
			http.StatusOK,
			res,
		)
	}
}

func handlerParentContextMiddeware(ctx context.Context) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func requestLogMiddleware(logger *slog.Logger) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			reqId := c.Request().Header.Get(echo.HeaderXRequestID)

			// Log request details
			logger.Info("request",
				slog.String("id", reqId),
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("remote_ip", c.RealIP()),
			)

			start := time.Now()

			err := next(c)

			end := time.Now()

			// Log response details
			if err != nil {
				logger.Error("request error",
					slog.String("id", reqId),
					slog.String("error", err.Error()),
					slog.Int64("timeTookMs", end.Sub(start).Milliseconds()),
					slog.Int("status", res.Status),
				)
			} else {
				logger.Info("response",
					slog.String("id", reqId),
					slog.Int("status", res.Status),
					slog.Int64("timeTookMs", end.Sub(start).Milliseconds()),
				)
			}

			return err
		}
	}
}

func checkFail(err error, msg string) {
	if err == nil {
		return
	}

	slog.Error(msg, slog.String("err", err.Error()))
	os.Exit(1)
}
