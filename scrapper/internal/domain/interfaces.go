package domain

import (
	"context"

	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/gen/go/vacancy"
)

// Parser умеет зайти на сайт и достать сырой текст/HTML.
type Parser interface {
	GetRawContent(ctx context.Context, url string) (string, error)
}

// AIProcessor умеет превратить неструктурированный текст в массив стажировок.
type AIProcessor interface {
	Process(ctx context.Context, text string, link string) ([]*vacancy.CompanyInternship, error)
}

// Publisher умеет отправить готовые структуры в Kafka.
type Publisher interface {
	Publish(ctx context.Context, internships []*vacancy.CompanyInternship) error
	Close() error
}
