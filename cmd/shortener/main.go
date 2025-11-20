package main

import (
	"context"
	"errors"
	"github.com/RussiaFPS/shortlink/internal/audit"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/handler"
	"github.com/RussiaFPS/shortlink/internal/logger"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const lenShortURL = 8

func main() {
	cfg := config.NewConfig()

	if err := logger.NewLogger(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	auditService := audit.NewAuditService()
	if cfg.AuditFile != "" {
		auditService.Register(audit.NewLogObserver(cfg.AuditFile))
	}
	if cfg.AuditURL != "" {
		auditService.Register(audit.NewHTTPObserver(cfg.AuditURL))
	}

	router := gin.Default()
	rep := repository.New(context.Background(), cfg)
	ser := service.NewURLShortener(cfg, lenShortURL, rep)
	handler.NewURLShortenerHandler(router, cfg, ser, auditService)

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started at %v", cfg.ServerAddr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
