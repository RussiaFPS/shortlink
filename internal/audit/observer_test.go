package audit

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

type MockObserver struct {
	mu    sync.Mutex
	event Event
}

func (m *MockObserver) Notify(event Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.event = event
}

func TestServAudit(t *testing.T) {
	auditService := NewAuditService()
	mockObserver := &MockObserver{}
	auditService.Register(mockObserver)

	action := "test_action"
	userID := "test_user"
	url := "test_url"

	auditService.NotifyAll(action, userID, url)

	mockObserver.mu.Lock()
	defer mockObserver.mu.Unlock()

	assert.Equal(t, action, mockObserver.event.Action)
	assert.Equal(t, userID, mockObserver.event.UserID)
	assert.Equal(t, url, mockObserver.event.URL)
	assert.WithinDuration(t, time.Now(), time.Unix(mockObserver.event.Timestamp, 0), time.Second)
}

func TestLogObserver(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "audit_test")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	logObserver := NewLogObserver(tmpfile.Name())
	event := Event{
		Timestamp: time.Now().Unix(),
		Action:    "test_log",
		UserID:    "user_log",
		URL:       "url_log",
	}
	logObserver.Notify(event)

	data, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(t, err)

	var loggedEvent Event
	err = json.Unmarshal(data, &loggedEvent)
	assert.NoError(t, err)

	assert.Equal(t, event.Action, loggedEvent.Action)
	assert.Equal(t, event.UserID, loggedEvent.UserID)
	assert.Equal(t, event.URL, loggedEvent.URL)
	assert.Equal(t, event.Timestamp, loggedEvent.Timestamp)
}

func TestHTTPObserver(t *testing.T) {
	var receivedEvent Event
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedEvent)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpObserver := NewHTTPObserver(server.URL)
	event := Event{
		Timestamp: time.Now().Unix(),
		Action:    "test_http",
		UserID:    "user_http",
		URL:       "url_http",
	}
	httpObserver.Notify(event)

	assert.Equal(t, event.Action, receivedEvent.Action)
	assert.Equal(t, event.UserID, receivedEvent.UserID)
	assert.Equal(t, event.URL, receivedEvent.URL)
	assert.Equal(t, event.Timestamp, receivedEvent.Timestamp)
}
