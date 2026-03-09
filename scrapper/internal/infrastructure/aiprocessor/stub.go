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

// нуу пока так чтоб грустно не было легендарному джависту
func (s *StubProcessor) Process(ctx context.Context, text string) (*vacancy.CompanyInternship, error) {
	if text == "qwe" {
		return &vacancy.CompanyInternship{
			CompanyName: "Тинькофф",
			SourceUrl:   "https://hh.ru/vacancy/123456",
			SourceSite:  "hh.ru",
			Tracks: []*vacancy.Track{
				{
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
				},
				{
					PositionName:           "Стажер-разработчик Java",
					TechStack:              []string{"Java", "Spring Boot", "PostgreSQL", "Kafka"},
					MinSalary:              70000,
					Location:               "Москва",
					InternshipDates:        "Июль-Сентябрь 2025",
					SelectionProcess:       "Тестовое задание → Собеседование с командой",
					Description:            "Разработка бэкенда на Spring Boot",
					ApplicationDeadline:    "2025-05-15",
					ContactInfo:            "hr@tinkoff.ru",
					ExperienceRequirements: "Java Core, SQL, основы Spring",
				},
			},
		}, nil
	}
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
