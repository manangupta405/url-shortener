package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
	"url-shortener/internal/models"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestURLRepositoryPostgresqlImpl_GetShortURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	originalURL := "https://www.example.com"

	rows := sqlmock.NewRows([]string{"short_path", "original_url", "expiry", "created_at", "created_by", "modified_at", "modified_by"}).
		AddRow("shortPath", originalURL, time.Now().Add(time.Minute*60), time.Now(), "user", time.Now(), "user")

	mock.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by FROM urls WHERE original_url = ?").WithArgs(originalURL).WillReturnRows(rows)

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

	rows := sqlmock.NewRows([]string{"short_path", "original_url", "expiry", "created_at", "created_by", "modified_at", "modified_by"}).
		AddRow(shortPath, "https://www.example.com", time.Now().Add(time.Minute*60), time.Now(), "user", time.Now(), "user")

	mock.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by FROM urls WHERE short_path = ?").WithArgs(shortPath).WillReturnRows(rows)

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
	currentTime := time.Now()

	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
		CreatedAt:   &currentTime,
		CreatedBy:   "system",
	}

	mock.ExpectExec("INSERT INTO urls \\(short_path, original_url, expiry, created_at, created_by\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").WithArgs(url.ShortPath, url.OriginalURL, url.Expiry, url.CreatedAt, url.CreatedBy).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.InsertShortURL(ctx, url)
	assert.Nil(t, err)
}

func TestURLRepositoryPostgresqlImpl_UpdateShortURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	currentTime := time.Now()
	modifiedBy := "system"

	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
		ModifiedAt:  &currentTime,
		ModifiedBy:  &modifiedBy,
	}

	mock.ExpectExec("UPDATE urls SET original_url = \\$1, expiry = \\$2, modified_at = \\$3, modified_by = \\$4 WHERE short_path = \\$5").WithArgs(url.OriginalURL, url.Expiry, url.ModifiedAt, url.ModifiedBy, url.ShortPath).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateShortURL(ctx, url)
	assert.Nil(t, err)
}
func TestURLRepositoryPostgresqlImpl_DeleteShortURL_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "testuser"

	mockDB.ExpectBegin()

	returnedRows := sqlmock.NewRows([]string{"short_path", "original_url", "expiry", "created_at", "created_by", "modified_at", "modified_by"}).
		AddRow(shortPath, "https://www.example.com", currentTime.Add(time.Minute*60), currentTime, "user", currentTime, "user")
	mockDB.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by FROM urls WHERE short_path = \\$1").WithArgs(shortPath).WillReturnRows(returnedRows)

	mockDB.ExpectExec("INSERT INTO urls_archive \\(short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7, \\$8, \\$9\\)").
		WithArgs(shortPath, "https://www.example.com", currentTime.Add(time.Minute*60), currentTime, "user", currentTime, "user", &currentTime, &deletedBy).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectExec("DELETE FROM urls WHERE short_path = \\$1").WithArgs(shortPath).WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectCommit()

	err = repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.NoError(t, err)
}

func TestURLRepositoryPostgresqlImpl_DeleteShortURL_BeginTxError(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "testuser"
	mockDB.ExpectBegin().WillReturnError(fmt.Errorf("begin error"))
	err = repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.NotNil(t, err)
}

func TestURLRepositoryPostgresqlImpl_DeleteShortURL_SelectError(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "testuser"
	mockDB.ExpectBegin()
	mockDB.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by FROM urls WHERE short_path = \\$1").WithArgs(shortPath).WillReturnError(fmt.Errorf("select error"))
	mockDB.ExpectRollback()
	err = repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.NotNil(t, err)
}

func TestURLRepositoryPostgresqlImpl_DeleteShortURL_InsertArchiveError(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "testuser"
	mockDB.ExpectBegin()
	mockDB.ExpectExec("INSERT INTO urls_archive (.+) VALUES (.+)").
		WithArgs(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything). // Match any args
		WillReturnError(fmt.Errorf("insert archive error"))
	mockDB.ExpectRollback()

	err = repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.NotNil(t, err)
}

func TestURLRepositoryPostgresqlImpl_DeleteShortURL_DeleteError(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	currentTime := time.Now()
	deletedBy := "testuser"
	mockDB.ExpectBegin()
	mockDB.ExpectExec("DELETE FROM urls WHERE short_path = \\$1").WithArgs(shortPath).WillReturnError(fmt.Errorf("delete error"))

	mockDB.ExpectRollback()
	err = repo.DeleteShortURL(ctx, shortPath, currentTime, deletedBy)
	assert.NotNil(t, err)
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
	currentTime := time.Now()
	createdBy := "system"
	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
		CreatedAt:   &currentTime,
		CreatedBy:   createdBy,
	}

	mock.ExpectExec("INSERT INTO urls \\(short_path, original_url, expiry, created_at, created_by\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").WithArgs(url.ShortPath, url.OriginalURL, url.Expiry, url.CreatedAt, url.CreatedBy).WillReturnError(fmt.Errorf("some error"))

	err = repo.InsertShortURL(ctx, url)
	assert.NotNil(t, err)
}

func TestURLRepositoryPostgresqlImpl_UpdateShortURL_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	currentTime := time.Now()
	modifiedBy := "system"
	url := &models.URL{
		ShortPath:   "shortPath",
		OriginalURL: "https://www.example.com",
		Expiry:      &time.Time{},
		ModifiedAt:  &currentTime,
		ModifiedBy:  &modifiedBy,
	}

	mock.ExpectExec("UPDATE urls SET original_url = \\$1, expiry = \\$2, modified_at = \\$3, modified_by = \\$4 WHERE short_path = \\$5").WithArgs(url.OriginalURL, url.Expiry, url.ModifiedAt, url.ModifiedBy, url.ShortPath).WillReturnError(fmt.Errorf("some error"))

	err = repo.UpdateShortURL(ctx, url)
	assert.NotNil(t, err)
}

func TestURLRepositoryPostgresqlImpl_GetShortURL_NoRowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	originalURL := "https://www.example.com"

	mock.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by FROM urls WHERE original_url = ?").WithArgs(originalURL).WillReturnError(sql.ErrNoRows)

	url, err := repo.GetShortURL(ctx, originalURL)
	assert.Nil(t, err)
	assert.Nil(t, url)
}

func TestURLRepositoryPostgresqlImpl_GetOriginalURL_NoRowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := NewURLRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"

	mock.ExpectQuery("SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by FROM urls WHERE short_path = ?").WithArgs(shortPath).WillReturnError(sql.ErrNoRows)

	url, err := repo.GetOriginalURL(ctx, shortPath)
	assert.NotNil(t, err)
	assert.Nil(t, url)
}
