package services

import (
	"context"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/repositories"
	"url-shortener/internal/utils"
)

//go:generate mockery --name=URLService --output=./mocks
type URLService interface {
	CreateShortURL(ctx context.Context, originalURL string, expiry *time.Time) (string, error)
	GetLongURL(ctx context.Context, shortPath string) (string, error)
}

type urlServiceImpl struct {
	repo         repositories.URLRepository
	idGenerator  utils.NanoIDGenerator
	timeProvider utils.TimeProvider
}

func NewURLService(repo repositories.URLRepository, idGenerator utils.NanoIDGenerator, timeProvider utils.TimeProvider) URLService {
	return &urlServiceImpl{repo: repo, idGenerator: idGenerator, timeProvider: timeProvider}
}

// CreateShortURL implements URLService.
func (s *urlServiceImpl) CreateShortURL(ctx context.Context, originalURL string, expiry *time.Time) (string, error) {
	shortUrl, err := s.repo.GetShortURL(ctx, originalURL)
	if err != nil {
		return "", err
	}
	if shortUrl != nil {
		return shortUrl.ShortPath, nil
	}
	currentTime := s.timeProvider.Now()
	shortURL := &models.URL{
		OriginalURL: originalURL,
		Expiry:      expiry,
		CreatedAt:   &currentTime,
		CreatedBy:   "system",
	}

	shortPath, err := s.idGenerator.Generate()
	if err != nil {
		return "", err
	}
	shortURL.ShortPath = shortPath

	err = s.repo.InsertShortURL(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return shortPath, nil
}

// GetLongURL implements URLService.
func (s *urlServiceImpl) GetLongURL(ctx context.Context, shortPath string) (string, error) {
	url, err := s.repo.GetOriginalURL(ctx, shortPath)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}
