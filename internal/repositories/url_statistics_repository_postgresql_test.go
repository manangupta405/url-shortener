package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestURLStatisticsRepositoryPostgresqlImpl_GetURLStatistics(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLStatisticsRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"

	rows := sqlmock.NewRows([]string{"last_24_hours", "past_week", "all_time"}).
		AddRow(10, 100, 1000)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FILTER \\(WHERE accessed_at >= NOW\\(\\) - INTERVAL '24 hours'\\) AS last_24_hours, COUNT\\(\\*\\) FILTER \\(WHERE accessed_at >= NOW\\(\\) - INTERVAL '7 days'\\) AS past_week, COUNT\\(\\*\\) AS all_time FROM url_access_logs WHERE short_path = \\$1;").WithArgs(shortPath).WillReturnRows(rows)

	stats, err := repo.GetURLStatistics(ctx, shortPath)

	assert.Nil(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, shortPath, stats.ShortPath)
	assert.Equal(t, int64(10), stats.Last24Hours)
	assert.Equal(t, int64(100), stats.PastWeek)
	assert.Equal(t, int64(1000), stats.AllTime)
}

func TestURLStatisticsRepositoryPostgresqlImpl_GetURLStatistics_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLStatisticsRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FILTER \\(WHERE accessed_at >= NOW\\(\\) - INTERVAL '24 hours'\\) AS last_24_hours, COUNT\\(\\*\\) FILTER \\(WHERE accessed_at >= NOW\\(\\) - INTERVAL '7 days'\\) AS past_week, COUNT\\(\\*\\) AS all_time FROM url_access_logs WHERE short_path = \\$1;").WithArgs(shortPath).WillReturnError(fmt.Errorf("some error"))

	stats, err := repo.GetURLStatistics(ctx, shortPath)

	assert.Error(t, err)
	assert.Nil(t, stats)
}

func TestURLStatisticsRepositoryPostgresqlImpl_InsertAccessLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLStatisticsRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	accessedAt := time.Now()

	mock.ExpectExec("INSERT INTO url_access_logs \\(short_path,accessed_at\\) VALUES \\(\\$1,\\$2\\);").WithArgs(shortPath, accessedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertAccessLog(ctx, shortPath, accessedAt)

	assert.NoError(t, err)
}

func TestURLStatisticsRepositoryPostgresqlImpl_InsertAccessLog_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewURLStatisticsRepositoryPostgresql(db)
	ctx := context.Background()
	shortPath := "shortPath"
	accessedAt := time.Now()

	mock.ExpectExec("INSERT INTO url_access_logs \\(short_path,accessed_at\\) VALUES \\(\\$1,\\$2\\);").WithArgs(shortPath, accessedAt).WillReturnError(fmt.Errorf("some error"))

	err = repo.InsertAccessLog(ctx, shortPath, accessedAt)

	assert.Error(t, err)
}
