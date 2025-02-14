package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"url-shortener/internal/models"

	"log"
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // URL not found
		}
		log.Printf("Error getting short URL from database: %v, originalURL: %s", err, originalURL)
		return nil, ErrDBError
	}
	if url.Expiry != nil && url.Expiry.Before(time.Now()) {
		return nil, nil // URL expired
	}
	return url, nil
}

// GetOriginalURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) GetOriginalURL(ctx context.Context, shortPath string) (*models.URL, error) {
	row := r.db.QueryRowContext(ctx, PG_GET_BY_SHORT_URL, shortPath)
	url := &models.URL{}
	err := row.Scan(&url.ShortPath, &url.OriginalURL, &url.Expiry, &url.CreatedAt, &url.CreatedBy, &url.ModifiedAt, &url.ModifiedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrURLNotFound //Specific error for no rows
		}
		log.Printf("Error getting original URL from database: %v, shortPath: %s", err, shortPath)
		return nil, ErrDBError
	}
	if url.Expiry != nil && url.Expiry.Before(time.Now()) {
		return nil, ErrURLExpired //Specific error for expired URL
	}
	return url, nil
}

// UpdateShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) UpdateShortURL(ctx context.Context, url *models.URL) error {
	_, err := r.db.ExecContext(ctx, PG_UPDATE_SHORT_URL, url.OriginalURL, url.Expiry, url.ModifiedAt, url.ModifiedBy, url.ShortPath)
	if err != nil {
		log.Printf("Error updating short URL in database: %v, url: %+v", err, url)
		return ErrDBError
	}
	return nil
}

// DeleteShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) DeleteShortURL(ctx context.Context, shortPath string, currentTime time.Time, deletedBy string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction: %v, shortPath: %s", err, shortPath)
		return ErrDBError
	}
	defer tx.Rollback() // Rollback on error

	// Delete from urls and retrieve the deleted row using RETURNING
	row := tx.QueryRowContext(ctx, PG_GET_BY_SHORT_URL, shortPath)

	urlArchive := &models.URLArchive{}
	err = row.Scan(&urlArchive.ShortPath, &urlArchive.OriginalURL, &urlArchive.Expiry, &urlArchive.CreatedAt, &urlArchive.CreatedBy, &urlArchive.ModifiedAt, &urlArchive.ModifiedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrURLNotFound // URL not found
		}
		log.Printf("Error retrieving URL before deletion: %v, shortPath: %s", err, shortPath)
		return ErrDBError
	}

	urlArchive.DeletedAt = &currentTime
	urlArchive.DeletedBy = &deletedBy

	// Insert into urls_archive
	_, err = tx.ExecContext(ctx, PG_INSERT_URL_ARCHIVE, urlArchive.ShortPath, urlArchive.OriginalURL, urlArchive.Expiry, urlArchive.CreatedAt, urlArchive.CreatedBy, urlArchive.ModifiedAt, urlArchive.ModifiedBy, urlArchive.DeletedAt, urlArchive.DeletedBy)
	if err != nil {
		log.Printf("Error inserting into url_archive: %v, urlArchive: %+v", err, urlArchive)
	}

	// Delete from urls
	_, err = tx.ExecContext(ctx, PG_DELETE_SHORT_URL, urlArchive.ShortPath)
	if err != nil {
		log.Printf("Error deleting URL from urls table: %v, shortPath: %v", err, urlArchive.ShortPath)
		return ErrDBError
	}

	return tx.Commit()
}

// InsertShortURL implements URLRepository.
func (r *urlRepositoryPostgresqlImpl) InsertShortURL(ctx context.Context, url *models.URL) error {
	_, err := r.db.ExecContext(ctx, PG_INSERT_SHORT_URL, url.ShortPath, url.OriginalURL, url.Expiry, url.CreatedAt, url.CreatedBy)
	if err != nil {
		log.Printf("Error inserting short URL into database: %v, url: %+v", err, url)
		return ErrDBError
	}
	return nil
}
