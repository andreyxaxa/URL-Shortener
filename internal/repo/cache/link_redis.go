package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	nredis "github.com/andreyxaxa/URL-Shortener/pkg/redis"
	"github.com/andreyxaxa/URL-Shortener/pkg/types/errs"
	"github.com/redis/go-redis/v9"
)

type LinkCache struct {
	c *nredis.Client
}

func New(c *nredis.Client) *LinkCache {
	return &LinkCache{c: c}
}

func (r *LinkCache) Get(ctx context.Context, key string) (string, error) {
	v, err := r.c.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("LinkCache - Get: %w", errs.ErrRecordNotFound)
		}
		return "", fmt.Errorf("LinkCache - Get - r.c.Client.Get: %w", err)
	}

	return v, nil
}

func (r *LinkCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.c.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("LinkCache - Set - r.c.Client.Set: %w", err)
	}

	return nil
}

func (r *LinkCache) Delete(ctx context.Context, key string) error {
	err := r.c.Client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("LinkCache - Delete - r.c.Client.Del: %w", err)
	}

	return nil
}

func (r *LinkCache) Increment(ctx context.Context, key string) (int64, error) {
	v, err := r.c.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("LinkCache - Increment - r.c.Client.Incr: %w", err)
	}

	return v, nil
}

func (r *LinkCache) IncrementWithExpiry(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	pipe := r.c.Client.Pipeline()

	v, err := pipe.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("LinkCache - IncrementWithExpiry - pipe.Incr: %w", err)
	}

	err = pipe.Expire(ctx, key, ttl).Err()
	if err != nil {
		return 0, fmt.Errorf("LinkCache - IncrementWithExpiry - pipe.Expire: %w", err)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("LinkCache - IncrementWithExpiry - pipe.Exec: %w", err)
	}

	return v, nil
}

func (r *LinkCache) GetInt(ctx context.Context, key string) (int64, error) {
	v, err := r.c.Client.Get(ctx, key).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, fmt.Errorf("LinkCache - GetInt: %w", errs.ErrRecordNotFound)
		}
		return 0, fmt.Errorf("LinkCache - GetInt - r.c.Client.Get: %w", err)
	}

	return v, nil
}
