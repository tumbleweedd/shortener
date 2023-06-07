package services

import (
	"context"
	"github.com/tumbleweedd/shortener/internal/repositories"
)

type ShortenerService struct {
	repo repositories.Shortener
}

func NewShortenerService(repo repositories.Shortener) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, longURL string) (string, error) {
	shortURL, err := s.repo.SaveURL(ctx, longURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *ShortenerService) GetLongURL(ctx context.Context, code string) (string, error) {
	longURL, err := s.repo.GetURL(ctx, code)
	if err != nil {
		return "", err
	}

	return longURL, err
}
