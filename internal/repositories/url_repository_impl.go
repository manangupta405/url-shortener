package repositories

import (
	"context"
	"errors"

	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/utils"

	"github.com/rs/zerolog/log"
)

type urlRepositoryImpl struct {
	redisRepo    URLRepository
	postgresRepo URLRepository
	timeProvider utils.TimeProvider
}

func NewURLRepository(redisRepo URLRepository, postgresRepo URLRepository, timeProvider utils.TimeProvider) URLRepository {
	return &urlRepositoryImpl{redisRepo: redisRepo, postgresRepo: postgresRepo, timeProvider: timeProvider}
}

func (r *urlRepositoryImpl) GetShortURL(ctx context.Context, originalURL string) (*models.URL, error) {
	url, err := r.postgresRepo.GetShortURL(ctx, originalURL)
	if err != nil {
		return nil, err
	}
	if url != nil {
		go func() {
			err = r.redisRepo.InsertShortURL(context.Background(), url)
			if err != nil {
				log.Print("Error inserting to redis" + err.Error())
			}
		}()

	}
	return url, nil
}

func (r *urlRepositoryImpl) GetOriginalURL(ctx context.Context, shortPath string) (*models.URL, error) {
	url, err := r.redisRepo.GetOriginalURL(ctx, shortPath)
	if err != nil && !errors.Is(err, ErrCacheMiss) {
		return nil, err
	}
	if url != nil {
		return url, nil
	}
	url, err = r.postgresRepo.GetOriginalURL(ctx, shortPath)
	if err != nil {
		return nil, err
	}
	if url != nil {
		err = r.redisRepo.InsertShortURL(ctx, url)
		if err != nil {
			return nil, err
		}
	}
	return url, nil
}

func (r *urlRepositoryImpl) UpdateShortURL(ctx context.Context, url *models.URL) error {
	err := r.postgresRepo.UpdateShortURL(ctx, url)
	if err != nil {
		return err
	}
	return r.redisRepo.DeleteShortURL(ctx, url.ShortPath, r.timeProvider.Now(), "system")
}

func (r *urlRepositoryImpl) DeleteShortURL(ctx context.Context, shortPath string, currentTime time.Time, deletedBy string) error {
	err := r.postgresRepo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	if err != nil {
		return err
	}
	return r.redisRepo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
}

func (r *urlRepositoryImpl) InsertShortURL(ctx context.Context, url *models.URL) error {
	err := r.postgresRepo.InsertShortURL(ctx, url)
	if err != nil {
		return err
	}
	return nil
}

var ErrCacheMiss = errors.New("cache miss")
