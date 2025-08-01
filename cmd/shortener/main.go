package main

import (
	"github.com/RussiaFPS/shortlink/internal/handler"
	"log"
	"net/http"
)

func main() {
	h := handler.NewURLShortenerHandler()

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
