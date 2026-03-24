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

// получает HTML, отправляет в Ollama, возвращает массив структур стажировок
func (p *aiProcessor) Process(ctx context.Context, html string, link string) ([]*vacancy.CompanyInternship, error) {
	log.Println("пошла ai-возня для " + link)

	cleanText := cleanHTML(html)
	log.Printf("HTML был %d символов, стал %d", len(html), len(cleanText))

	prompt := buildPrompt(cleanText, link)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  p.modelName,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.4,
			"top_p":       0.9,
		},
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	log.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama вернула статус %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа Ollama: %w", err)
	}

	var internships []*vacancy.CompanyInternship
	if err := json.Unmarshal([]byte(cleanModelResponse(ollamaResp.Response)), &internships); err != nil {
		var single vacancy.CompanyInternship
		if err := json.Unmarshal([]byte(cleanModelResponse(ollamaResp.Response)), &single); err != nil {
			return nil, fmt.Errorf("ошибка парсинга JSON от модели: %w\nОтвет модели: %s", err, ollamaResp.Response)
		}
		internships = []*vacancy.CompanyInternship{&single}
	}

	log.Printf("найдено %d стажировок", len(internships))
	log.Println("выхожу из функции Process")

	return internships, nil
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

// Удаляем ```json ... ``` обертки но вообще это не обязательно
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
