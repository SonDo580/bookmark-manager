package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLStorage interface {
	StoreURL(ctx context.Context, key string, url string, exp time.Duration) error
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}

func (r *urlStorage) StoreURL(ctx context.Context, key string, url string, exp time.Duration) error {
	// TODO: handle conflict
	return r.c.Set(ctx, key, url, exp).Err()
}

func (r *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	// TODO: does it return error when key expires
	return r.c.Get(ctx, key).Result()
}
