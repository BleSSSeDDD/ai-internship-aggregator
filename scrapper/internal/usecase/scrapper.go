package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/internal/domain"
)

// ScraperUsecase — конвейер: скачать страницу → извлечь стажировки через LLM → отправить в Kafka.
type ScraperUsecase struct {
	parser    domain.Parser
	ai        domain.AIProcessor
	publisher domain.Publisher
}

func NewScraperUsecase(p domain.Parser, ai domain.AIProcessor, pub domain.Publisher) *ScraperUsecase {
	return &ScraperUsecase{parser: p, ai: ai, publisher: pub}
}

func (u *ScraperUsecase) Run(ctx context.Context, link string) error {
	text, err := u.parser.GetRawContent(ctx, link)
	if err != nil {
		return fmt.Errorf("fetch page: %w", err)
	}

	data, err := u.ai.Process(ctx, text, link)
	if err != nil {
		return fmt.Errorf("extract internships: %w", err)
	}

	if len(data) == 0 {
		slog.Info("no internships found", "url", link)
		return nil
	}

	if err := u.publisher.Publish(ctx, data); err != nil {
		return fmt.Errorf("publish internships: %w", err)
	}

	return nil
}
