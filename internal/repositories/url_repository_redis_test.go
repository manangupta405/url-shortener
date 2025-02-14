package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
	dbMocks "url-shortener/internal/db/mocks"
	"url-shortener/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestURLRepositoryRedisImpl_InsertShortURL(t *testing.T) {
	client := &dbMocks.RedisClient{}
	defer client.AssertExpectations(t)
	repo := NewURLRepositoryRedis(client, time.Minute*60)
	ctx := context.Background()
	currentTime := time.Now()
	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
		CreatedAt:   &currentTime,
		CreatedBy:   "system",
	}
	data, _ := json.Marshal(url)
	client.On("Set", ctx, url.ShortPath, string(data), time.Minute*60).Return(redis.NewStatusCmd(ctx, "OK")).Once()
	err := repo.InsertShortURL(ctx, url)
	assert.Nil(t, err)
	client.AssertExpectations(t)
}

func TestURLRepositoryRedisImpl_InsertShortURL_Error(t *testing.T) {
	client := &dbMocks.RedisClient{}
	defer client.AssertExpectations(t)
	repo := NewURLRepositoryRedis(client, time.Minute*60)
	ctx := context.Background()
	currentTime := time.Now()
	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
		CreatedAt:   &currentTime,
		CreatedBy:   "system",
	}
	errorResponse := redis.NewStatusCmd(ctx)
	errorResponse.SetErr(errors.New("Internal"))
	client.On("Set", ctx, url.ShortPath, mock.Anything, time.Minute*60).Return(errorResponse).Once()
	err := repo.InsertShortURL(ctx, url)
	assert.NotNil(t, err)
	client.AssertExpectations(t)
}

func TestURLRepositoryRedisImpl_GetOriginalURL(t *testing.T) {
	client := &dbMocks.RedisClient{}
	defer client.AssertExpectations(t)
	repo := NewURLRepositoryRedis(client, time.Minute*60)
	ctx := context.Background()
	shortPath := "shortPath"
	originalURL := "https://www.example.com"
	url := &models.URL{
		OriginalURL: originalURL,
		ShortPath:   shortPath,
	}
	data, _ := json.Marshal(url)
	response := redis.NewStringCmd(ctx)
	response.SetVal(string(data))
	client.On("Get", ctx, shortPath).Return(response).Once()
	urlGot, err := repo.GetOriginalURL(ctx, shortPath)
	assert.Nil(t, err)
	assert.Equal(t, originalURL, urlGot.OriginalURL)
	client.AssertExpectations(t)
}

func TestURLRepositoryRedisImpl_GetOriginalURL_Nil(t *testing.T) {
	client := &dbMocks.RedisClient{}
	defer client.AssertExpectations(t)
	repo := NewURLRepositoryRedis(client, time.Minute*60)
	ctx := context.Background()
	shortPath := "shortPath"
	response := redis.NewStringCmd(ctx)
	response.SetErr(redis.Nil)
	client.On("Get", ctx, shortPath).Return(response).Once()
	urlGot, err := repo.GetOriginalURL(ctx, shortPath)
	assert.Nil(t, err)
	assert.Nil(t, urlGot)
	client.AssertExpectations(t)
}

func TestURLRepositoryRedisImpl_GetOriginalURL_Error(t *testing.T) {
	client := &dbMocks.RedisClient{}
	defer client.AssertExpectations(t)
	repo := NewURLRepositoryRedis(client, time.Minute*60)
	ctx := context.Background()
	shortPath := "shortPath"
	client.On("Get", ctx, shortPath).Return(redis.NewStringCmd(ctx, errors.New("Internal"))).Once()
	_, err := repo.GetOriginalURL(ctx, shortPath)
	assert.NotNil(t, err)
	client.AssertExpectations(t)
}

func TestURLRepositoryRedisImpl_DeleteShortURL(t *testing.T) {
	client := &dbMocks.RedisClient{}
	defer client.AssertExpectations(t)
	repo := NewURLRepositoryRedis(client, time.Minute*60)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "system"
	client.On("Del", ctx, shortPath).Return(redis.NewIntCmd(ctx, 1)).Once()
	err := repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.Nil(t, err)
	client.AssertExpectations(t)
}

func TestURLRepositoryRedisImpl_DeleteShortURL_Error(t *testing.T) {
	client := &dbMocks.RedisClient{}
	defer client.AssertExpectations(t)
	repo := NewURLRepositoryRedis(client, time.Minute*60)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "system"
	response := redis.NewIntCmd(ctx)
	response.SetErr(errors.New("Internal"))
	client.On("Del", ctx, shortPath).Return(response).Once()
	err := repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.NotNil(t, err)
	client.AssertExpectations(t)
}
