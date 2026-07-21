package aiprocessor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/gen/go/vacancy"
	"github.com/BleSSSeDDD/ai-internship-aggregator/scrapper/internal/domain"
	"github.com/microcosm-cc/bluemonday"
)

type ollamaResponse struct {
	Response string `json:"response"`
}

type aiProcessor struct {
	client    *http.Client
	ollamaURL string
	modelName string
}

func NewAiProcessor(ollamaURL, model string, timeout time.Duration) domain.AIProcessor {
	return &aiProcessor{
		client:    &http.Client{Timeout: timeout},
		ollamaURL: ollamaURL,
		modelName: model,
	}
}

// Process очищает HTML от разметки, отправляет текст в Ollama
// и разбирает ответ модели в структуры стажировок.
func (p *aiProcessor) Process(ctx context.Context, html string, link string) ([]*vacancy.CompanyInternship, error) {
	cleanText := cleanHTML(html)
	slog.Info("sending page to LLM",
		"url", link,
		"html_len", len(html),
		"text_len", len(cleanText),
	)

	requestBody, err := json.Marshal(map[string]any{
		"model":  p.modelName,
		"prompt": buildPrompt(cleanText, link),
		"stream": false,
		"options": map[string]any{
			"temperature": 0.4,
			"top_p":       0.9,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal ollama request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.ollamaURL+"/api/generate", bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("build ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read ollama response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("decode ollama response: %w", err)
	}

	internships, err := parseModelOutput(ollamaResp.Response)
	if err != nil {
		return nil, err
	}

	slog.Info("LLM extracted internships",
		"url", link,
		"count", len(internships),
		"took", time.Since(start).Round(time.Second),
	)

	return internships, nil
}

// parseModelOutput ожидает JSON-массив, но модель иногда возвращает одиночный объект.
func parseModelOutput(response string) ([]*vacancy.CompanyInternship, error) {
	cleaned := cleanModelResponse(response)

	var internships []*vacancy.CompanyInternship
	if err := json.Unmarshal([]byte(cleaned), &internships); err == nil {
		return internships, nil
	}

	var single vacancy.CompanyInternship
	if err := json.Unmarshal([]byte(cleaned), &single); err != nil {
		return nil, fmt.Errorf("model output is not valid JSON: %w\nmodel output: %s", err, response)
	}
	return []*vacancy.CompanyInternship{&single}, nil
}

func cleanHTML(html string) string {
	text := bluemonday.StrictPolicy().Sanitize(html)

	var result []string
	for _, line := range strings.Split(text, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\n")
}

func buildPrompt(cleanText, link string) string {
	return fmt.Sprintf(`Ты — строгий парсер направлений стажировок. Извлеки информацию о направлениях стажировок из HTML.

КРИТИЧЕСКИ ВАЖНО: Если на странице упоминаются разные направления (бэкенд, фронтенд, мобильная разработка, devops и т.д.),
ты ОБЯЗАН создать ОТДЕЛЬНЫЙ объект в массиве для КАЖДОГО направления.

Правила:
1. Верни ТОЛЬКО валидный JSON массив
2. Если есть несколько направлений — каждый в отдельном объекте
3. tech_stack должен содержать ТОЛЬКО технологии для этого конкретного направления
4. position_name должен отражать направление (например: "Стажер-бэкенд разработчик")
5. НЕ объединяй разные направления в один объект

Пример правильного ответа для страницы с бэкендом и фронтендом:
[
  {
    "company_name": "имя компании",
    "source_url": "%s",
    "source_site": "%s",
    "position_name": "должность",
    "tech_stack": ["...", "..."],
    "min_salary": если не указано, то 0,
    "location": "город или удаленно",
    "internship_dates": "",
    "selection_process": "строкой через запятую все этапы",
    "description": "описание",
    "application_deadline": "",
    "contact_info": "",
    "experience_requirements": ""
  },
  {
    "company_name": "...",
    "source_url": "%s",
    "source_site": "%s",
    "position_name": "...",
    "tech_stack": ["...", "..."],
    "min_salary": ...,
    "location": "...",
    "internship_dates": "",
    "selection_process": "",
    "description": "",
    "application_deadline": "",
    "contact_info": "",
    "experience_requirements": ""
  }
]

HTML для парсинга:
%s

ВЕРНИ ТОЛЬКО JSON МАССИВ, НИЧЕГО ДРУГОГО.`, link, link, link, link, cleanText)
}

// cleanModelResponse срезает обёртку ```json ... ```, которую модель иногда добавляет.
func cleanModelResponse(response string) string {
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
	}
	return strings.TrimSpace(response)
}
