package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/audit"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"log"
	"math/big"
	"strings"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// URLService is an interface for URL shortening services.
type URLService interface {
	Shorten(ctx context.Context, originalURL string, userID string) (string, bool, error)
	GetOriginal(ctx context.Context, shortURL string, userID string) (*model.URLStorage, bool)
	ShorterMulti(ctx context.Context, req []model.MultiReq, userID string) ([]model.MultiResp, error)
	PingDB(ctx context.Context) error
	randShortID() string
	generateShortID(ctx context.Context) string
	GetShortenedURLByUserID(ctx context.Context, userID string) ([]model.RespUserURL, error)
	DeleteURLs(req *[]string, userID string) error
}

// URLShortenerService is a struct that implements the URLService interface.
type URLShortenerService struct {
	cfg          *config.Config
	shortLen     int
	r            repository.URLRepository
	auditService audit.AuditService
}

// NewURLShortener is a constructor for URLShortenerService.
func NewURLShortener(cfg *config.Config, shortLen int, rep repository.URLRepository, auditService audit.AuditService) URLService {
	return &URLShortenerService{
		shortLen:     shortLen,
		cfg:          cfg,
		r:            rep,
		auditService: auditService,
	}
}

// PingDB checks the database connection.
func (s *URLShortenerService) PingDB(ctx context.Context) error {
	return s.r.Ping(ctx)
}

// generateShortID generates a random short ID and ensures it's unique.
func (s *URLShortenerService) generateShortID(ctx context.Context) string {
	var id string
	for {
		id = s.randShortID()
		if _, exists := s.r.GetOriginalURL(ctx, id); !exists {
			break
		}
	}
	return id
}

// randShortID generates a random short ID.
func (s *URLShortenerService) randShortID() string {
	id := make([]byte, s.shortLen)
	numChars := big.NewInt(int64(len(chars)))
	for i := range id {
		n, err := rand.Int(rand.Reader, numChars)
		if err != nil {
			log.Panicf("failed to generate random number for short ID: %v", err)
		}
		id[i] = chars[n.Int64()]
	}

	return string(id)
}

// Shorten shortens a URL.
func (s *URLShortenerService) Shorten(ctx context.Context, originalURL string, userID string) (string, bool, error) {
	log.Printf("URLShortenerService:Shorten with originalURL: %s", originalURL)

	s.auditService.NotifyAll("shorten", userID, originalURL)

	if shortURL, ok := s.r.GetShortURL(ctx, originalURL); ok {
		return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortURL), true, nil
	}

	URL, err := s.r.StorageURL(ctx, originalURL, s.generateShortID(ctx), userID)
	if err != nil {
		return "", false, err
	}

	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, URL), false, nil
}

// GetOriginal retrieves the original URL from a short ID.
func (s *URLShortenerService) GetOriginal(ctx context.Context, shortID string, userID string) (*model.URLStorage, bool) {
	log.Printf("URLShortenerService:GetOriginal with shortID: %s", shortID)
	url, ok := s.r.GetOriginalURL(ctx, shortID)
	if ok {
		s.auditService.NotifyAll("follow", userID, url.OriginalURL)
	}
	return url, ok
}

// ShorterMulti shortens multiple URLs at once.
func (s *URLShortenerService) ShorterMulti(ctx context.Context, req []model.MultiReq, userID string) ([]model.MultiResp, error) {
	data, oldData := make([]model.URLStorage, 0), make([]model.MultiResp, 0)

	for _, v := range req {
		d := model.URLStorage{
			UUID:        v.CorrID,
			OriginalURL: v.URL,
			UserID:      userID,
		}

		if shortURL, ok := s.r.GetShortURL(ctx, v.URL); ok {
			oldData = append(oldData, model.MultiResp{CorrID: v.CorrID, ShortURL: shortURL})
		} else {
			d.ShortURL = s.generateShortID(ctx)
			data = append(data, d)
		}
	}

	resp, err := s.r.StoreMultiURL(ctx, data)
	if err != nil {
		return nil, err
	}

	resp = append(resp, oldData...)
	for i := range resp {
		if !strings.Contains(resp[i].ShortURL, "http") {
			resp[i].ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, resp[i].ShortURL)
		}
	}

	return resp, nil
}

// GetShortenedURLByUserID retrieves all URLs for a user.
func (s *URLShortenerService) GetShortenedURLByUserID(ctx context.Context, userID string) ([]model.RespUserURL, error) {
	data, err := s.r.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range data {
		if !strings.Contains(data[i].ShortURL, "http") {
			data[i].ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, data[i].ShortURL)
		}
	}
	return data, nil
}

// DeleteURLs deletes multiple URLs for a user.
func (s *URLShortenerService) DeleteURLs(req *[]string, userID string) error {
	return s.r.DeleteRecords(*req, userID)
}
