package repositories

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Shortener interface {
	SaveURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, code string) (string, error)
}

type Repository struct {
	Shortener
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		Shortener: NewShortenerRepository(client),
	}
}
