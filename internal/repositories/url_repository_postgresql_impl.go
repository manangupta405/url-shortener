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
	err := row.Scan(&url.ShortPath, &url.OriginalURL, &url.Expiry, &url.CreatedAt, &url.CreatedBy, &url.ModifiedAt, &url.ModifiedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if url.Expiry != nil && url.Expiry.Before(time.Now()) {
		return nil, nil
	}
	return url, nil
}

// GetOriginalURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) GetOriginalURL(ctx context.Context, shortPath string) (*models.URL, error) {
	row := r.db.QueryRowContext(ctx, PG_GET_BY_SHORT_URL, shortPath)
	url := &models.URL{}
	err := row.Scan(&url.ShortPath, &url.OriginalURL, &url.Expiry, &url.CreatedAt, &url.CreatedBy, &url.ModifiedAt, &url.ModifiedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if url.Expiry != nil && url.Expiry.Before(time.Now()) {
		return nil, nil
	}
	return url, nil
}

// UpdateShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) UpdateShortURL(ctx context.Context, url *models.URL) error {
	_, err := r.db.ExecContext(ctx, PG_UPDATE_SHORT_URL, url.OriginalURL, url.Expiry, url.ModifiedAt, url.ModifiedBy, url.ShortPath)
	return err
}

// DeleteShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) DeleteShortURL(ctx context.Context, shortPath string, currentTime time.Time, deletedBy string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback on error

	// Delete from urls and retrieve the deleted row using RETURNING
	row := tx.QueryRowContext(ctx, PG_GET_BY_SHORT_URL, shortPath)

	urlArchive := &models.URLArchive{}
	err = row.Scan(&urlArchive.ShortPath, &urlArchive.OriginalURL, &urlArchive.Expiry, &urlArchive.CreatedAt, &urlArchive.CreatedBy, &urlArchive.ModifiedAt, &urlArchive.ModifiedBy)
	if err != nil {
		return err // Handle the error appropriately
	}

	urlArchive.DeletedAt = &currentTime
	urlArchive.DeletedBy = &deletedBy

	// Insert into urls_archive
	_, err = tx.ExecContext(ctx, PG_INSERT_URL_ARCHIVE, urlArchive.ShortPath, urlArchive.OriginalURL, urlArchive.Expiry, urlArchive.CreatedAt, urlArchive.CreatedBy, urlArchive.ModifiedAt, urlArchive.ModifiedBy, urlArchive.DeletedAt, urlArchive.DeletedBy)
	if err != nil {
		return err
	}

	// Insert into urls_archive
	_, err = tx.ExecContext(ctx, PG_DELETE_SHORT_URL, urlArchive.ShortPath)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// InsertShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) InsertShortURL(ctx context.Context, url *models.URL) error {
	_, err := r.db.ExecContext(ctx, PG_INSERT_SHORT_URL, url.ShortPath, url.OriginalURL, url.Expiry, url.CreatedAt, url.CreatedBy) // using system now, will be replaced by user
	return err
}
