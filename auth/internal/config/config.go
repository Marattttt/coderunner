package config

import "os"

const (
	ENV_DB_POSTGRESURI = "PG_URI"
)

type AppConfig struct {
	DB DBConfig
}

type DBConfig struct {
	PostgresURI string
}

func Config() (*AppConfig, error) {
	var conf AppConfig

	pguri := os.Getenv(ENV_DB_POSTGRESURI)
	if pguri == "" {
		return nil, EnvNotSetOrEmptyError{ENV_DB_POSTGRESURI}
	}

	return &conf, nil
}

type EnvNotSetOrEmptyError struct {
	EnvName string
}

func (e EnvNotSetOrEmptyError) Error() string {
	return "Env variable not set or empty: " + e.EnvName
}
