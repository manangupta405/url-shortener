package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"url-shortener/internal/db"
	"url-shortener/internal/models"

	"github.com/go-redis/redis/v8"
)

type urlRepositoryRedisImpl struct {
	client      db.RedisClient
	cacheExpiry time.Duration
}

func NewURLRepositoryRedis(client db.RedisClient, cacheExpiry time.Duration) URLRepository {
	return &urlRepositoryRedisImpl{client: client, cacheExpiry: cacheExpiry}
}

func (r *urlRepositoryRedisImpl) GetShortURL(ctx context.Context, originalURL string) (*models.URL, error) {
	return nil, errors.New("not implemented")
}

func (r *urlRepositoryRedisImpl) GetOriginalURL(ctx context.Context, shortPath string) (*models.URL, error) {
	val, err := r.client.Get(ctx, shortPath).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var url models.URL
	err = json.Unmarshal([]byte(val), &url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepositoryRedisImpl) UpdateShortURL(ctx context.Context, url *models.URL) error {
	return errors.New("not implemented")
}

func (r *urlRepositoryRedisImpl) DeleteShortURL(ctx context.Context, shortPath string, currentTime time.Time, deletedBy string) error {
	err := r.client.Del(ctx, shortPath).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *urlRepositoryRedisImpl) InsertShortURL(ctx context.Context, url *models.URL) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, url.ShortPath, string(data), r.cacheExpiry).Err()
	if err != nil {
		return err
	}
	return nil
}
