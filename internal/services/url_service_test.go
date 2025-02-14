package services

import (
	"context"
	"errors"
	"testing"
	"time"

	repoMocks "url-shortener/internal/repositories/mocks"
	utilsMocks "url-shortener/internal/utils/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestURLServiceImpl_CreateShortURL(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
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
	repo.On("InsertShortURL", ctx, mock.Anything).Return(nil).Once()

	service := NewURLService(repo, idGenerator, timeProvider)
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
	repo.On("InsertShortURL", ctx, mock.Anything).Return(errors.New("Internal")).Once()
	service := NewURLService(repo, idGenerator, timeProvider)
	_, err := service.CreateShortURL(ctx, originalURL, &expiry)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}

func TestURLServiceImpl_CreateShortURL_IDError(t *testing.T) {
	repo := &repoMocks.URLRepository{}
	defer repo.AssertExpectations(t)
	idGenerator := &utilsMocks.NanoIDGenerator{}
	defer idGenerator.AssertExpectations(t)
	timeProvider := &utilsMocks.TimeProvider{}
	defer timeProvider.AssertExpectations(t)
	ctx := context.Background()
	originalURL := "https://www.example.com"
	expiry := time.Now().Add(time.Minute * 60)
	currentTime := time.Now()
	idGenerator.On("Generate").Return("", errors.New("Internal")).Once()
	timeProvider.On("Now").Return(currentTime).Once()
	service := NewURLService(repo, idGenerator, timeProvider)
	_, err := service.CreateShortURL(ctx, originalURL, &expiry)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	idGenerator.AssertExpectations(t)
	timeProvider.AssertExpectations(t)
}
