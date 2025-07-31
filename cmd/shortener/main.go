package main

import (
	"github.com/RussiaFPS/shortlink/internal/handler"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	h := handler.NewURLShortenerHandler()

	mux.HandleFunc("/", h.GetShortURL)
	mux.HandleFunc("/{id}", h.GetOriginURL)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
