package aiprocessor

import (
	"context"

	vacancy "github.com/BleSSSeDDD/reviewer-assignment/generated"
	"github.com/BleSSSeDDD/reviewer-assignment/internal/domain"
)

type StubProcessor struct{}

func NewStub() domain.AIProcessor {
	return &StubProcessor{}
}

func (s *StubProcessor) Process(ctx context.Context, text string) (*vacancy.CompanyInternship, error) {
	return &vacancy.CompanyInternship{
		CompanyName: "Тестовая Компания",
		SourceUrl:   "https://example.com/vacancy",
		SourceSite:  "hh.ru",
		Tracks: []*vacancy.Track{
			{
				PositionName: "Стажер Go",
				TechStack:    []string{"Go", "PostgreSQL", "Kafka"},
				MinSalary:    50000,
				Location:     "Москва",
			},
		},
	}, nil
}
