package memory

import (
	"encoding/json"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"log"
	"sync"
)

type MemStorage struct {
	cfg        *config.Config
	mu         sync.RWMutex
	urls       map[string]string // {shortID: originalURL}
	urlToShort map[string]string // {originalURL: shortID}
	records    []model.URLStorage
	storage    IStorage
}

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

func (m *MemStorage) GetShortURL(originalURL string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if v, ok := m.urlToShort[originalURL]; ok {
		return fmt.Sprintf("%s/%s", m.cfg.BaseURL, v), true
	}

	return "", false
}

func (m *MemStorage) StorageURL(originalURL string, shortID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := model.URLStorage{
		UUID:        shortID,
		ShortURL:    shortID,
		OriginalURL: originalURL,
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

	return fmt.Sprintf("%s/%s", m.cfg.BaseURL, shortID), nil
}

func (m *MemStorage) GetOriginalURL(shortID string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if v, ok := m.urls[shortID]; ok {
		return v, true
	}
	return "", false
}

func (m *MemStorage) Ping() error {
	return fmt.Errorf("db not ready")
}
