package repositories

const (
	PG_GET_BY_SHORT_URL      = `SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by FROM urls WHERE short_path = $1`
	PG_GET_BY_ORIGINAL_URL   = `SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by FROM urls WHERE original_url = $1`
	PG_INSERT_SHORT_URL      = `INSERT INTO urls (short_path, original_url, expiry, created_at, created_by) VALUES ($1, $2, $3, $4, $5)`
	PG_UPDATE_SHORT_URL      = `UPDATE urls SET original_url = $1, expiry = $2, modified_at = $3, modified_by = $4 WHERE short_path = $5`
	PG_SOFT_DELETE_SHORT_URL = `UPDATE urls SET deleted_at = $1, deleted_by = $2 WHERE short_path = $3`
)
