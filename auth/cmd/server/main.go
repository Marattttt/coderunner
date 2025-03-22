package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/Marattttt/new_new_portfolio/auth/internal/db"
	"github.com/labstack/echo/v4"
)

func main() {
	ctx, cancel, conf := setup()
	defer cancel()

	logger := conf.MakeLogger()

	logger.Info("Created config", slog.Any("config", *conf))

	err := db.Migrate(ctx, &conf.DB)
	checkFatal(err, "Applying migrations")

	dbconn, err := db.Connect(&conf.DB)
	checkFatal(err, "Connecting to db")

	e := echo.New()
	applyMiddleware(ctx, conf, logger, e)
	applyRoutes(conf, e, UsersProviderFromDBConnn(dbconn))

	go func() {
		e.Start(conf.GetListenAddr())
	}()

	<-ctx.Done()

	downCtx, downCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer downCancel()

	logger.Info("Began shutdown")

	e.Shutdown(downCtx)
}

func setup() (context.Context, func(), *config.AppConfig) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)

	conf, err := config.Config()
	checkFatal(err, "Creating config")

	return ctx, cancel, conf
}

func checkFatal(err error, msg string, extra ...any) {
	if err == nil {
		return
	}

	slog.Error(
		msg,
		slog.String("error", err.Error()),
		slog.Any("extra", extra),
	)

	os.Exit(1)
}
