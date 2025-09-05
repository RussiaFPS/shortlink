package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/config/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type DbStorage struct {
	config  *config.Config
	pgxPool *pgxpool.Pool
}

func New(cfg *config.Config) (*DbStorage, error) {
	dbConn, err := postgres.NewPostgres(context.Background(), cfg.DSN)
	if err != nil {
		log.Printf("failed to initialize dbConn: %v", err)
		return nil, err
	}

	return &DbStorage{
		config:  cfg,
		pgxPool: dbConn,
	}, nil
}

func (db *DbStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return db.pgxPool.Ping(ctx)
}

func (db *DbStorage) GetShortURL(originalURL string) (string, bool) {
	var ur string
	err := db.pgxPool.QueryRow(context.Background(), findShortURL, originalURL).Scan(&ur)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", false
	}
	if err != nil {
		log.Fatalf("failed to fetch short url: %v", err)
	}
	return fmt.Sprintf("%s/%s", db.config.BaseURL, ur), true
}

func (db *DbStorage) StorageURL(originalURL string, shortID string) (string, error) {
	_, err := db.pgxPool.Exec(context.Background(), addURL, shortID, originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to add short url to db: %v", err)
	}
	return fmt.Sprintf("%s/%s", db.config.BaseURL, shortID), nil
}

func (db *DbStorage) GetOriginalURL(shortID string) (string, bool) {
	var ur string
	err := db.pgxPool.QueryRow(context.Background(), findLongURL, shortID).Scan(&ur)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", false
	}
	if err != nil {
		log.Fatalf("failed to fetch long url: %v", err)
	}
	return ur, true
}
