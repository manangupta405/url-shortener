package repositories

import (
	"context"
	"time"
	"url-shortener/internal/models"
)

//go:generate mockery --name=URLStatisticsRepository --output=./mocks
type URLStatisticsRepository interface {
	GetURLStatistics(ctx context.Context, shortPath string) (*models.URLStatistics, error)
	InsertAccessLog(ctx context.Context, shortPath string, accessedAt time.Time) error
}
