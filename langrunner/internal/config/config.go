package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type AppConfig struct {
	RunnerConig
	Port      int    `env:"PORT,default=8080"`
}

type RunnerConig struct {
	Go            *GoRunnerConfig `env:", prefix=GO_, noinit"`
	Python        *PyRunnerConfig `env:",prefix=PY_, noinit"`
	Enabled       []string        `env:"ENABLED_LANGS, required"`
	EnvVars       []string        `env:"RUNTIME_ENVS,required"`
	SafeUserNames []string        `env:"SAFE_USERNAMES"`
}

type GoRunnerConfig struct {
	RuntimeDir string `env:"RUNTIME_DIR"`
}

type PyRunnerConfig struct {
	RuntimeDir string `env:"RUNTIME_DIR"`
}

func Config(ctx context.Context) (*AppConfig, error) {
	godotenv.Load()

	var conf AppConfig

	if err := envconfig.Process(ctx, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
func (a *AppConfig) CreateLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Enable debug level logging
	})

	return slog.New(handler)
}
