package handler

import (
	"github.com/RussiaFPS/shortlink/internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetShortURL(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Errorf("Failed to create config: %v", err)
		return
	}

	h := NewURLShortenerHandler(cfg)

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "Valid POST request",
			method:         http.MethodPost,
			body:           "https://practicum.yandex.ru/",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid method (GET)",
			method:         http.MethodGet,
			body:           "https://practicum.yandex.ru/",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty URL",
			method:         http.MethodPost,
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid URL",
			method:         http.MethodPost,
			body:           "not-a-valid-url",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				if !strings.HasPrefix(rr.Body.String(), cfg.BaseURL) {
					t.Errorf("expected shortened URL, got %s", rr.Body.String())
				}
			}
		})
	}
}

func TestGetOriginURL(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Errorf("Failed to create config: %v", err)
		return
	}

	h := NewURLShortenerHandler(cfg)
	originalURL := "https://practicum.yandex.ru/"

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	shortID := strings.TrimPrefix(rr.Body.String(), cfg.BaseURL+"/")
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedLoc    string
	}{
		{
			name:           "Valid GET request",
			method:         http.MethodGet,
			path:           "/" + shortID,
			expectedStatus: http.StatusTemporaryRedirect,
			expectedLoc:    originalURL,
		},
		{
			name:           "Invalid method (POST)",
			method:         http.MethodPost,
			path:           "/" + shortID,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Non-existent short ID",
			method:         http.MethodGet,
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty short ID",
			method:         http.MethodGet,
			path:           "/",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req = httptest.NewRequest(tt.method, tt.path, nil)
			rr = httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusTemporaryRedirect {
				loc := rr.Header().Get("Location")
				if loc != tt.expectedLoc {
					t.Errorf("expected Location %s, got %s", tt.expectedLoc, loc)
				}
			}
		})
	}
}
