package handler

import (
	"encoding/json"
	"github.com/RussiaFPS/shortlink/internal/audit"
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

type URLHandler interface {
	GetAPIShortURL(c *gin.Context)
	GetShortURL(c *gin.Context)
	GetOriginURL(c *gin.Context)
	PingDB(c *gin.Context)
	ShorterMulti(c *gin.Context)
	GetUserURL(c *gin.Context)
	DeleteURLs(c *gin.Context)
}

type URLShortenerHandler struct {
	mux          *gin.Engine
	s            service.URLService
	cfg          *config.Config
	auditService *audit.AuditService
}

func NewURLShortenerHandler(router *gin.Engine, cfg *config.Config, ser service.URLService, auditService *audit.AuditService) URLHandler {
	h := &URLShortenerHandler{
		mux:          router,
		s:            ser,
		cfg:          cfg,
		auditService: auditService,
	}

	h.mux.Use(logger.ReqLogger()).Use(GzipMiddleware()).Use(AuthenticationMiddleware(&cfg.SecretKey))
	h.mux.POST("/api/shorten/batch", h.ShorterMulti)
	h.mux.POST("/api/shorten", h.GetAPIShortURL)
	h.mux.GET("/api/user/urls", h.GetUserURL)
	h.mux.POST("/", h.GetShortURL)
	h.mux.GET("/:id", h.GetOriginURL)
	h.mux.GET("/ping", h.PingDB)
	h.mux.DELETE("/api/user/urls", h.DeleteURLs)

	return h
}

func (h *URLShortenerHandler) GetUserURL(c *gin.Context) {
	uid := c.GetString("uid")
	if uid == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	originalURL, err := h.s.GetShortenedURLByUserID(c, uid)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(originalURL) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, originalURL)
}

func (h *URLShortenerHandler) ShorterMulti(c *gin.Context) {
	buffer := make([]model.MultiReq, 0)

	if err := c.BindJSON(&buffer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.s.ShorterMulti(c, buffer, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *URLShortenerHandler) PingDB(c *gin.Context) {
	if err := h.s.PingDB(c); err != nil {
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

	uid := c.GetString("uid")
	result, ok, err := h.s.Shorten(c, req.URL, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.auditService.NotifyAll("shorten", uid, req.URL)

	if ok {
		resp.Result = result
		c.JSON(http.StatusConflict, resp)
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

	uid := c.GetString("uid")
	shortURL, ok, err := h.s.Shorten(c, originalURL, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.auditService.NotifyAll("shorten", uid, originalURL)

	if ok {
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusConflict, shortURL)
		return
	}

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, shortURL)
}

func (h *URLShortenerHandler) GetOriginURL(c *gin.Context) {
	id := c.Param("id")

	originalURL, exists := h.s.GetOriginal(c, id)
	if !exists {
		log.Printf("not exists: %s", id)
		c.Status(http.StatusNotFound)
		return
	}

	if originalURL.DeletedFlag {
		log.Printf("not exists: %s", id)
		c.AbortWithStatus(http.StatusGone)
		return
	}

	uid := c.GetString("uid")
	h.auditService.NotifyAll("follow", uid, originalURL.OriginalURL)

	log.Printf("find shortId: %v and redirectURL: %v", id, originalURL)
	c.Redirect(http.StatusTemporaryRedirect, originalURL.OriginalURL)
}

func (h *URLShortenerHandler) DeleteURLs(c *gin.Context) {
	uid := c.GetString("uid")
	if uid == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	var req []string
	if err := c.BindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.s.DeleteURLs(&req, uid); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusAccepted, nil)
}
