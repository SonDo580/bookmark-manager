package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=URLStorage --filename url_storage.go --outpkg repository_mocks
type URLStorage interface {
	StoreURL(ctx context.Context, key string, url string, exp time.Duration) (bool, error)
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}

func (r *urlStorage) StoreURL(ctx context.Context, key string, url string, exp time.Duration) (bool, error) {
	return r.c.SetNX(ctx, key, url, exp).Result()
}

func (r *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	return r.c.Get(ctx, key).Result()
}
