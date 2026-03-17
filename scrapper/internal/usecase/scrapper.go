package usecase

import (
	"context"

	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/domain"
)

type ScraperUsecase struct {
	parser    domain.Parser
	ai        domain.AIProcessor
	publisher domain.Publisher
}

func NewScraperUsecase(p domain.Parser, ai domain.AIProcessor, pub domain.Publisher) *ScraperUsecase {
	return &ScraperUsecase{p, ai, pub}
}

func (u *ScraperUsecase) Run(ctx context.Context, link string) error {
	text, err := u.parser.GetRawContent(ctx, link)
	if err != nil {
		return err
	}

	data, err := u.ai.Process(ctx, text, link)
	if err != nil {
		return err
	}

	return u.publisher.Publish(ctx, data)
}
