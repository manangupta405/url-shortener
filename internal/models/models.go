package models

import "time"

/*
CREATE TABLE IF NOT EXISTS urls (

	short_path VARCHAR(255) PRIMARY KEY NOT NULL,
	original_url TEXT NOT NULL,
	expiry TIMESTAMP,
	created_at TIMESTAMP NOT NULL,
	created_by VARCHAR(255) NOT NULL,
	modified_at TIMESTAMP,
	modified_by VARCHAR(255),
	deleted_at TIMESTAMP,
	deleted_by VARCHAR(255)

);
*/
type URL struct {
	ShortPath   string     `json:"short_path"`
	OriginalURL string     `json:"original_url"`
	Expiry      *time.Time `json:"expiry"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   string     `json:"created_by"`
	ModifiedAt  *time.Time `json:"modified_at"`
	ModifiedBy  *string    `json:"modified_by"`
	DeletedAt   *time.Time `json:"deleted_at"`
	DeletedBy   *string    `json:"deleted_by"`
}
