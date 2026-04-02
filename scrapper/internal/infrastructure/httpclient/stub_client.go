package httpclient

import (
	"context"

	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/domain"
)

// для мока вот такое есть:

type stubClient struct{}

func NewStubClient() domain.Parser {
	return &stubClient{}
}

func (p *stubClient) GetRawContent(ctx context.Context, url string) (string, error) {
	return string(`Стажировка в компании ООО "ИВАН"

Направления стажировки:

1. Бэкенд-разработка:
   - Технологии: Go, Kafka, gRPC
   - Зарплата: 50000 рублей
   - Локация: любой город или удаленно

2. Фронтенд-разработка:
   - Технологии: HTML, CSS, JavaScript
   - Зарплата: 50000 рублей
   - Локация: любой город или удаленно

3. Мобильная разработка:
   - Технологии: Kotlin
   - Зарплата: 50000 рублей
   - Локация: любой город или удаленно

Этапы отбора: тестовое задание, станцевать танец на шоу талантов`), nil
}
