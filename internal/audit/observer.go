package audit

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Event is a struct that represents an audit event.
type Event struct {
	Timestamp int64  `json:"ts"`
	Action    string `json:"action"`
	UserID    string `json:"user_id"`
	URL       string `json:"url"`
}

// Observer is an interface for an audit observer.
type Observer interface {
	Notify(event Event)
}

// AuditService is an interface for an audit service.
type AuditService interface {
	Register(observer Observer)
	NotifyAll(action, userID, url string)
}

// LogObserver is a struct that implements the Observer interface and logs events to a file.
type LogObserver struct {
	filePath string
}

// NewLogObserver creates a new LogObserver.
func NewLogObserver(filePath string) *LogObserver {
	return &LogObserver{filePath: filePath}
}

// Notify logs an audit event to a file.
func (o *LogObserver) Notify(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal audit event: %v", err)
		return
	}

	f, err := os.OpenFile(o.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("failed to open audit file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		log.Printf("failed to write to audit file: %v", err)
	}
}

// HTTPObserver is a struct that implements the Observer interface and sends events to a URL.
type HTTPObserver struct {
	url string
}

// NewHTTPObserver creates a new HTTPObserver.
func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{url: url}
}

// Notify sends an audit event to a URL.
func (o *HTTPObserver) Notify(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal audit event: %v", err)
		return
	}

	resp, err := http.Post(o.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("failed to send audit event: %v", err)
	}
	defer resp.Body.Close()
}

// ServAudit is a struct that manages audit observers.
type ServAudit struct {
	observers []Observer
}

// NewAuditService creates a new ServAudit.
func NewAuditService() AuditService {
	return &ServAudit{observers: []Observer{}}
}

// Register registers a new observer.
func (s *ServAudit) Register(observer Observer) {
	s.observers = append(s.observers, observer)
}

// NotifyAll notifies all registered observers of an event.
func (s *ServAudit) NotifyAll(action, userID, url string) {
	event := Event{
		Timestamp: time.Now().Unix(),
		Action:    action,
		UserID:    userID,
		URL:       url,
	}

	for _, observer := range s.observers {
		observer.Notify(event)
	}
}
