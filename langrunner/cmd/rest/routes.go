package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Marattttt/coderunner/langrunner/internal/config"
	"github.com/Marattttt/coderunner/langrunner/internal/models"
	"github.com/Marattttt/coderunner/langrunner/internal/runner"
	"github.com/labstack/echo/v4"
)

func applyRoutes(e *echo.Echo, conf *config.AppConfig, logger *slog.Logger, runtimeManager *runner.RuntimeManager) {
	slog.Info("Creating routes", slog.Any("langs", conf.Enabled))

	for _, lang := range conf.Enabled {
		switch lang {
		case "go":
			e.POST("/api/run/go", runGoHandler(conf, logger, runtimeManager))
		case "py":
			e.POST("/api/run/py", runPyHandler(conf, logger, runtimeManager))
		default:
			panic(fmt.Sprintf("Cannot create route for language %s", lang))
		}
	}
}

func runPyHandler(conf *config.AppConfig, logger *slog.Logger, man *runner.RuntimeManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		pyrunner := runner.NewPyRunner(logger, *conf.Python, man)

		var req models.CodeReq

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				map[string]string{"error": "bad request"},
			)
		}

		res, err := pyrunner.RunCode(ctx, []byte(req.Code))
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

func runGoHandler(conf *config.AppConfig, logger *slog.Logger, man *runner.RuntimeManager) echo.HandlerFunc {
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
