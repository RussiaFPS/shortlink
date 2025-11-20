package memory

import (
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"os"
)

// IStorage is an interface for file storage.
type IStorage interface {
	// GetDataFromFile reads data from the storage file.
	GetDataFromFile() ([]byte, error)
	// SaveDataToFile saves data to the storage file.
	SaveDataToFile(data []byte) error
}

// FileStorage is a file storage implementation.
type FileStorage struct {
	cfg *config.Config
}

// NewFileStorage creates a new file storage.
func NewFileStorage(cfg *config.Config) IStorage {
	return &FileStorage{cfg: cfg}
}

// GetDataFromFile reads data from the storage file.
func (f *FileStorage) GetDataFromFile() ([]byte, error) {
	file, err := os.OpenFile(f.cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage file: %v", err)
	}
	defer file.Close()

	return os.ReadFile(f.cfg.FileStoragePath)
}

// SaveDataToFile saves data to the storage file.
func (f *FileStorage) SaveDataToFile(data []byte) error {
	file, err := os.OpenFile(f.cfg.FileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open storage file: %v", err)
	}
	defer file.Close()

	if err = os.WriteFile(f.cfg.FileStoragePath, data, 0666); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}
