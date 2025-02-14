package repositories

import (
	"context"
	"time"
	"url-shortener/internal/models"
)

type URLRepository interface {
	GetShortURL(ctx context.Context, originalURL string) (*models.URL, error)
	GetOriginalURL(ctx context.Context, shortPath string) (*models.URL, error)
	UpdateShortURL(ctx context.Context, url *models.URL, currentTime time.Time) error
	DeleteShortURL(ctx context.Context, shortPath string, currentTime time.Time) error
	InsertShortURL(ctx context.Context, url *models.URL, currentTime time.Time) error
}
