package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func newTestParser() *vacancyParser {
	return &vacancyParser{
		client:     &http.Client{Timeout: time.Second},
		limiter:    time.NewTicker(time.Millisecond).C,
		retryDelay: time.Millisecond,
	}
}

func TestGetRawContentSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("page content"))
	}))
	defer srv.Close()

	got, err := newTestParser().GetRawContent(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("GetRawContent() returned error: %v", err)
	}
	if got != "page content" {
		t.Errorf("GetRawContent() = %q, want %q", got, "page content")
	}
}

func TestGetRawContentRetriesOn5xx(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	got, err := newTestParser().GetRawContent(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("GetRawContent() returned error: %v", err)
	}
	if got != "ok" {
		t.Errorf("GetRawContent() = %q, want %q", got, "ok")
	}
	if calls.Load() != 3 {
		t.Errorf("server got %d requests, want 3", calls.Load())
	}
}

func TestGetRawContentDoesNotRetryOn404(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := newTestParser().GetRawContent(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("GetRawContent() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "non-retryable") {
		t.Errorf("error = %v, want non-retryable", err)
	}
	if calls.Load() != 1 {
		t.Errorf("server got %d requests, want 1", calls.Load())
	}
}

func TestGetRawContentGivesUpAfterMaxRetries(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	_, err := newTestParser().GetRawContent(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("GetRawContent() expected error, got nil")
	}
	if calls.Load() != maxRetries {
		t.Errorf("server got %d requests, want %d", calls.Load(), maxRetries)
	}
}

func TestGetRawContentRespectsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := newTestParser().GetRawContent(ctx, "http://127.0.0.1:0")
	if err == nil {
		t.Fatal("GetRawContent() expected error for cancelled context, got nil")
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		want   bool
	}{
		{"network error", context.DeadlineExceeded, 0, true},
		{"500", nil, 500, true},
		{"503", nil, 503, true},
		{"429", nil, 429, true},
		{"404", nil, 404, false},
		{"403", nil, 403, false},
		{"200", nil, 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRetry(tt.err, tt.status); got != tt.want {
				t.Errorf("shouldRetry(%v, %d) = %v, want %v", tt.err, tt.status, got, tt.want)
			}
		})
	}
}
