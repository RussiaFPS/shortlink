package memory

import (
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"os"
)

type IStorage interface {
	GetDataFromFile() ([]byte, error)
	SaveDataToFile(data []byte) error
}

type FileStorage struct {
	cfg *config.Config
}

func NewFileStorage(cfg *config.Config) IStorage {
	return &FileStorage{cfg: cfg}
}

func (f *FileStorage) GetDataFromFile() ([]byte, error) {
	file, err := os.OpenFile(f.cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage file: %v", err)
	}
	defer file.Close()

	return os.ReadFile(f.cfg.FileStoragePath)
}

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
