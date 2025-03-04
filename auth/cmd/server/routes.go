package main

import (
	"log/slog"
	"net/http"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/labstack/echo/v4"
)

func applyRoutes(conf *config.AppConfig, e *echo.Echo) {
	e.GET("/auth/login/google", handleLoginGoogle(conf))
	e.GET("/auth/callback/google", handleCallbackGoogle(conf))
}

func handleLoginGoogle(conf *config.AppConfig) func(echo.Context) error {
	return func(c echo.Context) error {
		msg := "handle login"

		logger := c.Get("logger").(*slog.Logger)
		logger.Info(msg)

		return c.String(http.StatusOK, msg)
	}
}

func handleCallbackGoogle(conf *config.AppConfig) func(echo.Context) error {
	return func(c echo.Context) error {
		msg := "handle callback"

		logger := c.Get("logger").(*slog.Logger)
		logger.Info(msg)

		return c.String(http.StatusOK, msg)
	}
}
