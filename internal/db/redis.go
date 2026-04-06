package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(url string) *redis.Client {
	// Parse URL and create client
	opt, err := redis.ParseURL(url)
	if err != nil {
		// Fallback to default localhost
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	return redis.NewClient(opt)
}

// RedisWrapper provides convenience methods for Redis operations
type RedisWrapper struct {
	client *redis.Client
}

// NewRedisWrapper creates a new Redis wrapper
func NewRedisWrapper(client *redis.Client) *RedisWrapper {
	return &RedisWrapper{client: client}
}

// Set stores a key-value pair
func (r *RedisWrapper) Set(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

// Get retrieves a value by key
func (r *RedisWrapper) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete removes a key
func (r *RedisWrapper) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
