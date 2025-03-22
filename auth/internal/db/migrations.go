package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
)

func Migrate(ctx context.Context, conf *config.DBConfig) error {
	defer func(from time.Time) {
		slog.Info("Finished migrations", slog.Duration("took", time.Now().Sub(from)))
	}(time.Now())

	m, err := migrate.New(conf.MigrationsURI, conf.PostgresURI)
	if err != nil {
		return fmt.Errorf("initializing: %w", err)
	}
	defer m.Close()

	// Use a separate goroutine to run migration
	errCh := make(chan error, 1)
	go func() {
		errCh <- m.Up() // Run all up migrations
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context done: %w", ctx.Err())
	case err := <-errCh:
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
