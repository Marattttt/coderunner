package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB                 DBConfig
	OAuth              OAuthConfig
	Port               int    `env:"PORT" envDefault:"8080"`
	Mode               string `env:"MODE" envDefault:"development"`
	JWTSecret          string `env:"JWT_SECRET,unset"`
	JWTAccessExprirySecs int    `env:"JWT_ACCESS_EXPIRY" envDefault:"1440"`
}

func (a *AppConfig) GetListenAddr() string {
	return fmt.Sprintf(":%d", a.Port)
}

func (a *AppConfig) MakeLogger() *slog.Logger {
	var handler slog.Handler

	if a.Mode == "development" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	}

	return slog.New(handler)
}

type OAuthConfig struct {
	Google GoogleAuthConfig
}

type GoogleAuthConfig struct {
	ClientID     string   `env:"G_CLIENT_ID,unset"`
	Redirect     string   `env:"G_REDIRECT_URL"`
	ClientSecret string   `env:"G_CLIENT_SECRET,unset"`
	Scopes       []string `env:"G_OAUTH_SCOPES"`
}

type DBConfig struct {
	PostgresURI   string `env:"PG_URI"`
	RedisTokenURI string `env:"REDIS_TOKEN_URI"`
	MigrationsURI string `env:"MIGRATIONS_URI"`
}

func Config() (*AppConfig, error) {
	_ = godotenv.Load()

	var conf AppConfig

	err := env.Parse(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
