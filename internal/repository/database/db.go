package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/config/db/postgres"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type DBStorage struct {
	config  *config.Config
	pgxPool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*DBStorage, error) {
	dbConn, err := postgres.NewPostgres(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dbConn: %v", err)
	}

	if err = initTable(ctx, dbConn); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &DBStorage{
		config:  cfg,
		pgxPool: dbConn,
	}, nil
}

func initTable(ctx context.Context, db *pgxpool.Pool) error {
	if _, err := db.Exec(ctx, createTable); err != nil {
		return err
	}
	return nil
}

func (db *DBStorage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	return db.pgxPool.Ping(ctx)
}

func (db *DBStorage) GetShortURL(ctx context.Context, originalURL string) (string, bool) {
	var ur string

	log.Printf("DBStorage:GetShortURL with originalURL: %s", originalURL)

	err := db.pgxPool.QueryRow(ctx, findShortURL, originalURL).Scan(&ur)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", false
	}
	if err != nil {
		log.Fatalf("failed to fetch short url: %v", err)
	}
	return ur, true
}

func (db *DBStorage) StorageURL(ctx context.Context, originalURL string, shortID string, userID string) (string, error) {
	var shortURL string
	log.Printf("DBStorage:StorageURL with originalURL: %s,shortID: %s", originalURL, shortID)

	err := db.pgxPool.QueryRow(ctx, addURL, shortID, originalURL, userID).Scan(&shortURL)
	if err != nil {
		return "", fmt.Errorf("failed to add short url to db: %v", err)
	}
	return shortURL, nil
}

func (db *DBStorage) GetOriginalURL(ctx context.Context, shortID string) (string, bool) {
	var ur string

	log.Printf("DBStorage:GetOriginalURL with shortID: %s", shortID)

	err := db.pgxPool.QueryRow(ctx, findLongURL, shortID).Scan(&ur)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", false
	}
	if err != nil {
		log.Fatalf("failed to fetch long url: %v", err)
	}
	return ur, true
}

func (db *DBStorage) StoreMultiURL(ctx context.Context, req []model.URLStorage) ([]model.MultiResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := db.pgxPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]model.MultiResp, 0)
	for _, v := range req {
		var short string

		if err = tx.QueryRow(ctx, addURL, v.ShortURL, v.OriginalURL, v.UserID).Scan(&short); err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		response = append(response, model.MultiResp{
			CorrID:   v.UUID,
			ShortURL: short,
		})
	}

	if err = tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	log.Printf("Resp multi add: %#v\n , with req: %#v\n", response, req)
	return response, nil
}

func (db *DBStorage) FindAllByUserID(ctx context.Context, userID string) ([]model.RespUserURL, error) {
	rows, err := db.pgxPool.Query(ctx, selectFindAllByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userURLs []model.RespUserURL
	for rows.Next() {
		u := model.RespUserURL{}
		if err = rows.Scan(&u.ShortURL, &u.OriginalURL); err != nil {
			return nil, err
		}
		userURLs = append(userURLs, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userURLs, nil
}
