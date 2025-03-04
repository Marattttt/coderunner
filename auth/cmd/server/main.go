package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/Marattttt/new_new_portfolio/auth/internal/db"
)

func main() {
	start := time.Now()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer cancel()

	conf, err := config.Config()
	checkFatal(err, "Creating config")

	err = db.Migrate(ctx, &conf.DB)
	checkFatal(err, "Applying migrations")

	slog.Info(
		"Migrations applied successfully",
		slog.Duration("took", time.Now().Sub(start)),
	)
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
