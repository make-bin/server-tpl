package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/make-bin/server-tpl/pkg/utils/config"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// RedisClient wraps the redis client with additional functionality
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client instance
func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Info("Connected to Redis at %s", addr)

	return &RedisClient{client: rdb}, nil
}

// Set stores a key-value pair with optional expiration
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete removes a key
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire sets an expiration time for a key
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// HealthCheck checks the Redis connection
func (r *RedisClient) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
