package audit

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Event struct {
	Timestamp int64  `json:"ts"`
	Action    string `json:"action"`
	UserID    string `json:"user_id"`
	URL       string `json:"url"`
}

type Observer interface {
	Notify(event Event)
}

type LogObserver struct {
	filePath string
}

func NewLogObserver(filePath string) *LogObserver {
	return &LogObserver{filePath: filePath}
}

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

type HTTPObserver struct {
	url string
}

func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{url: url}
}

func (o *HTTPObserver) Notify(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal audit event: %v", err)
		return
	}

	_, err = http.Post(o.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("failed to send audit event: %v", err)
	}
}

type AuditService struct {
	observers []Observer
}

func NewAuditService() *AuditService {
	return &AuditService{observers: []Observer{}}
}

func (s *AuditService) Register(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *AuditService) NotifyAll(action, userID, url string) {
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
