package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/gen/go/vacancy"
)

type fakeParser struct {
	content string
	err     error
}

func (f *fakeParser) GetRawContent(ctx context.Context, url string) (string, error) {
	return f.content, f.err
}

type fakeAI struct {
	result  []*vacancy.CompanyInternship
	err     error
	gotText string
}

func (f *fakeAI) Process(ctx context.Context, text, link string) ([]*vacancy.CompanyInternship, error) {
	f.gotText = text
	return f.result, f.err
}

type fakePublisher struct {
	published []*vacancy.CompanyInternship
	err       error
	calls     int
}

func (f *fakePublisher) Publish(ctx context.Context, in []*vacancy.CompanyInternship) error {
	f.calls++
	f.published = in
	return f.err
}

func (f *fakePublisher) Close() error { return nil }

func TestRunPublishesExtractedInternships(t *testing.T) {
	internships := []*vacancy.CompanyInternship{
		{CompanyName: "Acme", PositionName: "Go Intern"},
		{CompanyName: "Acme", PositionName: "Frontend Intern"},
	}
	ai := &fakeAI{result: internships}
	pub := &fakePublisher{}

	u := NewScraperUsecase(&fakeParser{content: "<html>raw</html>"}, ai, pub)

	if err := u.Run(context.Background(), "https://example.com"); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}
	if ai.gotText != "<html>raw</html>" {
		t.Errorf("AI got %q, want raw page content", ai.gotText)
	}
	if pub.calls != 1 {
		t.Fatalf("Publish called %d times, want 1", pub.calls)
	}
	if len(pub.published) != 2 {
		t.Errorf("published %d internships, want 2", len(pub.published))
	}
}

func TestRunParserErrorSkipsPipeline(t *testing.T) {
	pub := &fakePublisher{}
	u := NewScraperUsecase(&fakeParser{err: errors.New("boom")}, &fakeAI{}, pub)

	if err := u.Run(context.Background(), "https://example.com"); err == nil {
		t.Fatal("Run() expected error, got nil")
	}
	if pub.calls != 0 {
		t.Errorf("Publish called %d times, want 0", pub.calls)
	}
}

func TestRunAIErrorSkipsPublish(t *testing.T) {
	pub := &fakePublisher{}
	u := NewScraperUsecase(&fakeParser{content: "x"}, &fakeAI{err: errors.New("bad json")}, pub)

	if err := u.Run(context.Background(), "https://example.com"); err == nil {
		t.Fatal("Run() expected error, got nil")
	}
	if pub.calls != 0 {
		t.Errorf("Publish called %d times, want 0", pub.calls)
	}
}

func TestRunNoInternshipsIsNotAnError(t *testing.T) {
	pub := &fakePublisher{}
	u := NewScraperUsecase(&fakeParser{content: "x"}, &fakeAI{result: nil}, pub)

	if err := u.Run(context.Background(), "https://example.com"); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}
	if pub.calls != 0 {
		t.Errorf("Publish called %d times for empty result, want 0", pub.calls)
	}
}

func TestRunPublisherErrorIsReturned(t *testing.T) {
	pub := &fakePublisher{err: errors.New("kafka down")}
	ai := &fakeAI{result: []*vacancy.CompanyInternship{{CompanyName: "Acme"}}}
	u := NewScraperUsecase(&fakeParser{content: "x"}, ai, pub)

	if err := u.Run(context.Background(), "https://example.com"); err == nil {
		t.Fatal("Run() expected error, got nil")
	}
}
