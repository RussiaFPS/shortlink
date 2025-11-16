package memory

import (
	"context"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestMemStorage_New(t *testing.T) {
	cfg := &config.Config{}
	storage, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	tmpfile, err := ioutil.TempFile("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	cfg.FileStoragePath = tmpfile.Name()
	storage, err = New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}

func TestMemStorage_StorageAndGet(t *testing.T) {
	storage, _ := New(&config.Config{})
	ctx := context.Background()
	originalURL := "https://example.com"
	shortID := "short"
	userID := "user1"

	storedID, err := storage.StorageURL(ctx, originalURL, shortID, userID)
	assert.NoError(t, err)
	assert.Equal(t, shortID, storedID)

	res, ok := storage.GetOriginalURL(ctx, shortID)
	assert.True(t, ok)
	assert.Equal(t, originalURL, res.OriginalURL)

	resShort, ok := storage.GetShortURL(ctx, originalURL)
	assert.True(t, ok)
	assert.Equal(t, shortID, resShort)
}

func TestMemStorage_StoreMultiURL(t *testing.T) {
	storage, _ := New(&config.Config{})
	ctx := context.Background()
	req := []model.URLStorage{
		{UUID: "1", OriginalURL: "https://ex1.com", ShortURL: "s1", UserID: "u1"},
		{UUID: "2", OriginalURL: "https://ex2.com", ShortURL: "s2", UserID: "u1"},
	}

	resp, err := storage.StoreMultiURL(ctx, req)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "s1", resp[0].ShortURL)

	res, ok := storage.GetOriginalURL(ctx, "s1")
	assert.True(t, ok)
	assert.Equal(t, "https://ex1.com", res.OriginalURL)
}

func TestMemStorage_FindAllByUserID(t *testing.T) {
	storage, _ := New(&config.Config{})
	ctx := context.Background()
	storage.StorageURL(ctx, "https://ex1.com", "s1", "u1")
	storage.StorageURL(ctx, "https://ex2.com", "s2", "u2")
	storage.StorageURL(ctx, "https://ex3.com", "s3", "u1")

	res, err := storage.FindAllByUserID(ctx, "u1")
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestMemStorage_Ping(t *testing.T) {
	storage, _ := New(&config.Config{})
	err := storage.Ping(context.Background())
	assert.Error(t, err)
}

func TestMemStorage_loadRecords(t *testing.T) {
	content := `[
		{"uuid": "s1", "short_url": "s1", "original_url": "https://ex1.com", "user_id": "u1"},
		{"uuid": "s2", "short_url": "s2", "original_url": "https://ex2.com", "user_id": "u1"}
	]`
	tmpfile, err := ioutil.TempFile("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	assert.NoError(t, err)
	tmpfile.Close()

	cfg := &config.Config{FileStoragePath: tmpfile.Name()}
	storage, err := New(cfg)
	assert.NoError(t, err)

	res, ok := storage.GetOriginalURL(context.Background(), "s1")
	assert.True(t, ok)
	assert.Equal(t, "https://ex1.com", res.OriginalURL)
}
