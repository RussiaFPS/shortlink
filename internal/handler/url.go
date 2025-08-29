package handler

import (
	"encoding/json"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/logger"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const lenShortURL = 8

type IURLShortenerHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	GetAPIShortURL(c *gin.Context)
	GetShortURL(c *gin.Context)
	GetOriginURL(c *gin.Context)
}

type URLShortenerHandler struct {
	mux gin.IRouter
	s   service.IURLShortenerService
	cfg *config.Config
}

func NewURLShortenerHandler(cfg *config.Config) IURLShortenerHandler {
	h := &URLShortenerHandler{
		mux: gin.Default(),
		s:   service.NewURLShortener(cfg, lenShortURL),
		cfg: cfg,
	}

	h.mux.Use(logger.ReqLogger()).Use(GzipMiddleware())
	h.mux.POST("/api/shorten", h.GetAPIShortURL)
	h.mux.POST("/", h.GetShortURL)
	h.mux.GET("/:id", h.GetOriginURL)

	return h
}

func (h *URLShortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := h.mux.(http.Handler); ok {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *URLShortenerHandler) GetAPIShortURL(c *gin.Context) {
	req, resp := model.RequestGetAPIShortURL{}, model.ResponseGetAPIShortURL{}

	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err = url.ParseRequestURI(req.URL); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.s.Shorten(req.URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp.Result = result
	c.JSON(http.StatusCreated, resp)
}

func (h *URLShortenerHandler) GetShortURL(c *gin.Context) {
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	originalURL := strings.TrimSpace(string(body))
	if _, err = url.ParseRequestURI(originalURL); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL, err := h.s.Shorten(originalURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
