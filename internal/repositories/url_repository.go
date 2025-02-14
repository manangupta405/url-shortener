package repositories

import (
	"context"
	"time"
	"url-shortener/internal/models"
)

//go:generate mockery --name=URLRepository --output=./mocks
type URLRepository interface {
	GetShortURL(ctx context.Context, originalURL string) (*models.URL, error)
	GetOriginalURL(ctx context.Context, shortPath string) (*models.URL, error)
	UpdateShortURL(ctx context.Context, url *models.URL) error
	DeleteShortURL(ctx context.Context, shortPath string, currentTime time.Time, deletedBy string) error
	InsertShortURL(ctx context.Context, url *models.URL) error
}
