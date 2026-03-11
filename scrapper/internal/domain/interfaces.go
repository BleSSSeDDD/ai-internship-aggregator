package domain

import (
	"context"

	"github.com/BleSSSeDDD/ai-internship-aggregator/gen/go/vacancy"
)

// Parser - умеет зайти на сайт и достать сырой текст/HTML
type Parser interface {
	GetRawContent(ctx context.Context, url string) (string, error)
}

// AIProcessor - умеет превратить мусорный текст в структуру Internship
type AIProcessor interface {
	Process(ctx context.Context, text string) (*vacancy.CompanyInternship, error)
}

// Publisher - умеет отправить готовую структуру в Кафку
type Publisher interface {
	Publish(ctx context.Context, internship *vacancy.CompanyInternship) error
	Close() error
}
