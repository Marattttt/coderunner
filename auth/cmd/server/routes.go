package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Marattttt/new_new_portfolio/auth/internal/auth"
	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/Marattttt/new_new_portfolio/auth/internal/db"
	"github.com/Marattttt/new_new_portfolio/auth/pkg/models"
	"github.com/labstack/echo/v4"
)

func applyRoutes(
	conf *config.AppConfig,
	e *echo.Echo,
	usersProvider UsersProvider,
	tokensProvider TokensProvider,
	codesProvider CodesProvider,
) {
	stateProvider := auth.NewInMemoryStateProvider()

	e.GET("/login", handleLogin(conf, tokensProvider, codesProvider, usersProvider))
	e.GET("/login/google", handleLoginGoogle(conf, stateProvider, codesProvider))
	e.GET("/callback/google", handleCallbackGoogle(conf, stateProvider, usersProvider))
}

func handleLogin(
	conf *config.AppConfig,
	tokensProvider TokensProvider,
	codesProvider CodesProvider,
	usersProvider UsersProvider,
) func(echo.Context) error {
	return func(c echo.Context) error {
		logger := c.Get("logger").(*slog.Logger)

		code := c.QueryParams().Get("code")
		logger.Debug("Got code param from request", slog.String("code", code))
		if code == "" {
			return c.JSON(
				http.StatusBadRequest,
				objectError("Code query param is required"),
			)
		}

		logger = logger.With(slog.String("code", code))

		userIdStr := c.QueryParams().Get("userId")
		logger.Debug("Got userId param from request", slog.String("userid", userIdStr))
		if userIdStr == "" {
			return c.JSON(
				http.StatusBadRequest,
				objectError("userId query param is required"),
			)
		}

		expectedUserId, err := strconv.Atoi(userIdStr)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				objectError("userid is not an integer"),
			)
		}

		codesRepo := codesProvider(logger)

		userId, err := codesRepo.Consume(c.Request().Context(), code)
		if err != nil && err != db.ErrNotFound {
			logger.Error("Error consuming userid", slog.String("err", err.Error()))
			return c.JSON(
				http.StatusBadGateway,
				objectError("Something went wrong"),
			)
		}

		if err == db.ErrNotFound || expectedUserId != *userId {
			logger.Warn("User id not found by code")
			return c.JSON(
				http.StatusBadRequest,
				objectError("User id not associated with code"),
			)
		}

		usersRepo := usersProvider(logger)

		user, err := usersRepo.GetUser(c.Request().Context(), *userId)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			logger.Error("Could not get user from db", slog.String("error", err.Error()))
			return c.JSON(
				http.StatusBadGateway,
				objectError("Something went wrong"),
			)
		}

		if errors.Is(err, db.ErrNotFound) {
			logger.Warn("User not found by id", slog.Int("userId", *userId))
			return c.JSON(
				http.StatusBadRequest,
				objectError("User not found by id"),
			)
		}

		accessTok, err := auth.JWTGenerateAccess(conf, user)
		if err != nil {
			logger.Error("Could not generate jwt", slog.String("error", err.Error()))
			return c.JSON(
				http.StatusBadGateway,
				objectError("Something went wrong"),
			)
		}

		tokRepo := tokensProvider(logger)

		err = tokRepo.SaveAccess(c.Request().Context(), accessTok)
		if err != nil {
			logger.Error("Could not save access token", slog.String("error", err.Error()))
			return c.JSON(
				http.StatusBadGateway,
				objectError("Something went wrong"),
			)
		}

		refreshTok, err := auth.JWTGenerateRefresh(conf, user)
		if err != nil {
			logger.Error("Could not generate jwt", slog.String("error", err.Error()))
			return c.JSON(
				http.StatusBadGateway,
				objectError("Something went wrong"),
			)
		}

		err = tokRepo.SaveRefresh(c.Request().Context(), refreshTok)
		if err != nil {
			logger.Error("Could not save refresh token", slog.String("error", err.Error()))
			return c.JSON(
				http.StatusBadGateway,
				objectError("Something went wrong"),
			)
		}

		return c.JSON(
			http.StatusOK,
			models.AuthResponce{
				AccessTok:  accessTok,
				RefreshTok: refreshTok,
			},
		)
	}
}

//	func handleLogin(conf *config.AppConfig) func(echo.Context) error {
//		return func(c echo.Context) error {
//			logger := c.Get("logger").(*slog.Logger)
//
//			authH := c.Request().Header["Authorization"]
//			if len(authH) == 0 || !strings.HasPrefix(authH[0], "Bearer ") {
//				logger.Debug("Bearer token not found", slog.Any("authorizationHeader", authH))
//				return c.JSON(
//					http.StatusUnauthorized,
//					objectError("Authentication required"),
//				)
//			}
//
//			bearerToken := strings.TrimPrefix(authH[0], "Bearer ")
//
//			logger.Info("Found bearer token", slog.String("jwt", bearerToken))
//
//			return c.JSON(
//				http.StatusNotImplemented,
//				objectError("Login is not implemented"),
//			)
//		}
//	}
func handleLoginGoogle(
	conf *config.AppConfig,
	stateProvider auth.StateProvider,
	codesProvider CodesProvider,
) func(echo.Context) error {
	return func(c echo.Context) error {
		logger := c.Get("logger").(*slog.Logger)

		userIdStr := c.QueryParams().Get("userId")
		userId, err := strconv.Atoi(userIdStr)

		logger.Info("Login attempt with google", slog.String("userId", userIdStr))

		if userIdStr == "" || err != nil {
			logger.Warn("Could not parse user id")

			return c.JSON(
				http.StatusBadRequest,
				objectError("userId query param not present or invalid"),
			)

		}

		authService := auth.NewService(
			&conf.OAuth,
			stateProvider,
			logger,
		)

		codesRepo := codesProvider(logger)

		url := authService.GenerateLoginURL(auth.GoogleResourceServer)

		code := auth.GenerateAccessCode()

		go codesRepo.Save(
			context.TODO(),
			code,
			userId,
			time.Duration(int64(conf.OAuth.AccessCodeExpirySeconds))*1e9,
		)

		logger.Info("Generated login URL", slog.String("resourceProvider", "google"))

		return c.JSON(http.StatusOK, models.OAuthRedirect{
			LoginURL:        url,
			TokenAccessCode: code,
		})
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

		return c.JSON(http.StatusOK, map[string]int{"id": u.ID})
	}
}

func objectError(msg string) map[string]string {
	return map[string]string{"error": msg}
}

func createUserIfNotExists(ctx context.Context, logger *slog.Logger, provider UsersRepo, u *models.User) {
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

type UsersProvider func(*slog.Logger) UsersRepo

type UsersRepo interface {
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUserEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

func UsersProviderFromDBConnn(conn *db.DBConn) UsersProvider {
	return func(logger *slog.Logger) UsersRepo {
		return db.NewUsersRepository(conn, logger)
	}
}

type TokensProvider func(*slog.Logger) TokensRepo
type TokensRepo interface {
	SaveRefresh(ctx context.Context, tok string) error
	SaveAccess(ctx context.Context, tok string) error
	ConsumeRefresh(ctx context.Context, tok string) (bool, error)
	ConsumeAccess(ctx context.Context, tok string) (bool, error)
}

func TokensProviderFromRedisConn(conf *config.AppConfig, conn *db.RedisConnn) TokensProvider {
	return func(logger *slog.Logger) TokensRepo {
		return db.NewRedisTokenStore(conf, conn, logger)
	}
}

type CodesProvider func(*slog.Logger) CodesRepo
type CodesRepo interface {
	Consume(ctx context.Context, code string) (*int, error)
	Save(ctx context.Context, code string, userId int, ttl time.Duration) error
}

func CodesProviderFromRedisConn(conn *db.RedisConnn) CodesProvider {
	return func(logger *slog.Logger) CodesRepo {
		return db.NewRedisCodeStore(conn, logger)
	}
}
