package service

import (
	"context"
	"errors"
	"time"

	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/repository"
)

//go:generate mockery --name=ShortenURL --filename shorten_url.go --outpkg service_mocks
type ShortenURL interface {
	CreateShortenedLink(ctx context.Context, url string, exp time.Duration) (string, error)
	GetURLFromCode(ctx context.Context, code string) (string, error)
}

type shortenURL struct {
	r       repository.URLStorage
	genCode GenCode
}

func NewShortenURL(r repository.URLStorage, genCode GenCode) ShortenURL {
	return &shortenURL{r: r, genCode: genCode}
}

const MaxRetries = 5 // arbitrary

var ErrMaxRetriesExceeded = errors.New("failed to generate code after all attempts")

func (s *shortenURL) CreateShortenedLink(ctx context.Context, url string, exp time.Duration) (string, error) {
	for range MaxRetries {
		code, err := s.genCode.GenerateCode(DefaultCodeLength)
		if err != nil {
			return "", err
		}

		success, err := s.r.StoreURL(ctx, code, url, exp)
		if err != nil {
			return "", err
		}

		if success {
			return code, nil
		}
	}

	return "", ErrMaxRetriesExceeded
}

var ErrURLNotFound = errors.New(dto.MsgURLNotFound)

func (s *shortenURL) GetURLFromCode(ctx context.Context, code string) (string, error) {
	url, err := s.r.GetURL(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			return "", ErrURLNotFound
		}
		return "", err
	}

	return url, nil
}
