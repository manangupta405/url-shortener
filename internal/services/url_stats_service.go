package services

import (
	"context"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/repositories"
)

//go:generate mockery --name=URLStatsService --output=./mocks
type URLStatsService interface {
	GetURLStatistics(ctx context.Context, shortPath string) (*models.URLStatistics, error)
	InsertAccessLog(ctx context.Context, shortPath string, accessedAt time.Time) error
}

type urlStatsServiceImpl struct {
	repo repositories.URLStatisticsRepository
}

func NewURLStatsService(repo repositories.URLStatisticsRepository) URLStatsService {
	return &urlStatsServiceImpl{repo: repo}
}

func (s *urlStatsServiceImpl) GetURLStatistics(ctx context.Context, shortPath string) (*models.URLStatistics, error) {
	urlStats, err := s.repo.GetURLStatistics(ctx, shortPath)
	if err != nil {
		return nil, err
	}
	return urlStats, nil
}

func (s *urlStatsServiceImpl) InsertAccessLog(ctx context.Context, shortPath string, accessedAt time.Time) error {
	err := s.repo.InsertAccessLog(ctx, shortPath, accessedAt)
	if err != nil {
		return err
	}
	return nil
}
