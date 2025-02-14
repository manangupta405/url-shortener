package repositories

import (
	"context"
	"testing"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/repositories/mocks"

	utilMocks "url-shortener/internal/utils/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRepository() (*mocks.URLRepository, *mocks.URLRepository, *utilMocks.TimeProvider, URLRepository) {
	redisRepo := mocks.URLRepository{}
	postgresRepo := mocks.URLRepository{}
	timeProvider := utilMocks.TimeProvider{}
	repo := NewURLRepository(&redisRepo, &postgresRepo, &timeProvider)
	return &redisRepo, &postgresRepo, &timeProvider, repo
}

func TestGetShortURL_Success(t *testing.T) {
	redisRepo, postgresRepo, _, repo := setupRepository()
	mockURL := &models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	postgresRepo.On("GetShortURL", mock.Anything, "https://example.com").Return(mockURL, nil).Once()
	redisRepo.On("InsertShortURL", mock.Anything, mockURL).Return(nil).Once()
	url, err := repo.GetShortURL(context.Background(), "https://example.com")

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, "https://example.com", url.OriginalURL)
	postgresRepo.AssertExpectations(t)
}

func TestGetShortURL_Error(t *testing.T) {
	_, postgresRepo, _, repo := setupRepository()
	postgresRepo.On("GetShortURL", mock.Anything, "https://example.com").Return(nil, assert.AnError).Once()
	url, err := repo.GetShortURL(context.Background(), "https://example.com")

	assert.Error(t, err)
	assert.Nil(t, url)
	postgresRepo.AssertExpectations(t)
}

func TestGetOriginalURL_SuccessRedis(t *testing.T) {
	redisRepo, _, _, repo := setupRepository()
	mockURL := &models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	redisRepo.On("GetOriginalURL", mock.Anything, "shortpath").Return(mockURL, nil).Once()
	url, err := repo.GetOriginalURL(context.Background(), "shortpath")

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, "https://example.com", url.OriginalURL)
	redisRepo.AssertExpectations(t)
}

func TestGetOriginalURL_SuccessPostgres(t *testing.T) {
	redisRepo, postgresRepo, _, repo := setupRepository()
	mockURL := &models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	redisRepo.On("GetOriginalURL", mock.Anything, "shortpath").Return(nil, ErrCacheMiss).Once()
	postgresRepo.On("GetOriginalURL", mock.Anything, "shortpath").Return(mockURL, nil).Once()
	redisRepo.On("InsertShortURL", mock.Anything, mockURL).Return(nil).Once()
	url, err := repo.GetOriginalURL(context.Background(), "shortpath")

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, "https://example.com", url.OriginalURL)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestGetOriginalURL_ErrorRedis(t *testing.T) {
	redisRepo, _, _, repo := setupRepository()
	redisRepo.On("GetOriginalURL", mock.Anything, "shortpath").Return(nil, assert.AnError).Once()
	url, err := repo.GetOriginalURL(context.Background(), "shortpath")

	assert.Error(t, err)
	assert.Nil(t, url)
	redisRepo.AssertExpectations(t)
}

func TestGetOriginalURL_ErrorPostgres(t *testing.T) {
	redisRepo, postgresRepo, _, repo := setupRepository()
	redisRepo.On("GetOriginalURL", mock.Anything, "shortpath").Return(nil, ErrCacheMiss).Once()
	postgresRepo.On("GetOriginalURL", mock.Anything, "shortpath").Return(nil, assert.AnError).Once()
	url, err := repo.GetOriginalURL(context.Background(), "shortpath")

	assert.Error(t, err)
	assert.Nil(t, url)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestUpdateShortURL_Success(t *testing.T) {
	redisRepo, postgresRepo, timeProvider, repo := setupRepository()
	currTime := time.Now()
	mockURL := &models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath", Expiry: &currTime}
	postgresRepo.On("UpdateShortURL", mock.Anything, mockURL).Return(nil).Once()
	redisRepo.On("DeleteShortURL", mock.Anything, "shortpath", mock.Anything, "system").Return(nil).Once()
	timeProvider.On("Now").Return(currTime).Once()
	err := repo.UpdateShortURL(context.Background(), mockURL)

	assert.NoError(t, err)
	postgresRepo.AssertExpectations(t)
	redisRepo.AssertExpectations(t)
}

func TestUpdateShortURL_Error(t *testing.T) {
	redisRepo, postgresRepo, timeProvider, repo := setupRepository()
	mockURL := &models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	postgresRepo.On("UpdateShortURL", mock.Anything, mockURL).Return(assert.AnError).Once()

	timeProvider.On("Now").Return(time.Now()).Once()
	redisRepo.On("DeleteShortURL", mock.Anything, "shortpath", mock.Anything, "system").Return(nil).Once()
	err := repo.UpdateShortURL(context.Background(), mockURL)

	assert.Error(t, err)
	postgresRepo.AssertExpectations(t)
}

func TestDeleteShortURL_Success(t *testing.T) {
	redisRepo, postgresRepo, _, repo := setupRepository()
	postgresRepo.On("DeleteShortURL", mock.Anything, "shortpath", mock.Anything, "testuser").Return(nil).Once()
	redisRepo.On("DeleteShortURL", mock.Anything, "shortpath", mock.Anything, "testuser").Return(nil).Once()
	err := repo.DeleteShortURL(context.Background(), "shortpath", time.Now(), "testuser")

	assert.NoError(t, err)
	postgresRepo.AssertExpectations(t)
	redisRepo.AssertExpectations(t)
}

func TestDeleteShortURL_Error(t *testing.T) {
	_, postgresRepo, _, repo := setupRepository()
	postgresRepo.On("DeleteShortURL", mock.Anything, "shortpath", mock.Anything, "testuser").Return(assert.AnError).Once()
	err := repo.DeleteShortURL(context.Background(), "shortpath", time.Now(), "testuser")

	assert.Error(t, err)
	postgresRepo.AssertExpectations(t)
}

func TestInsertShortURL_Success(t *testing.T) {
	_, postgresRepo, _, repo := setupRepository()
	mockURL := &models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	postgresRepo.On("InsertShortURL", mock.Anything, mockURL).Return(nil).Once()
	err := repo.InsertShortURL(context.Background(), mockURL)

	assert.NoError(t, err)
	postgresRepo.AssertExpectations(t)
}

func TestInsertShortURL_Error(t *testing.T) {
	_, postgresRepo, _, repo := setupRepository()
	mockURL := &models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	postgresRepo.On("InsertShortURL", mock.Anything, mockURL).Return(assert.AnError).Once()
	err := repo.InsertShortURL(context.Background(), mockURL)

	assert.Error(t, err)
	postgresRepo.AssertExpectations(t)
}
