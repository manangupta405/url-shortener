package repositories

import "errors"

const (
	PG_GET_BY_SHORT_URL    = `SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by FROM urls WHERE short_path = $1`
	PG_GET_BY_ORIGINAL_URL = `SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by FROM urls WHERE original_url = $1`
	PG_INSERT_SHORT_URL    = `INSERT INTO urls (short_path, original_url, expiry, created_at, created_by) VALUES ($1, $2, $3, $4, $5)`
	PG_UPDATE_SHORT_URL    = `UPDATE urls SET original_url = $1, expiry = $2, modified_at = $3, modified_by = $4 WHERE short_path = $5`
	PG_DELETE_SHORT_URL    = `DELETE FROM urls WHERE short_path = $1`
	PG_INSERT_URL_ARCHIVE  = `INSERT INTO urls_archive (short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	PG_GET_URL_STATISTICS = `SELECT COUNT(*) FILTER (WHERE accessed_at >= NOW() - INTERVAL '24 hours') AS last_24_hours,
    								COUNT(*) FILTER (WHERE accessed_at >= NOW() - INTERVAL '7 days') AS past_week,
    								COUNT(*) AS all_time
							FROM url_access_logs
							WHERE short_path = $1;`
	PG_INSERT_ACCESS_LOG = `INSERT INTO url_access_logs (short_path,accessed_at) VALUES ($1,$2);`
)

var (
	ErrInvalidShortURL          = errors.New("invalid short url")
	ErrShortURLAlreadyExists    = errors.New("short url already exists")
	ErrOriginalURLAlreadyExists = errors.New("original url already exists")
	ErrURLNotFound              = errors.New("url not found")
	ErrInternalServerError      = errors.New("internal server error")
	ErrRequestTimeout           = errors.New("request timeout")
	ErrInvalidPayload           = errors.New("invalid payload")
	ErrDBError                  = errors.New("database error")
	ErrRedisError               = errors.New("redis error")
	ErrCacheMiss                = errors.New("cache miss")
	ErrShortURLNotFound         = errors.New("short url not found")
	ErrURLStatisticsNotFound    = errors.New("url statistics not found")
	ErrURLExpired               = errors.New("url expired")
)
