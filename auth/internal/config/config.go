package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB    DBConfig
	OAuth OAuthConfig
	Port  int `env:"PORT" envDefault:"8080"`
}

func (a AppConfig) GetListenAddr() string {
	return fmt.Sprintf(":%d", a.Port)
}

type OAuthConfig struct {
	Google GoogleAuthConfig
}

type GoogleAuthConfig struct {
	ClientID     string   `env:"G_CLIENT_ID"`
	Redirect     string   `env:"G_REDIRECT_URL"`
	ClientSecret string   `env:"G_CLIENT_SECRET"`
	Scopes       []string `env:"G_OAUTH_SCOPES"`
}

type DBConfig struct {
	PostgresURI   string `env:"PG_URI"`
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
