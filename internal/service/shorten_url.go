package service

import (
	"context"
	"time"

	"github.com/SonDo580/bookmark-manager/internal/repository"
)

type ShortenURL interface {
	CreateShortenedLink(ctx context.Context, url string, exp time.Duration) (string, error)
}

type shortenURL struct {
	r       repository.URLStorage
	genCode GenCode
}

func NewShortenURL(r repository.URLStorage, genCode GenCode) ShortenURL {
	return &shortenURL{r: r, genCode: genCode}
}

func (s *shortenURL) CreateShortenedLink(ctx context.Context, url string, exp time.Duration) (string, error) {
	code, err := s.genCode.GenerateCode(DefaultCodeLength)
	if err != nil {
		return "", err
	}

	err = s.r.StoreURL(ctx, code, url, exp)
	if err != nil {
		return "", err
	}

	return code, nil
}
