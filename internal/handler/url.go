package handler

import (
	"github.com/RussiaFPS/shortlink/internal/service"
	"io"
	"net/http"
	"strings"
)

const lenShortURL = 8

type URLShortenerHandler struct {
	mux *http.ServeMux
	s   *service.URLShortenerService
}

func NewURLShortenerHandler() *URLShortenerHandler {
	h := &URLShortenerHandler{
		mux: http.NewServeMux(),
		s:   service.NewURLShortener(lenShortURL),
	}

	h.mux.HandleFunc("/", h.GetShortURL)
	h.mux.HandleFunc("/{id}", h.GetOriginURL)

	return h
}

func (h *URLShortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *URLShortenerHandler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	shortURL, err := h.s.Shorten(originalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *URLShortenerHandler) GetOriginURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing short URL ID", http.StatusBadRequest)
		return
	}

	originalURL, exists := h.s.GetOriginal(id)
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
