package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisConnn struct{ client *redis.Client }

func ConnectRedis(conf *config.DBConfig) (*RedisConnn, error) {
	opts, err := redis.ParseURL(conf.RedisURI)
	if err != nil {
		return nil, fmt.Errorf("parsing url: %w", err)
	}
	return &RedisConnn{client: redis.NewClient(opts)}, nil
}

const (
	refreshPrefix = "refresh"
	accessPrefix  = "access"
)

type RedisTokenStore struct {
	conf   *config.AppConfig
	client *redis.Client
	logger *slog.Logger
}

func NewRedisTokenStore(conf *config.AppConfig, conn *RedisConnn, logger *slog.Logger) RedisTokenStore {
	return RedisTokenStore{
		conf:   conf,
		client: conn.client,
		logger: logger,
	}
}

func (r RedisTokenStore) SaveRefresh(ctx context.Context, tok string) error {
	return r.saveTok(ctx, tok, refreshPrefix, time.Duration(r.conf.JWTRefreshExprirySecs))
}

func (r RedisTokenStore) SaveAccess(ctx context.Context, tok string) error {
	return r.saveTok(ctx, tok, accessPrefix,  time.Duration(r.conf.JWTAccessExprirySecs))
}

// If the token exists, it gets deleted and true is returned
func (r RedisTokenStore) ConsumeRefresh(ctx context.Context, tok string) (bool, error) {
	return r.consumeIfExists(ctx, tok, refreshPrefix)
}

// If the token exists, it gets deleted and true is returned
func (r RedisTokenStore) ConsumeAccess(ctx context.Context, tok string) (bool, error) {
	return r.consumeIfExists(ctx, tok, accessPrefix)
}

// Uses a default timeout of 20 seconds
func (r RedisTokenStore) consumeIfExists(ctx context.Context, tok string, tokType string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	s := tokType + ":" + tok

	res := r.client.GetDel(ctx, s)

	if res.Err() != nil {
		return false, res.Err()
	}

	return res.Val() != "", nil
}

func (r RedisTokenStore) saveTok(ctx context.Context, tok string, tokType string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	s := tokType + ":" + tok

	res := r.client.Set(ctx, s, 1, ttl)

	return res.Err()
}
