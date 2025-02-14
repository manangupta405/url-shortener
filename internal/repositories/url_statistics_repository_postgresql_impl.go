package repositories

import (
	"context"
	"database/sql"
	"time"

	"url-shortener/internal/models"
)

type urlStatisticsRepositoryPostgresqlImpl struct {
	db *sql.DB
}

func NewURLStatisticsRepositoryPostgresql(db *sql.DB) URLStatisticsRepository {
	return &urlStatisticsRepositoryPostgresqlImpl{db: db}
}

func (r *urlStatisticsRepositoryPostgresqlImpl) GetURLStatistics(ctx context.Context, shortPath string) (*models.URLStatistics, error) {
	row := r.db.QueryRowContext(ctx, PG_GET_URL_STATISTICS, shortPath)
	var statistics models.URLStatistics
	err := row.Scan(&statistics.Last24Hours, &statistics.PastWeek, &statistics.AllTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrURLStatisticsNotFound
		}
		return nil, ErrInternalServerError
	}
	statistics.ShortPath = shortPath
	return &statistics, nil
}

func (r *urlStatisticsRepositoryPostgresqlImpl) InsertAccessLog(ctx context.Context, shortPath string, accessedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, PG_INSERT_ACCESS_LOG, shortPath, accessedAt)
	return err
}
