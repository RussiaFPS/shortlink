package handler_test

import (
	"context"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/audit"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/handler"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Example demonstrates the usage of the GetShortURL endpoint.
func ExampleURLShortenerHandler_GetShortURL() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	cfg := config.NewConfig()
	cfg.BaseURL = "http://localhost:8080"
	cfg.FileStoragePath = ""

	rep := repository.New(context.Background(), cfg)
	auditService := audit.NewAuditService()
	ser := service.NewURLShortener(cfg, 8, rep, auditService)
	handler.NewURLShortenerHandler(router, cfg, ser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader("https://yandex.ru"))
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(strings.HasPrefix(w.Body.String(), "http://localhost:8080/"))

	// Output:
	// 201
	// true
}
