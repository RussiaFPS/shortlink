package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	err := NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, Log)
}

func TestReqLogger(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	Log = zap.New(core)

	r := gin.New()
	r.Use(ReqLogger())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, len(recorded.All()))

	log := recorded.All()[0]
	assert.Equal(t, "got incoming HTTP request", log.Message)

	fields := log.ContextMap()
	assert.Equal(t, http.MethodGet, fields["method"])
	assert.Equal(t, "", fields["uri"])
	assert.Equal(t, int64(http.StatusOK), fields["status"])
	assert.Equal(t, int64(2), fields["size"])
	_, ok := fields["time"].(time.Duration)
	assert.True(t, ok)
}
