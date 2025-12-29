package service

import (
	"context"
	"errors"
	"github.com/RussiaFPS/shortlink/internal/audit"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Close() {
	m.Called()
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

func (m *MockURLRepository) GetStats(ctx context.Context) *model.Stats {
	args := m.Called(ctx)
	return args.Get(0).(*model.Stats)
}

func (m *MockURLRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockAuditService struct {
	mock.Mock
}

func (m *MockAuditService) NotifyAll(eventType, userID, data string) {
	m.Called(eventType, userID, data)
}

func (m *MockAuditService) Register(observer audit.Observer) {
	m.Called(observer)
}

func TestURLShortenerService_Shorten(t *testing.T) {
	t.Run("shorten new url", func(t *testing.T) {
		mockRepo := new(MockURLRepository)
		mockAudit := new(MockAuditService)
		cfg := &config.Config{BaseURL: "http://localhost:8080"}
		service := NewURLShortener(cfg, 8, mockRepo, mockAudit)

		ctx := context.Background()
		originalURL := "https://example.com"
		userID := "user1"
		shortID := "shortID"

		mockAudit.On("NotifyAll", "shorten", userID, originalURL).Return()
		mockRepo.On("GetShortURL", ctx, originalURL).Return("", false).Once()
		mockRepo.On("GetOriginalURL", ctx, mock.Anything).Return(nil, false)
		mockRepo.On("StorageURL", ctx, originalURL, mock.Anything, userID).Return(shortID, nil).Once()

		shortURL, exists, err := service.Shorten(ctx, originalURL, userID)
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.Contains(t, shortURL, cfg.BaseURL)
		mockRepo.AssertExpectations(t)
		mockAudit.AssertExpectations(t)
	})

	t.Run("shorten existing url", func(t *testing.T) {
		mockRepo := new(MockURLRepository)
		mockAudit := new(MockAuditService)
		cfg := &config.Config{BaseURL: "http://localhost:8080"}
		service := NewURLShortener(cfg, 8, mockRepo, mockAudit)

		ctx := context.Background()
		originalURL := "https://example.com"
		userID := "user1"
		shortID := "shortID"

		mockAudit.On("NotifyAll", "shorten", userID, originalURL).Return()
		mockRepo.On("GetShortURL", ctx, originalURL).Return(shortID, true).Once()
		shortURL, exists, err := service.Shorten(ctx, originalURL, userID)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, "http://localhost:8080/shortID", shortURL)
		mockRepo.AssertExpectations(t)
		mockAudit.AssertExpectations(t)
	})
}

func TestURLShortenerService_GetOriginal(t *testing.T) {
	t.Run("get original url found", func(t *testing.T) {
		mockRepo := new(MockURLRepository)
		mockAudit := new(MockAuditService)
		service := NewURLShortener(&config.Config{}, 8, mockRepo, mockAudit)
		ctx := context.Background()
		shortID := "shortID"
		userID := "user1"
		storage := &model.URLStorage{OriginalURL: "https://example.com"}

		mockAudit.On("NotifyAll", "follow", userID, storage.OriginalURL).Return()
		mockRepo.On("GetOriginalURL", ctx, shortID).Return(storage, true).Once()
		res, found := service.GetOriginal(ctx, shortID, userID)
		assert.True(t, found)
		assert.Equal(t, storage, res)
		mockRepo.AssertExpectations(t)
		mockAudit.AssertExpectations(t)
	})

	t.Run("get original url not found", func(t *testing.T) {
		mockRepo := new(MockURLRepository)
		mockAudit := new(MockAuditService)
		service := NewURLShortener(&config.Config{}, 8, mockRepo, mockAudit)
		ctx := context.Background()
		shortID := "shortID"
		userID := "user1"

		mockRepo.On("GetOriginalURL", ctx, shortID).Return(nil, false).Once()
		res, found := service.GetOriginal(ctx, shortID, userID)
		assert.False(t, found)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
		mockAudit.AssertNotCalled(t, "NotifyAll", "follow", userID, mock.Anything)
	})
}

func TestURLShortenerService_PingDB(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLShortener(&config.Config{}, 8, mockRepo, nil)
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
	service := NewURLShortener(cfg, 8, mockRepo, nil)
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
	service := NewURLShortener(&config.Config{}, 8, mockRepo, nil)
	urls := []string{"url1", "url2"}
	userID := "user1"

	mockRepo.On("DeleteRecords", urls, userID).Return(nil).Once()
	err := service.DeleteURLs(&urls, userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func BenchmarkRandShortID(b *testing.B) {
	s := &URLShortenerService{
		shortLen: 8,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.randShortID()
	}
}
