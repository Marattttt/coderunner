package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	ENV_DB_POSTGRESURI   = "PG_URI"
	ENV_DB_MIGRATIONSURI = "MIGRATIONS_URI"
)

type AppConfig struct {
	DB DBConfig
}

type DBConfig struct {
	PostgresURI   string
	MigrationsURI string
}

func Config() (*AppConfig, error) {
	_ = godotenv.Load()

	var conf AppConfig

	pguri := os.Getenv(ENV_DB_POSTGRESURI)
	if pguri == "" {
		return nil, EnvNotSetOrEmptyError{ENV_DB_POSTGRESURI}
	}
	conf.DB.PostgresURI = pguri

	migrations := os.Getenv(ENV_DB_MIGRATIONSURI)
	if migrations == "" {
		return nil, EnvNotSetOrEmptyError{ENV_DB_MIGRATIONSURI}
	}
	conf.DB.MigrationsURI = migrations

	return &conf, nil
}

type EnvNotSetOrEmptyError struct {
	EnvName string
}

func (e EnvNotSetOrEmptyError) Error() string {
	return "Env variable not set or empty: " + e.EnvName
}
