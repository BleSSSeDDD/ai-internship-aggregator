package aiprocessor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func fakeOllama(t *testing.T, modelResponse string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("request is not valid JSON: %v", err)
		}
		if req["model"] != "test-model" {
			t.Errorf("model = %v, want test-model", req["model"])
		}
		json.NewEncoder(w).Encode(map[string]string{"response": modelResponse})
	}))
}

func TestProcessParsesArrayResponse(t *testing.T) {
	srv := fakeOllama(t, `[
		{"company_name": "Acme", "position_name": "Go Intern", "tech_stack": ["Go", "Kafka"]},
		{"company_name": "Acme", "position_name": "Frontend Intern", "tech_stack": ["React"]}
	]`)
	defer srv.Close()

	p := NewAiProcessor(srv.URL, "test-model", time.Second)

	got, err := p.Process(context.Background(), "<html>vacancies</html>", "https://example.com")
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("Process() returned %d internships, want 2", len(got))
	}
	if got[0].CompanyName != "Acme" || got[0].PositionName != "Go Intern" {
		t.Errorf("unexpected first internship: %+v", got[0])
	}
}

func TestProcessAcceptsSingleObjectFallback(t *testing.T) {
	srv := fakeOllama(t, `{"company_name": "Acme", "position_name": "Go Intern"}`)
	defer srv.Close()

	p := NewAiProcessor(srv.URL, "test-model", time.Second)

	got, err := p.Process(context.Background(), "page", "https://example.com")
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}
	if len(got) != 1 || got[0].CompanyName != "Acme" {
		t.Errorf("Process() = %+v, want single Acme internship", got)
	}
}

func TestProcessStripsMarkdownFence(t *testing.T) {
	srv := fakeOllama(t, "```json\n[{\"company_name\": \"Acme\"}]\n```")
	defer srv.Close()

	p := NewAiProcessor(srv.URL, "test-model", time.Second)

	got, err := p.Process(context.Background(), "page", "https://example.com")
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("Process() returned %d internships, want 1", len(got))
	}
}

func TestProcessRejectsInvalidJSON(t *testing.T) {
	srv := fakeOllama(t, "Вот список стажировок: ...")
	defer srv.Close()

	p := NewAiProcessor(srv.URL, "test-model", time.Second)

	if _, err := p.Process(context.Background(), "page", "https://example.com"); err == nil {
		t.Fatal("Process() expected error for non-JSON output, got nil")
	}
}

func TestProcessOllamaErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "model not loaded", http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := NewAiProcessor(srv.URL, "test-model", time.Second)

	if _, err := p.Process(context.Background(), "page", "https://example.com"); err == nil {
		t.Fatal("Process() expected error for 500 from ollama, got nil")
	}
}

func TestCleanModelResponse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain json", `[{"a":1}]`, `[{"a":1}]`},
		{"json fence", "```json\n[{\"a\":1}]\n```", `[{"a":1}]`},
		{"bare fence", "```\n[{\"a\":1}]\n```", `[{"a":1}]`},
		{"surrounding whitespace", "  [{\"a\":1}]\n\n", `[{"a":1}]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanModelResponse(tt.in); got != tt.want {
				t.Errorf("cleanModelResponse(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCleanHTMLStripsTagsAndBlankLines(t *testing.T) {
	in := "<html><body>\n<h1>Стажировка</h1>\n\n  <p>Go, Kafka</p>\n<script>alert(1)</script>\n</body></html>"
	got := cleanHTML(in)

	if strings.Contains(got, "<") {
		t.Errorf("cleanHTML() left tags in output: %q", got)
	}
	if strings.Contains(got, "\n\n") {
		t.Errorf("cleanHTML() left blank lines: %q", got)
	}
	for _, want := range []string{"Стажировка", "Go, Kafka"} {
		if !strings.Contains(got, want) {
			t.Errorf("cleanHTML() lost text %q: %q", want, got)
		}
	}
}

func TestWaitForModelReturnsWhenModelPresent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(ollamaTagsResponse{Models: []ollamaModel{{Name: "test-model"}}})
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := WaitForModel(ctx, srv.URL, "test-model"); err != nil {
		t.Fatalf("WaitForModel() returned error: %v", err)
	}
}

func TestWaitForModelStopsOnContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ollamaTagsResponse{}) // модели нет
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := WaitForModel(ctx, srv.URL, "test-model"); err == nil {
		t.Fatal("WaitForModel() expected error for cancelled context, got nil")
	}
}
