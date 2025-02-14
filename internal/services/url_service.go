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
	DeleteURL(ctx context.Context, shortPath string) error
	UpdateShortURL(ctx context.Context, originalUrl string, shortUrl string, expiry *time.Time) error
	GetURLDetails(ctx context.Context, shortPath string) (*models.URL, error)
}

type urlServiceImpl struct {
	repo         repositories.URLRepository
	statRepo     repositories.URLStatisticsRepository
	idGenerator  utils.NanoIDGenerator
	timeProvider utils.TimeProvider
}

func NewURLService(repo repositories.URLRepository, statRepo repositories.URLStatisticsRepository, idGenerator utils.NanoIDGenerator, timeProvider utils.TimeProvider) URLService {
	return &urlServiceImpl{repo: repo, statRepo: statRepo, idGenerator: idGenerator, timeProvider: timeProvider}
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
	if url == nil {
		return "", nil
	}
	go func() {
		err := s.statRepo.InsertAccessLog(context.Background(), shortPath, s.timeProvider.Now())
		if err != nil {
			// Log the error.  Consider adding more robust error handling here.
			println("Error inserting access log:", err)
		}
	}()

	return url.OriginalURL, nil
}

// DeleteURL implements URLService.
func (s *urlServiceImpl) DeleteURL(ctx context.Context, shortPath string) error {
	err := s.repo.DeleteShortURL(ctx, shortPath, s.timeProvider.Now(), "system")
	if err != nil {
		return err
	}
	return nil
}

// UpdateShortURL implements URLService
func (s *urlServiceImpl) UpdateShortURL(ctx context.Context, originalUrl string, shortUrl string, expiry *time.Time) error {
	currentTime := s.timeProvider.Now()
	modifiedBy := "system"
	urlUpdate := &models.URL{
		OriginalURL: originalUrl,
		ShortPath:   shortUrl,
		Expiry:      expiry,
		ModifiedAt:  &currentTime,
		ModifiedBy:  &modifiedBy,
	}
	err := s.repo.UpdateShortURL(ctx, urlUpdate)
	if err != nil {
		return err
	}
	return nil
}

// GetURLDetails implements URLService
func (s *urlServiceImpl) GetURLDetails(ctx context.Context, shortPath string) (*models.URL, error) {
	url, err := s.repo.GetOriginalURL(ctx, shortPath)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, nil
	}
	return url, nil
}
