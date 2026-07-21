package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/BleSSSeDDD/ai-internship-aggregator/manual-control/gen/go/vacancy"
	"github.com/gin-gonic/gin"
)

type fakePublisher struct {
	sent  *vacancy.CompanyInternship
	err   error
	calls int
}

func (f *fakePublisher) SendInternship(in *vacancy.CompanyInternship) (int32, int64, error) {
	f.calls++
	f.sent = in
	return 0, 42, f.err
}

func (f *fakePublisher) Close() error { return nil }

func newTestRouter(pub *fakePublisher) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.SetHTMLTemplate(template.Must(template.New("index.html").Parse("ok")))

	h := NewHandlers(pub)
	router.GET("/", h.Index)
	router.POST("/submit", h.Submit)
	router.GET("/health", h.Health)
	return router
}

func postForm(router *gin.Engine, form url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func validForm() url.Values {
	return url.Values{
		"company_name":  {"Acme"},
		"position_name": {"Go Intern"},
		"tech_stack":    {" Go , Kafka ,, gRPC "},
		"min_salary":    {"50000"},
		"location":      {"Remote"},
	}
}

func TestSubmitSendsInternshipToKafka(t *testing.T) {
	pub := &fakePublisher{}
	w := postForm(newTestRouter(pub), validForm())

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if pub.calls != 1 {
		t.Fatalf("SendInternship called %d times, want 1", pub.calls)
	}
	if pub.sent.CompanyName != "Acme" || pub.sent.MinSalary != 50000 {
		t.Errorf("unexpected internship sent: %+v", pub.sent)
	}
	wantStack := []string{"Go", "Kafka", "gRPC"}
	if !reflect.DeepEqual(pub.sent.TechStack, wantStack) {
		t.Errorf("TechStack = %v, want %v", pub.sent.TechStack, wantStack)
	}
}

func TestSubmitValidation(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(url.Values)
	}{
		{"missing company", func(f url.Values) { f.Del("company_name") }},
		{"missing position", func(f url.Values) { f.Del("position_name") }},
		{"empty tech stack", func(f url.Values) { f.Set("tech_stack", " , , ") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pub := &fakePublisher{}
			form := validForm()
			tt.mutate(form)

			w := postForm(newTestRouter(pub), form)

			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", w.Code)
			}
			if pub.calls != 0 {
				t.Errorf("SendInternship called %d times, want 0", pub.calls)
			}
		})
	}
}

func TestSubmitKafkaErrorReturns500(t *testing.T) {
	pub := &fakePublisher{err: errors.New("kafka down")}
	w := postForm(newTestRouter(pub), validForm())

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	newTestRouter(&fakePublisher{}).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Errorf("body = %s, want status ok", w.Body.String())
	}
}
