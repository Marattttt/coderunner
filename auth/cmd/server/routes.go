package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Marattttt/new_new_portfolio/auth/internal/auth"
	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/Marattttt/new_new_portfolio/auth/internal/db"
	"github.com/Marattttt/new_new_portfolio/auth/pkg/models"
	"github.com/labstack/echo/v4"
)

type UsersProvider func(*slog.Logger) UserRepo

type UserRepo interface {
	GetUserEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

func UsersProviderFromDBConnn(conn *db.DBConn) UsersProvider {
	return func(logger *slog.Logger) UserRepo {
		return db.NewUsersRepository(conn, logger)
	}
}

func applyRoutes(conf *config.AppConfig, e *echo.Echo, userProvider UsersProvider) {
	stateProvider := auth.NewInMemoryStateProvider()

	e.GET("/login", handleLogin(conf))
	e.GET("/login/google", handleLoginGoogle(conf, stateProvider))
	e.GET("/callback/google", handleCallbackGoogle(conf, stateProvider, userProvider))
}

func handleLogin(conf *config.AppConfig) func(echo.Context) error {
	return func(c echo.Context) error {
		logger := c.Get("logger").(*slog.Logger)

		authH := c.Request().Header["Authorization"]
		if len(authH) == 0 || !strings.HasPrefix(authH[0], "Bearer ") {
			logger.Debug("Bearer token not found", slog.Any("authorizationHeader", authH))
			return c.JSON(
				http.StatusUnauthorized,
				objectError("Authentication required"),
			)
		}

		bearerToken := strings.TrimPrefix(authH[0], "Bearer ")

		logger.Info("Found bearer token", slog.String("jwt", bearerToken))

		return c.JSON(
			http.StatusNotImplemented,
			objectError("Login is not implemented"),
		)
	}
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

func handleCallbackGoogle(
	conf *config.AppConfig,
	stateProvider auth.StateProvider,
	users UsersProvider,
) func(echo.Context) error {
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

		users := users(logger)
		// TODO: use app context instead of todo
		go createUserIfNotExists(context.TODO(), logger, users, u)

		return c.JSON(http.StatusOK, u)
	}
}

func objectError(msg string) map[string]string {
	return map[string]string{"error": msg}
}

func createUserIfNotExists(ctx context.Context, logger *slog.Logger, provider UserRepo, u *models.User) {
	var err error
	defer func() {
		if err != nil {
			logger.Warn(
				"Error creating user if it does not exist",
				slog.String("email",
					u.Email),
				slog.String("error", err.Error()),
			)
		}
	}()

	_, err = provider.GetUserEmail(ctx, u.Email)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return
	}

	if errors.Is(err, db.ErrNotFound) {
		logger.Info(
			"User with provided email does not exist, creating a new one",
			slog.String("email", u.Email),
		)

		err = provider.CreateUser(ctx, u)
		return
	}
}
