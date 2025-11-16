package service

import (
	"context"
	"errors"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) GetShortURL(ctx context.Context, originalURL string) (string, bool) {
	args := m.Called(ctx, originalURL)
	return args.String(0), args.Bool(1)
}

func (m *MockURLRepository) GetOriginalURL(ctx context.Context, shortURL string) (*model.URLStorage, bool) {
	args := m.Called(ctx, shortURL)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*model.URLStorage), args.Bool(1)
}

func (m *MockURLRepository) StorageURL(ctx context.Context, originalURL, shortURL, userID string) (string, error) {
	args := m.Called(ctx, originalURL, shortURL, userID)
	return args.String(0), args.Error(1)
}

func (m *MockURLRepository) StoreMultiURL(ctx context.Context, data []model.URLStorage) ([]model.MultiResp, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.MultiResp), args.Error(1)
}

func (m *MockURLRepository) FindAllByUserID(ctx context.Context, userID string) ([]model.RespUserURL, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RespUserURL), args.Error(1)
}

func (m *MockURLRepository) DeleteRecords(urls []string, userID string) error {
	args := m.Called(urls, userID)
	return args.Error(0)
}

func (m *MockURLRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestURLShortenerService_Shorten(t *testing.T) {
	mockRepo := new(MockURLRepository)
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortener(cfg, 8, mockRepo)

	ctx := context.Background()
	originalURL := "https://example.com"
	userID := "user1"
	shortID := "shortID"

	mockRepo.On("GetShortURL", ctx, originalURL).Return("", false).Once()
	mockRepo.On("GetOriginalURL", ctx, mock.Anything).Return(nil, false)
	mockRepo.On("StorageURL", ctx, originalURL, mock.Anything, userID).Return(shortID, nil).Once()

	shortURL, exists, err := service.Shorten(ctx, originalURL, userID)
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Contains(t, shortURL, cfg.BaseURL)
	mockRepo.AssertExpectations(t)

	mockRepo.On("GetShortURL", ctx, originalURL).Return(shortID, true).Once()
	shortURL, exists, err = service.Shorten(ctx, originalURL, userID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "http://localhost:8080/shortID", shortURL)
	mockRepo.AssertExpectations(t)
}

func TestURLShortenerService_GetOriginal(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLShortener(&config.Config{}, 8, mockRepo)
	ctx := context.Background()
	shortID := "shortID"
	storage := &model.URLStorage{OriginalURL: "https://example.com"}

	mockRepo.On("GetOriginalURL", ctx, shortID).Return(storage, true).Once()
	res, found := service.GetOriginal(ctx, shortID)
	assert.True(t, found)
	assert.Equal(t, storage, res)
	mockRepo.AssertExpectations(t)

	mockRepo.On("GetOriginalURL", ctx, shortID).Return(nil, false).Once()
	res, found = service.GetOriginal(ctx, shortID)
	assert.False(t, found)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestURLShortenerService_PingDB(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLShortener(&config.Config{}, 8, mockRepo)
	ctx := context.Background()

	mockRepo.On("Ping", ctx).Return(nil).Once()
	err := service.PingDB(ctx)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	expectedErr := errors.New("db error")
	mockRepo.On("Ping", ctx).Return(expectedErr).Once()
	err = service.PingDB(ctx)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestURLShortenerService_GetShortenedURLByUserID(t *testing.T) {
	mockRepo := new(MockURLRepository)
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortener(cfg, 8, mockRepo)
	ctx := context.Background()
	userID := "user1"
	expected := []model.RespUserURL{
		{ShortURL: "short1", OriginalURL: "https://ex1.com"},
		{ShortURL: "short2", OriginalURL: "https://ex2.com"},
	}

	mockRepo.On("FindAllByUserID", ctx, userID).Return(expected, nil).Once()
	result, err := service.GetShortenedURLByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "http://localhost:8080/short1", result[0].ShortURL)
	mockRepo.AssertExpectations(t)
}

func TestURLShortenerService_DeleteURLs(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLShortener(&config.Config{}, 8, mockRepo)
	urls := []string{"url1", "url2"}
	userID := "user1"

	mockRepo.On("DeleteRecords", urls, userID).Return(nil).Once()
	err := service.DeleteURLs(&urls, userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
