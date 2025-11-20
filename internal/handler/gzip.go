package handler

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
)

// GzipWriter is a custom ResponseWriter that compresses the response body using gzip.
type GzipWriter struct {
	gin.ResponseWriter
	ginContext *gin.Context
	writer     *gzip.Writer
	once       sync.Once
}

// initGzip initializes the gzip writer.
func (g *GzipWriter) initGzip() {
	if g.writer == nil {
		g.writer = gzip.NewWriter(g.ResponseWriter)
		g.Header().Set("Content-Encoding", "gzip")
	}
}

// Write writes the compressed data to the response.
func (g *GzipWriter) Write(data []byte) (int, error) {
	if !g.checkContentType(g.Header().Get("Content-Type")) {
		return g.ResponseWriter.Write(data)
	}

	g.initGzip()
	return g.writer.Write(data)
}

// checkContentType checks if the content type is compressible.
func (g *GzipWriter) checkContentType(contentType string) bool {
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/html")
}

// WriteString writes a string to the response.
func (g *GzipWriter) WriteString(s string) (int, error) {
	return g.Write([]byte(s))
}

// Close closes the gzip writer.
func (g *GzipWriter) Close() {
	if g.writer != nil {
		g.writer.Close()
	}
}

// WriteHeader writes the HTTP header.
func (g *GzipWriter) WriteHeader(code int) {
	if g.checkContentType(g.Header().Get("Content-Type")) {
		g.initGzip()
	}
	g.ResponseWriter.WriteHeader(code)
}

// GzipMiddleware is a middleware that handles gzip compression and decompression.
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to decompress request: " + err.Error()})
				return
			}
			defer gz.Close()
			c.Request.Body = gz
		}

		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gzipWriter := &GzipWriter{
			ResponseWriter: c.Writer,
			ginContext:     c,
		}
		c.Writer = gzipWriter

		defer func() {
			gzipWriter.Close()
			c.Writer = gzipWriter.ResponseWriter
		}()

		c.Next()
	}
}
