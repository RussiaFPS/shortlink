package memory

import (
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileStorage_SaveAndGetData(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test_file_storage")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	cfg := &config.Config{FileStoragePath: tmpfile.Name()}
	fileStorage := NewFileStorage(cfg)

	dataToSave := []byte("hello world")
	err = fileStorage.SaveDataToFile(dataToSave)
	assert.NoError(t, err)

	readData, err := fileStorage.GetDataFromFile()
	assert.NoError(t, err)
	assert.Equal(t, dataToSave, readData)
}

func TestFileStorage_GetDataFromNonExistentFile(t *testing.T) {
	cfg := &config.Config{FileStoragePath: "non_existent_file.txt"}
	fileStorage := NewFileStorage(cfg)

	data, err := fileStorage.GetDataFromFile()
	assert.NoError(t, err)
	assert.Empty(t, data)

	os.Remove("non_existent_file.txt")
}
