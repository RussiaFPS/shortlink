package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func ExampleURLShortenerHandler_GetAPIShortURL() {
	cfg := config.NewConfig()
	gin.SetMode(gin.TestMode)

	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, 8, rep)
	r := gin.New()
	NewURLShortenerHandler(r, cfg, ser, nil)

	reqBody, _ := json.Marshal(model.RequestGetAPIShortURL{URL: "https://example.com"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(reqBody))

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 201
}

func ExampleURLShortenerHandler_GetShortURL() {
	cfg := config.NewConfig()
	gin.SetMode(gin.TestMode)

	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, 8, rep)
	r := gin.New()
	NewURLShortenerHandler(r, cfg, ser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("https://example.com")))

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 201
}

func ExampleURLShortenerHandler_GetOriginURL() {
	cfg := config.NewConfig()
	gin.SetMode(gin.TestMode)

	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, 8, rep)
	r := gin.New()
	NewURLShortenerHandler(r, cfg, ser, nil)

	shortURL, _, _ := ser.Shorten(context.Background(), "https://example.com", "")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, shortURL, nil)

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 307
}

func ExampleURLShortenerHandler_PingDB() {
	cfg := config.NewConfig()
	gin.SetMode(gin.TestMode)

	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, 8, rep)
	r := gin.New()
	NewURLShortenerHandler(r, cfg, ser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 200
}

func ExampleURLShortenerHandler_ShorterMulti() {
	cfg := config.NewConfig()
	gin.SetMode(gin.TestMode)

	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, 8, rep)
	r := gin.New()
	NewURLShortenerHandler(r, cfg, ser, nil)

	reqBody, _ := json.Marshal([]model.MultiReq{
		{CorrID: "1", URL: "https://example.com"},
		{CorrID: "2", URL: "https://example.org"},
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(reqBody))

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 201
}

func ExampleURLShortenerHandler_GetUserURL() {
	cfg := config.NewConfig()
	cfg.SecretKey = "test"
	gin.SetMode(gin.TestMode)

	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, 8, rep)
	r := gin.New()
	NewURLShortenerHandler(r, cfg, ser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/urls", nil)

	// Это пример, поэтому мы можем просто установить заголовок напрямую
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJ0ZXN0In0.pA_HS224b61-4HIp3a_2b-zo3a-v2w-pA_HS224b61")

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 204
}

func ExampleURLShortenerHandler_DeleteURLs() {
	cfg := config.NewConfig()
	cfg.SecretKey = "test"
	gin.SetMode(gin.TestMode)

	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, 8, rep)
	r := gin.New()
	NewURLShortenerHandler(r, cfg, ser, nil)

	shortURL, _, _ := ser.Shorten(context.Background(), "https://example.com", "test")

	reqBody, _ := json.Marshal([]string{shortURL})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(reqBody))

	// Это пример, поэтому мы можем просто установить заголовок напрямую
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJ0ZXN0In0.pA_HS224b61-4HIp3a_2b-zo3a-v2w-pA_HS224b61")

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 202
}
