package main

import (
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/handler"
	"github.com/RussiaFPS/shortlink/internal/logger"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()

	if err := logger.NewLogger(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	h := handler.NewURLShortenerHandler(cfg)

	log.Printf("Server started at %v", cfg.ServerAddr)
	log.Fatal(http.ListenAndServe(cfg.ServerAddr, h))
}
