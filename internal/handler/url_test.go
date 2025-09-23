package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const lenShortURL = 8

func TestGetAPIShortURL(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           model.RequestGetAPIShortURL
		expectedStatus int
	}{
		{
			name:           "Valid POST request",
			method:         http.MethodPost,
			body:           model.RequestGetAPIShortURL{URL: "https://practicum.yandex.ru"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid method (GET)",
			method:         http.MethodGet,
			body:           model.RequestGetAPIShortURL{URL: "https://practicum.yandex.ru"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty URL",
			method:         http.MethodPost,
			body:           model.RequestGetAPIShortURL{URL: ""},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid URL",
			method:         http.MethodPost,
			body:           model.RequestGetAPIShortURL{URL: "not-a-valid-url"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalFlags := flag.CommandLine
			flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
			defer func() {
				flag.CommandLine = originalFlags
			}()
			os.Args = []string{"test"}

			router := gin.Default()
			cfg := config.NewConfig()
			rep := repository.New(context.Background(), cfg)
			ser := service.NewURLShortener(cfg, lenShortURL, rep)
			NewURLShortenerHandler(router, cfg, ser)
			srv := &http.Server{
				Addr:    cfg.ServerAddr,
				Handler: router.Handler(),
			}
			log.Printf("Server started at %v", srv.Addr)

			resp := model.ResponseGetAPIShortURL{}

			jsonBody, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("failed marshaled body: %v", err)
			}

			req := httptest.NewRequest(tt.method, "/api/shorten", bytes.NewReader(jsonBody))

			rr := httptest.NewRecorder()
			srv.Handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				if err = json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed unmarshaled body: %v", err)
				}

				if !strings.HasPrefix(resp.Result, cfg.BaseURL) {
					t.Errorf("expected shortened URL, got %s", rr.Body.String())
				}
			}
		})
	}
}

func TestGetShortURL(t *testing.T) {
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
			originalFlags := flag.CommandLine
			flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
			defer func() {
				flag.CommandLine = originalFlags
			}()
			os.Args = []string{"test"}

			router := gin.Default()
			cfg := config.NewConfig()
			rep := repository.New(context.Background(), cfg)
			ser := service.NewURLShortener(cfg, lenShortURL, rep)
			NewURLShortenerHandler(router, cfg, ser)
			srv := &http.Server{
				Addr:    cfg.ServerAddr,
				Handler: router.Handler(),
			}
			log.Printf("Server started at %v", srv.Addr)

			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))

			rr := httptest.NewRecorder()
			srv.Handler.ServeHTTP(rr, req)

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
	originalFlags := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	defer func() {
		flag.CommandLine = originalFlags
	}()
	os.Args = []string{"test"}

	router := gin.Default()
	cfg := config.NewConfig()
	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, lenShortURL, rep)
	NewURLShortenerHandler(router, cfg, ser)
	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router.Handler(),
	}
	log.Printf("Server started at %v", srv.Addr)
	originalURL := "https://practicum.yandex.ru/"

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalURL))
	rr := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rr, req)

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

			srv.Handler.ServeHTTP(rr, req)

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
