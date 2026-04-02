package aiprocessor

import (
	"context"

	"github.com/BleSSSeDDD/ai-internship-aggregator/gen/go/vacancy"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/domain"
)

type StubProcessor struct{}

func NewStub() domain.AIProcessor {
	return &StubProcessor{}
}

// Process возвращает один трек стажировки (теперь это CompanyInternship)
func (s *StubProcessor) Process(ctx context.Context, text string, link string) ([]*vacancy.CompanyInternship, error) {
	if text == "qwe" {

		return []*vacancy.CompanyInternship{&vacancy.CompanyInternship{
			CompanyName:            "Тинькофф",
			SourceUrl:              link,
			SourceSite:             "hh.ru",
			PositionName:           "Стажер-разработчик Go",
			TechStack:              []string{"Go", "PostgreSQL", "Kafka", "gRPC"},
			MinSalary:              70000,
			Location:               "Москва",
			InternshipDates:        "Июль-Сентябрь 2025",
			SelectionProcess:       "Тестовое задание → Собеседование с тимлидом",
			Description:            "Разработка высоконагруженных микросервисов",
			ApplicationDeadline:    "2025-05-15",
			ContactInfo:            "hr@tinkoff.ru",
			ExperienceRequirements: "Знание Go, понимание многопоточности, SQL",
		}}, nil
	}

	return []*vacancy.CompanyInternship{&vacancy.CompanyInternship{
		CompanyName:  "Тестовая Компания",
		SourceUrl:    "https://example.com/vacancy",
		SourceSite:   "hh.ru",
		PositionName: "Стажер Go",
		TechStack:    []string{"Go", "PostgreSQL", "Kafka"},
		MinSalary:    50000,
		Location:     "Москва",
	}}, nil
}
