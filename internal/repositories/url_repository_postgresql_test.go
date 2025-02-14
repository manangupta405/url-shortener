package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"
	"url-shortener/internal/models"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestURLRepositoryPostgresqlImpl_GetShortURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	originalURL := "https://www.example.com"

	rows := sqlmock.NewRows([]string{"short_path", "original_url", "expiry", "created_at", "created_by", "modified_at", "modified_by", "deleted_at", "deleted_by"}).
		AddRow("shortPath", originalURL, time.Now().Add(time.Minute*60), time.Now(), "user", time.Now(), "user", nil, nil)

	mock.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by FROM urls WHERE original_url = ?").WithArgs(originalURL).WillReturnRows(rows)

	url, err := repo.GetShortURL(ctx, originalURL)
	assert.Nil(t, err)
	assert.Equal(t, "shortPath", url.ShortPath)
	assert.Equal(t, originalURL, url.OriginalURL)
}
func TestURLRepositoryPostgresqlImpl_GetOriginalURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"

	rows := sqlmock.NewRows([]string{"short_path", "original_url", "expiry", "created_at", "created_by", "modified_at", "modified_by", "deleted_at", "deleted_by"}).
		AddRow(shortPath, "https://www.example.com", time.Now().Add(time.Minute*60), time.Now(), "user", time.Now(), "user", nil, nil)

	mock.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by FROM urls WHERE short_path = ?").WithArgs(shortPath).WillReturnRows(rows)

	url, err := repo.GetOriginalURL(ctx, shortPath)
	assert.Nil(t, err)
	assert.Equal(t, shortPath, url.ShortPath)
	assert.Equal(t, "https://www.example.com", url.OriginalURL)
}

func TestURLRepositoryPostgresqlImpl_InsertShortURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
	}
	currentTime := time.Now()

	mock.ExpectExec("INSERT INTO urls \\(short_path, original_url, expiry, created_at, created_by\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").WithArgs(url.ShortPath, url.OriginalURL, url.Expiry, currentTime, "system").WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertShortURL(ctx, url, currentTime)
	assert.Nil(t, err)
}
func TestURLRepositoryPostgresqlImpl_UpdateShortURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
	}
	currentTime := time.Now()

	mock.ExpectExec("UPDATE urls SET original_url = \\$1, expiry = \\$2, modified_at = \\$3, modified_by = \\$4 WHERE short_path = \\$5").WithArgs(url.OriginalURL, url.Expiry, currentTime, "system", url.ShortPath).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateShortURL(ctx, url, currentTime)
	assert.Nil(t, err)
}

func TestURLRepositoryPostgresqlImpl_DeleteShortURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()

	mock.ExpectExec("UPDATE urls SET deleted_at = \\$1, deleted_by = \\$2 WHERE short_path = \\$3").WithArgs(currentTime, "system", shortPath).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteShortURL(ctx, shortPath, currentTime)
	assert.Nil(t, err)
}
func TestURLRepositoryPostgresqlImpl_GetShortURL_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	originalURL := "https://www.example.com"

	mock.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by FROM urls WHERE original_url = ?").WithArgs(originalURL).WillReturnError(fmt.Errorf("some error"))

	_, err = repo.GetShortURL(ctx, originalURL)
	assert.NotNil(t, err)
}
func TestURLRepositoryPostgresqlImpl_GetOriginalURL_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"

	mock.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by FROM urls WHERE short_path = ?").WithArgs(shortPath).WillReturnError(fmt.Errorf("some error"))

	_, err = repo.GetOriginalURL(ctx, shortPath)
	assert.NotNil(t, err)
}

func TestURLRepositoryPostgresqlImpl_InsertShortURL_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
	}
	currentTime := time.Now()

	mock.ExpectExec("INSERT INTO urls \\(short_path, original_url, expiry, created_at, created_by\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").WithArgs(url.ShortPath, url.OriginalURL, url.Expiry, currentTime, "system").WillReturnError(fmt.Errorf("some error"))

	err = repo.InsertShortURL(ctx, url, currentTime)
	assert.NotNil(t, err)
}
func TestURLRepositoryPostgresqlImpl_UpdateShortURL_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
	}
	currentTime := time.Now()

	mock.ExpectExec("UPDATE urls SET original_url = \\$1, expiry = \\$2, modified_at = \\$3, modified_by = \\$4 WHERE short_path = \\$5").WithArgs(url.OriginalURL, url.Expiry, currentTime, "system", url.ShortPath).WillReturnError(fmt.Errorf("some error"))

	err = repo.UpdateShortURL(ctx, url, currentTime)
	assert.NotNil(t, err)
}

func TestURLRepositoryPostgresqlImpl_DeleteShortURL_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()

	mock.ExpectExec("UPDATE urls SET deleted_at = \\$1, deleted_by = \\$2 WHERE short_path = \\$3").WithArgs(currentTime, "system", shortPath).WillReturnError(fmt.Errorf("some error"))

	err = repo.DeleteShortURL(ctx, shortPath, currentTime)
	assert.NotNil(t, err)
}
