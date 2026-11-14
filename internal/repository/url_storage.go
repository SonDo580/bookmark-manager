package repository

import (
	"context"
	"errors"
	"time"

	"github.com/SonDo580/bookmark-manager/internal/dto"
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

var ErrURLNotFound = errors.New(dto.MsgURLNotFound)

func (r *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	url, err := r.c.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrURLNotFound
		}
		return "", err
	}

	return url, nil
}
