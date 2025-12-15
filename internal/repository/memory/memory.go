package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"log"
	"sync"
)

// MemStorage is a struct for in-memory URL storage.
type MemStorage struct {
	cfg        *config.Config
	mu         sync.RWMutex
	urls       map[string]string // {shortID: originalURL}
	urlToShort map[string]string // {originalURL: shortID}
	records    []model.URLStorage
	storage    IStorage
}

// New creates a new MemStorage.
func New(cfg *config.Config) (*MemStorage, error) {
	s := &MemStorage{
		urls:       make(map[string]string),
		urlToShort: make(map[string]string),
		records:    make([]model.URLStorage, 0),
		cfg:        cfg,
		storage:    NewFileStorage(cfg),
	}

	if cfg.FileStoragePath != "" {
		if err := s.loadRecords(); err != nil {
			return nil, fmt.Errorf("failed to load records: %v", err)
		}
	}
	return s, nil
}

// loadRecords loads records from a file.
func (m *MemStorage) loadRecords() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := m.storage.GetDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to load data from file: %v", err)
	}

	if len(data) != 0 {
		if err = json.Unmarshal(data, &m.records); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	for _, record := range m.records {
		m.urls[record.ShortURL] = record.OriginalURL
		m.urlToShort[record.OriginalURL] = record.ShortURL
	}

	log.Printf("Loaded %d URLs from storage file", len(m.records))
	return nil
}

// GetShortURL retrieves a short URL from memory.
func (m *MemStorage) GetShortURL(ctx context.Context, originalURL string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if v, ok := m.urlToShort[originalURL]; ok {
		return v, true
	}

	return "", false
}

// StorageURL stores a URL in memory.
func (m *MemStorage) StorageURL(ctx context.Context, originalURL string, shortID string, userID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := model.URLStorage{
		UUID:        shortID,
		ShortURL:    shortID,
		OriginalURL: originalURL,
		UserID:      userID,
	}

	m.records = append(m.records, s)
	if m.cfg.FileStoragePath != "" {
		newData, err := json.MarshalIndent(m.records, "", "   ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal data: %v", err)
		}

		if err = m.storage.SaveDataToFile(newData); err != nil {
			return "", fmt.Errorf("failed to save data to file: %v", err)
		}
	}

	m.urls[shortID] = originalURL
	m.urlToShort[originalURL] = shortID

	return shortID, nil
}

// GetOriginalURL retrieves an original URL from memory.
func (m *MemStorage) GetOriginalURL(ctx context.Context, shortID string) (*model.URLStorage, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if v, ok := m.urls[shortID]; ok {
		return &model.URLStorage{
			OriginalURL: v,
		}, true
	}
	return nil, false
}

// Ping checks the database connection.
func (m *MemStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("db not ready")
}

// StoreMultiURL stores multiple URLs in memory.
func (m *MemStorage) StoreMultiURL(ctx context.Context, req []model.URLStorage) ([]model.MultiResp, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	resp := make([]model.MultiResp, 0)
	m.records = append(m.records, req...)
	if m.cfg.FileStoragePath != "" {
		newData, err := json.MarshalIndent(m.records, "", "   ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data: %v", err)
		}

		if err = m.storage.SaveDataToFile(newData); err != nil {
			return nil, fmt.Errorf("failed to save data to file: %v", err)
		}
	}

	for _, r := range req {
		m.urls[r.ShortURL] = r.OriginalURL
		m.urlToShort[r.OriginalURL] = r.ShortURL

		resp = append(resp, model.MultiResp{
			CorrID:   r.UUID,
			ShortURL: r.ShortURL,
		})
	}

	return resp, nil
}

// FindAllByUserID finds all URLs for a user.
func (m *MemStorage) FindAllByUserID(ctx context.Context, userID string) ([]model.RespUserURL, error) {
	userURLs := make([]model.RespUserURL, 0)
	for _, v := range m.records {
		if v.UserID == userID {
			userURLs = append(userURLs, model.RespUserURL{
				ShortURL:    v.ShortURL,
				OriginalURL: v.OriginalURL,
			})
		}
	}

	return userURLs, nil
}

// DeleteRecords deletes records from memory.
func (m *MemStorage) DeleteRecords(ids []string, userID string) error {
	return nil
}

// Close closes memory storage.
func (m *MemStorage) Close() {}
