package db

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	codePrefix = "code"
)

type RedisCodeStore struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisCodeStore(conn *RedisConnn, logger *slog.Logger) RedisCodeStore {
	return RedisCodeStore{
		client: conn.client,
		logger: logger,
	}
}

// Uses a default timeout of 20 seconds
func (r RedisCodeStore) Consume(ctx context.Context, code string) (*int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	s := codePrefix + ":" + code

	val, err := r.client.GetDel(ctx, s).Result()

	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	if err == redis.Nil  {
		return nil, ErrNotFound
	}

	parsedUserId, err := strconv.Atoi(val)
	if err != nil {
		r.logger.Error("Could not parse value as integer", slog.String("value", val))
		return nil, fmt.Errorf("parsing userid: %w", err)
	}

	return &parsedUserId, nil
}

func (r RedisCodeStore) Save(ctx context.Context, code string, userId int, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	s := codePrefix + ":" + code

	res := r.client.Set(ctx, s, userId, ttl)

	return res.Err()
}
