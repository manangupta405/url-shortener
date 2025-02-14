package services

import (
	"context"
	"testing"
	"time"

	"url-shortener/internal/models"
	repoMocks "url-shortener/internal/repositories/mocks"

	"github.com/stretchr/testify/assert"
)

func TestURLStatsServiceImpl_GetURLStatistics(t *testing.T) {
	repo := &repoMocks.URLStatisticsRepository{}
	defer repo.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	mockStats := &models.URLStatistics{ShortPath: shortPath, Last24Hours: 5, PastWeek: 5, AllTime: 5}
	repo.On("GetURLStatistics", ctx, shortPath).Return(mockStats, nil).Once()
	service := NewURLStatsService(repo)
	urlStats, err := service.GetURLStatistics(ctx, shortPath)
	assert.Nil(t, err)
	assert.Equal(t, mockStats, urlStats)
	repo.AssertExpectations(t)
}

func TestURLStatsServiceImpl_GetURLStatistics_RepoError(t *testing.T) {
	repo := &repoMocks.URLStatisticsRepository{}
	defer repo.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	repo.On("GetURLStatistics", ctx, shortPath).Return(nil, assert.AnError).Once()
	service := NewURLStatsService(repo)
	_, err := service.GetURLStatistics(ctx, shortPath)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
}

func TestURLStatsServiceImpl_InsertAccessLog(t *testing.T) {
	repo := &repoMocks.URLStatisticsRepository{}
	defer repo.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	accessedAt := time.Now()
	repo.On("InsertAccessLog", ctx, shortPath, accessedAt).Return(nil).Once()
	service := NewURLStatsService(repo)
	err := service.InsertAccessLog(ctx, shortPath, accessedAt)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestURLStatsServiceImpl_InsertAccessLog_RepoError(t *testing.T) {
	repo := &repoMocks.URLStatisticsRepository{}
	defer repo.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	accessedAt := time.Now()
	repo.On("InsertAccessLog", ctx, shortPath, accessedAt).Return(assert.AnError).Once()
	service := NewURLStatsService(repo)
	err := service.InsertAccessLog(ctx, shortPath, accessedAt)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
}
