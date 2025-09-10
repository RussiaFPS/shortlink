package handler

import (
	"encoding/json"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/logger"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"log"
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
	PingDB(c *gin.Context)
	ShorterMulti(c *gin.Context)
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
	h.mux.POST("/api/shorten/batch", h.ShorterMulti)
	h.mux.POST("/api/shorten", h.GetAPIShortURL)
	h.mux.POST("/", h.GetShortURL)
	h.mux.GET("/:id", h.GetOriginURL)
	h.mux.GET("/ping", h.PingDB)

	return h
}

func (h *URLShortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := h.mux.(http.Handler); ok {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *URLShortenerHandler) ShorterMulti(c *gin.Context) {
	buffer := make([]model.MultiReq, 0)

	if err := c.BindJSON(&buffer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.s.ShorterMulti(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *URLShortenerHandler) PingDB(c *gin.Context) {
	if err := h.s.PingDB(); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (h *URLShortenerHandler) GetAPIShortURL(c *gin.Context) {
	req, resp := model.RequestGetAPIShortURL{}, model.ResponseGetAPIShortURL{}

	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err = url.ParseRequestURI(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.s.Shorten(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp.Result = result
	c.JSON(http.StatusCreated, resp)
}

func (h *URLShortenerHandler) GetShortURL(c *gin.Context) {
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	originalURL := strings.TrimSpace(string(body))
	if _, err = url.ParseRequestURI(originalURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL, err := h.s.Shorten(originalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, shortURL)
}

func (h *URLShortenerHandler) GetOriginURL(c *gin.Context) {
	id := c.Param("id")

	originalURL, exists := h.s.GetOriginal(id)
	if !exists {
		log.Printf("not exists: %s", id)
		c.Status(http.StatusNotFound)
		return
	}

	log.Printf("find shortId: %v and redirectURL: %v", id, originalURL)

	c.Writer.Header().Set("Location", originalURL)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
	//c.Redirect(http.StatusTemporaryRedirect, originalURL)
}
