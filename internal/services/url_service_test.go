package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"url-shortener/internal/models"
	repoMocks "url-shortener/internal/repositories/mocks"
	utilsMocks "url-shortener/internal/utils/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestURLServiceImpl_CreateShortURL(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	expiry := time.Now().Add(time.Minute * 60)
	shortPath := "shortPath"
	currentTime := time.Now()

	idGenerator.On("Generate").Return(shortPath, nil).Once()
	timeProvider.On("Now").Return(currentTime).Once()
	repo.On("GetShortURL", ctx, originalURL).Return(nil, nil).Once()
	repo.On("InsertShortURL", ctx, mock.Anything).Return(nil).Once()

	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	shortPathGenerated, err := service.CreateShortURL(ctx, originalURL, &expiry)
	assert.Nil(t, err)
	assert.Equal(t, shortPath, shortPathGenerated)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_CreateShortURL_DBError(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	expiry := time.Now().Add(time.Minute * 60)
	shortPath := "shortPath"
	idGenerator.On("Generate").Return(shortPath, nil).Once()
	timeProvider.On("Now").Return(time.Now()).Once()
	repo.On("GetShortURL", ctx, originalURL).Return(nil, nil).Once()
	repo.On("InsertShortURL", ctx, mock.Anything).Return(errors.New("Internal")).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	_, err := service.CreateShortURL(ctx, originalURL, &expiry)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_CreateShortURL_IDError(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	expiry := time.Now().Add(time.Minute * 60)
	currentTime := time.Now()
	repo.On("GetShortURL", ctx, originalURL).Return(nil, nil).Once()
	idGenerator.On("Generate").Return("", errors.New("Internal")).Once()
	timeProvider.On("Now").Return(currentTime).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	_, err := service.CreateShortURL(ctx, originalURL, &expiry)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_CreateShortURL_URLFound(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	expiry := time.Now().Add(time.Minute * 60)
	shortPath := "shortPath"
	shortURL := &models.URL{
		OriginalURL: originalURL,
		ShortPath:   shortPath,
		Expiry:      &expiry,
	}
	repo.On("GetShortURL", ctx, originalURL).Return(shortURL, nil).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	shortPathGenerated, err := service.CreateShortURL(ctx, originalURL, &expiry)
	assert.Nil(t, err)
	assert.Equal(t, shortPath, shortPathGenerated)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_CreateShortURL_RepoError(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	expiry := time.Now().Add(time.Minute * 60)
	repo.On("GetShortURL", ctx, originalURL).Return(nil, errors.New("Internal")).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	_, err := service.CreateShortURL(ctx, originalURL, &expiry)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_GetLongURL(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	originalURL := "https://www.example.com"
	url := &models.URL{
		OriginalURL: originalURL,
		ShortPath:   shortPath,
	}

	currentTime := time.Now()

	repo.On("GetOriginalURL", ctx, shortPath).Return(url, nil).Once()
	timeProvider.On("Now").Return(currentTime).Once()
	statRepo.On("InsertAccessLog", mock.Anything, shortPath, currentTime).Return(nil).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	longURL, err := service.GetLongURL(ctx, shortPath)
	assert.Eventually(t, func() bool {
		return err == nil
	}, time.Millisecond*100, time.Millisecond*10)
	assert.Equal(t, originalURL, longURL)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_GetLongURL_RepoError(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	// currentTime := time.Now()
	shortPath := "shortPath"
	repo.On("GetOriginalURL", ctx, shortPath).Return(nil, errors.New("Internal")).Once()
	// timeProvider.On("Now").Return(currentTime).Once()
	// statRepo.On("InsertAccessLog", mock.Anything, shortPath, currentTime).Return(nil).Once()

	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	_, err := service.GetLongURL(ctx, shortPath)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_DeleteURL(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "system"
	repo.On("DeleteShortURL", ctx, shortPath, currentTime, deletedBy).Return(nil).Once()
	timeProvider.On("Now").Return(currentTime).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	err := service.DeleteURL(ctx, shortPath)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_DeleteURL_RepoError(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "system"
	repo.On("DeleteShortURL", ctx, shortPath, currentTime, deletedBy).Return(errors.New("Internal")).Once()
	timeProvider.On("Now").Return(currentTime).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	err := service.DeleteURL(ctx, shortPath)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_UpdateShortURL(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	shortPath := "shortPath"
	expiry := time.Now().Add(time.Minute * 60)
	currentTime := time.Now()
	modifiedBy := "system"
	urlUpdate := &models.URL{
		OriginalURL: originalURL,
		ShortPath:   shortPath,
		Expiry:      &expiry,
		ModifiedAt:  &currentTime,
		ModifiedBy:  &modifiedBy,
	}

	repo.On("UpdateShortURL", ctx, urlUpdate).Return(nil).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	timeProvider.On("Now").Return(currentTime).Once()

	err := service.UpdateShortURL(ctx, originalURL, shortPath, &expiry)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_UpdateShortURL_RepoError(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	shortPath := "shortPath"
	expiry := time.Now().Add(time.Minute * 60)
	currentTime := time.Now()
	modifiedBy := "system"
	urlUpdate := &models.URL{
		OriginalURL: originalURL,
		ShortPath:   shortPath,
		Expiry:      &expiry,
		ModifiedAt:  &currentTime,
		ModifiedBy:  &modifiedBy,
	}
	timeProvider.On("Now").Return(currentTime).Once()
	repo.On("UpdateShortURL", ctx, urlUpdate).Return(errors.New("Internal")).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	err := service.UpdateShortURL(ctx, originalURL, shortPath, &expiry)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_GetURLDetails(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	originalURL := "https://www.example.com"
	url := &models.URL{
		OriginalURL: originalURL,
		ShortPath:   shortPath,
	}
	repo.On("GetOriginalURL", ctx, shortPath).Return(url, nil).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	urlDetails, err := service.GetURLDetails(ctx, shortPath)
	assert.Nil(t, err)
	assert.Equal(t, originalURL, urlDetails.OriginalURL)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_GetURLDetails_RepoError(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	statRepo := &repoMocks.URLStatisticsRepository{}
	defer statRepo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	shortPath := "shortPath"
	repo.On("GetOriginalURL", ctx, shortPath).Return(nil, errors.New("Internal")).Once()
	service := NewURLService(repo, statRepo, idGenerator, timeProvider)
	_, err := service.GetURLDetails(ctx, shortPath)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}
