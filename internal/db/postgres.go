package db

import (
	"context"
	"database/sql"
	"fmt"
	"url-shortener/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(config *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName))
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}
	return db, nil
}
