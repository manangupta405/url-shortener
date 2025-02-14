package repositories

import (
	"context"
	"errors"
	"testing"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/repositories/mocks"
	utilsMocks "url-shortener/internal/utils/mocks"

	"github.com/stretchr/testify/assert"
)

func TestURLRepositoryImpl_GetShortURL(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	shortPath := "shortPath"
	url := &models.URL{
		OriginalURL: originalURL,
		ShortPath:   shortPath,
	}
	redisRepo.On("GetShortURL", ctx, originalURL).Return(nil, ErrCacheMiss).Once()
	postgresRepo.On("GetShortURL", ctx, originalURL).Return(url, nil).Once()
	redisRepo.On("InsertShortURL", context.Background(), url).Return(nil).Once()
	urlGot, err := repo.GetShortURL(ctx, originalURL)
	assert.Eventually(t, func() bool {
		return err == nil
	}, time.Millisecond*100, time.Millisecond*10)
	assert.Equal(t, shortPath, urlGot.ShortPath)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_GetShortURL_RedisError(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	redisRepo.On("GetShortURL", ctx, originalURL).Return(nil, errors.New("Internal")).Once()
	_, err := repo.GetShortURL(ctx, originalURL)
	assert.NotNil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_GetShortURL_PostgresError(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	redisRepo.On("GetShortURL", ctx, originalURL).Return(nil, ErrCacheMiss).Once()
	postgresRepo.On("GetShortURL", ctx, originalURL).Return(nil, errors.New("Internal")).Once()
	_, err := repo.GetShortURL(ctx, originalURL)
	assert.NotNil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_GetOriginalURL(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	shortPath := "shortPath"
	originalURL := "https://www.example.com"
	url := &models.URL{
		OriginalURL: originalURL,
		ShortPath:   shortPath,
	}
	redisRepo.On("GetOriginalURL", ctx, shortPath).Return(nil, ErrCacheMiss).Once()
	postgresRepo.On("GetOriginalURL", ctx, shortPath).Return(url, nil).Once()
	redisRepo.On("InsertShortURL", ctx, url).Return(nil).Once()
	urlGot, err := repo.GetOriginalURL(ctx, shortPath)
	assert.Nil(t, err)
	assert.Equal(t, originalURL, urlGot.OriginalURL)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_GetOriginalURL_RedisError(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	shortPath := "shortPath"
	redisRepo.On("GetOriginalURL", ctx, shortPath).Return(nil, errors.New("Internal")).Once()
	_, err := repo.GetOriginalURL(ctx, shortPath)
	assert.NotNil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_GetOriginalURL_PostgresError(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	shortPath := "shortPath"
	redisRepo.On("GetOriginalURL", ctx, shortPath).Return(nil, ErrCacheMiss).Once()
	postgresRepo.On("GetOriginalURL", ctx, shortPath).Return(nil, errors.New("Internal")).Once()
	_, err := repo.GetOriginalURL(ctx, shortPath)
	assert.NotNil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_UpdateShortURL(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	currentTime := time.Now()
	ctx := context.Background()
	shortPath := "shortPath"
	url := &models.URL{
		ShortPath: shortPath,
	}
	postgresRepo.On("UpdateShortURL", ctx, url).Return(nil).Once()
	redisRepo.On("DeleteShortURL", ctx, shortPath, currentTime, "system").Return(nil).Once()
	timeProvider.On("Now").Return(currentTime).Once()
	err := repo.UpdateShortURL(ctx, url)
	assert.Eventually(t, func() bool {
		return err == nil
	}, time.Millisecond*100, time.Millisecond*10)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_UpdateShortURL_Error(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	shortPath := "shortPath"
	url := &models.URL{
		ShortPath: shortPath,
	}
	postgresRepo.On("UpdateShortURL", ctx, url).Return(errors.New("Internal")).Once()
	err := repo.UpdateShortURL(ctx, url)
	assert.NotNil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_DeleteShortURL(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "system"
	postgresRepo.On("DeleteShortURL", ctx, shortPath, currentTime, deletedBy).Return(nil).Once()
	redisRepo.On("DeleteShortURL", ctx, shortPath, currentTime, deletedBy).Return(nil).Once()
	err := repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.Nil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_DeleteShortURL_Error(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "system"
	postgresRepo.On("DeleteShortURL", ctx, shortPath, currentTime, deletedBy).Return(errors.New("Internal")).Once()
	err := repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.NotNil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_InsertShortURL(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	url := &models.URL{}
	postgresRepo.On("InsertShortURL", ctx, url).Return(nil).Once()
	err := repo.InsertShortURL(ctx, url)
	assert.Nil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestURLRepositoryImpl_InsertShortURL_Error(t *testing.T) {
	redisRepo := &mocks.URLRepository{}
	defer redisRepo.AssertExpectations(t)
	postgresRepo := &mocks.URLRepository{}
	defer postgresRepo.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	repo := NewURLRepository(redisRepo, postgresRepo, timeProvider)
	ctx := context.Background()
	url := &models.URL{}
	postgresRepo.On("InsertShortURL", ctx, url).Return(errors.New("Internal")).Once()
	err := repo.InsertShortURL(ctx, url)
	assert.NotNil(t, err)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}
