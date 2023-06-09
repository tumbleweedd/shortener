package services

import (
	"context"
	"github.com/tumbleweedd/shortener/internal/repositories"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Shortener interface {
	ShortenURL(ctx context.Context, longURL string) (string, error)
	GetLongURL(ctx context.Context, code string) (string, error)
}

type Service struct {
	Shortener
}

func NewService(r *repositories.Repository) *Service {
	return &Service{
		Shortener: NewShortenerService(r),
	}
}
