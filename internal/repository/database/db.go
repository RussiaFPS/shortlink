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

type DBStorage struct {
	config  *config.Config
	pgxPool *pgxpool.Pool
}

func New(cfg *config.Config) (*DBStorage, error) {
	dbConn, err := postgres.NewPostgres(context.Background(), cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dbConn: %v", err)
	}

	if err = initTable(dbConn); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &DBStorage{
		config:  cfg,
		pgxPool: dbConn,
	}, nil
}

func initTable(db *pgxpool.Pool) error {
	if _, err := db.Exec(context.Background(), createTable); err != nil {
		return err
	}
	return nil
}

func (db *DBStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return db.pgxPool.Ping(ctx)
}

func (db *DBStorage) GetShortURL(originalURL string) (string, bool) {
	var ur string

	log.Printf("DBStorage:GetShortURL with originalURL: %s", originalURL)

	err := db.pgxPool.QueryRow(context.Background(), findShortURL, originalURL).Scan(&ur)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", false
	}
	if err != nil {
		log.Fatalf("failed to fetch short url: %v", err)
	}
	return fmt.Sprintf("%s/%s", db.config.BaseURL, ur), true
}

func (db *DBStorage) StorageURL(originalURL string, shortID string) (string, error) {
	log.Printf("DBStorage:StorageURL with originalURL: %s,shortID: %s", originalURL, shortID)

	_, err := db.pgxPool.Exec(context.Background(), addURL, shortID, originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to add short url to db: %v", err)
	}
	return fmt.Sprintf("%s/%s", db.config.BaseURL, shortID), nil
}

func (db *DBStorage) GetOriginalURL(shortID string) (string, bool) {
	var ur string

	log.Printf("DBStorage:GetOriginalURL with shortID: %s", shortID)

	err := db.pgxPool.QueryRow(context.Background(), findLongURL, shortID).Scan(&ur)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", false
	}
	if err != nil {
		log.Fatalf("failed to fetch long url: %v", err)
	}
	return ur, true
}
