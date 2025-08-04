package main

import (
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/handler"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	h := handler.NewURLShortenerHandler(cfg)

	log.Printf("Server started at %v", cfg.ServerAddr)
	log.Fatal(http.ListenAndServe(cfg.ServerAddr, h))
}
