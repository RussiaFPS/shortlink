package main

import (
	"github.com/RussiaFPS/shortlink/internal/handler"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	h := handler.NewURLShortenerHandler()

	mux.HandleFunc("/", h.GetShortUrl)
	mux.HandleFunc("/{id}", h.GetOriginUrl)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
