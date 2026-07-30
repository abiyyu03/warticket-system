package pkg

import (
	"context"
	"fmt"
	"time"

	"go-projects/hexagonal-example/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
	Config config.RedisConfig
}

func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// test ping
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

func NewRedis(cfg config.RedisConfig) (*Redis, error) {
	client, err := InitRedis(cfg)
	if err != nil {
		return nil, err
	}
	return &Redis{
		Client: client,
		Config: cfg,
	}, nil
}
