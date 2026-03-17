package aiprocessor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BleSSSeDDD/ai-internship-aggregator/gen/go/vacancy"
	"github.com/BleSSSeDDD/ai-internship-aggregator/internal/domain"
	"github.com/microcosm-cc/bluemonday"
)

const AI_RESPONSE_TIMEOUT = 60 * time.Minute

type ollamaResponse struct {
	Response string `json:"response"`
}

type aiProcessor struct {
	client    *http.Client
	ollamaURL string
	modelName string
}

func NewAiProcessor(ollamaURL, model string) domain.AIProcessor {
	return &aiProcessor{
		client:    &http.Client{Timeout: AI_RESPONSE_TIMEOUT},
		ollamaURL: ollamaURL,
		modelName: model,
	}
}

// получает HTML, отправляет в Ollama, возвращает структуру стажировки
func (p *aiProcessor) Process(ctx context.Context, html string, link string) (*vacancy.CompanyInternship, error) {
	log.Println("зашел в Process")

	cleanText := cleanHTML(html)
	log.Printf("HTML был %d символов, стал %d", len(html), len(cleanText))

	prompt := fmt.Sprintf(`Ты — парсер вакансий. Извлеки из HTML информацию о стажировке.

HTML:
%s

Верни ТОЛЬКО JSON в таком формате (без пояснений, только сам JSON):
{
  "company_name": "название компании",
  "source_url": "ссылка на страницу",
  "source_site": "%s",
  "position_name": "название позиции",
  "tech_stack": ["технология1", "технология2"],
  "min_salary": число (если не указано, ставь 0),
  "location": "город или Remote",
  "internship_dates": "сроки стажировки",
  "selection_process": "этапы отбора (одной строкой, через запятые)",
  "description": "описание задач",
  "application_deadline": "дедлайн подачи",
  "contact_info": "контакты",
  "experience_requirements": "требования к кандидату"
}`, cleanText, link)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  p.modelName,
		"prompt": prompt,
		"stream": false,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования запроса: %w", err)
	}

	log.Println("щас буду отправлять запрос")

	req, err := http.NewRequestWithContext(ctx, "POST", p.ollamaURL+"/api/generate", bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("ОШИБКА client.Do: %v", err)
		return nil, fmt.Errorf("ошибка вызова Ollama: %w", err)
	}
	defer resp.Body.Close()

	log.Println("p.client.Do(req) не вернул ошибку")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama вернула статус %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа Ollama: %w", err)
	}

	var internship vacancy.CompanyInternship
	if err := json.Unmarshal([]byte(ollamaResp.Response), &internship); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON от модели: %w\nОтвет модели: %s", err, ollamaResp.Response)
	}

	log.Println("выхожу из функции Process")

	return &internship, nil
}

func cleanHTML(html string) string {
	p := bluemonday.StrictPolicy()
	text := p.Sanitize(html)
	lines := strings.Split(text, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\n")
}
