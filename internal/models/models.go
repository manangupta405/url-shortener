package models

import "time"

type URL struct {
	ShortPath   string    `json:"short_path"`
	OriginalURL string    `json:"original_url"`
	Expiry      time.Time `json:"expiry"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	ModifiedAt  time.Time `json:"modified_at"`
	ModifiedBy  string    `json:"modified_by"`
	DeletedAt   time.Time `json:"deleted_at"`
	DeletedBy   string    `json:"deleted_by"`
}
