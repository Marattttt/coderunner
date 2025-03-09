package main

import (
	"log/slog"
	"net/http"

	"github.com/Marattttt/new_new_portfolio/auth/internal/auth"
	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/labstack/echo/v4"
)

func applyRoutes(conf *config.AppConfig, e *echo.Echo) {
	stateProvider := auth.NewInMemoryStateProvider()

	e.GET("/login/google", handleLoginGoogle(conf, stateProvider))
	e.GET("/callback/google", handleCallbackGoogle(conf, stateProvider))
}

func handleLoginGoogle(conf *config.AppConfig, stateProvider auth.StateProvider) func(echo.Context) error {
	return func(c echo.Context) error {
		logger := c.Get("logger").(*slog.Logger)
		authService := auth.NewService(
			&conf.OAuth,
			stateProvider,
			logger,
		)

		url := authService.GenerateLoginURL(auth.GoogleResourceServer)

		logger.Info("Generated login URL", slog.String("resourceProvider", "google"))

		return c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func handleCallbackGoogle(conf *config.AppConfig, stateProvider auth.StateProvider) func(echo.Context) error {
	return func(c echo.Context) error {
		logger := c.Get("logger").(*slog.Logger)
		authService := auth.NewService(
			&conf.OAuth,
			stateProvider,
			logger,
		)

		state := c.FormValue("state")

		if err := stateProvider.InvalidateState(state); err != nil {
			logger.Warn(
				"Could not invalidate state",
				slog.String("error", err.Error()),
			)
			return c.JSON(http.StatusBadRequest, objectError("Invalid auth state"))
		}

		code := c.FormValue("code")
		if code == "" {
			logger.Warn("Code form value not found in form values")
			return c.JSON(http.StatusBadRequest, objectError("Code not in form values"))
		}

		u, err := authService.HandleCallback(
			c.Request().Context(),
			auth.GoogleResourceServer,
			code,
		)

		if err != nil {
			logger.Warn(
				"Error handling auth callback",
				slog.String("error", err.Error()),
			)
			return c.JSON(http.StatusBadGateway, objectError("Something went wrong"))
		}

		return c.JSON(http.StatusOK, u)
	}
}

func objectError(msg string) map[string]string {
	return map[string]string{"error": msg}
}
