package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func applyMiddleware(
	ctx context.Context,
	conf *config.AppConfig,
	logger *slog.Logger,
	e *echo.Echo,
) {
	e.Use(
		handlerParentContextMiddeware(ctx),
		slogProviderMiddleware(logger),
		middleware.RequestID(),
		requestLogMiddleware(),
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAccept,
				echo.HeaderContentType,
			},
			AllowMethods:     []string{http.MethodGet, http.MethodPost},
			AllowCredentials: true,
		}),
	)
}

func handlerParentContextMiddeware(ctx context.Context) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func slogProviderMiddleware(logger *slog.Logger) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", logger)
			return next(c)
		}
	}
}

func requestLogMiddleware() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			reqId := c.Request().Header.Get(echo.HeaderXRequestID)

			logger := c.Get("logger").(*slog.Logger)

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
