package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
	"url-shortener/internal/models"

	dbMocks "url-shortener/internal/db/mocks"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRedisRepository() (*dbMocks.RedisClient, *urlRepositoryRedisImpl) {
	mockClient := &dbMocks.RedisClient{}
	repo := NewURLRepositoryRedis(mockClient, 10*time.Minute)
	return mockClient, repo.(*urlRepositoryRedisImpl)
}

func TestRedisGetOriginalURL_Success(t *testing.T) {
	mockClient, repo := setupRedisRepository()
	mockURL := models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	data, _ := json.Marshal(mockURL)
	expectedOut := &redis.StringCmd{}
	expectedOut.SetVal(string(data))
	expectedOut.SetErr(nil)
	mockClient.On("Get", mock.Anything, "shortpath").Return(expectedOut, nil).Once()

	url, err := repo.GetOriginalURL(context.Background(), "shortpath")

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, "https://example.com", url.OriginalURL)
	mockClient.AssertExpectations(t)
}

func TestRedisGetOriginalURL_NotFound(t *testing.T) {
	mockClient, repo := setupRedisRepository()
	expectedOut := &redis.StringCmd{}
	expectedOut.SetVal("")
	expectedOut.SetErr(redis.Nil)
	mockClient.On("Get", mock.Anything, "shortpath").Return(expectedOut, redis.Nil).Once()

	url, err := repo.GetOriginalURL(context.Background(), "shortpath")

	assert.NoError(t, err)
	assert.Nil(t, url)
	mockClient.AssertExpectations(t)
}

func TestRedisGetOriginalURL_Error(t *testing.T) {
	mockClient, repo := setupRedisRepository()
	expectedOut := &redis.StringCmd{}
	expectedOut.SetVal("")
	expectedOut.SetErr(errors.New("redis error"))
	mockClient.On("Get", mock.Anything, "shortpath").Return(expectedOut).Once()

	url, err := repo.GetOriginalURL(context.Background(), "shortpath")

	assert.Error(t, err)
	assert.Nil(t, url)
	mockClient.AssertExpectations(t)
}
func TestRedisInsertShortURL_Success(t *testing.T) {
	mockClient, repo := setupRedisRepository()
	expectedOut := &redis.StatusCmd{}
	expectedOut.SetVal("")
	expectedOut.SetErr(nil)
	mockURL := models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	mockClient.On("Set", mock.Anything, "shortpath", mock.Anything, repo.cacheExpiry).Return(expectedOut).Once()
	err := repo.InsertShortURL(context.Background(), &mockURL)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisInsertShortURL_Error(t *testing.T) {
	mockClient, repo := setupRedisRepository()
	mockURL := models.URL{OriginalURL: "https://example.com", ShortPath: "shortpath"}
	expectedOut := &redis.StatusCmd{}
	expectedOut.SetVal("")
	expectedOut.SetErr(errors.New("redis error"))
	mockClient.On("Set", mock.Anything, "shortpath", mock.Anything, repo.cacheExpiry).Return(expectedOut).Once()
	err := repo.InsertShortURL(context.Background(), &mockURL)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisDeleteShortURL_Success(t *testing.T) {
	mockClient, repo := setupRedisRepository()
	expectedOut := &redis.IntCmd{}
	expectedOut.SetVal(0)
	expectedOut.SetErr(nil)
	mockClient.On("Del", mock.Anything, "shortpath").Return(expectedOut).Once()
	err := repo.DeleteShortURL(context.Background(), "shortpath", time.Now(), "testuser")
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisDeleteShortURL_Error(t *testing.T) {
	mockClient, repo := setupRedisRepository()
	expectedOut := &redis.IntCmd{}
	expectedOut.SetVal(0)
	expectedOut.SetErr(errors.New("redis error"))
	mockClient.On("Del", mock.Anything, "shortpath").Return(expectedOut).Once()
	err := repo.DeleteShortURL(context.Background(), "shortpath", time.Now(), "testuser")
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}
