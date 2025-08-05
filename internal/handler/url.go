package handler

import (
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const lenShortURL = 8

type URLShortenerHandler struct {
	mux *gin.Engine
	s   *service.URLShortenerService
	cfg *config.Config
}

func NewURLShortenerHandler(cfg *config.Config) *URLShortenerHandler {
	h := &URLShortenerHandler{
		mux: gin.Default(),
		s:   service.NewURLShortener(cfg, lenShortURL),
		cfg: cfg,
	}

	h.mux.POST("/", h.GetShortURL)
	h.mux.GET("/:id", h.GetOriginURL)

	return h
}

func (h *URLShortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *URLShortenerHandler) GetShortURL(c *gin.Context) {
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	if _, err = url.ParseRequestURI(originalURL); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	shortURL, err := h.s.Shorten(originalURL)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, shortURL)
}

func (h *URLShortenerHandler) GetOriginURL(c *gin.Context) {
	id := c.Param("id")

	originalURL, exists := h.s.GetOriginal(id)
	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Header("Location", originalURL)
	c.Redirect(http.StatusTemporaryRedirect, originalURL)
}
