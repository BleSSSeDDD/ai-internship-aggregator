package domain

import "context"

// Parser - умеет зайти на сайт и достать сырой текст/HTML
type Parser interface {
	GetRawContent(ctx context.Context, url string) (string, error)
}

// AIProcessor - умеет превратить мусорный текст в структуру Internship
type AIProcessor interface {
	Process(ctx context.Context, text string) (*Internship, error)
}

// Publisher - умеет отправить готовую структуру в Кафку
type Publisher interface {
	Publish(ctx context.Context, internship *Internship) error
}
