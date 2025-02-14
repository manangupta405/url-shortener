package repositories

import (
	"context"
	"database/sql"
	"time"
	"url-shortener/internal/models"
)

type urlRepositoryPostgresqlImpl struct {
	db *sql.DB
}

func NewURLRepositoryPostgresql(db *sql.DB) URLRepository {
	return &urlRepositoryPostgresqlImpl{db: db}
}

// GetShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) GetShortURL(ctx context.Context, originalURL string) (*models.URL, error) {
	row := r.db.QueryRowContext(ctx, PG_GET_BY_ORIGINAL_URL, originalURL)
	url := &models.URL{}
	err := row.Scan(&url.ShortPath, &url.OriginalURL, &url.Expiry, &url.CreatedAt, &url.CreatedBy, &url.ModifiedAt, &url.ModifiedBy, &url.DeletedAt, &url.DeletedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return url, nil
}

// GetOriginalURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) GetOriginalURL(ctx context.Context, shortPath string) (*models.URL, error) {
	row := r.db.QueryRowContext(ctx, PG_GET_BY_SHORT_URL, shortPath)
	url := &models.URL{}
	err := row.Scan(&url.ShortPath, &url.OriginalURL, &url.Expiry, &url.CreatedAt, &url.CreatedBy, &url.ModifiedAt, &url.ModifiedBy, &url.DeletedAt, &url.DeletedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return url, nil
}

// UpdateShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) UpdateShortURL(ctx context.Context, url *models.URL) error {
	_, err := r.db.ExecContext(ctx, PG_UPDATE_SHORT_URL, url.OriginalURL, url.Expiry, url.ModifiedAt, url.ModifiedBy, url.ShortPath)
	return err
}

// DeleteShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) DeleteShortURL(ctx context.Context, shortPath string, currentTime time.Time) error {
	_, err := r.db.ExecContext(ctx, PG_SOFT_DELETE_SHORT_URL, currentTime, "system", shortPath)
	return err
}

// InsertShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) InsertShortURL(ctx context.Context, url *models.URL) error {
	_, err := r.db.ExecContext(ctx, PG_INSERT_SHORT_URL, url.ShortPath, url.OriginalURL, url.Expiry, url.CreatedAt, url.CreatedBy) // using system now, will be replaced by user
	return err
}
